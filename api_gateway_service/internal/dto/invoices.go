package dto

import coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"

type InvoiceHandlerDto struct {
	ID          string `json:"id,omitempty"`
	Version     string `json:"version,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type InvoiceHandlersListReqDto struct{}

func InvoiceHandlerDtoFromGrpc(in *coreService.InvoiceHandler) *InvoiceHandlerDto {
	return &InvoiceHandlerDto{
		ID:          in.GetID(),
		Title:       in.GetTitle(),
		Description: in.GetDescription(),
		Version:     in.GetVersion(),
	}
}

func InvoiceHandlerDtoListFromGrpc(resp *coreService.InvoiceHandlersListRes) []*InvoiceHandlerDto {
	list := make([]*InvoiceHandlerDto, 0, len(resp.GetHandlers()))
	for _, v := range resp.GetHandlers() {
		list = append(list, InvoiceHandlerDtoFromGrpc(v))
	}
	return list
}
