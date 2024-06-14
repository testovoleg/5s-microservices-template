package commands

import (
	"context"

	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/repository"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
)

type UpdateProductCmdHandler interface {
	Handle(ctx context.Context, command *UpdateProductCommand) error
}

type updateProductCmdHandler struct {
	log       logger.Logger
	cfg       *config.Config
	redisRepo repository.CacheRepository
}

func NewUpdateProductCmdHandler(log logger.Logger, cfg *config.Config, redisRepo repository.CacheRepository) *updateProductCmdHandler {
	return &updateProductCmdHandler{log: log, cfg: cfg, redisRepo: redisRepo}
}

func (c *updateProductCmdHandler) Handle(ctx context.Context, command *UpdateProductCommand) error {
	// span, ctx := opentracing.StartSpanFromContext(ctx, "updateProductCmdHandler.Handle")
	// defer span.Finish()

	// product := &models.Product{
	// 	ProductID:   command.ProductID,
	// 	Name:        command.Name,
	// 	Description: command.Description,
	// 	Price:       command.Price,
	// 	UpdatedAt:   command.UpdatedAt,
	// }

	// c.redisRepo.PutProduct(ctx, updated.ProductID, updated)
	return nil
}
