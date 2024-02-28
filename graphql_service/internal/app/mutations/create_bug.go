package mutations

import (
	"context"

	"github.com/opentracing/opentracing-go"
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/graphql_service/config"
	model "github.com/testovoleg/5s-microservice-template/graphql_service/internal/graph_model"
	kafkaClient "github.com/testovoleg/5s-microservice-template/pkg/kafka"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

type CreateBugHandler interface {
	Handle(ctx context.Context, command *CreateBugCommand) (*model.Bug, error)
}

type createBugHandler struct {
	log           logger.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
	coreClient    coreService.CoreServiceClient
}

func NewCreateBugHandler(
	log logger.Logger,
	cfg *config.Config,
	kafkaProducer kafkaClient.Producer,
	coreClient coreService.CoreServiceClient,
) *createBugHandler {
	return &createBugHandler{
		log:           log,
		cfg:           cfg,
		kafkaProducer: kafkaProducer,
		coreClient:    coreClient,
	}
}

func (c *createBugHandler) Handle(ctx context.Context, command *CreateBugCommand) (*model.Bug, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createBugHandler.Handle")
	defer span.Finish()

	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.Context())
	// res, err := c.coreClient.CreateBug(ctx, &coreService.CreateBugReq{
	// 	AccessToken:        command.AccessToken,
	// 	Name:               command.Name,
	// 	Description:        command.Description,
	// 	CreateForReleaseID: command.CreateForReleaseID,
	// 	FileID:             command.Files,
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// bug := res.GetBug()

	// if len(command.Files) == 0 {
	// 	return model.BugFromGrpcMessage(bug), nil
	// }

	// createPreviewDto := &kafkaMessages.UploadedFiles{
	// 	BugID:     bug.BugID,
	// 	CommentID: "",
	// 	FileID:    command.Files,
	// }

	// dtoBytes, err := proto.Marshal(createPreviewDto)
	// if err != nil {
	// 	return nil, err
	// }

	// err = c.kafkaProducer.PublishMessage(ctx, kafka.Message{
	// 	Topic:   c.cfg.KafkaTopics.FileUploaded.TopicName,
	// 	Value:   dtoBytes,
	// 	Time:    time.Now().UTC(),
	// 	Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	// })

	// return model.BugFromGrpcMessage(bug), nil
	return nil, nil
}
