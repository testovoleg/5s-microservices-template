package v1

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/config"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/app/commands"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/app/service"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/dto"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/metrics"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/middlewares"
	httpErrors "github.com/testovoleg/5s-microservice-template/pkg/http_errors"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

type appHandlers struct {
	group   *echo.Group
	log     logger.Logger
	mw      middlewares.MiddlewareManager
	cfg     *config.Config
	svc     *service.Service
	v       *validator.Validate
	metrics *metrics.ApiGatewayMetrics
}

func NewAppHandlers(
	group *echo.Group,
	log logger.Logger,
	mw middlewares.MiddlewareManager,
	cfg *config.Config,
	svc *service.Service,
	v *validator.Validate,
	metrics *metrics.ApiGatewayMetrics,
) *appHandlers {
	return &appHandlers{group: group, log: log, mw: mw, cfg: cfg, svc: svc, v: v, metrics: metrics}
}

// InvoiceHandlerDto
// @Tags Invoices
// @Summary List of current invoice handlers
// @Description List of current invoice handlers
// @Accept json
// @Produce json
// @Success 200 {array} dto.InvoiceHandlerDto
// @Router /invoice/handler/list [get]
func (h *appHandlers) InvoiceHandlerList() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.InvoiceHandlersListHttpRequests.Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "appHandlers.InvoiceHandlersList")
		defer span.Finish()

		reqDto := &dto.InvoiceHandlersListReqDto{}
		if err := c.Bind(reqDto); err != nil {
			h.log.WarnMsg("Bind", err)
			h.traceErr(span, err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.v.StructCtx(ctx, reqDto); err != nil {
			h.log.WarnMsg("validate", err)
			h.traceErr(span, err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		res, err := h.svc.Commands.InvoiceHandlersList.Handle(ctx, commands.NewInvoiceHandlersListCommand(reqDto))
		if err != nil {
			h.log.WarnMsg("CreateProduct", err)
			h.metrics.ErrorHttpRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, res)
	}
}

func (h *appHandlers) traceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_code", err.Error())
	h.metrics.ErrorHttpRequests.Inc()
}
