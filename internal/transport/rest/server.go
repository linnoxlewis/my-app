package rest

import (
	"applicationDesignTest/internal/transport/rest/controllers"
	v1 "applicationDesignTest/internal/transport/rest/route/v1"
	"applicationDesignTest/pkg/log"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
)

type RestServer struct {
	server *http.Server
	logger log.Logging
}

func NewRestServer(address string,
	orderController *controllers.OrderController,
	logger log.Logging) *RestServer {
	router := chi.NewRouter()
	router.Get("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(map[string]any{"success": "ok"})

		return
	})

	v1.RegisterOrderRoutes(router, orderController)

	srv := &http.Server{
		Addr:    address,
		Handler: router,
	}

	return &RestServer{
		server: srv,
		logger: logger,
	}
}

func (r *RestServer) StartServer() {
	r.logger.LogInfo("Server is starting")
	err := r.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		r.logger.LogErrorf("err listen: %s\n", err)
		os.Exit(1)
	}
}

func (r *RestServer) StopServer(ctx context.Context) {
	if err := r.server.Shutdown(ctx); err != nil {
		r.logger.LogErrorf("err shutdown: %s\n", err)
		os.Exit(1)
	}
	r.logger.LogInfo("Server closed")
}
