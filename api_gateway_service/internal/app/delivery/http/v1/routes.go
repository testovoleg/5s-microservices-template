package v1

import (
	"github.com/labstack/echo/v4"
)

func (h *appHandlers) MapRoutes() {
	h.group.POST("/admin/api", h.AddApi())
	h.group.GET("/admin/api", h.GetApi())
	h.group.GET("/admin/api/:API_UUID", h.GetFullDataApi())
	h.group.PATCH("/admin/api/:API_UUID", h.UpdateApi())
	h.group.DELETE("/admin/api/:API_UUID", h.DeleteApi())

	h.group.POST("/example", h.Example())

	h.group.GET("/common/event", h.EventGet())

	h.group.POST("/webhook", h.Webhook())

	h.group.GET("/swagger/*", h.swaggerHandler())

	h.group.GET("/doc", func(c echo.Context) error {
		c.HTML(200, h.redocly.Html)
		return nil
	})
	// Kubernetes probes (/ready, /live) listen on probes.port — not on business HTTP routes.
}
