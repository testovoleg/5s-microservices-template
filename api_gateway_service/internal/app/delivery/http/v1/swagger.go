package v1

import (
	_ "embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type swaggerUIData struct {
	AppName string
}

// swaggerHandler — stock echo-swagger со своей index-страницей (swagger_ui.html):
// заголовок из http.title и плагин, который не строит curl-сниппет для больших тел.
// doc.json, js и css отдаёт echo-swagger.
func (h *appHandlers) swaggerHandler() echo.HandlerFunc {
	indexTpl := template.Must(template.New("swagger_ui").Parse(swaggerUIHTML))

	handler := echoSwagger.EchoWrapHandler(
		echoSwagger.PersistAuthorization(true),
	)

	return func(c echo.Context) error {
		switch c.Param("*") {
		case "", "index.html":
			var buf strings.Builder
			if err := indexTpl.Execute(&buf, swaggerUIData{AppName: h.cfg.Http.Title}); err != nil {
				return err
			}
			return c.HTML(http.StatusOK, buf.String())
		default:
			return handler(c)
		}
	}
}

// swaggerUIHTML — stock echo-swagger index + плагин для больших тел запроса.
// Это Go-шаблон (html/template), а не готовый HTML.
//
//go:embed swagger_ui.html
var swaggerUIHTML string
