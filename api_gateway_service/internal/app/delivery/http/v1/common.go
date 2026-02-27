package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/config"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/app/service"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/middlewares"
	"github.com/testovoleg/5s-microservice-template/docs"
	httpErrors "github.com/testovoleg/5s-microservice-template/pkg/http_errors"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/metrics"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type appHandlers struct {
	group   *echo.Group
	log     logger.Logger
	mw      middlewares.MiddlewareManager
	cfg     *config.Config
	svc     *service.GatewayService
	v       *validator.Validate
	metrics *metrics.MetricsManager
	redocly *docs.Redocly
}

func NewAppHandlers(
	group *echo.Group,
	log logger.Logger,
	mw middlewares.MiddlewareManager,
	cfg *config.Config,
	svc *service.GatewayService,
	v *validator.Validate,
	metrics *metrics.MetricsManager,
	redocly *docs.Redocly,
) *appHandlers {
	return &appHandlers{group: group, log: log, mw: mw, cfg: cfg, svc: svc, v: v, metrics: metrics, redocly: redocly}
}

// Example
// @Tags 2. Данные
// @Summary Получить данные по запросу
// @Description Пример запроса для получения данных
// @Accept json
// @Produce json
// @Param company_uuid query string	false	"Идентификатор компании"
// @Param api_uuid query string true "Идентификатор интеграции"
// @Success 200
// @Failure      401  {object}  httpErrors.RestError "Авторизационные данные неверны или устарели"
// @Failure      500  {object}  httpErrors.RestError "Внутренняя ошибка сервера"
// @Router /example [post]
func (h *appHandlers) Example() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.Get("Example", metrics.HTTP).Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "appHandlers.Example")
		defer span.End()

		token := c.Request().Header.Get("Authorization")
		if strings.HasPrefix(strings.ToLower(token), "bearer ") {
			token = token[7:]
		}
		if token == "" {
			return h.traceErr(c, span, "security", errors.New("invalid token"))
		}

		companyUuid := c.QueryParam("company_uuid")
		apiUuid := c.QueryParam("api_uuid")
		if companyUuid == "" || apiUuid == "" {
			return h.traceErr(c, span, "params", errors.New("company_uuid or api_uuid is empty"))
		}

		// _, err := h.svc.Commands.Example.Handle(ctx, mappers.NewExampleCommand(token, companyUuid, apiUuid))
		// if err != nil {
		// 	return h.traceErr(c, span, "SendMessage", err)
		// }

		_ = ctx // only for example

		h.metrics.Get("Success", metrics.HTTP).Inc()
		return c.NoContent(http.StatusOK)
	}
}

func (h *appHandlers) traceWebhookErr(c echo.Context, span trace.Span, title string, err error) error {
	h.log.WarnMsg(title, err)
	span.SetStatus(codes.Error, "operation in api_gateway failed")
	span.RecordError(err)
	h.metrics.Get("WebhookError", metrics.HTTP).Inc()
	return c.NoContent(http.StatusOK)
}

func (h *appHandlers) traceErr(c echo.Context, span trace.Span, title string, err error) error {
	h.log.WarnMsg(title, err)
	span.SetStatus(codes.Error, "operation in api_gateway failed")
	span.RecordError(err)
	h.metrics.Get("Error", metrics.HTTP).Inc()
	return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
}
