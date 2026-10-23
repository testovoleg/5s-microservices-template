package server

import (
	"context"
	"errors"

	"github.com/testovoleg/5s-microservice-template/pkg/probes"
)

func (s *server) runHealthCheck(ctx context.Context) {
	probes.ListenReadiness(ctx, s.log, s.cfg.ServiceName, s.cfg.Probes.Port, s.cfg.Probes.CheckIntervalSeconds, func() error {
		if s.cfg != nil {
			return nil
		}
		return errors.New("Config not loaded")
	})
}
