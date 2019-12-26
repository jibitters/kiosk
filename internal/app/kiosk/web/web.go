// Copyright 2019 The JIBit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE file.

package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jibitters/kiosk/internal/app/kiosk/configuration"
	"github.com/jibitters/kiosk/internal/app/kiosk/services"
	"github.com/jibitters/kiosk/internal/pkg/logging"
	natsclient "github.com/nats-io/nats.go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// API version number.
const version = "/v1"

// Global json marshaler and unmarshaler.
var (
	marshaler   = &jsonpb.Marshaler{OrigName: true, EmitDefaults: true}
	unmarshaler = &jsonpb.Unmarshaler{}
)

// Route is the abstraction over an HTTP route and its handler.
type Route struct {
	Method  string
	Path    string
	Handler func(w http.ResponseWriter, request *http.Request)
}

// ListenWeb creates a new HTTP server and listens on provided host and port.
func ListenWeb(logger *logging.Logger, config *configuration.Config, db *pgxpool.Pool, nats *natsclient.Conn) *http.Server {
	echoController := NewEchoController(services.NewEchoService())
	ticketController := NewTicketController(services.NewTicketService(config, logger, db, nats))
	commentController := NewCommentController(services.NewCommentService(config, logger, db, nats))

	router := mux.NewRouter()
	router = router.PathPrefix(version).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete).Subrouter()
	for _, route := range echoController.Routes() {
		router.Methods(route.Method).Path(route.Path).HandlerFunc(route.Handler)
	}

	for _, route := range ticketController.Routes() {
		router.Methods(route.Method).Path(route.Path).HandlerFunc(route.Handler)
	}

	for _, route := range commentController.Routes() {
		router.Methods(route.Method).Path(route.Path).HandlerFunc(route.Handler)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.WEB.Host, config.WEB.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go server.ListenAndServe()
	return server
}

func handleError(w http.ResponseWriter, err error) {
	status, ok := status.FromError(err)
	switch status.Code() {
	case codes.Unknown:
		w.WriteHeader(http.StatusInternalServerError)
	case codes.InvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	case codes.DeadlineExceeded:
		w.WriteHeader(http.StatusRequestTimeout)
	case codes.NotFound:
		w.WriteHeader(http.StatusNotFound)
	case codes.AlreadyExists:
		w.WriteHeader(http.StatusConflict)
	case codes.PermissionDenied:
		w.WriteHeader(http.StatusForbidden)
	case codes.FailedPrecondition:
		w.WriteHeader(http.StatusPreconditionFailed)
	case codes.Unimplemented:
		w.WriteHeader(http.StatusNotImplemented)
	case codes.Internal:
		w.WriteHeader(http.StatusInternalServerError)
	case codes.Unavailable:
		w.WriteHeader(http.StatusUnavailableForLegalReasons)
	case codes.Unauthenticated:
		w.WriteHeader(http.StatusUnauthorized)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	es := &Errors{}
	if ok {
		es.Errors = append(es.Errors, Error{Code: status.Message()})
	} else {
		es.Errors = append(es.Errors, Error{Message: status.Message()})
	}

	responseBody, _ := json.Marshal(es)
	_, _ = w.Write(responseBody)
}
