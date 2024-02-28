package queries

import (
	"strconv"

	model "github.com/testovoleg/5s-microservice-template/graphql_service/internal/graph_model"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"
)

type Queries struct {
	GetBugs GetBugsHandler
}

func NewQueries(
	getBugs GetBugsHandler,
) *Queries {
	return &Queries{
		GetBugs: getBugs,
	}
}

type GetBugsQuery struct {
	Pagination *utils.Pagination `json:"pagination"`
	ProductID  string            `json:"productID"`
	State      *model.BugState   `json:"state"`
	BugID      string            `json:"bugID"`
	ReleaseID  string            `json:"releaseID"`
}

func NewGetBugsQuery(pagination *utils.Pagination, productID int, state *model.BugState, bugID *int, releaseID *int) *GetBugsQuery {
	var bugTmp, releaseTmp string
	if bugID != nil {
		bugTmp = strconv.Itoa(*bugID)
	}
	if releaseID != nil {
		releaseTmp = strconv.Itoa(*releaseID)
	}

	return &GetBugsQuery{Pagination: pagination, ProductID: strconv.Itoa(productID), State: state, BugID: bugTmp, ReleaseID: releaseTmp}
}
