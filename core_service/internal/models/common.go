package models

type Company struct {
	Uuid   string `json:"uuid"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

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

type FileContentType string

const (
	SOMETHING FileContentType = "something"
	IMAGE     FileContentType = "image"
	DOC       FileContentType = "doc"
	AUDIO     FileContentType = "audio"
	VIDEO     FileContentType = "video"
)

type File struct {
	Filename    string          `json:"filename"`
	FileId      string          `json:"file_id"`
	Url         string          `json:"url"`
	GreenApiUrl string          `json:"user_api_url"`
	MimeType    *string         `json:"mime_type"`
	ContentType FileContentType `json:"content_type"`
	Size        int64           `json:"size"`
	Data        []byte          `json:"data"`
	IsPreview   bool            `json:"is_preview"`
}

type PresignUrl struct {
	Method      string `json:"method"`
	ObjectId    string `json:"object_id"`
	Url         string `json:"url"`
	XAmzTagging string `json:"x_amz_tagging"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ExtSession struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
}
