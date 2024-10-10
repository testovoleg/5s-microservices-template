В репозиториях ( а желательно во всех ), которые могут вызываться желательно сделать проверку того, начат ли span \
Начало трассировки без проверки выглядит вот так:
```go
span, ctx := opentracing.StartSpanFromContext(ctx, "someRepository.SomeFunc")
defer span.Finish()
```
С проверкой начата ли трассирвока вот так:
```go
if opentracing.SpanFromContext(ctx) != nil { // if have tracing, start new span
    var span opentracing.Span
    span, ctx = opentracing.StartSpanFromContext(ctx, "someRepository.SomeFunc")
    defer span.Finish()
}
```

Замену подобного желательно делать в каждом нужном файликах, а не глобально везде. \
Строки vscode для замены подобных строчек ( должен быть включен поиск по regexp ) \
Строка для поиска
```go
ctx, span := tracing.StartSpan\(ctx, "(.+)"\)
	defer span.End\(\)
```
Строка для замены
```go
if trace.SpanContextFromContext(ctx).IsValid() { // if have tracing, start new span
    var span trace.Span
    ctx, span = tracing.StartSpan(ctx, "$1")
    defer span.End()
}
```