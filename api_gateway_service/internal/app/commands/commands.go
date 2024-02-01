package commands

import (
	uuid "github.com/satori/go.uuid"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/dto"
)

type Commands struct {
	InvoiceHandlersList InvoiceHandlersListCmdHandler
	UpdateProduct       UpdateProductCmdHandler
	DeleteProduct       DeleteProductCmdHandler
}

func NewCommands(invoiceHandlersList InvoiceHandlersListCmdHandler, updateProduct UpdateProductCmdHandler, deleteProduct DeleteProductCmdHandler) *Commands {
	return &Commands{InvoiceHandlersList: invoiceHandlersList, UpdateProduct: updateProduct, DeleteProduct: deleteProduct}
}

type InvoiceHandlersListCommand struct {
	InvoiceHandlerListReqDto *dto.InvoiceHandlersListReqDto
}

func NewInvoiceHandlersListCommand(invoiceHandlerListReqDto *dto.InvoiceHandlersListReqDto) *InvoiceHandlersListCommand {
	return &InvoiceHandlersListCommand{InvoiceHandlerListReqDto: invoiceHandlerListReqDto}
}

type UpdateProductCommand struct {
	UpdateDto *dto.UpdateProductDto
}

func NewUpdateProductCommand(updateDto *dto.UpdateProductDto) *UpdateProductCommand {
	return &UpdateProductCommand{UpdateDto: updateDto}
}

type DeleteProductCommand struct {
	ProductID uuid.UUID `json:"productId" validate:"required"`
}

func NewDeleteProductCommand(productID uuid.UUID) *DeleteProductCommand {
	return &DeleteProductCommand{ProductID: productID}
}
