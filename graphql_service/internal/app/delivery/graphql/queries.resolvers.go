package graph_resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.44

import (
	"context"
	"strconv"

	"github.com/testovoleg/5s-microservice-template/graphql_service/internal/app/queries"
	model "github.com/testovoleg/5s-microservice-template/graphql_service/internal/graph_model"
	graph "github.com/testovoleg/5s-microservice-template/graphql_service/schema"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"
)

// Bugs is the resolver for the bugs field.
func (r *queryResolver) Bugs(ctx context.Context, productID int, state *model.BugState, bugID *int, solvedInReleaseID *int, page int, size int, orderBy *model.OrderBy) (*model.BugsResponse, error) {
	r.metrics.GetBugsGraphQLQueries.Inc()

	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "queryResolver.Bugs")
	defer span.End()

	pq := utils.NewPaginationFromQueryParams(strconv.Itoa(size), strconv.Itoa(page))
	query := queries.NewGetBugsQuery(pq, productID, state, bugID, solvedInReleaseID)

	response, err := r.bs.Queries.GetBugs.Handle(ctx, query)
	if err != nil {
		r.log.WarnMsg("queryResolver.Bugs", err)
		r.metrics.ErrorHttpRequests.Inc()
		return nil, err
	}

	r.metrics.SuccessHttpRequests.Inc()
	return response, nil
}

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
