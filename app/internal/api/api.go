package api

import (
	"context"
	"net"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"go.uber.org/zap"

	"github.com/emorikvendy/url-shortener/internal/controllers"
	"github.com/emorikvendy/url-shortener/internal/resources"
)

// Api represents a api server.
type Api struct {
	server http.Server
	errors chan error
	logger *zap.SugaredLogger
}

// New returns a new instance of the Api server.
func New(logger *zap.SugaredLogger, port int, r resources.R) *Api {
	router := mux.NewRouter()
	controllers.AddURLRoutes(router, logger, r.Adapters.URL, r.Config.HashLen)
	return &Api{
		server: http.Server{
			Addr:    net.JoinHostPort("", strconv.Itoa(port)),
			Handler: router,
		},
		errors: make(chan error, 1),
		logger: logger,
	}
}

// Start api server.
func (d *Api) Start() {
	go func() {
		d.errors <- d.server.ListenAndServe()
		close(d.errors)
	}()
}

// Stop api server.
func (d *Api) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return d.server.Shutdown(ctx)
}

// Notify returns a channel to notify the caller about errors.
// If you receive an error from the channel diagnostic you should stop the application.
func (d *Api) Notify() <-chan error {
	return d.errors
}
