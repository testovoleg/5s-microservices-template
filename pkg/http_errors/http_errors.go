package httpErrors

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	ErrBadRequest          = "Bad request"
	ErrNotFound            = "Not Found"
	ErrUnauthorized        = "Unauthorized"
	ErrRequestTimeout      = "Request Timeout"
	ErrInvalidEmail        = "Invalid email"
	ErrInvalidPassword     = "Invalid password"
	ErrInvalidField        = "Invalid field"
	ErrInternalServerError = "Internal Server Error"
)

var (
	BadRequest          = errors.New("Bad request")
	WrongCredentials    = errors.New("Wrong Credentials")
	NotFound            = errors.New("Not Found")
	Unauthorized        = errors.New("Unauthorized")
	Forbidden           = errors.New("Forbidden")
	InternalServerError = errors.New("Internal Server Error")
)

// RestErr Rest error interface
type RestErr interface {
	Status() int
	Error() string
	Causes() interface{}
	ErrBody() RestError
}

// RestError Rest error struct
type RestError struct {
	ErrStatus  int         `json:"status,omitempty"`
	ErrError   string      `json:"error,omitempty"`
	ErrMessage interface{} `json:"message,omitempty"`
	Timestamp  time.Time   `json:"timestamp,omitempty"`
}

// ErrBody Error body
func (e RestError) ErrBody() RestError {
	return e
}

// Error  Error() interface method
func (e RestError) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - causes: %v", e.ErrStatus, e.ErrError, e.ErrMessage)
}

// Status Error status
func (e RestError) Status() int {
	return e.ErrStatus
}

// Causes RestError Causes
func (e RestError) Causes() interface{} {
	return e.ErrMessage
}

// NewRestError New Rest Error
func NewRestError(status int, err string, causes interface{}, debug bool) RestErr {
	restError := RestError{
		ErrStatus: status,
		ErrError:  err,
		Timestamp: time.Now().UTC(),
	}
	if debug {
		restError.ErrMessage = causes
	}
	return restError
}

// NewRestErrorWithMessage New Rest Error With Message
func NewRestErrorWithMessage(status int, err string, causes interface{}) RestErr {
	return RestError{
		ErrStatus:  status,
		ErrError:   err,
		ErrMessage: causes,
		Timestamp:  time.Now().UTC(),
	}
}

// NewRestErrorFromBytes New Rest Error From Bytes
func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
	var apiErr RestError
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return apiErr, nil
}

// NewBadRequestError New Bad Request Error
func NewBadRequestError(ctx echo.Context, causes interface{}, debug bool) error {
	restError := RestError{
		ErrStatus: http.StatusBadRequest,
		ErrError:  BadRequest.Error(),
		Timestamp: time.Now().UTC(),
	}
	if debug {
		restError.ErrMessage = causes
	}
	return ctx.JSON(http.StatusBadRequest, restError)
}

// NewNotFoundError New Not Found Error
func NewNotFoundError(ctx echo.Context, causes interface{}, debug bool) error {
	restError := RestError{
		ErrStatus: http.StatusNotFound,
		ErrError:  NotFound.Error(),
		Timestamp: time.Now().UTC(),
	}
	if debug {
		restError.ErrMessage = causes
	}
	return ctx.JSON(http.StatusNotFound, restError)
}

// NewUnauthorizedError New Unauthorized Error
func NewUnauthorizedError(ctx echo.Context, causes interface{}, debug bool) error {

	restError := RestError{
		ErrStatus: http.StatusUnauthorized,
		ErrError:  Unauthorized.Error(),
		Timestamp: time.Now().UTC(),
	}
	if debug {
		restError.ErrMessage = causes
	}
	return ctx.JSON(http.StatusUnauthorized, restError)
}

// NewForbiddenError New Forbidden Error
func NewForbiddenError(ctx echo.Context, causes interface{}, debug bool) error {

	restError := RestError{
		ErrStatus: http.StatusForbidden,
		ErrError:  Forbidden.Error(),
		Timestamp: time.Now().UTC(),
	}
	if debug {
		restError.ErrMessage = causes
	}
	return ctx.JSON(http.StatusForbidden, restError)
}

// NewInternalServerError New Internal Server Error
func NewInternalServerError(ctx echo.Context, causes interface{}, debug bool) error {

	restError := RestError{
		ErrStatus: http.StatusInternalServerError,
		ErrError:  InternalServerError.Error(),
		Timestamp: time.Now().UTC(),
	}
	if debug {
		restError.ErrMessage = causes
	}
	return ctx.JSON(http.StatusInternalServerError, restError)
}

// ParseErrors Parser of error string messages returns RestError
func ParseErrors(err error, debug bool) RestErr {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return NewRestError(http.StatusNotFound, ErrNotFound, err.Error(), debug)
	case errors.Is(err, context.DeadlineExceeded):
		return NewRestError(http.StatusRequestTimeout, ErrRequestTimeout, err.Error(), debug)
	case errors.Is(err, Unauthorized):
		return NewRestError(http.StatusUnauthorized, ErrUnauthorized, err.Error(), debug)
	case errors.Is(err, WrongCredentials):
		return NewRestError(http.StatusUnauthorized, ErrUnauthorized, err.Error(), debug)
	case strings.Contains(strings.ToLower(err.Error()), "sqlstate"):
		return parseSqlErrors(err, debug)
	case strings.Contains(strings.ToLower(err.Error()), "field validation"):
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return NewRestError(http.StatusBadRequest, ErrBadRequest, validationErrors.Error(), debug)
		}
		return parseValidatorError(err, debug)
	case strings.Contains(strings.ToLower(err.Error()), "required header"):
		return NewRestError(http.StatusBadRequest, ErrBadRequest, err.Error(), debug)
	case strings.Contains(strings.ToLower(err.Error()), "base64"):
		return NewRestError(http.StatusBadRequest, ErrBadRequest, err.Error(), debug)
	case strings.Contains(strings.ToLower(err.Error()), "unmarshal"):
		return NewRestError(http.StatusBadRequest, ErrBadRequest, err.Error(), debug)
	case strings.Contains(strings.ToLower(err.Error()), "uuid"):
		return NewRestError(http.StatusBadRequest, ErrBadRequest, err.Error(), debug)
	case strings.Contains(strings.ToLower(err.Error()), "cookie"):
		return NewRestError(http.StatusUnauthorized, ErrUnauthorized, err.Error(), debug)
	case strings.Contains(strings.ToLower(err.Error()), "token"):
		return NewRestError(http.StatusUnauthorized, ErrUnauthorized, err.Error(), debug)
	case strings.Contains(strings.ToLower(err.Error()), "bcrypt"):
		return NewRestError(http.StatusBadRequest, ErrBadRequest, err.Error(), debug)
	case strings.Contains(strings.ToLower(err.Error()), "no documents in result"):
		return NewRestError(http.StatusNotFound, ErrNotFound, err.Error(), debug)
	default:
		if restErr, ok := err.(*RestError); ok {
			return restErr
		}
		return NewRestError(http.StatusInternalServerError, ErrInternalServerError, errors.Cause(err).Error(), debug)
	}
}

// clientCausedSqlStateCodes are Postgres error classes caused by bad client input
// (integrity constraints, malformed data). Everything else is a server-side problem.
var clientCausedSqlStateCodes = map[string]bool{
	"23505": true, // unique_violation
	"23503": true, // foreign_key_violation
	"23502": true, // not_null_violation
	"23514": true, // check_violation
	"22001": true, // string_data_right_truncation
	"22P02": true, // invalid_text_representation
}

func parseSqlErrors(err error, debug bool) RestErr {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && clientCausedSqlStateCodes[pgErr.Code] {
		return NewRestError(http.StatusBadRequest, ErrBadRequest, err, debug)
	}
	return NewRestError(http.StatusInternalServerError, ErrInternalServerError, err, debug)
}

func parseValidatorError(err error, debug bool) RestErr {
	if strings.Contains(err.Error(), "Password") {
		return NewRestError(http.StatusBadRequest, ErrInvalidPassword, err, debug)
	}

	if strings.Contains(err.Error(), "Email") {
		return NewRestError(http.StatusBadRequest, ErrInvalidEmail, err, debug)
	}

	return NewRestError(http.StatusBadRequest, ErrInvalidField, err, debug)
}

// ErrorResponse Error response
func ErrorResponse(err error, debug bool) (int, interface{}) {
	return ParseErrors(err, debug).Status(), ParseErrors(err, debug)
}

// ErrorCtxResponse Error response object and status code
func ErrorCtxResponse(ctx echo.Context, err error, debug bool) error {
	restErr := ParseErrors(err, debug)
	return ctx.JSON(restErr.Status(), restErr)
}

func ErrorWrap(err error, span trace.Span, opt ...string) error {
	if err != nil && len(opt) > 0 && opt[0] != "" {
		err = errors.Wrap(err, opt[0])
	}
	return UpdateSpan(span, err)
}

func HttpResponse(res *http.Response, span trace.Span, opt ...string) error {
	return ErrorWrap(ParseHttpResponse(res), span, opt...)
}

// HttpBody is like HttpResponse but for a body that was already read from res.
func HttpBody(status string, body []byte, span trace.Span, opt ...string) error {
	return ErrorWrap(ParseHttpBody(status, body), span, opt...)
}

func UpdateSpan(span trace.Span, err error) error {
	if span != nil && err != nil {
		span.SetStatus(codes.Error, "operation in core_service failed")
		span.RecordError(err)
	}
	return AddTraceID(err, span)
}

func AddTraceID(err error, s trace.Span) error {
	if err == nil {
		return nil
	}

	if s == nil {
		return err
	}

	if strings.Contains(err.Error(), "TraceID:") {
		return err
	}

	return errors.New(err.Error() + " TraceID: " + s.SpanContext().TraceID().String())
}

func ParseHttpResponse(res *http.Response) error {
	if res == nil || res.Body == nil {
		return errors.New("empty HTTP response")
	}

	body, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		return err
	}

	return ParseHttpBody(res.Status, body)
}

const maxHTTPErrorBodyLen = 2048

// ParseHttpBody builds an error from an already-read HTTP response body.
// It understands common upstream shapes (Admin/Auth/Keycloak/5s APIs) and
// falls back to a truncated raw body when the payload is not structured.
func ParseHttpBody(status string, body []byte) error {
	body = bytes.TrimSpace(body)
	if len(body) == 0 {
		if status != "" {
			return errors.New(status)
		}
		return errors.New("empty HTTP response")
	}

	if detail := extractHTTPErrorDetail(body); detail != "" {
		return formatHTTPError(status, detail)
	}

	return formatHTTPError(status, truncateRunes(string(body), maxHTTPErrorBodyLen))
}

func formatHTTPError(status, detail string) error {
	detail = strings.TrimSpace(detail)
	status = strings.TrimSpace(status)
	switch {
	case status == "" && detail == "":
		return errors.New("empty HTTP response")
	case status == "":
		return errors.New(detail)
	case detail == "" || detail == status:
		return errors.New(status)
	case strings.HasPrefix(detail, status):
		return errors.New(detail)
	default:
		return errors.New(status + ": " + detail)
	}
}

// extractHTTPErrorDetail pulls a human-readable detail from known JSON error shapes.
// Returns "" when body is not JSON or contains no recognizable error fields.
func extractHTTPErrorDetail(body []byte) string {
	var payload struct {
		Message          any    `json:"message"`
		Error            any    `json:"error"`
		ErrError         any    `json:"err_error"`
		ErrorDescription any    `json:"error_description"`
		StatusCode       any    `json:"statusCode"`
		Path             string `json:"path"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return ""
	}

	parts := make([]string, 0, 4)
	appendUnique := func(s string) {
		s = strings.TrimSpace(s)
		if s == "" {
			return
		}
		for _, p := range parts {
			if p == s {
				return
			}
		}
		parts = append(parts, s)
	}

	// Prefer descriptive fields first, then codes/path for context.
	for _, v := range []any{
		payload.Message,
		payload.ErrorDescription,
		payload.Error,
		payload.ErrError,
	} {
		appendUnique(anyToString(v))
	}
	appendUnique(anyToString(payload.StatusCode))
	appendUnique(payload.Path)

	return strings.Join(parts, "; ")
}

func anyToString(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(t)
	case float64:
		if t == float64(int64(t)) {
			return fmt.Sprintf("%d", int64(t))
		}
		return fmt.Sprintf("%v", t)
	case json.Number:
		return t.String()
	case bool:
		return fmt.Sprintf("%t", t)
	case []any:
		items := make([]string, 0, len(t))
		for _, item := range t {
			if s := anyToString(item); s != "" {
				items = append(items, s)
			}
		}
		return strings.Join(items, ", ")
	case map[string]any:
		// Nested objects like {"message":"..."}; prefer common keys.
		for _, key := range []string{"message", "error", "error_description", "description"} {
			if s := anyToString(t[key]); s != "" {
				return s
			}
		}
		b, err := json.Marshal(t)
		if err != nil {
			return ""
		}
		return truncateRunes(string(b), 512)
	default:
		return strings.TrimSpace(fmt.Sprintf("%v", t))
	}
}

func truncateRunes(s string, max int) string {
	if max <= 0 || utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	return string(runes[:max]) + "…"
}
