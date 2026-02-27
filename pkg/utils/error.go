package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ErrorResponseBody struct {
	StatusCode string `json:"statusCode"`
	Timestamp  string `json:"timestamp"`
	Path       string `json:"path"`
	Message    string `json:"message"`
}

// how to use
//
//	if res.StatusCode != http.StatusOK {
//		return utils.NewError(res, span, "Error was during on Admin API")
//	}
// or
//	if res.StatusCode != http.StatusOK {
// 		return nil, utils.GetErrorMessage(res, span)
// 	}

func NewError(res *http.Response, span trace.Span, opt ...string) error {
	message := ""
	if len(opt) > 0 && opt[0] != "" {
		message = opt[0]
	}
	err := parseErrorMessage(res)
	TraceErr(span, err)
	if message != "" {
		return errors.New(message)
	}

	return err
}

func parseErrorMessage(res *http.Response) error {
	if res == nil || res.Body == nil {
		return errors.New("empty HTTP response")
	}

	body, err := io.ReadAll(res.Body)
	_ = res.Body.Close()

	if err != nil {
		return err
	}

	if len(body) == 0 {
		return errors.New(res.Status)
	}

	var customErr ErrorResponseBody
	if json.Unmarshal(body, &customErr) == nil {
		switch {
		case customErr.Message != "":
			return errors.New(customErr.Message)
		case customErr.StatusCode != "":
			return errors.New(customErr.StatusCode)
		case customErr.Path != "":
			return errors.New(customErr.Path)
		}
	}

	var defaultErr struct {
		ErrError string `json:"err_error"`
	}
	if json.Unmarshal(body, &defaultErr) == nil && defaultErr.ErrError != "" {
		return errors.New(defaultErr.ErrError)
	}

	return errors.New(string(body))
}

func TraceErr(span trace.Span, err error) error {
	if span != nil && err != nil {
		span.SetStatus(codes.Error, "operation in core_service failed")
		span.RecordError(err)
	}
	return err
}
