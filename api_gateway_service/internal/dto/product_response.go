package dto

import (
	"time"

	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
)

type ProductResponse struct {
	ProductID   string    `json:"productId"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

func ProductResponseFromGrpc(product *coreService.Product) *ProductResponse {
	return &ProductResponse{
		ProductID:   product.GetProductID(),
		Name:        product.GetName(),
		Description: product.GetDescription(),
		Price:       product.GetPrice(),
		CreatedAt:   product.GetCreatedAt().AsTime(),
		UpdatedAt:   product.GetUpdatedAt().AsTime(),
	}
}
