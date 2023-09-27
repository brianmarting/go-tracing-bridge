package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	restc "go-tracing-bridge/c-service/rest"
	"go-tracing-bridge/common"
	"net/http"
	"os"
)

func main() {
	common.RunMain("Service C", func() {
		h := restc.NewHandler()
		h.CreateAllRoutes()
		if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("SERVICE_PORT")), h); err != nil {
			log.Fatal().Err(err).Msg("failed to start c-service ")
		}
	})
}
