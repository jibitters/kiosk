// Copyright 2019 The Jibit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE.md file.

package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/fasthttp/router"
	"github.com/golang/protobuf/jsonpb"
	"github.com/jackc/pgx/v4/pgxpool"
	rpc "github.com/jibitters/kiosk/g/rpc/kiosk"
	"github.com/jibitters/kiosk/internal/app/kiosk/configuration"
	"github.com/jibitters/kiosk/internal/app/kiosk/services"
	"github.com/jibitters/kiosk/internal/pkg/logging"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const version = "/v1"

type handler struct {
	echoService   *services.EchoService
	ticketService *services.TicketService
	marshaler     *jsonpb.Marshaler
	unmarshaler   *jsonpb.Unmarshaler
}

// ListenWeb creates a new HTTP server and listens on provided host and port.
func ListenWeb(config *configuration.Config, logger *logging.Logger, db *pgxpool.Pool) error {
	handler := setup(config, logger, db)

	router := router.New()
	router.POST(version+"/echo", handler.echo)

	router.POST(version+"/tickets", handler.createTicket)
	router.GET(version+"/tickets/:id", handler.readTicket)

	go fasthttp.ListenAndServe(fmt.Sprintf("%s:%d", config.WEB.Host, config.WEB.Port), router.Handler)
	return nil
}

func setup(config *configuration.Config, logger *logging.Logger, db *pgxpool.Pool) *handler {
	return &handler{
		echoService:   services.NewEchoService(),
		ticketService: services.NewTicketService(logger, db),
		marshaler:     &jsonpb.Marshaler{OrigName: true, EmitDefaults: true},
		unmarshaler:   &jsonpb.Unmarshaler{},
	}
}

func (h *handler) echo(context *fasthttp.RequestCtx) {
	message := &rpc.Message{}

	if err := h.unmarshaler.Unmarshal(bytes.NewReader(context.Request.Body()), message); err != nil {
		handleError(err, context)
		return
	}

	response, err := h.echoService.Echo(context, message)
	if err != nil {
		handleError(err, context)
		return
	}

	responseBody := new(bytes.Buffer)
	h.marshaler.Marshal(responseBody, response)
	context.Response.Header.Add("Content-Type", "application/json; application/json; charset=utf-8")
	context.Write(responseBody.Bytes())
}

func (h *handler) createTicket(context *fasthttp.RequestCtx) {
	ticket := &rpc.Ticket{}

	if err := h.unmarshaler.Unmarshal(bytes.NewReader(context.Request.Body()), ticket); err != nil {
		handleError(err, context)
		return
	}

	response, err := h.ticketService.Create(context, ticket)
	if err != nil {
		handleError(err, context)
		return
	}

	responseBody := new(bytes.Buffer)
	h.marshaler.Marshal(responseBody, response)
	context.Response.Header.Add("Content-Type", "application/json; application/json; charset=utf-8")
	context.Write(responseBody.Bytes())
}

func (h *handler) readTicket(context *fasthttp.RequestCtx) {
	path := context.UserValue("id").(string)
	id, err := strconv.Atoi(path)
	if err != nil {
		handleError(err, context)
		return
	}

	response, err := h.ticketService.Read(context, &rpc.Id{Id: int64(id)})
	if err != nil {
		handleError(err, context)
		return
	}

	responseBody := new(bytes.Buffer)
	h.marshaler.Marshal(responseBody, response)
	context.Response.Header.Add("Content-Type", "application/json; application/json; charset=utf-8")
	context.Write(responseBody.Bytes())
}

func handleError(err error, context *fasthttp.RequestCtx) {
	context.Response.Header.Add("Content-Type", "application/json; charset=utf-8")

	status, ok := status.FromError(err)
	switch status.Code() {
	case codes.Unknown:
		context.Response.Header.SetStatusCode(http.StatusInternalServerError)
	case codes.InvalidArgument:
		context.Response.Header.SetStatusCode(http.StatusBadRequest)
	case codes.DeadlineExceeded:
		context.Response.Header.SetStatusCode(http.StatusRequestTimeout)
	case codes.NotFound:
		context.Response.Header.SetStatusCode(http.StatusNotFound)
	case codes.AlreadyExists:
		context.Response.Header.SetStatusCode(http.StatusConflict)
	case codes.PermissionDenied:
		context.Response.Header.SetStatusCode(http.StatusForbidden)
	case codes.FailedPrecondition:
		context.Response.Header.SetStatusCode(http.StatusPreconditionFailed)
	case codes.Unimplemented:
		context.Response.Header.SetStatusCode(http.StatusNotImplemented)
	case codes.Internal:
		context.Response.Header.SetStatusCode(http.StatusInternalServerError)
	case codes.Unavailable:
		context.Response.Header.SetStatusCode(http.StatusUnavailableForLegalReasons)
	case codes.Unauthenticated:
		context.Response.Header.SetStatusCode(http.StatusUnauthorized)
	default:
		context.Response.Header.SetStatusCode(http.StatusInternalServerError)
	}

	es := &Errors{}
	if ok {
		es.Errors = append(es.Errors, Error{Code: status.Message()})
	} else {
		es.Errors = append(es.Errors, Error{Message: status.Message()})
	}

	responseBody, _ := json.Marshal(es)
	context.Write(responseBody)
}
