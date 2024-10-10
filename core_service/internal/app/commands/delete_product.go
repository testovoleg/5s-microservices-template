package commands

import (
	"context"

	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/repository"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

type DeleteProductCmdHandler interface {
	Handle(ctx context.Context, command *DeleteProductCommand) error
}

type deleteProductCmdHandler struct {
	log       logger.Logger
	cfg       *config.Config
	redisRepo repository.CacheRepository
}

func NewDeleteProductCmdHandler(log logger.Logger, cfg *config.Config, redisRepo repository.CacheRepository) *deleteProductCmdHandler {
	return &deleteProductCmdHandler{log: log, cfg: cfg, redisRepo: redisRepo}
}

func (c *deleteProductCmdHandler) Handle(ctx context.Context, command *DeleteProductCommand) error {
	ctx, span := tracing.StartSpan(ctx, "deleteProductCmdHandler.Handle")
	defer span.End()

	c.redisRepo.DelProduct(ctx, command.ProductID.String())
	return nil
}
