package commands

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Commands struct {
	InvoiceHandlersList InvoiceHandlersListCmdHandler
	UpdateProduct       UpdateProductCmdHandler
	DeleteProduct       DeleteProductCmdHandler
}

func NewCommands(
	invoiceHandlersList InvoiceHandlersListCmdHandler,
	updateProduct UpdateProductCmdHandler,
	deleteProduct DeleteProductCmdHandler,
) *Commands {
	return &Commands{InvoiceHandlersList: invoiceHandlersList, UpdateProduct: updateProduct, DeleteProduct: deleteProduct}
}

type InvoiceHandlersListCommand struct{}

func NewInvoiceHandlersListCommand() *InvoiceHandlersListCommand {
	return &InvoiceHandlersListCommand{}
}

type UpdateProductCommand struct {
	ProductID   string    `json:"productId" bson:"_id,omitempty"`
	Name        string    `json:"name,omitempty" bson:"name,omitempty" validate:"required,min=3,max=250"`
	Description string    `json:"description,omitempty" bson:"description,omitempty" validate:"required,min=3,max=500"`
	Price       float64   `json:"price,omitempty" bson:"price,omitempty" validate:"required"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

func NewUpdateProductCommand(productID string, name string, description string, price float64, updatedAt time.Time) *UpdateProductCommand {
	return &UpdateProductCommand{ProductID: productID, Name: name, Description: description, Price: price, UpdatedAt: updatedAt}
}

type DeleteProductCommand struct {
	ProductID uuid.UUID `json:"productId" bson:"_id,omitempty"`
}

func NewDeleteProductCommand(productID uuid.UUID) *DeleteProductCommand {
	return &DeleteProductCommand{ProductID: productID}
}
