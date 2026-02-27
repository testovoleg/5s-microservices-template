package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	"github.com/testovoleg/5s-microservice-template/pkg/constants"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"
)

type storageRepository struct {
	log logger.Logger
	cfg *config.Config
}

func NewStorageRepository(log logger.Logger, cfg *config.Config) *storageRepository {
	return &storageRepository{log: log, cfg: cfg}
}

func (r *storageRepository) PutFile(ctx context.Context, accessToken, companyUuid string, file *models.File) (*models.File, error) {
	ctx, span := tracing.StartSpan(ctx, "storageRepository.PutFile")
	defer span.End()

	if file.Data == nil {
		data, filename, err := r.downloadFile(ctx, file.Url)
		if err != nil {
			return nil, errors.Wrap(err, "r.downloadFile")
		}

		file.Filename = filename
		file.Data = data
		file.Size = int64(len(data))
		if file.Filename == "" {
			file.Filename = filename
		}
	}

	mimetype := utils.MimeType(file.Data, file.Filename)
	file.MimeType = &mimetype.MimeType

	if file.Filename == "" {
		file.Filename = string(file.ContentType) + mimetype.Extension
	}

	var tags []*models.Tag
	tags = append(
		tags,
		&models.Tag{Key: constants.FileStorageContentTypeTag, Value: string(file.ContentType)},
		&models.Tag{Key: constants.FileStorageCompanyUuidTag, Value: companyUuid},
	)

	presignUrl, err := r.GetPresignUploadUrl(ctx, accessToken, companyUuid, file.Filename, tags)
	if err != nil {
		return nil, errors.Wrap(err, "r.GetPresignUploadUrl")
	}

	err = r.uploadFile(ctx, presignUrl, &file.Data)
	if err != nil {
		return nil, errors.Wrap(err, "r.uploadFile")
	}

	err = r.DeleteTempTag(ctx, accessToken, companyUuid, presignUrl.ObjectId)
	if err != nil {
		return nil, errors.Wrap(err, "r.deleteTempTag")
	}

	file.FileId = presignUrl.ObjectId
	file.Url = ""

	return file, nil
}

func (r *storageRepository) GetPresignDownloadUrl(ctx context.Context, accessToken, companyUuid, fileId string) (*models.PresignUrl, error) {
	_, span := tracing.StartSpan(ctx, "storageRepository.GetPresignDownloadUrl")
	defer span.End()

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	requestUrl := r.cfg.API.StorageApiUrl + "/object/url_download"

	pattern := regexp.MustCompile(`(.+)\/(.+)$`)
	finded := pattern.FindStringSubmatch(fileId)
	if len(finded) == 3 && len(finded[1]) == 36 {
		companyUuid = finded[1]
	}

	reqBody := map[string]interface{}{}
	reqBody["bucket"] = r.cfg.API.StorageBucket
	reqBody["company_uuid"] = companyUuid
	reqBody["object_id"] = fileId

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, requestUrl, bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest")
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client.Do")
	}

	if res.StatusCode != 200 {
		return nil, utils.NewError(res, span, "Error was during on Storage API")
	}

	b, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "io.ReadAll")
	}

	type ResponseBody struct {
		Url *models.PresignUrl `json:"url"`
	}
	var body ResponseBody

	err = json.Unmarshal(b, &body)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}

	if body.Url == nil || (body.Url != nil && body.Url.Url == "") {
		return nil, errors.New("url is empty")
	}

	return body.Url, nil
}

func (r *storageRepository) GetPresignUploadUrl(ctx context.Context, accessToken, companyUuid, filename string, tags []*models.Tag) (*models.PresignUrl, error) {
	_, span := tracing.StartSpan(ctx, "storageRepository.GetPresignUploadUrl")
	defer span.End()

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	requestUrl := r.cfg.API.StorageApiUrl + "/object/url_upload"

	reqBody := map[string]interface{}{}
	reqBody["bucket"] = r.cfg.API.StorageBucket
	reqBody["company_uuid"] = companyUuid
	reqBody["filename"] = filename
	reqBody["prefixs"] = []string{constants.ShortMicroserviceName}

	tagsBody := []map[string]interface{}{}

	for _, t := range tags {
		tag := map[string]interface{}{}
		tag["key"] = t.Key
		tag["value"] = t.Value
		tagsBody = append(tagsBody, tag)
	}

	reqBody["tags"] = tagsBody

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, requestUrl, bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest")
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client.Do")
	}

	b, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "io.ReadAll")
	}

	if res.StatusCode != 200 {
		return nil, utils.NewError(res, span, "Error was during on Storage API")
	}

	type ResponseBody struct {
		Url *models.PresignUrl `json:"url"`
	}
	var body ResponseBody

	err = json.Unmarshal(b, &body)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}

	return body.Url, nil
}

func (r *storageRepository) GetFileContentType(ctx context.Context, accessToken, companyUuid, fileId string) (models.FileContentType, error) {
	_, span := tracing.StartSpan(ctx, "storageRepository.GetFileContentType")
	defer span.End()

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	encodedFileId := url.QueryEscape(fileId)

	requestUrl := r.cfg.API.StorageApiUrl + "/object/" + encodedFileId + "?bucket=" + r.cfg.API.StorageBucket + "&company_uuid=" + companyUuid

	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		return "", errors.Wrap(err, "http.NewRequest")
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "client.Do")
	}

	if res.StatusCode != 200 {
		return "", utils.NewError(res, span, "Error was during on Storage API")
	}

	b, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return "", errors.Wrap(err, "io.ReadAll")
	}

	type ResponseBody struct {
		Tags []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"tags"`
	}
	var body ResponseBody

	err = json.Unmarshal(b, &body)
	if err != nil {
		return "", errors.Wrap(err, "json.Unmarshal")
	}

	for _, t := range body.Tags {
		if t.Key == constants.FileStorageContentTypeTag {
			return models.FileContentType(t.Value), nil
		}
	}

	return "", errors.New("Content type of file not defined. FileId: " + fileId)
}

func (r *storageRepository) getFilename(ctx context.Context, accessToken, companyUuid, fileId string) (string, error) {
	_, span := tracing.StartSpan(ctx, "storageRepository.GetFilename")
	defer span.End()

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	encodedFileId := url.QueryEscape(fileId)

	requestUrl := r.cfg.API.StorageApiUrl + "/object/" + encodedFileId + "?bucket=" + r.cfg.API.StorageBucket + "&company_uuid=" + companyUuid

	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		return "", errors.Wrap(err, "http.NewRequest")
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "client.Do")
	}

	if res.StatusCode != 200 {
		return "", utils.NewError(res, span, "Error was during on Storage API")
	}

	b, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return "", errors.Wrap(err, "io.ReadAll")
	}

	type ResponseBody struct {
		Tags []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"tags"`
	}
	var body ResponseBody

	err = json.Unmarshal(b, &body)
	if err != nil {
		return "", errors.Wrap(err, "json.Unmarshal")
	}

	for _, t := range body.Tags {
		if t.Key == constants.FileStorageFilenameTag {
			return t.Value, nil
		}
	}

	return "", errors.New("Content type of file not defined. FileId: " + fileId)
}

func (r *storageRepository) uploadFile(ctx context.Context, presignUrl *models.PresignUrl, data *[]byte) error {
	_, span := tracing.StartSpan(ctx, "storageRepository.uploadFile")
	defer span.End()

	if data == nil {
		return errors.New("empty file data")
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	bodyReader := bytes.NewReader(*data)

	req, err := http.NewRequest(http.MethodPut, presignUrl.Url, bodyReader)
	if err != nil {
		return errors.Wrap(err, "http.NewRequest")
	}

	req.Header.Set("X-Amz-Tagging", presignUrl.XAmzTagging)
	req.Header.Set("Content-Length", fmt.Sprint(len(*data)))

	res, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "client.Do")
	}

	if res.StatusCode != 200 {
		return utils.NewError(res, span, "Error was during on Storage API")
	}

	return nil
}

func (r *storageRepository) downloadFile(ctx context.Context, url string) ([]byte, string, error) {
	_, span := tracing.StartSpan(ctx, "storageRepository.downloadFile")
	defer span.End()

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, "", errors.Wrap(err, "http.NewRequest")
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, "", errors.Wrap(err, "client.Do")
	}

	if res.StatusCode != 200 {
		return nil, "", utils.NewError(res, span, "Error was during on Storage API")
	}

	b, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, "", errors.Wrap(err, "io.ReadAll")
	}

	var filename string
	cd := res.Header.Get("Content-Disposition")
	if cd != "" {
		filename = r.extractFilenameFromContentDisposition(cd)
	}

	return b, filename, nil
}

func (r *storageRepository) extractFilenameFromContentDisposition(cd string) string {
	if cd == "" {
		return ""
	}

	filename := ""

	if idx := strings.Index(cd, "filename*="); idx != -1 {
		start := idx + len("filename*=")
		end := strings.Index(cd[start:], "\"")
		if end != -1 {
			filename = cd[start : start+end]
		}
	}

	if filename == "" {
		if idx := strings.Index(cd, "filename="); idx != -1 {
			start := idx + len("filename=")
			if cd[start] == '"' {
				start++
				end := strings.Index(cd[start:], "\"")
				if end != -1 {
					filename = cd[start : start+end]
				}
			} else {
				end := strings.IndexAny(cd[start:], " ;")
				if end != -1 {
					filename = cd[start : start+end]
				} else {
					filename = cd[start:]
				}
			}
		}
	}

	if filename == "" {
		return ""
	}

	if len(filename) > 25 {
		ext := path.Ext(filename)
		nameWithoutExt := strings.TrimSuffix(filename, ext)

		maxNameLen := 25 - len(ext)
		if len(nameWithoutExt) > maxNameLen {
			return "file" + ext
		}

		truncatedName := nameWithoutExt[:maxNameLen]
		filename = truncatedName + ext
	}

	return filename
}

func (r *storageRepository) DownloadFile(ctx context.Context, accessToken, companyUuid, fileId string) ([]byte, string, error) {
	if fileId == "" {
		return nil, "", errors.New("empty input")
	}

	presignURL, err := r.GetPresignDownloadUrl(ctx, accessToken, companyUuid, fileId)
	if err != nil {
		return nil, "", errors.Wrap(err, "r.GetPresignDownloadURL")
	}

	fileData, _, err := r.downloadFile(ctx, presignURL.Url)
	if err != nil {
		return nil, "", errors.Wrap(err, "r.downloadFile")
	}

	filename, _ := r.getFilename(ctx, accessToken, companyUuid, fileId)

	return fileData, filename, nil
}

func (r *storageRepository) DeleteTempTag(ctx context.Context, accessToken, companyUuid, objectId string) error {
	_, span := tracing.StartSpan(ctx, "storageRepository.DeleteTempTag")
	defer span.End()

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	requestUrl := r.cfg.API.StorageApiUrl + "/object/delete_temp_tag"

	reqBody := map[string]interface{}{}
	reqBody["bucket"] = r.cfg.API.StorageBucket
	reqBody["company_uuid"] = companyUuid
	reqBody["object_id"] = objectId

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPatch, requestUrl, bodyReader)
	if err != nil {
		return errors.Wrap(err, "http.NewRequest")
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "client.Do")
	}

	if res.StatusCode != 200 {
		return utils.NewError(res, span, "Error was during on Storage API")
	}

	return nil
}

func (r *storageRepository) searchFileInBucket(ctx context.Context, accessToken, companyUuid, bucket string, tags []*models.Tag) (*models.File, error) {
	_, span := tracing.StartSpan(ctx, "storageRepository.searchFileInBucket")
	defer span.End()

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	requestUrl := r.cfg.API.StorageApiUrl + "/object/search"

	reqBody := map[string]interface{}{}
	reqBody["bucket"] = bucket
	reqBody["company_uuid"] = companyUuid

	tagsBody := []map[string]interface{}{}

	for _, t := range tags {
		tag := map[string]interface{}{}
		tag["key"] = t.Key
		tag["value"] = t.Value
		tagsBody = append(tagsBody, tag)
	}

	reqBody["tags"] = tagsBody

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, requestUrl, bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest")
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client.Do")
	}

	if res.StatusCode != 200 {
		return nil, utils.NewError(res, span, "Error was during on Storage API")
	}

	b, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "io.ReadAll")
	}

	type ResponseBody struct {
		Id             *string `json:"id,omitempty"`
		LastModifiedAt *string `json:"last_modified_at,omitempty"`
		Size           *int    `json:"size_bytes,omitempty"`
	}

	var body []*ResponseBody

	err = json.Unmarshal(b, &body)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}

	if len(body) > 0 {
		return &models.File{
			FileId: utils.DerefString(body[0].Id),
			Size:   int64(utils.DerefInt(body[0].Size)),
		}, nil
	}

	return nil, nil
}

func (r *storageRepository) findLink(ctx context.Context, accessToken, companyUuid, sourceBucket, sourceObjectId string) (targetBucket string, targetObjectId string, err error) {
	_, span := tracing.StartSpan(ctx, "storageRepository.findLink")
	defer span.End()

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	requestUrl := r.cfg.API.StorageApiUrl + "/link/target?company_uuid=" + companyUuid + "&bucket=" + sourceBucket + "&object_id=" + sourceObjectId

	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		return "", "", errors.Wrap(err, "http.NewRequest")
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", "", errors.Wrap(err, "client.Do")
	}

	if res.StatusCode != 200 {
		return "", "", utils.NewError(res, span, "Error was during on Storage API")
	}

	b, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return "", "", errors.Wrap(err, "io.ReadAll")
	}

	type ResponseBody struct {
		TargetList []struct {
			Bucket   string `json:"bucket"`
			ObjectId string `json:"object_id"`
		} `json:"target_list"`
	}
	var body ResponseBody

	err = json.Unmarshal(b, &body)
	if err != nil {
		return "", "", errors.Wrap(err, "json.Unmarshal")
	}

	for _, v := range body.TargetList {
		if v.Bucket == r.cfg.API.StorageBucket {
			return v.Bucket, v.ObjectId, nil
		}
	}

	return "", "", errors.New("Link to file not found")
}

func (r *storageRepository) CreateLink(ctx context.Context, accessToken, companyUuid, bucket, objectId string) error {
	_, span := tracing.StartSpan(ctx, "storageRepository.CreateLink")
	defer span.End()

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	requestUrl := r.cfg.API.StorageApiUrl + "/link"

	reqBody := map[string]interface{}{}
	reqBody["company_uuid"] = companyUuid

	sourceBody := map[string]interface{}{}
	sourceBody["bucket"] = r.cfg.API.StorageBucket
	sourceBody["object_id"] = objectId
	reqBody["source"] = sourceBody

	targetBody := map[string]interface{}{}
	targetBody["bucket"] = bucket
	targetBody["object_id"] = objectId
	reqBody["target"] = targetBody

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, requestUrl, bodyReader)
	if err != nil {
		return errors.Wrap(err, "http.NewRequest")
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "client.Do")
	}

	if res.StatusCode != 200 {
		return utils.NewError(res, span, "Error was during on Storage API")
	}

	return nil
}
