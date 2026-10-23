package httpUtils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/testovoleg/5s-microservice-template/pkg/constants"
	httpErrors "github.com/testovoleg/5s-microservice-template/pkg/http_errors"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"go.opentelemetry.io/otel/trace"
)

// SendOpts configures auth and timeouts for SendHTTPRequest.
type SendOpts struct {
	BearerToken   string
	BasicUser     string
	BasicPassword string
	// ErrorPrefix is appended to HTTP error wrapping (optional).
	ErrorPrefix string
	// Timeout overrides the default HttpClientTimeout when > 0.
	Timeout time.Duration
	// ContentType sets Content-Type. If empty and body != nil, application/json is used.
	ContentType string
	// OmitContentType skips setting Content-Type (e.g. S3 presigned PUT).
	OmitContentType bool
	// Headers are extra request headers (e.g. X-Amz-Tagging).
	Headers map[string]string
}

// SendHTTPRequest performs an HTTP call with tracing injection and optional auth.
// Responses outside 2xx are turned into errors via httpErrors.HttpBody.
//
// Dataset-style Basic Auth to company.ApiData:
//
//	httpUtils.SendHTTPRequest(ctx, method, url, body, httpUtils.SendOpts{
//	    BasicUser: company.ApiData.Username, BasicPassword: company.ApiData.Password,
//	})
func SendHTTPRequest(ctx context.Context, method, requestURL string, body []byte, opts SendOpts) ([]byte, error) {
	ctx, span := tracing.StartSpan(ctx, "httpUtils.SendHTTPRequest")
	defer span.End()

	status, respBody, _, err := doHTTP(ctx, span, method, requestURL, body, opts)
	if err != nil {
		return nil, err
	}
	if status < 200 || status >= 300 {
		return nil, httpErrors.HttpBody(formatStatus(status), respBody, span, opts.ErrorPrefix)
	}
	return respBody, nil
}

// SendHTTPRequestRaw is like SendHTTPRequest but does not fail on non-OK status.
// Caller inspects status (e.g. special-case empty Admin property).
func SendHTTPRequestRaw(ctx context.Context, method, requestURL string, body []byte, opts SendOpts) (status int, respBody []byte, err error) {
	ctx, span := tracing.StartSpan(ctx, "httpUtils.SendHTTPRequest")
	defer span.End()
	status, respBody, _, err = doHTTP(ctx, span, method, requestURL, body, opts)
	return status, respBody, err
}

// SendHTTPRequestFull returns status, body and response headers. Non-OK status is not an error.
func SendHTTPRequestFull(ctx context.Context, method, requestURL string, body []byte, opts SendOpts) (status int, respBody []byte, header http.Header, err error) {
	ctx, span := tracing.StartSpan(ctx, "httpUtils.SendHTTPRequest")
	defer span.End()
	return doHTTP(ctx, span, method, requestURL, body, opts)
}

func doHTTP(ctx context.Context, span trace.Span, method, requestURL string, body []byte, opts SendOpts) (int, []byte, http.Header, error) {
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = constants.HttpClientTimeout * time.Second
	}
	client := http.Client{Timeout: timeout}

	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL, reader)
	if err != nil {
		return 0, nil, nil, httpErrors.ErrorWrap(err, span, "http.NewRequest")
	}

	if opts.BearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+opts.BearerToken)
	}
	if opts.BasicUser != "" || opts.BasicPassword != "" {
		req.SetBasicAuth(opts.BasicUser, opts.BasicPassword)
	}
	if !opts.OmitContentType {
		switch {
		case opts.ContentType != "":
			req.Header.Set("Content-Type", opts.ContentType)
		case body != nil:
			req.Header.Set("Content-Type", "application/json")
		}
	}
	for k, v := range opts.Headers {
		req.Header.Set(k, v)
	}

	if headers, injErr := tracing.InjectTextMapCarrier(ctx); injErr == nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	if tid := span.SpanContext().TraceID(); tid.IsValid() {
		req.Header.Set("x-trace-id", tid.String())
	}

	res, err := client.Do(req)
	if err != nil {
		return 0, nil, nil, httpErrors.ErrorWrap(err, span, "client.Do")
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, nil, res.Header, httpErrors.ErrorWrap(err, span, "io.ReadAll")
	}

	return res.StatusCode, respBody, res.Header, nil
}

func formatStatus(code int) string {
	if text := http.StatusText(code); text != "" {
		return fmt.Sprintf("%d %s", code, text)
	}
	return fmt.Sprintf("%d", code)
}
