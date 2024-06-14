package commands

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/repository"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "deleteProductCmdHandler.Handle")
	defer span.Finish()

	c.redisRepo.DelProduct(ctx, command.ProductID.String())
	return nil
}
