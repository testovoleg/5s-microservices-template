package models

type Company struct {
	Uuid    string   `json:"uuid"`
	Name    string   `json:"name"`
	Active  bool     `json:"active"`
	ApiData *ApiData `json:"api_data"`
}

type ApiData struct {
	ApiUrl   string `json:"api_url"`
	Version  string `json:"version"`
	Username string `json:"username"`
	Password string `json:"password"`
}
