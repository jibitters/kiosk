package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jibitters/kiosk/web/handlers"
	"github.com/lireza/lib/configuring"
	nc "github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

const (
	v1       = "/v1"
	echo     = "/echo"
	tickets  = "/tickets"
	comments = "/comments"
	metrics  = "/metrics"
)

// StartServer setups and then runs an HTTP server.
func StartServer(logger *zap.SugaredLogger, config *configuring.Config, natsClient *nc.Conn) *http.Server {
	host := config.Get("web.server.host").StringOrElse("localhost")
	port := config.Get("web.server.port").UintOrElse(8080)
	readTimeout := config.Get("web.server.read_timeout").DurationOrElse(10 * time.Second)
	readHeaderTimeout := config.Get("web.server.read_header_timeout").DurationOrElse(5 * time.Second)
	writeTimeout := config.Get("web.server.write_timeout").DurationOrElse(10 * time.Second)
	idleTimeout := config.Get("web.server.idle_timeout").DurationOrElse(30 * time.Second)

	logger.Info("web.server.host -> ", host)
	logger.Info("web.server.port -> ", port)
	logger.Info("web.server.read_timeout -> ", readTimeout)
	logger.Info("web.server.read_header_timeout -> ", readHeaderTimeout)
	logger.Info("web.server.write_timeout -> ", writeTimeout)
	logger.Info("web.server.idle_timeout -> ", idleTimeout)

	router := setupRoutes(logger, natsClient)

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

func setupRoutes(logger *zap.SugaredLogger, natsClient *nc.Conn) *mux.Router {
	// Router
	router := mux.NewRouter().
		PathPrefix(v1).
		Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete).
		Subrouter()

	// Meddlers
	meddlers := handlers.NewMeddlers()
	router.Use(meddlers.JSONContentTypeHeaderMiddleware)

	// Echo handler
	echoHandler := handlers.NewEchoHandler(logger)
	router.Methods(http.MethodPost).PathPrefix(echo).HandlerFunc(echoHandler.Echo())

	// Ticket handler
	ticketHandler := handlers.NewTicketHandler(logger, natsClient)
	router.Methods(http.MethodPost).PathPrefix(tickets).HandlerFunc(ticketHandler.Create())
	router.Methods(http.MethodGet).PathPrefix(tickets).HandlerFunc(ticketHandler.Filter())

	// Comment handler
	commentHandler := handlers.NewCommentHandler(logger, natsClient)
	router.Methods(http.MethodPost).PathPrefix(comments).HandlerFunc(commentHandler.Create())

	// Metrics handler
	router.Handle(metrics, promhttp.Handler())

	return router
}
