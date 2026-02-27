package v1

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/mappers"
	"github.com/testovoleg/5s-microservice-template/pkg/constants"
	"github.com/testovoleg/5s-microservice-template/pkg/metrics"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

// Webhook
// @Tags 3. Уведомления
// @Summary Уведомление
// @Description Вызывается внешним сервисом для уведомления о событиях, на которые осуществлена подписка
// @Accept json
// @Produce json
// @Param input body string true "Тело уведомления"
// @Success 200
// @Router /webhook [post]
func (h *appHandlers) Webhook() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.Get("Webhook", metrics.HTTP).Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "appHandlers.Webhook")
		defer span.End()

		//test webhook
		var bodyBytes []byte
		if c.Request().Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request().Body)
		}
		// span.SetAttributes(attribute.String("reqDto", string(bodyBytes)))
		// Restore the io.ReadCloser to its original state
		c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		//test webhook
		defer c.Request().Body.Close()

		secretToken := c.Request().Header.Get(constants.HeaderSecretToken)
		if strings.HasPrefix(strings.ToLower(secretToken), "bearer ") {
			secretToken = secretToken[7:]
		}

		companyUuid := ""
		apiUuid := ""
		if secretToken != "" {
			pattern := regexp.MustCompile(`(.+)_(.+)`)
			finded := pattern.FindStringSubmatch(secretToken)
			if len(finded) == 0 {
				desc := "CompanyUuid_ApiUuid not found in secret token"
				return h.traceWebhookErr(c, span, desc, errors.New(desc))
			}

			companyUuid = finded[1]
			if len(finded) >= 2 {
				apiUuid = finded[2]
			}
		}

		if companyUuid == "" || apiUuid == "" {
			desc := "COMPANY_UUID or API_UUID is empty"
			return h.traceWebhookErr(c, span, desc, errors.New(desc))
		}

		err := h.svc.Commands.Webhook.Handle(ctx, mappers.NewWebhookCommand(companyUuid, apiUuid, bodyBytes))
		if err != nil {
			return h.traceWebhookErr(c, span, "Webhook", err)
		}

		h.metrics.Get("Success", metrics.HTTP).Inc()
		return c.NoContent(http.StatusOK)
	}
}
