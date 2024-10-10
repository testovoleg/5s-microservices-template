
В Dockerfile необходимо найти строчку где загружается golang
```go
FROM --platform=linux/amd64 golang:1.21 AS builder
```
И заменить на последнюю версию
```go
FROM --platform=linux/amd64 golang:1.22 AS builder
```

В некоторых случаях, необходимо так же обновить go.mod файл. Сделав команду go mod tidy или поправив руками версию.
