package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (h *appHandlers) MapRoutes() {
	h.group.GET("/invoice/handler/list", h.InvoiceHandlerList())

	h.group.GET("/swagger/*", echoSwagger.WrapHandler)

	h.group.Any("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
}
