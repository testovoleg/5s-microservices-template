package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/dto"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/mappers"
	"github.com/testovoleg/5s-microservice-template/pkg/metrics"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

// EventGet
// @Tags 4. Общее
// @Summary Получить статус ивента
// @Description Получение статуса асинхронного ивента по его идентификатору
// @Security BearerToken
// @Accept json
// @Produce json
// @Param company_uuid query string true "Идентификатор компании"
// @Param event_uuid query string true "Идентификатор ивента"
// @Success 200 {object} dto.EventDto
// @Failure 401 {object} httpErrors.RestError "Авторизационные данные неверны"
// @Failure 500 {object} httpErrors.RestError "Внутренняя ошибка сервера"
// @Router /common/event [get]
func (h *appHandlers) EventGet() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.Get("EventGet", metrics.HTTP).Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "appHandlers.EventGet")
		defer span.End()

		token, err := extractBearerToken(c)
		if err != nil {
			return h.traceErr(c, span, "security", err)
		}

		params := bindApiParams(c, token)

		reqDto := &dto.EventGetReqDto{}
		if err := c.Bind(reqDto); err != nil {
			return h.traceErr(c, span, "bind req", err)
		}

		if params.CompanyUuid == "" {
			params.CompanyUuid = reqDto.CompanyUuid
		}

		command := mappers.NewEventGetCommand(params, reqDto)
		response, err := h.svc.Commands.CommonMethods.EventGet(ctx, command)
		if err != nil {
			return h.traceErr(c, span, "EventGet", err)
		}

		h.metrics.Get("Success", metrics.HTTP).Inc()
		return c.JSON(http.StatusOK, response)
	}
}
