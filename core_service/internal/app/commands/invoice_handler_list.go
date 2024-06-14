package commands

import (
	"context"

	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/repository"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/models"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
)

type InvoiceHandlersListCmdHandler interface {
	Handle(ctx context.Context, command *InvoiceHandlersListCommand) ([]*models.InvoiceHandler, error)
}

type invoiceHandlersListHandler struct {
	log       logger.Logger
	cfg       *config.Config
	redisRepo repository.CacheRepository
}

func NewInvoiceHandlersListHandler(log logger.Logger, cfg *config.Config, redisRepo repository.CacheRepository) *invoiceHandlersListHandler {
	return &invoiceHandlersListHandler{log: log, cfg: cfg, redisRepo: redisRepo}
}

func (c *invoiceHandlersListHandler) Handle(ctx context.Context, command *InvoiceHandlersListCommand) ([]*models.InvoiceHandler, error) {
	// span, ctx := opentracing.StartSpanFromContext(ctx, "invoiceHandlersListHandler.Handle")
	// defer span.Finish()

	// product := &models.Product{
	// 	ProductID:   command.ProductID,
	// 	Name:        command.Name,
	// 	Description: command.Description,
	// 	Price:       command.Price,
	// 	CreatedAt:   command.CreatedAt,
	// 	UpdatedAt:   command.UpdatedAt,
	// }

	// created, err := c.mongoRepo.CreateProduct(ctx, product)
	// if err != nil {
	// 	return err
	// }

	// c.redisRepo.PutProduct(ctx, created.ProductID, created)
	return nil, nil
}
