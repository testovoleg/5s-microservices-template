# 10. Bulk-батчи и переполнение очереди трассировки (потерянные parent-спаны)

На большом батч-цикле (тысячи items → десятки тысяч спанов в одном трейсе) Tempo
может показать детей без родителя. Хелпер и лимиты SDK лежат в `pkg/tracing`, чтобы
сервис из шаблона не повторил эту картину.

## Симптом

В Tempo/Grafana часть поддеревьев «висела» без родителя: сотни/тысячи спанов
репозитория (`redisRepository.PutEvent`, `dbRepository.*`) были в трейсе, а их
родитель (command-хендлер / Kafka `ProcessMessages`) — нет.

## Причина

- Каждый repo/redis/http-вызов автоматически инструментирован через `tracing.StartSpan`.
  На один item в цикле уходит несколько таких вызовов.
- При батче в тысячи items это десятки тысяч спанов в одном трейсе.
- `BatchSpanProcessor` (`pkg/tracing/opentelemetry.go`) с дефолтными лимитами SDK
  (`MaxQueueSize=2048`, `MaxExportBatchSize=512`) переполняется, и SDK **молча дропает**
  новый спан, если очередь полна (штатное поведение, не баг SDK).
- Спан-родитель цикла вызывает `span.End()` (→ попадает в очередь на экспорт) **после всех
  своих детей** — поэтому именно родительские спаны batch-цикла систематически чаще
  попадают под дроп. Результат — «дети есть, родителя нет».

## Фикс в шаблоне

1. `pkg/tracing/opentelemetry.go` — увеличены лимиты `BatchSpanProcessor`
   (`trace.WithMaxQueueSize(20000)`, `trace.WithMaxExportBatchSize(2000)`) как страховка
   от переполнения на «обычных» путях.
2. `pkg/tracing/utils.go` — `tracing.WithBulkTracing(ctx)`: помечает `ctx` так, что все
   `tracing.StartSpan` под ним (в т.ч. через repository-слой) становятся no-op (через
   `otel/trace/noop`), вместо создания реальных спанов. Это устраняет проблему в источнике,
   а не просто расширяет очередь.

В demo-коде шаблона (CRUD `/admin/api`, webhook) хелпер **не вызывается** — там нет
цикла по тысячам items. Подключать его нужно в конкретном сервисе, когда появляется
bulk-обработка.

## Как использовать в сервисе из шаблона

```go
ctx, span := tracing.StartSpan(ctx, "handler.ProcessBatch")
defer span.End()

if total > 200 { // порог — по объёму своего цикла
    utils.Attr(span, "bulk_tracing", "true")
    ctx = tracing.WithBulkTracing(ctx)
}
// дальше цикл по items: repo-вызовы больше не плодят спаны
```

На внешнем span пишите агрегаты (`total`, `accepted`, `skipped`) через `utils.Attr` —
детальная видимость по каждому item теряется, но трейс не взрывается и родитель не дропается.

## Что смотреть при переносе

- `pkg/tracing/utils.go` — `WithBulkTracing`/`isBulkTracing`/`noopTracer`, обёртка в `StartSpan`.
- `pkg/tracing/opentelemetry.go` — лимиты `BatchSpanProcessor`.
- Любой bulk-цикл (обработка батча из N item с repo-вызовами внутри) — тот же паттерн.
