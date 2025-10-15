Чтобы добавить в header "x-trace-id" нужно в файле "pkg\tracing\utils.go" в методе StartHttpServerTracerSpan рядом с строчкой `c.Set("traceid", traceID.String())` добавить следующую строчку:
```
c.Response().Header().Set("x-trace-id", traceID.String())
```