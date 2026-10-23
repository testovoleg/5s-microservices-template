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

type PropertyData struct {
	List []*Api `json:"api_list"`
}

type PropertyValueReq struct {
	Value string `json:"value"`
}

type AdminUserRes struct {
	Id      string   `json:"uuid"`
	Company *Company `json:"company"`
	Active  bool     `json:"active"`
	Roles   []string `json:"roles"`
	Name    string   `json:"username"`
}

type AdminCompanyRes struct {
	Uuid   string `json:"uuid"`
	Name   string `json:"name"`
	Active *bool  `json:"active"`
}

// AdminCompanyList mirrors the pagination envelope returned by GET /company/list, e.g.:
//
//	{"companies": [{"active": true, "name": "...", "uuid": "..."}],
//	 "hasMore": true, "page": 1, "size": 100, "totalCount": 0, "totalPages": 0}
type AdminCompanyList struct {
	Companies  []AdminCompanyRes `json:"companies"`
	HasMore    bool              `json:"hasMore"`
	TotalPages int               `json:"totalPages"`
}

const (
	AdminCompanyListPageSize = 100
	// AdminCompanyListMaxPages caps the walk so a misbehaving Admin API (hasMore stuck at
	// true) can't turn this into an infinite loop.
	AdminCompanyListMaxPages = 1000
)
