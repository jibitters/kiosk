package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/lireza/lib/configuring"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

const (
	metrics = "/metrics"
)

// StartServer setups and then runs an HTTP server.
func StartServer(logger *zap.SugaredLogger, config *configuring.Config) *http.Server {

	host := config.Get("web.server.host").StringOrElse("localhost")
	port := config.Get("web.server.port").UintOrElse(8080)
	readTimeout := config.Get("web.server.read_timeout").DurationOrElse(5 * time.Second)
	readHeaderTimeout := config.Get("web.server.read_header_timeout").DurationOrElse(2 * time.Second)
	writeTimeout := config.Get("web.server.write_timeout").DurationOrElse(10 * time.Second)
	idleTimeout := config.Get("web.server.idle_timeout").DurationOrElse(30 * time.Second)

	logger.Debug("web.server.host -> ", host)
	logger.Debug("web.server.port -> ", port)
	logger.Debug("web.server.read_timeout -> ", readTimeout)
	logger.Debug("web.server.read_header_timeout -> ", readHeaderTimeout)
	logger.Debug("web.server.write_timeout -> ", writeTimeout)
	logger.Debug("web.server.idle_timeout -> ", idleTimeout)

	router := setupRoutes()

	server := &http.Server{
		Addr:              fmt.Sprintf("%v:%v", host, port),
		Handler:           router,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	go func() { _ = server.ListenAndServe() }()

	logger.Info("Web server started successfully and listening on ", host, ":", port)
	return server
}

func setupRoutes() *mux.Router {
	// Router
	router := mux.NewRouter().
		Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete).
		Subrouter()

	// Metrics handler
	router.Handle(metrics, promhttp.Handler())

	return router
}
