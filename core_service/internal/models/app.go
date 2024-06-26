package models

import (
	"time"

	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type InvoiceHandler struct {
	ID          string `json:"id,omitempty"`
	Version     string `json:"version,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

func InvoiceHandlerToGrpcMessage(in *InvoiceHandler) *coreService.InvoiceHandler {
	return &coreService.InvoiceHandler{
		ID:          in.ID,
		Title:       in.Title,
		Description: in.Description,
		Version:     in.Version,
	}
}

func InvoiceHandlersListToGrpc(arr []*InvoiceHandler) []*coreService.InvoiceHandler {
	list := make([]*coreService.InvoiceHandler, 0, len(arr))
	for _, v := range arr {
		list = append(list, InvoiceHandlerToGrpcMessage(v))
	}
	return list
}

type Product struct {
	ProductID   string    `json:"productId" bson:"_id,omitempty"`
	Name        string    `json:"name,omitempty" bson:"name,omitempty" validate:"required,min=3,max=250"`
	Description string    `json:"description,omitempty" bson:"description,omitempty" validate:"required,min=3,max=500"`
	Price       float64   `json:"price,omitempty" bson:"price,omitempty" validate:"required"`
	CreatedAt   time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

// ProductsList products list response with pagination
type ProductsList struct {
	TotalCount int64      `json:"totalCount" bson:"totalCount"`
	TotalPages int64      `json:"totalPages" bson:"totalPages"`
	Page       int64      `json:"page" bson:"page"`
	Size       int64      `json:"size" bson:"size"`
	HasMore    bool       `json:"hasMore" bson:"hasMore"`
	Products   []*Product `json:"products" bson:"products"`
}

func NewProductListWithPagination(products []*Product, count int64, pagination *utils.Pagination) *ProductsList {
	return &ProductsList{
		TotalCount: count,
		TotalPages: int64(pagination.GetTotalPages(int(count))),
		Page:       int64(pagination.GetPage()),
		Size:       int64(pagination.GetSize()),
		HasMore:    pagination.GetHasMore(int(count)),
		Products:   products,
	}
}

func ProductToGrpcMessage(product *Product) *coreService.Product {
	return &coreService.Product{
		ProductID:   product.ProductID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   timestamppb.New(product.CreatedAt),
		UpdatedAt:   timestamppb.New(product.UpdatedAt),
	}
}
