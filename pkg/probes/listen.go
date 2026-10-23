package probes

import (
	"context"
	"net/http"
	"time"

	"github.com/heptiolabs/healthcheck"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
)

// ListenReadiness starts a Kubernetes readiness HTTP listener on port.
func ListenReadiness(ctx context.Context, log logger.Logger, serviceName, port string, intervalSec int, ready func() error) {
	health := healthcheck.NewHandler()
	health.AddReadinessCheck(serviceName, healthcheck.AsyncWithContext(ctx, ready, time.Duration(intervalSec)*time.Second))
	go func() {
		log.Infof("%s Kubernetes probes listening on port: %s", serviceName, port)
		if err := http.ListenAndServe(port, health); err != nil {
			log.WarnMsg("ListenAndServe", err)
		}
	}()
}
