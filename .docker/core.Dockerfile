FROM --platform=linux/amd64 golang:1.22 AS builder

LABEL stage=gobuilder
ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /build/core core_service/cmd/main.go 

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /build/core /app/core
COPY core_service/config/config.yaml /app/config.yml

EXPOSE 5001

ENTRYPOINT [ "/app/core", "-config=/app/config.yml"]