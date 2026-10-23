package models

type ExtSession struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
}

type ExtSessionReq struct {
	CompanyUuid string `json:"company_uuid"`
	ClientId    string `json:"client_id"`
}

type OidcLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type OidcLoginRes struct {
	AccessToken string `json:"access_token"`
}

type ExtSelfDataRes struct {
	UserId      string `json:"client_id"`
	CompanyUuid string `json:"company_uuid"`
}
