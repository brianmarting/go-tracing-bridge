package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"go-tracing-bridge/common"
	in_out_handler "go-tracing-bridge/common/in-out-handler"
	"net/http"
	"os"
)

func main() {
	common.RunMain("Service B", func() {
		h := in_out_handler.NewHandler()
		h.CreateAllRoutes()
		if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("SERVICE_PORT")), h); err != nil {
			log.Fatal().Err(err).Msg("failed to start b-service ")
		}
	})
}
