package models

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

type StorageObjectReq struct {
	Bucket      string `json:"bucket"`
	CompanyUuid string `json:"company_uuid"`
	ObjectId    string `json:"object_id"`
}

type StoragePresignUploadReq struct {
	Bucket      string   `json:"bucket"`
	CompanyUuid string   `json:"company_uuid"`
	Filename    string   `json:"filename"`
	Prefixs     []string `json:"prefixs"`
	Tags        []*Tag   `json:"tags"`
}

type StorageSearchReq struct {
	Bucket      string `json:"bucket"`
	CompanyUuid string `json:"company_uuid"`
	Tags        []*Tag `json:"tags"`
}

type StorageObjectRef struct {
	Bucket   string `json:"bucket"`
	ObjectId string `json:"object_id"`
}

type StorageCreateLinkReq struct {
	CompanyUuid string           `json:"company_uuid"`
	Source      StorageObjectRef `json:"source"`
	Target      StorageObjectRef `json:"target"`
}

type StoragePresignUrlRes struct {
	Url *PresignUrl `json:"url"`
}

type StorageObjectTagsRes struct {
	Tags []*Tag `json:"tags"`
}

type StorageSearchItemRes struct {
	Id             *string `json:"id,omitempty"`
	LastModifiedAt *string `json:"last_modified_at,omitempty"`
	Size           *int    `json:"size_bytes,omitempty"`
}

type StorageLinkTargetRes struct {
	TargetList []StorageObjectRef `json:"target_list"`
}
