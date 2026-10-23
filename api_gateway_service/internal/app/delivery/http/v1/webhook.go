package v1

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/mappers"
	"github.com/testovoleg/5s-microservice-template/pkg/constants"
	httpUtils "github.com/testovoleg/5s-microservice-template/pkg/http_utils"
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

		var bodyBytes []byte
		if c.Request().Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request().Body)
		}
		// span.SetAttributes(attribute.String("reqDto", string(bodyBytes)))
		// Restore the io.ReadCloser to its original state
		c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		defer c.Request().Body.Close()

		secretToken := httpUtils.StripBearer(c.Request().Header.Get(constants.HeaderSecretToken))

		companyUuid := ""
		apiUuid := ""
		if secretToken != "" {
			left, right, ok := strings.Cut(secretToken, "_")
			if !ok {
				msg := "CompanyUuid_ApiUuid not found in secret token"
				return h.traceWebhookErr(c, span, msg, errors.New(msg))
			}
			companyUuid = left
			apiUuid = right
		}

		if companyUuid == "" || apiUuid == "" {
			msg := "COMPANY_UUID or API_UUID is empty"
			return h.traceWebhookErr(c, span, msg, errors.New(msg))
		}

		err := h.svc.Commands.WebhookMethods.Webhook(ctx, mappers.NewWebhookCommand(companyUuid, apiUuid, bodyBytes))
		if err != nil {
			return h.traceWebhookErr(c, span, "Webhook", err)
		}

		h.metrics.Get("Success", metrics.HTTP).Inc()
		return c.NoContent(http.StatusOK)
	}
}
