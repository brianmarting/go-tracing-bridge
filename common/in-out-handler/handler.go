package in_out_handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"go-tracing-bridge/common/observability/tracing"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
	"os"
	"time"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)

	CreateAllRoutes()
}

type handler struct {
	*chi.Mux

	client http.Client
}

func NewHandler() Handler {
	return &handler{
		Mux: chi.NewMux(),
		client: http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

func (h *handler) CreateAllRoutes() {
	h.Post("/in", post(h.client))
}

func post(client http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		spanCtx, span := tracing.GetTracer().Start(r.Context(), "sending request")
		defer span.End()

		ctx, cancel := context.WithTimeout(spanCtx, time.Second*10)
		defer cancel()

		endpoint := fmt.Sprintf("http://%s:%s/in", os.Getenv("NEXT_SERVICE"), os.Getenv("NEXT_SERVICE_PORT"))
		log.Info().Msg("Forwarding to " + endpoint)
		request, err := http.NewRequestWithContext(ctx, "POST", endpoint, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response, err := client.Do(request)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				http.Error(w, "timed out when validating token", http.StatusBadGateway)
				return
			}

			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		if response.StatusCode != 200 {
			http.Error(w, "err", http.StatusBadGateway)
			return
		}

		log.Info().Msg("Forwarded to " + endpoint)

		w.WriteHeader(200)
	}
}
