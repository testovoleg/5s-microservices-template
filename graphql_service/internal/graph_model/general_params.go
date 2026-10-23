package model

// GeneralParamsInput holds common GraphQL request params.
// AccessToken is filled from the Authorization header only (not part of the GraphQL schema).
type GeneralParamsInput struct {
	AccessToken *string `json:"-"`
	CompanyUUID string  `json:"companyUuid"`
	APIUUID     *string `json:"apiUuid,omitempty"`
	IdmUserUUID *string `json:"idmUserUuid,omitempty"`
}
