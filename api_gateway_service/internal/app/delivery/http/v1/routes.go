package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (h *appHandlers) MapRoutes() {
	h.group.POST("/admin/api", h.AddApi())
	h.group.GET("/admin/api", h.GetApi())
	h.group.GET("/admin/api/:API_UUID", h.GetFullDataApi())
	h.group.PATCH("/admin/api/:API_UUID", h.UpdateApi())
	h.group.DELETE("/admin/api/:API_UUID", h.DeleteApi())

	h.group.POST("/example", h.Example())

	h.group.POST("/webhook", h.Webhook())

	h.group.GET("/swagger/*", echoSwagger.WrapHandler)

	h.group.GET("/doc", func(c echo.Context) error {
		c.HTML(200, h.redocly.Html)
		return nil
	})

	h.group.Any("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
}
