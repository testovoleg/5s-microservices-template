package repository

import (
	"context"
	"encoding/json"
	"fmt"
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
	httpErrors "github.com/testovoleg/5s-microservice-template/pkg/http_errors"
	httpUtils "github.com/testovoleg/5s-microservice-template/pkg/http_utils"
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

func storageHTTPOpts(accessToken string) httpUtils.SendOpts {
	return httpUtils.SendOpts{
		BearerToken: accessToken,
		Timeout:     30 * time.Second,
		ErrorPrefix: "Error was during on Storage API",
	}
}

func (r *storageRepository) PutFile(ctx context.Context, accessToken, companyUuid string, file *models.File) (*models.File, error) {
	ctx, span := tracing.StartSpan(ctx, "storageRepository.PutFile")
	defer span.End()

	if file.Data == nil {
		data, filename, err := r.downloadFile(ctx, file.Url)
		if err != nil {
			return nil, httpErrors.ErrorWrap(err, span, "r.downloadFile")
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
		return nil, httpErrors.ErrorWrap(err, span, "r.GetPresignUploadUrl")
	}

	err = r.uploadFile(ctx, presignUrl, &file.Data)
	if err != nil {
		return nil, httpErrors.ErrorWrap(err, span, "r.uploadFile")
	}

	err = r.DeleteTempTag(ctx, accessToken, companyUuid, presignUrl.ObjectId)
	if err != nil {
		return nil, httpErrors.ErrorWrap(err, span, "r.deleteTempTag")
	}

	file.FileId = presignUrl.ObjectId
	file.Url = ""

	return file, nil
}

func (r *storageRepository) GetPresignDownloadUrl(ctx context.Context, accessToken, companyUuid, fileId string) (*models.PresignUrl, error) {
	ctx, span := tracing.StartSpan(ctx, "storageRepository.GetPresignDownloadUrl")
	defer span.End()

	requestUrl := r.cfg.API.StorageApiUrl + "/object/url_download"

	pattern := regexp.MustCompile(`(.+)\/(.+)$`)
	finded := pattern.FindStringSubmatch(fileId)
	if len(finded) == 3 && len(finded[1]) == 36 {
		companyUuid = finded[1]
	}

	jsonBody, err := json.Marshal(models.StorageObjectReq{
		Bucket:      r.cfg.API.StorageBucket,
		CompanyUuid: companyUuid,
		ObjectId:    fileId,
	})
	if err != nil {
		return nil, httpErrors.ErrorWrap(err, span, "json.Marshal")
	}

	b, err := httpUtils.SendHTTPRequest(ctx, http.MethodPost, requestUrl, jsonBody, storageHTTPOpts(accessToken))
	if err != nil {
		return nil, err
	}

	var body models.StoragePresignUrlRes
	if err = json.Unmarshal(b, &body); err != nil {
		return nil, httpErrors.ErrorWrap(err, span, "json.Unmarshal")
	}

	if body.Url == nil || body.Url.Url == "" {
		return nil, errors.New("url is empty")
	}

	return body.Url, nil
}

func (r *storageRepository) GetPresignUploadUrl(ctx context.Context, accessToken, companyUuid, filename string, tags []*models.Tag) (*models.PresignUrl, error) {
	ctx, span := tracing.StartSpan(ctx, "storageRepository.GetPresignUploadUrl")
	defer span.End()

	requestUrl := r.cfg.API.StorageApiUrl + "/object/url_upload"

	jsonBody, err := json.Marshal(models.StoragePresignUploadReq{
		Bucket:      r.cfg.API.StorageBucket,
		CompanyUuid: companyUuid,
		Filename:    filename,
		Prefixs:     []string{constants.ShortMicroserviceName},
		Tags:        tags,
	})
	if err != nil {
		return nil, httpErrors.ErrorWrap(err, span, "json.Marshal")
	}

	b, err := httpUtils.SendHTTPRequest(ctx, http.MethodPost, requestUrl, jsonBody, storageHTTPOpts(accessToken))
	if err != nil {
		return nil, err
	}

	var body models.StoragePresignUrlRes
	if err = json.Unmarshal(b, &body); err != nil {
		return nil, httpErrors.ErrorWrap(err, span, "json.Unmarshal")
	}

	return body.Url, nil
}

func (r *storageRepository) GetFileContentType(ctx context.Context, accessToken, companyUuid, fileId string) (models.FileContentType, error) {
	ctx, span := tracing.StartSpan(ctx, "storageRepository.GetFileContentType")
	defer span.End()

	encodedFileId := url.QueryEscape(fileId)
	requestUrl, err := httpUtils.BuildQueryURL(
		r.cfg.API.StorageApiUrl+"/object/"+encodedFileId,
		companyUuid, "", "",
		&struct {
			Bucket string `url:"bucket"`
		}{Bucket: r.cfg.API.StorageBucket},
	)
	if err != nil {
		return "", httpErrors.ErrorWrap(err, span, "httpUtils.BuildQueryURL")
	}

	b, err := httpUtils.SendHTTPRequest(ctx, http.MethodGet, requestUrl, nil, storageHTTPOpts(accessToken))
	if err != nil {
		return "", err
	}

	var body models.StorageObjectTagsRes
	if err = json.Unmarshal(b, &body); err != nil {
		return "", httpErrors.ErrorWrap(err, span, "json.Unmarshal")
	}

	for _, t := range body.Tags {
		if t.Key == constants.FileStorageContentTypeTag {
			return models.FileContentType(t.Value), nil
		}
	}

	return "", errors.New("Content type of file not defined. FileId: " + fileId)
}

func (r *storageRepository) getFilename(ctx context.Context, accessToken, companyUuid, fileId string) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "storageRepository.GetFilename")
	defer span.End()

	encodedFileId := url.QueryEscape(fileId)
	requestUrl, err := httpUtils.BuildQueryURL(
		r.cfg.API.StorageApiUrl+"/object/"+encodedFileId,
		companyUuid, "", "",
		&struct {
			Bucket string `url:"bucket"`
		}{Bucket: r.cfg.API.StorageBucket},
	)
	if err != nil {
		return "", httpErrors.ErrorWrap(err, span, "httpUtils.BuildQueryURL")
	}

	b, err := httpUtils.SendHTTPRequest(ctx, http.MethodGet, requestUrl, nil, storageHTTPOpts(accessToken))
	if err != nil {
		return "", err
	}

	var body models.StorageObjectTagsRes
	if err = json.Unmarshal(b, &body); err != nil {
		return "", httpErrors.ErrorWrap(err, span, "json.Unmarshal")
	}

	for _, t := range body.Tags {
		if t.Key == constants.FileStorageFilenameTag {
			return t.Value, nil
		}
	}

	return "", errors.New("Content type of file not defined. FileId: " + fileId)
}

func (r *storageRepository) uploadFile(ctx context.Context, presignUrl *models.PresignUrl, data *[]byte) error {
	ctx, span := tracing.StartSpan(ctx, "storageRepository.uploadFile")
	defer span.End()

	if data == nil {
		return errors.New("empty file data")
	}

	_, err := httpUtils.SendHTTPRequest(ctx, http.MethodPut, presignUrl.Url, *data, httpUtils.SendOpts{
		Timeout:         30 * time.Second,
		ErrorPrefix:     "Error was during on Storage API",
		OmitContentType: true,
		Headers: map[string]string{
			"X-Amz-Tagging":  presignUrl.XAmzTagging,
			"Content-Length": fmt.Sprint(len(*data)),
		},
	})
	return err
}

func (r *storageRepository) downloadFile(ctx context.Context, url string) ([]byte, string, error) {
	ctx, span := tracing.StartSpan(ctx, "storageRepository.downloadFile")
	defer span.End()

	status, b, header, err := httpUtils.SendHTTPRequestFull(ctx, http.MethodGet, url, nil, httpUtils.SendOpts{
		Timeout: 30 * time.Second,
	})
	if err != nil {
		return nil, "", err
	}
	if status != http.StatusOK {
		return nil, "", httpErrors.HttpBody(fmt.Sprintf("%d %s", status, http.StatusText(status)), b, span, "Error was during on Storage API")
	}

	filename := truncateDownloadFilename(httpUtils.FilenameFromContentDisposition(header.Get("Content-Disposition")))
	return b, filename, nil
}

func truncateDownloadFilename(filename string) string {
	if filename == "" || len(filename) <= 25 {
		return filename
	}

	ext := path.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	maxNameLen := 25 - len(ext)
	if len(nameWithoutExt) > maxNameLen {
		return "file" + ext
	}
	return nameWithoutExt[:maxNameLen] + ext
}

func (r *storageRepository) DownloadFile(ctx context.Context, accessToken, companyUuid, fileId string) ([]byte, string, error) {
	if fileId == "" {
		return nil, "", errors.New("empty input")
	}

	presignURL, err := r.GetPresignDownloadUrl(ctx, accessToken, companyUuid, fileId)
	if err != nil {
		return nil, "", httpErrors.ErrorWrap(err, nil, "r.GetPresignDownloadURL")
	}

	fileData, _, err := r.downloadFile(ctx, presignURL.Url)
	if err != nil {
		return nil, "", httpErrors.ErrorWrap(err, nil, "r.downloadFile")
	}

	filename, _ := r.getFilename(ctx, accessToken, companyUuid, fileId)

	return fileData, filename, nil
}

func (r *storageRepository) DeleteTempTag(ctx context.Context, accessToken, companyUuid, objectId string) error {
	ctx, span := tracing.StartSpan(ctx, "storageRepository.DeleteTempTag")
	defer span.End()

	requestUrl := r.cfg.API.StorageApiUrl + "/object/delete_temp_tag"

	jsonBody, err := json.Marshal(models.StorageObjectReq{
		Bucket:      r.cfg.API.StorageBucket,
		CompanyUuid: companyUuid,
		ObjectId:    objectId,
	})
	if err != nil {
		return httpErrors.ErrorWrap(err, span, "json.Marshal")
	}

	_, err = httpUtils.SendHTTPRequest(ctx, http.MethodPatch, requestUrl, jsonBody, storageHTTPOpts(accessToken))
	return err
}

func (r *storageRepository) searchFileInBucket(ctx context.Context, accessToken, companyUuid, bucket string, tags []*models.Tag) (*models.File, error) {
	ctx, span := tracing.StartSpan(ctx, "storageRepository.searchFileInBucket")
	defer span.End()

	requestUrl := r.cfg.API.StorageApiUrl + "/object/search"

	jsonBody, err := json.Marshal(models.StorageSearchReq{
		Bucket:      bucket,
		CompanyUuid: companyUuid,
		Tags:        tags,
	})
	if err != nil {
		return nil, httpErrors.ErrorWrap(err, span, "json.Marshal")
	}

	b, err := httpUtils.SendHTTPRequest(ctx, http.MethodPost, requestUrl, jsonBody, storageHTTPOpts(accessToken))
	if err != nil {
		return nil, err
	}

	var body []*models.StorageSearchItemRes
	if err = json.Unmarshal(b, &body); err != nil {
		return nil, httpErrors.ErrorWrap(err, span, "json.Unmarshal")
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
	ctx, span := tracing.StartSpan(ctx, "storageRepository.findLink")
	defer span.End()

	requestUrl, err := httpUtils.BuildQueryURL(
		r.cfg.API.StorageApiUrl+"/link/target",
		companyUuid, "", "",
		&struct {
			Bucket   string `url:"bucket"`
			ObjectId string `url:"object_id"`
		}{Bucket: sourceBucket, ObjectId: sourceObjectId},
	)
	if err != nil {
		return "", "", httpErrors.ErrorWrap(err, span, "httpUtils.BuildQueryURL")
	}

	opts := storageHTTPOpts(accessToken)
	opts.ContentType = "application/json"
	b, err := httpUtils.SendHTTPRequest(ctx, http.MethodGet, requestUrl, nil, opts)
	if err != nil {
		return "", "", err
	}

	var body models.StorageLinkTargetRes
	if err = json.Unmarshal(b, &body); err != nil {
		return "", "", httpErrors.ErrorWrap(err, span, "json.Unmarshal")
	}

	for _, v := range body.TargetList {
		if v.Bucket == r.cfg.API.StorageBucket {
			return v.Bucket, v.ObjectId, nil
		}
	}

	return "", "", errors.New("Link to file not found")
}

func (r *storageRepository) CreateLink(ctx context.Context, accessToken, companyUuid, bucket, objectId string) error {
	ctx, span := tracing.StartSpan(ctx, "storageRepository.CreateLink")
	defer span.End()

	requestUrl := r.cfg.API.StorageApiUrl + "/link"

	jsonBody, err := json.Marshal(models.StorageCreateLinkReq{
		CompanyUuid: companyUuid,
		Source: models.StorageObjectRef{
			Bucket:   r.cfg.API.StorageBucket,
			ObjectId: objectId,
		},
		Target: models.StorageObjectRef{
			Bucket:   bucket,
			ObjectId: objectId,
		},
	})
	if err != nil {
		return httpErrors.ErrorWrap(err, span, "json.Marshal")
	}

	_, err = httpUtils.SendHTTPRequest(ctx, http.MethodPost, requestUrl, jsonBody, storageHTTPOpts(accessToken))
	return err
}
