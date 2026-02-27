package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/dto"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/mappers"
	"github.com/testovoleg/5s-microservice-template/pkg/metrics"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

// AddApi
// @Tags 1. Администрирование
// @Summary Добавить новую интеграцию
// @Description Добавление новой интеграции
// @Security BearerToken
// @Accept json
// @Produce json
// @Param company_uuid query string	false	"Идентификатор компании"
// @Param input body dto.AddApiReqDto true "Тело запроса"
// @Success 200 {object} dto.ApiDto
// @Failure      401  {object}  httpErrors.RestError "Авторизационные данные неверны или устарели"
// @Failure      500  {object}  httpErrors.RestError "Внутренняя ошибка сервера"
// @Router /admin/api [post]
func (h *appHandlers) AddApi() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.Get("AddApi", metrics.HTTP).Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "appHandlers.AddApi")
		defer span.End()

		token := c.Request().Header.Get("Authorization")
		if strings.HasPrefix(strings.ToLower(token), "bearer ") {
			token = token[7:]
		}

		if token == "" {
			return h.traceErr(c, span, "security", errors.New("invalid token"))
		}

		companyUuid := c.QueryParam("company_uuid")
		if companyUuid == "" {
			return h.traceErr(c, span, "params", errors.New("company_uuid is empty"))
		}

		reqDto := &dto.AddApiReqDto{}
		if err := c.Bind(reqDto); err != nil {
			return h.traceErr(c, span, "Bind", err)
		}

		if err := h.v.StructCtx(ctx, reqDto); err != nil {
			return h.traceErr(c, span, "validate", err)
		}

		res, err := h.svc.Commands.AddApi.Handle(ctx, mappers.NewAddApiCommand(token, companyUuid, reqDto))
		if err != nil {
			return h.traceErr(c, span, "AddApi", err)
		}

		h.metrics.Get("Success", metrics.HTTP).Inc()
		return c.JSON(http.StatusOK, res)
	}
}

// GetApi
// @Tags 1. Администрирование
// @Summary Получить список интеграций
// @Description Получение списка интеграций по компании
// @Security BearerToken
// @Accept json
// @Produce json
// @Param company_uuid query string	false	"Идентификатор компании"
// @Success 200 {array} dto.ApiDto
// @Failure      401  {object}  httpErrors.RestError "Авторизационные данные неверны или устарели"
// @Failure      500  {object}  httpErrors.RestError "Внутренняя ошибка сервера"
// @Router /admin/api [get]
func (h *appHandlers) GetApi() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.Get("GetApi", metrics.HTTP).Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "appHandlers.GetApi")
		defer span.End()

		token := c.Request().Header.Get("Authorization")
		if strings.HasPrefix(strings.ToLower(token), "bearer ") {
			token = token[7:]
		}

		if token == "" {
			return h.traceErr(c, span, "security", errors.New("invalid token"))
		}

		reqDto := &dto.GetApiReqDto{}
		if err := c.Bind(reqDto); err != nil {
			return h.traceErr(c, span, "Bind", err)
		}

		if err := h.v.StructCtx(ctx, reqDto); err != nil {
			return h.traceErr(c, span, "validate", err)
		}

		res, err := h.svc.Commands.GetApi.Handle(ctx, mappers.NewGetApiCommand(token, reqDto))
		if err != nil {
			return h.traceErr(c, span, "GetApi", err)
		}

		h.metrics.Get("Success", metrics.HTTP).Inc()

		return c.JSON(http.StatusOK, res)
	}
}

// GetFullDataApi
// @Tags 1. Администрирование
// @Summary Получить данные по интеграции
// @Description Получение подробных данных по интеграции
// @Security BearerToken
// @Accept json
// @Produce json
// @Param company_uuid query string	false	"Идентификатор компании"
// @Param API_UUID path string true "Идентификатор интеграции"
// @Success 200 {object} dto.ApiFullDto
// @Failure      401  {object}  httpErrors.RestError "Авторизационные данные неверны или устарели"
// @Failure      500  {object}  httpErrors.RestError "Внутренняя ошибка сервера"
// @Router /admin/api/{API_UUID} [get]
func (h *appHandlers) GetFullDataApi() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.Get("GetFullDataApi", metrics.HTTP).Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "appHandlers.GetFullDataApi")
		defer span.End()

		token := c.Request().Header.Get("Authorization")
		if strings.HasPrefix(strings.ToLower(token), "bearer ") {
			token = token[7:]
		}

		if token == "" {
			return h.traceErr(c, span, "security", errors.New("invalid token"))
		}

		reqDto := &dto.GetFullApiReqDto{}
		if err := c.Bind(reqDto); err != nil {
			return h.traceErr(c, span, "Bind", err)
		}

		if err := h.v.StructCtx(ctx, reqDto); err != nil {
			return h.traceErr(c, span, "validate", err)
		}

		res, err := h.svc.Commands.GetFullApi.Handle(ctx, mappers.NewGetFullApiCommand(token, reqDto))
		if err != nil {
			return h.traceErr(c, span, "GetFullDataApi", err)
		}

		h.metrics.Get("Success", metrics.HTTP).Inc()
		return c.JSON(http.StatusOK, res)
	}
}

// UpdateApi
// @Tags 1. Администрирование
// @Summary Обновить интеграцию
// @Description Обновление данных интеграции
// @Security BearerToken
// @Accept json
// @Produce json
// @Param company_uuid query string	false	"Идентификатор компании"
// @Param API_UUID path string true "Идентификатор интеграции"
// @Param input body dto.UpdateApiReqDto true "Тело запроса"
// @Success 200 {object} dto.ApiDto
// @Failure      401  {object}  httpErrors.RestError "Авторизационные данные неверны или устарели"
// @Failure      500  {object}  httpErrors.RestError "Внутренняя ошибка сервера"
// @Router /admin/api/{API_UUID} [patch]
func (h *appHandlers) UpdateApi() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.Get("UpdateApi", metrics.HTTP).Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "appHandlers.UpdateApi")
		defer span.End()

		token := c.Request().Header.Get("Authorization")
		if strings.HasPrefix(strings.ToLower(token), "bearer ") {
			token = token[7:]
		}

		if token == "" {
			return h.traceErr(c, span, "security", errors.New("invalid token"))
		}

		companyUuid := c.QueryParam("company_uuid")
		apiUuid := c.Param("API_UUID")
		if companyUuid == "" || apiUuid == "" {
			return h.traceErr(c, span, "params", errors.New("company_uuid or API_UUID is empty"))
		}

		reqDto := &dto.UpdateApiReqDto{}
		if err := c.Bind(reqDto); err != nil {
			return h.traceErr(c, span, "Bind", err)
		}

		if err := h.v.StructCtx(ctx, reqDto); err != nil {
			return h.traceErr(c, span, "validate", err)
		}

		res, err := h.svc.Commands.UpdateApi.Handle(ctx, mappers.NewUpdateApiCommand(token, companyUuid, apiUuid, reqDto))
		if err != nil {
			return h.traceErr(c, span, "UpdateApi", err)
		}

		h.metrics.Get("Success", metrics.HTTP).Inc()
		return c.JSON(http.StatusOK, res)
	}
}

// DeleteApi
// @Tags 1. Администрирование
// @Summary Удалить интеграцию
// @Description Удаление интеграции
// @Security BearerToken
// @Accept json
// @Produce json
// @Param company_uuid query string	false	"Идентификатор компании"
// @Param API_UUID path string true "Идентификатор интеграции"
// @Success 200
// @Failure      401  {object}  httpErrors.RestError "Авторизационные данные неверны или устарели"
// @Failure      500  {object}  httpErrors.RestError "Внутренняя ошибка сервера"
// @Router /admin/api/{API_UUID} [delete]
func (h *appHandlers) DeleteApi() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.Get("DeleteApi", metrics.HTTP).Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "appHandlers.DeleteApi")
		defer span.End()

		token := c.Request().Header.Get("Authorization")
		if strings.HasPrefix(strings.ToLower(token), "bearer ") {
			token = token[7:]
		}

		if token == "" {
			return h.traceErr(c, span, "security", errors.New("invalid token"))
		}

		reqDto := &dto.DeleteApiReqDto{}
		if err := c.Bind(reqDto); err != nil {
			return h.traceErr(c, span, "Bind", err)
		}

		if err := h.v.StructCtx(ctx, reqDto); err != nil {
			return h.traceErr(c, span, "validate", err)
		}

		err := h.svc.Commands.DeleteApi.Handle(ctx, mappers.NewDeleteApiCommand(token, reqDto))
		if err != nil {
			return h.traceErr(c, span, "DeleteApi", err)
		}

		h.metrics.Get("Success", metrics.HTTP).Inc()
		return c.NoContent(http.StatusOK)
	}
}
