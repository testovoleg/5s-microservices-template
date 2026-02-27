package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"

	"github.com/Nerzal/gocloak/v13"
	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	"github.com/testovoleg/5s-microservice-template/pkg/constants"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"
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
	if trace.SpanContextFromContext(ctx).IsValid() {
		var span trace.Span
		ctx, span = tracing.StartSpan(ctx, "idmRepository.UserData")
		defer span.End()
	}
	if user, err := p.decodeExtToken(ctx, accessToken); err == nil {
		return user, nil
	}

	result, err := p.client.RetrospectToken(ctx, accessToken, p.cfg.Keycloak.ClientID, p.cfg.Keycloak.ClientSecret, p.cfg.Keycloak.Realm)
	if err != nil {
		return nil, errors.New("Invalid token")
	}
	if !*result.Active {
		return nil, errors.New("Invalid or expired token")
	}

	user, err := p.DecodeToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (p *idmRepository) DecodeToken(ctx context.Context, accessToken string) (*models.User, error) {
	if trace.SpanContextFromContext(ctx).IsValid() {
		var span trace.Span
		ctx, span = tracing.StartSpan(ctx, "idmRepository.DecodeToken")
		defer span.End()
	}
	tokenData, claims, err := p.client.DecodeAccessToken(ctx, accessToken, p.cfg.Keycloak.Realm)
	if err != nil {
		return nil, errors.Wrap(err, "Decode token failed")
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

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	requestURL := r.cfg.API.AuthApiUrl + "/oidc/login"

	reqBody := map[string]interface{}{}
	reqBody["username"] = r.cfg.API.ApiUsername
	reqBody["password"] = r.cfg.API.ApiPassword

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", errors.Wrap(err, "json.Marshal")
	}
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, requestURL, bodyReader)
	if err != nil {
		return "", errors.Wrap(err, "http.NewRequest")
	}

	req.Header.Set("Content-Type", "application/json")
	tracingHeaders, _ := tracing.InjectTextMapCarrier(ctx)
	for i, k := range tracingHeaders {
		req.Header.Set(i, k)
	}

	res, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "client.Do")
	}

	if res.StatusCode != http.StatusOK {
		return "", utils.NewError(res, span, "Error was during on Admin API")
	}

	b, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return "", errors.Wrap(err, "io.ReadAll")
	}

	type ResponseBody struct {
		AccessToken string `json:"access_token"`
	}
	var body ResponseBody

	err = json.Unmarshal(b, &body)
	if err != nil {
		return "", errors.Wrap(err, "json.Unmarshal")
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
	if trace.SpanContextFromContext(ctx).IsValid() {
		var span trace.Span
		ctx, span = tracing.StartSpan(ctx, "idmRepository.getExtUser")
		defer span.End()
	}
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	requestURL := p.cfg.API.AuthApiUrl + "/ext/self-data"

	req_token := token
	if key != "" {
		req_token = req_token + ":" + key
	}

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest")
	}

	req.Header.Set("Authorization", "Bearer "+req_token)
	tracingHeaders, _ := tracing.InjectTextMapCarrier(ctx)
	for i, k := range tracingHeaders {
		req.Header.Set(i, k)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client.Do")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("Error was during on get self-data token. The token may have expired")
	}

	b, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadAll")
	}

	type ResponseBody struct {
		UserId      string `json:"client_id"`
		CompanyUuid string `json:"company_uuid"`
	}
	var body ResponseBody

	err = json.Unmarshal(b, &body)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal(")
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
