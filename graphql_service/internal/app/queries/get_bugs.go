package queries

import (
	"context"

	"github.com/opentracing/opentracing-go"
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/graphql_service/config"
	model "github.com/testovoleg/5s-microservice-template/graphql_service/internal/graph_model"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
)

type GetBugsHandler interface {
	Handle(ctx context.Context, query *GetBugsQuery) (*model.BugsResponse, error)
}

type getBugsHandler struct {
	log        logger.Logger
	cfg        *config.Config
	coreClient coreService.CoreServiceClient
}

func NewGetBugsHandler(log logger.Logger, cfg *config.Config, coreClient coreService.CoreServiceClient) *getBugsHandler {
	return &getBugsHandler{log: log, cfg: cfg, coreClient: coreClient}
}

func (s *getBugsHandler) Handle(ctx context.Context, query *GetBugsQuery) (*model.BugsResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getBugsHandler.Handle")
	defer span.Finish()

	// ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.Context())
	// res, err := s.coreClient.GetBugList(ctx, &coreService.GetBugListReq{
	// 	ProductID: query.ProductID,
	// 	State:     model.GetBugStateToGrpc(query.State),
	// 	BugID:     query.BugID,
	// 	ReleaseID: query.ReleaseID,
	// 	Page:      int64(query.Pagination.Page),
	// 	Size:      int64(query.Pagination.Size),
	// })
	// if err != nil {
	// 	return nil, err
	// }
	// bugs := model.BugListResponseFromGrpcMessage(res)

	// res2, err := s.rlClient.GetReleasesList(ctx, &coreService.GetReleasesListReq{ProductID: query.ProductID})
	// if err != nil {
	// 	return nil, err
	// }

	// rmap := make(map[string]*model.Release)
	// for _, v := range model.ReleasesListFromGrpcMessage(res2.GetRelease()) {
	// 	rmap[v.ID] = v
	// }

	// for _, v := range bugs.Bugs {
	// 	if v.CreatedForRelease != nil && v.CreatedForRelease.ID != "" {
	// 		v.CreatedForRelease = rmap[v.CreatedForRelease.ID]
	// 	}
	// 	if v.SolvedInRelease != nil && v.SolvedInRelease.ID != "" {
	// 		v.SolvedInRelease = rmap[v.SolvedInRelease.ID]
	// 	}
	// }

	// return bugs, nil
	return nil, nil
}
