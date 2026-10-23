package dto

type ApiDto struct {
	Uuid        string `json:"api_uuid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Token       bool   `json:"token"`
}

type ApiFullDto struct {
	Uuid          string   `json:"api_uuid"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Token         string   `json:"token"`
	Subscriptions []string `json:"subscriptions"`
}

type ApiParamsDto struct {
	AccessToken string `json:"access_token"`
	CompanyUuid string `json:"company_uuid"`
	ApiUuid     string `json:"api_uuid"`
	IdmUserUuid string `json:"idm_user_uuid"`
}

type AddApiReqDto struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Token       string `json:"token"`
}

type GetApiReqDto struct {
	CompanyUuid string `query:"company_uuid"`
}

type GetFullApiReqDto struct {
	CompanyUuid string `query:"company_uuid"`
	ApiUuid     string `param:"API_UUID"`
}

type UpdateApiReqDto struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Token       *string `json:"token"`
}

type DeleteApiReqDto struct {
	CompanyUuid string `query:"company_uuid"`
	ApiUuid     string `param:"API_UUID"`
}

type UpdateWebhookResDto struct {
	Message string `json:"message"`
}

type EventResDto struct {
	EventUuid string `json:"event_uuid,omitempty"`
}

/*
	Swagger/JSON: protobuf empty repeated fields are nil and encode as null.
	Keep omitempty off the slice and map with utils.NonNilSlice:

		Subscriptions: utils.NonNilSlice(in.GetSubscriptions()),
*/
