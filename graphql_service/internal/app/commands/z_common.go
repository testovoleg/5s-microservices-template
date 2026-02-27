package commands

import (
	model "github.com/testovoleg/5s-microservice-template/graphql_service/internal/graph_model"
)

type GetTipicalDataCommand struct {
	Params *model.GeneralParamsInput
}

type PostTipicalMutationCommand struct {
	Params *model.GeneralParamsInput
}
