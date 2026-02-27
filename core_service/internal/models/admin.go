package models

type ApiParams struct {
	AccessToken string `json:"access_token"`
	CompanyUuid string `json:"company_uuid"`
	ApiUuid     string `json:"api_uuid"`
	IdmUserUuid string `json:"idm_user_uuid"`
}

type Api struct {
	Uuid        string `json:"api_uuid"`
	CompanyUuid string `json:"company_uuid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Token       string `json:"token"`
}

type ApiFull struct {
	Uuid          string   `json:"api_uuid"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Token         string   `json:"token"`
	Subscriptions []string `json:"subscriptions"`
}

type UpdateApi struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Token       *string `json:"token"`
}
