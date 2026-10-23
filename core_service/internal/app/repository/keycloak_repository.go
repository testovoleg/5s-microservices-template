package repository

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"

	"github.com/Nerzal/gocloak/v13"
	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	"github.com/testovoleg/5s-microservice-template/pkg/constants"
	httpErrors "github.com/testovoleg/5s-microservice-template/pkg/http_errors"
	httpUtils "github.com/testovoleg/5s-microservice-template/pkg/http_utils"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

type idmRepository struct {
	log    logger.Logger
	cfg    *config.Config
	client *gocloak.GoCloak
}

func NewIDMRepository(log logger.Logger, cfg *config.Config, keycloakClient *gocloak.GoCloak) *idmRepository {
	return &idmRepository{log: log, cfg: cfg, client: keycloakClient}
}

func (p *idmRepository) UserData(ctx context.Context, accessToken string) (*models.User, error) {
	ctx, span := tracing.StartSpan(ctx, "idmRepository.UserData")
	defer span.End()
	if user, err := p.decodeExtToken(ctx, accessToken); err == nil {
		return user, nil
	}

	active, err := p.introspectAccessToken(ctx, accessToken)
	if err != nil {
		return nil, errors.New("Invalid token")
	}
	if !active {
		return nil, errors.New("Invalid or expired token")
	}

	user, err := p.DecodeToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// introspectAccessToken calls Keycloak's introspect endpoint for a regular OAuth2 access token.
// gocloak's RetrospectToken hardcodes token_type_hint=requesting_party_token (UMA RPT), which makes
// Keycloak report active:false for plain access tokens — introspect manually with the right hint.
func (p *idmRepository) introspectAccessToken(ctx context.Context, accessToken string) (bool, error) {
	tokenURL := strings.TrimRight(p.cfg.Keycloak.Host, "/") +
		"/realms/" + p.cfg.Keycloak.Realm + "/protocol/openid-connect/token/introspect"

	form := url.Values{
		"token":           {accessToken},
		"token_type_hint": {"access_token"},
	}
	body, err := httpUtils.SendHTTPRequest(ctx, http.MethodPost, tokenURL, []byte(form.Encode()), httpUtils.SendOpts{
		BasicUser:     p.cfg.Keycloak.ClientID,
		BasicPassword: p.cfg.Keycloak.ClientSecret,
		ContentType:   "application/x-www-form-urlencoded",
		ErrorPrefix:   "Keycloak introspect",
	})
	if err != nil {
		return false, err
	}

	var result struct {
		Active bool `json:"active"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, err
	}
	if !result.Active {
		p.log.Warnf("Keycloak introspect inactive: url=%s raw=%s claims=%s", tokenURL, string(body), unverifiedTokenClaims(accessToken))
	}
	return result.Active, nil
}

// unverifiedTokenClaims decodes a JWT payload WITHOUT verifying the signature — diagnostics only,
// to compare iss/azp/sid/exp against the realm/client used for introspection.
func unverifiedTokenClaims(token string) string {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "not a JWT"
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "decode error: " + err.Error()
	}
	var claims struct {
		Iss string `json:"iss"`
		Azp string `json:"azp"`
		Sid string `json:"sid"`
		Exp int64  `json:"exp"`
		Iat int64  `json:"iat"`
		Aud any    `json:"aud"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "unmarshal error: " + err.Error()
	}
	b, _ := json.Marshal(claims)
	return string(b)
}

func (p *idmRepository) DecodeToken(ctx context.Context, accessToken string) (*models.User, error) {
	ctx, span := tracing.StartSpan(ctx, "idmRepository.DecodeToken")
	defer span.End()

	tokenData, claims, err := p.client.DecodeAccessToken(ctx, accessToken, p.cfg.Keycloak.Realm)
	if err != nil {
		return nil, httpErrors.ErrorWrap(err, span, "Decode token failed")
	}
	if !tokenData.Valid {
		return nil, errors.New("Invalid token")
	}

	userUuid, err := getClaim(claims, "user_uuid")
	if err != nil {
		return nil, err
	}
	if userUuid == "" {
		return nil, errors.New("Invalid token")
	}

	cloakId, err := getClaim(claims, "sub")
	if err != nil {
		return nil, err
	}
	if cloakId == "" {
		return nil, errors.New("Invalid token")
	}

	name, err := getClaim(claims, "preferred_username")
	if err != nil {
		return nil, err
	}

	if name == "" {
		name, err = getClaim(claims, "name")
		if err != nil {
			return nil, err
		}
	}

	userEmail, err := getClaim(claims, "email")
	if err != nil {
		return nil, err
	}

	companyUuid, _ := getClaim(claims, "company_uuid")
	// if err != nil {
	// 	return nil, err
	// }

	userRoles, err := getRoles(claims)
	if err != nil {
		return nil, err
	}
	if len(userRoles) == 0 {
		return nil, errors.New("")
	}

	return &models.User{
		Id:      userUuid,
		CloakId: cloakId,
		Name:    models.UserName{FullName: name},
		Contacts: models.UserContacts{
			Email: []string{userEmail},
		},
		Roles:   userRoles,
		Company: &models.Company{Uuid: companyUuid},
	}, nil
}

func (p *idmRepository) getAttributeInKeycloakUser(user *gocloak.User, attribute_name string) []string {
	if user == nil || user.Attributes == nil {
		return nil
	}

	val, ok := (*user.Attributes)[attribute_name]
	if !ok { // if not finded with this attribute name
		return nil
	} else {
		if len(val) > 0 { // if have some variables, get first
			return val
		} else { // if is just empty array
			return nil
		}
	}
}

func (p *idmRepository) IsExternalUser(u *models.User) bool {
	if u == nil {
		return false
	}
	return u.External
}

func (p *idmRepository) IsWebservice(u *models.User) bool {
	if u == nil {
		return false
	}
	for _, v := range u.Roles {
		if v.Name == constants.RoleWebservice {
			return true
		}
	}
	return false
}

func (p *idmRepository) IsAdministrator(u *models.User) bool {
	if u == nil {
		return false
	}
	for _, v := range u.Roles {
		if v.Name == constants.RoleAdministrator {
			return true
		}
	}
	return false
}

func (p *idmRepository) IsSuperuser(u *models.User) bool {
	if u == nil {
		return false
	}
	for _, v := range u.Roles {
		if v.Name == constants.RoleSuperuser {
			return true
		}
	}
	return false
}

func (r *idmRepository) GetAdminToken(ctx context.Context) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "idmRepository.GetAdminToken")
	defer span.End()

	requestURL := r.cfg.API.AuthApiUrl + "/oidc/login"

	jsonBody, err := json.Marshal(models.OidcLoginReq{
		Username: r.cfg.API.ApiUsername,
		Password: r.cfg.API.ApiPassword,
	})
	if err != nil {
		return "", httpErrors.ErrorWrap(err, span, "json.Marshal")
	}

	b, err := httpUtils.SendHTTPRequest(ctx, http.MethodPost, requestURL, jsonBody, httpUtils.SendOpts{
		Timeout:     30 * time.Second,
		ErrorPrefix: "Error was during on Admin API",
	})
	if err != nil {
		return "", err
	}

	var body models.OidcLoginRes
	if err = json.Unmarshal(b, &body); err != nil {
		return "", httpErrors.ErrorWrap(err, span, "json.Unmarshal")
	}

	return body.AccessToken, nil
}

func getClaim(c *jwt.MapClaims, name string) (string, error) {
	if c == nil {
		return "", errors.New("Invalid token")
	}

	var res string

	v, ok := (*c)[name]
	if !ok {
		return "", errors.New("Invalid token, not field 'sub'")
	}
	res, ok = v.(string)
	if !ok {
		return "", errors.New("Invalid token, uncorrect field 'sub'")
	}
	if res == "" {
		return "", errors.New("Invalid token")
	}

	return res, nil
}

func getRoles(c *jwt.MapClaims) ([]models.UserRole, error) {
	var res []models.UserRole

	if c == nil {
		return nil, errors.New("Invalid token")
	}

	v, ok := (*c)["realm_access"]
	if !ok {
		return nil, errors.New("Invalid token, not field 'sub'")
	}
	v2, ok := v.(map[string]interface{})
	if !ok {
		return nil, errors.New("Invalid token, uncorrect field 'sub'")
	}
	v3, ok := v2["roles"]
	if !ok {
		return nil, errors.New("Invalid token, uncorrect field 'sub'")
	}
	v4, ok := v3.([]interface{})
	if !ok {
		return nil, errors.New("Invalid token, uncorrect field 'sub'")
	}
	if len(v4) == 0 {
		return nil, errors.New("Invalid token, roles are ampty")
	}
	for _, v5 := range v4 {
		v6, ok := v5.(string)
		if !ok {
			return nil, errors.New("Invalid token")
		}
		r := models.UserRole{Name: v6}
		res = append(res, r)
	}

	if len(res) == 0 {
		return nil, errors.New("Invalid token")
	}

	return res, nil
}

func (p *idmRepository) decodeExtToken(ctx context.Context, accessToken string) (user *models.User, err error) {
	pattern := regexp.MustCompile(`(^.{36}$)|(^.{36}):(.+)`)
	finded := pattern.FindStringSubmatch(accessToken)
	if len(finded) == 0 {
		return nil, errors.New("Not external user")
	}
	var token, key string
	token = finded[1]
	if len(finded) >= 2 {
		key = finded[2]
	}
	return p.getExtUser(ctx, token, key)
}

func (p *idmRepository) getExtUser(ctx context.Context, token, key string) (*models.User, error) {
	ctx, span := tracing.StartSpan(ctx, "idmRepository.getExtUser")
	defer span.End()

	requestURL := p.cfg.API.AuthApiUrl + "/ext/self-data"

	reqToken := token
	if key != "" {
		reqToken = reqToken + ":" + key
	}

	status, b, err := httpUtils.SendHTTPRequestRaw(ctx, http.MethodGet, requestURL, nil, httpUtils.SendOpts{
		BearerToken: reqToken,
		Timeout:     30 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, errors.New("Error was during on get self-data token. The token may have expired")
	}

	var body models.ExtSelfDataRes
	if err = json.Unmarshal(b, &body); err != nil {
		return nil, httpErrors.ErrorWrap(err, span, "json.Unmarshal(")
	}

	user := &models.User{
		Id: body.UserId,
		Company: &models.Company{
			Uuid: body.CompanyUuid,
		},
		External: true,
	}
	return user, nil
}
