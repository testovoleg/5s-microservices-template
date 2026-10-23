package models

type UserRole struct {
	Name string `json:"name"`
}

type UserName struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	FullName   string `json:"full_name"`
}

type UserContacts struct {
	Phone []string `json:"phone"`
	MaxId []string `json:"max_id"`
	Email []string `json:"email"`
}

type User struct {
	Id       string       `json:"id"`
	Company  *Company     `json:"-"`
	Active   bool         `json:"active"`
	Roles    []UserRole   `json:"roles"`
	CloakId  string       `json:"cloak_id"`
	External bool         `json:"external"`
	Name     UserName     `json:"name"`
	Birthday string       `json:"birthday"`
	Contacts UserContacts `json:"contacts"`
	Sex      string       `json:"sex"`
	Employee bool         `json:"employee"`
	Photo    []string     `json:"photo"`
}
