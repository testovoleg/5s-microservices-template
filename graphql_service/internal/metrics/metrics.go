package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/testovoleg/5s-microservice-template/graphql_service/config"
)

type ApiGatewayMetrics struct {
	SuccessHttpRequests           prometheus.Counter
	ErrorHttpRequests             prometheus.Counter
	GetProductsGraphQLQueries     prometheus.Counter
	GetReleasesGraphQLQueries     prometheus.Counter
	GetBugsGraphQLQueries         prometheus.Counter
	GetCommentsGraphQLQueries     prometheus.Counter
	GetVotesGraphQLQueries        prometheus.Counter
	SearchBugsGraphQLQueries      prometheus.Counter
	GetFavouritesGraphQLQueries   prometheus.Counter
	CreateProductGraphQLQueries   prometheus.Counter
	UpdateProductGraphQLQueries   prometheus.Counter
	CreateReleaseGraphQLQueries   prometheus.Counter
	UpdateReleaseGraphQLQueries   prometheus.Counter
	CreateBugGraphQLQueries       prometheus.Counter
	UpdateBugGraphQLQueries       prometheus.Counter
	DeleteBugGraphQLQueries       prometheus.Counter
	SetStateBugGraphQLQueries     prometheus.Counter
	AddCommentGraphQLQueries      prometheus.Counter
	AddVoteGraphQLQueries         prometheus.Counter
	AddFavouritesGraphQLQueries   prometheus.Counter
	DelFavouritesGraphQLQueries   prometheus.Counter
	CreateUploadURLGraphQLQueries prometheus.Counter
	CreateFileURLGraphQLQueries   prometheus.Counter
	UserGraphQLQueries            prometheus.Counter
	SelfUserGraphQLQueries        prometheus.Counter
	CreateBugHttpRequests         prometheus.Counter
	SetInworkStatusHttpRequests   prometheus.Counter
	SetDoneStatusHttpRequests     prometheus.Counter
	SetRejectedStatusRequests     prometheus.Counter
}

func NewApiGatewayMetrics(cfg *config.Config) *ApiGatewayMetrics {
	return &ApiGatewayMetrics{
		SuccessHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_http_requests_total", cfg.ServiceName),
			Help: "The total number of success http requests",
		}),
		ErrorHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_error_http_requests_total", cfg.ServiceName),
			Help: "The total number of error http requests",
		}),
		GetProductsGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_products_graphql_queries_total", cfg.ServiceName),
			Help: "The total number of get products graphql queries",
		}),
		GetReleasesGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_releases_graphql_queries_total", cfg.ServiceName),
			Help: "The total number of get releases graphql queries",
		}),
		GetBugsGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_bugs_graphql_queries_total", cfg.ServiceName),
			Help: "The total number of get bugs graphql queries",
		}),
		GetCommentsGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_comments_graphql_queries_total", cfg.ServiceName),
			Help: "The total number of get comments graphql queries",
		}),
		GetVotesGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_votes_graphql_queries_total", cfg.ServiceName),
			Help: "The total number of get votes graphql queries",
		}),
		SearchBugsGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_search_bugs_graphql_queries_total", cfg.ServiceName),
			Help: "The total number of search bugs graphql queries",
		}),
		GetFavouritesGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_favourites_graphql_queries_total", cfg.ServiceName),
			Help: "The total number of get favourites graphql queries",
		}),
		CreateProductGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_product_graphql_queries_total", cfg.ServiceName),
			Help: "The total number of create product graphql queries",
		}),
		UpdateProductGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_product_graphql_queries_total", cfg.ServiceName),
			Help: "The total number of update product graphql queries",
		}),
		CreateReleaseGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_release_graphql_queries_total", cfg.ServiceName),
			Help: "The total number of create release graphql queries",
		}),
		UpdateReleaseGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_release_graphql_queries_total", cfg.ServiceName),
			Help: "The total number of update release graphql queries",
		}),
		CreateBugGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_bug_graphql_queries_total", cfg.ServiceName),
			Help: "The total number of create bug graphql queries",
		}),
		UpdateBugGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_bug_graphql_queries_total", cfg.ServiceName),
			Help: "The total number of update bug graphql queries",
		}),
		DeleteBugGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_delete_bug_graphql_queries_total", cfg.ServiceName),
			Help: "The total number of delete bug graphql queries",
		}),
		SetStateBugGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_set_state_bug_graphql_queries_total", cfg.ServiceName),
			Help: "The total number of set state bug graphql queries",
		}),
		AddCommentGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_add_comment_graphql_queries_total", cfg.ServiceName),
			Help: "The total number of add comment graphql queries",
		}),
		AddVoteGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_add_vote_graphql_queries_total", cfg.ServiceName),
			Help: "The total number of add vote graphql queries",
		}),
		AddFavouritesGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_add_favourites_graphql_mutation_total", cfg.ServiceName),
			Help: "The total number of add favourites graphql mutation",
		}),
		DelFavouritesGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_del_favourites_graphql_mutation_total", cfg.ServiceName),
			Help: "The total number of del favourites graphql mutation",
		}),
		CreateUploadURLGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_upload_url_graphql_mutation_total", cfg.ServiceName),
			Help: "The total number of create upload url graphql mutation",
		}),
		CreateFileURLGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_file_url_graphql_mutation_total", cfg.ServiceName),
			Help: "The total number of create file url graphql mutation",
		}),
		UserGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_user_data_graphql_mutation_total", cfg.ServiceName),
			Help: "The total number of get user data graphql mutation",
		}),
		SelfUserGraphQLQueries: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_self_user_data_graphql_mutation_total", cfg.ServiceName),
			Help: "The total number of get self user data graphql mutation",
		}),
		CreateBugHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_bug_http_requests_total", cfg.ServiceName),
			Help: "The total number of create bug http requests",
		}),
		SetInworkStatusHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_set_inwork_status_http_requests_total", cfg.ServiceName),
			Help: "The total number of set inwork status http requests",
		}),
		SetDoneStatusHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_set_done_status_http_requests_total", cfg.ServiceName),
			Help: "The total number of set done status http requests",
		}),
		SetRejectedStatusRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_set_rejected_sratus_http_requests_total", cfg.ServiceName),
			Help: "The total number of set rejected status http requests",
		}),
	}
}
