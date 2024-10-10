
# Добавление traceID в контекст
В pkg/tracing/utils.go в функции StartHttpServerTracerSpan добавить следующее перед возвратом функции ( возможно у вас уже это сделано ):
```go
	traceID := serverSpan.SpanContext().TraceID()
	if traceID.IsValid() {
		c.Set("traceid", traceID.String())
	}
```

# Добавление константы
В pkg/constants/constants.go добавить константу
```go
TRACEID = "TRACEID"
```

# Добавление переменной в функцию логирования
В pkg/logger/logger.go добавить функции HttpMiddlewareAccessLogger переменную traceID \
В итоге у меня получилось так
```go 
func (l *appLogger) HttpMiddlewareAccessLogger(method, uri string, status int, size int64, time time.Duration, traceID string) {
	l.logger.Info(
		constants.HTTP,
		zap.String(constants.METHOD, method),
		zap.String(constants.URI, uri),
		zap.Int(constants.STATUS, status),
		zap.Int64(constants.SIZE, size),
		zap.Duration(constants.TIME, time),
        zap.String(constants.TRACEID, traceID),
	)
}
```
И так же добавить в интерфейс Logger у этого же метода переменную traceID. Получилось так:
```go
HttpMiddlewareAccessLogger(method string, uri string, status int, size int64, time time.Duration, traceID string)
```

# Добавляем передачу traceid при логировании
В методах где используется логирование http, а это api_gateway ( api_gateway_service/internal/middlewares/middlewares.go ) и graphql ( graphql_service/internal/middlewares/middlewares.go ) находим строку где используется HttpMiddlewareAccessLogger и добавляем туда traceid ( и заодно её там создаем )
```go
traceID, _ := ctx.Get("traceid").(string)

if !mw.checkIgnoredURI(ctx.Request().RequestURI, mw.cfg.Http.IgnoreLogUrls) {
    mw.log.HttpMiddlewareAccessLogger(req.Method, req.URL.String(), status, size, s, traceID)
}
```
