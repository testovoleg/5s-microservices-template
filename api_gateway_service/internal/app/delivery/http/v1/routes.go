package v1

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (h *appHandlers) MapRoutes() {
	h.group.GET("/invoice/handler/list", h.InvoiceHandlerList())

	h.group.GET("/swagger/*", echoSwagger.WrapHandler)

	h.group.GET("/doc", func(c echo.Context) error {
		c.HTML(200, fmt.Sprintf(redoclyHTML, h.cfg.Resources.REDOCLY_JSON))
		return nil
	})

	h.group.Any("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
}
