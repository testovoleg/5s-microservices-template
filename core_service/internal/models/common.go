package models

type Company struct {
	Name        string `json:"name"`
	CompanyUUID string `json:"uuid"`
	Active      bool   `json:"active"`
}

type User struct {
	UserUUID  string
	CloakUUID string
	Username  string
	External  bool
	Company   *Company
	Roles     []UserRole
	Active    bool
}

type UserRole struct {
	Name string `json:"name"`
}

type Client struct {
	ID         string
	FirstName  string
	MiddleName string
	LastName   string
	FullName   string
	Phone      []string
	VkID       []string
}

type ExtSession struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
}
