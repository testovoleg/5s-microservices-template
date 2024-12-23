
В методах которые используют http при обращении к другим нашим микросервисам, перед вызовом http метода, в месте где устанавливаются header запроса, надо добавить следующее:
```go
tracingHeaders, _ := tracing.InjectTextMapCarrier(ctx)
for i, k := range tracingHeaders {
    req.Header.Set(i, k)
}
```