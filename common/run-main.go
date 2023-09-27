package common

import (
	"context"
	"go-tracing-bridge/common/observability/tracing"
	"os"
	"os/signal"

	"github.com/rs/zerolog/log"
)

func RunMain(serviceName string, fn func()) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	tracer := tracing.InitTracerProvider(serviceName)
	defer func() {
		cancel()
		if err := tracer.Shutdown(context.Background()); err != nil {
			log.Info().Err(err).Msg("failed to shut down tracer provider")
		}
	}()

	go func() {
		log.Info().Msg("starting" + serviceName)
		fn()
		log.Info().Msg("finished" + serviceName)
	}()

	<-ctx.Done()
	log.Info().Msg("stopped" + serviceName)
}
