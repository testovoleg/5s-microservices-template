FROM --platform=linux/amd64 golang:1.22 AS builder

LABEL stage=gobuilder
ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /build/api_gateway api_gateway_service/cmd/main.go 

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates && apk add mailcap && rm /var/cache/apk/*

WORKDIR /app
COPY --from=builder /build/api_gateway /app/api_gateway
COPY api_gateway_service/config/config.yaml /app/config.yml

EXPOSE 5001

ENTRYPOINT [ "/app/api_gateway", "-config=/app/config.yml"]