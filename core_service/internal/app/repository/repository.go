package repository

import (
	"context"

	uuid "github.com/satori/go.uuid"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"
)

type Repository interface {
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	DeleteProduct(ctx context.Context, uuid uuid.UUID) error

	GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error)
	Search(ctx context.Context, search string, pagination *utils.Pagination) (*models.ProductsList, error)
}

type CacheRepository interface {
	PutProduct(ctx context.Context, key string, product *models.Product)
	GetProduct(ctx context.Context, key string) (*models.Product, error)
	DelProduct(ctx context.Context, key string)
	DelAllProducts(ctx context.Context)
}
