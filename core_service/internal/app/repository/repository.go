package repository

import (
	"context"

	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
)

type Repository interface {
}

type CacheRepository interface {
	PutAdminToken(ctx context.Context, token string) error
	GetAdminToken(ctx context.Context) (string, error)

	GetApiData(ctx context.Context, company *models.Company) (*models.ApiData, error)
	PutApiData(ctx context.Context, company *models.Company, apiData *models.ApiData) error

	PutApiList(ctx context.Context, companyUuid string, apiList []*models.Api) error
	GetApiList(ctx context.Context, companyUuid string) ([]*models.Api, error)
	GetApi(ctx context.Context, companyUuid, apiUuid string) (*models.Api, error)
	DeleteApiList(ctx context.Context, company *models.Company) error

	GetCompany(ctx context.Context, companyUuid string) (*models.Company, error)
	PutCompany(ctx context.Context, company *models.Company) error
	// GetCompanyList/PutCompanyList — кэш полного справочника компаний (короткий TTL,
	// см. redisCompanyFullListLife); (nil, nil) при промахе.
	GetCompanyList(ctx context.Context) ([]*models.Company, error)
	PutCompanyList(ctx context.Context, companies []*models.Company) error

	PutEvent(ctx context.Context, companyUuid string, event *models.Event) error
	GetEvent(ctx context.Context, companyUuid, eventUuid string) (*models.Event, error)
}

type IDMRepository interface {
	UserData(ctx context.Context, accessToken string) (*models.User, error)
	IsAdministrator(u *models.User) bool
	IsSuperuser(u *models.User) bool
	IsWebservice(u *models.User) bool
	IsExternalUser(u *models.User) bool
	GetAdminToken(ctx context.Context) (string, error)
}

type AdminRepository interface {
	AddProperty(ctx context.Context, accessToken, companyUuid string, apiList []*models.Api) error
	GetProperty(ctx context.Context, accessToken, companyUuid string) ([]*models.Api, error)
	GetUserData(ctx context.Context, accessToken, userUuid string) (*models.User, error)
	GetCompany(ctx context.Context, accessToken, companyUuid string) (*models.Company, error)
	// GetCompanyList — весь справочник компаний Admin API (постранично). Только для
	// администраторских сценариев; на горячем пути — через CacheRepository.GetCompanyList.
	GetCompanyList(ctx context.Context, accessToken string) ([]*models.Company, error)
	GetApiData(ctx context.Context, accessToken string, company *models.Company) error
}

type StorageRepository interface {
	GetPresignUploadUrl(ctx context.Context, accessToken, companyUuid, filename string, tags []*models.Tag) (*models.PresignUrl, error)
	GetPresignDownloadUrl(ctx context.Context, accessToken, companyUuid, fileId string) (*models.PresignUrl, error)
	PutFile(ctx context.Context, accessToken, companyUuid string, file *models.File) (*models.File, error)
	GetFileContentType(ctx context.Context, accessToken, companyUuid, fileId string) (models.FileContentType, error)
	DeleteTempTag(ctx context.Context, accessToken, companyUuid, objectId string) error
	CreateLink(ctx context.Context, accessToken, companyUuid, bucket, objectId string) error
	DownloadFile(ctx context.Context, accessToken, companyUuid, fileId string) ([]byte, string, error)
}
