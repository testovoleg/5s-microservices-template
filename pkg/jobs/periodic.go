package jobs

import (
	"context"
	"math/rand"
	"runtime/debug"
	"time"

	"github.com/testovoleg/5s-microservice-template/pkg/logger"
)

// RunPeriodic — общий цикл фоновых джоб core_service: тик, пауза interval + случайный
// джиттер до jitter (чтобы поды не били в базу синхронно), остановка по ctx. Каждый тик
// обёрнут в recover: паника в логике одной джобы пишется в лог со стеком и не роняет
// весь процесс. Если startImmediately — первый тик сразу при старте, иначе после паузы.
//
// Использование (server.Run, после подключения репозиториев):
//
//	go jobs.RunPeriodic(ctx, s.log, "cleanupJob", 10*time.Minute, time.Minute, true, func(ctx context.Context) {
//		if err := s.svc.Commands.CleanupMethods.Cleanup(ctx); err != nil {
//			s.log.WarnMsg("cleanupJob", err)
//		}
//	})
func RunPeriodic(ctx context.Context, log logger.Logger, name string, interval, jitter time.Duration, startImmediately bool, tick func(ctx context.Context)) {
	safeTick := func() {
		defer func() {
			if rec := recover(); rec != nil {
				log.Errorf("%s: recovered from panic: %v\n%s", name, rec, debug.Stack())
			}
		}()
		tick(ctx)
	}
	if startImmediately {
		select {
		case <-ctx.Done():
			return
		default:
		}
		safeTick()
	}
	for {
		sleep := interval
		if jitter > 0 {
			sleep += time.Duration(rand.Int63n(int64(jitter)))
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(sleep):
		}
		safeTick()
	}
}
