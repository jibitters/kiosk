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
	natsclient "github.com/nats-io/nats.go"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const version = "/v1"

type handler struct {
	echoService    *services.EchoService
	ticketService  *services.TicketService
	commentService *services.CommentService
	marshaler      *jsonpb.Marshaler
	unmarshaler    *jsonpb.Unmarshaler
}

// ListenWeb creates a new HTTP server and listens on provided host and port.
func ListenWeb(config *configuration.Config, logger *logging.Logger, db *pgxpool.Pool, nats *natsclient.Conn) error {
	handler := setup(config, logger, db, nats)

	router := router.New()
	router.POST(version+"/echo", handler.echo)

	router.POST(version+"/tickets", handler.createTicket)
	router.GET(version+"/tickets/:id", handler.readTicket)
	router.PUT(version+"/tickets", handler.updateTicket)
	router.DELETE(version+"/tickets/:id", handler.deleteTicket)
	router.GET(version+"/tickets", handler.filterTickets)

	router.POST(version+"/comments", handler.createComment)
	router.PUT(version+"/comments", handler.updateComment)
	router.DELETE(version+"/comments/:id", handler.deleteComment)

	go fasthttp.ListenAndServe(fmt.Sprintf("%s:%d", config.WEB.Host, config.WEB.Port), router.Handler)
	return nil
}

func setup(config *configuration.Config, logger *logging.Logger, db *pgxpool.Pool, nats *natsclient.Conn) *handler {
	return &handler{
		echoService:    services.NewEchoService(),
		ticketService:  services.NewTicketService(config, logger, db, nats),
		commentService: services.NewCommentService(config, logger, db, nats),
		marshaler:      &jsonpb.Marshaler{OrigName: true, EmitDefaults: true},
		unmarshaler:    &jsonpb.Unmarshaler{},
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
	pathSegment := context.UserValue("id").(string)
	id, err := strconv.ParseInt(pathSegment, 10, 64)
	if err != nil {
		handleError(err, context)
		return
	}

	response, err := h.ticketService.Read(context, &rpc.Id{Id: id})
	if err != nil {
		handleError(err, context)
		return
	}

	responseBody := new(bytes.Buffer)
	h.marshaler.Marshal(responseBody, response)
	context.Response.Header.Add("Content-Type", "application/json; application/json; charset=utf-8")
	context.Write(responseBody.Bytes())
}

func (h *handler) updateTicket(context *fasthttp.RequestCtx) {
	ticket := &rpc.Ticket{}

	if err := h.unmarshaler.Unmarshal(bytes.NewReader(context.Request.Body()), ticket); err != nil {
		handleError(err, context)
		return
	}

	response, err := h.ticketService.Update(context, ticket)
	if err != nil {
		handleError(err, context)
		return
	}

	responseBody := new(bytes.Buffer)
	h.marshaler.Marshal(responseBody, response)
	context.Response.Header.Add("Content-Type", "application/json; application/json; charset=utf-8")
	context.Write(responseBody.Bytes())
}

func (h *handler) deleteTicket(context *fasthttp.RequestCtx) {
	pathSegment := context.UserValue("id").(string)
	id, err := strconv.ParseInt(pathSegment, 10, 64)
	if err != nil {
		handleError(err, context)
		return
	}

	response, err := h.ticketService.Delete(context, &rpc.Id{Id: id})
	if err != nil {
		handleError(err, context)
		return
	}

	responseBody := new(bytes.Buffer)
	h.marshaler.Marshal(responseBody, response)
	context.Response.Header.Add("Content-Type", "application/json; application/json; charset=utf-8")
	context.Write(responseBody.Bytes())
}

func (h *handler) filterTickets(context *fasthttp.RequestCtx) {
	issuer := string(context.QueryArgs().Peek("issuer"))
	owner := string(context.QueryArgs().Peek("owner"))
	ticketImportanceLevel := string(context.QueryArgs().Peek("ticket_importance_level"))
	ticketStatus := string(context.QueryArgs().Peek("ticket_status"))
	fromDate := string(context.QueryArgs().Peek("from_date"))
	toDate := string(context.QueryArgs().Peek("to_data"))
	pageNumber, _ := context.QueryArgs().GetUint("page_number")
	pageSize, _ := context.QueryArgs().GetUint("page_size")

	response, err := h.ticketService.Filter(context, &rpc.FilterTicketsRequest{
		Issuer:                issuer,
		Owner:                 owner,
		TicketImportanceLevel: rpc.TicketImportanceLevel(rpc.TicketImportanceLevel_value[ticketImportanceLevel]),
		TicketStatus:          rpc.TicketStatus(rpc.TicketImportanceLevel_value[ticketStatus]),
		FromDate:              fromDate,
		ToDate:                toDate,
		Page:                  &rpc.Page{Number: int32(pageNumber), Size: int32(pageSize)},
	})
	if err != nil {
		handleError(err, context)
		return
	}

	responseBody := new(bytes.Buffer)
	h.marshaler.Marshal(responseBody, response)
	context.Response.Header.Add("Content-Type", "application/json; application/json; charset=utf-8")
	context.Write(responseBody.Bytes())
}

func (h *handler) createComment(context *fasthttp.RequestCtx) {
	comment := &rpc.Comment{}

	if err := h.unmarshaler.Unmarshal(bytes.NewReader(context.Request.Body()), comment); err != nil {
		handleError(err, context)
		return
	}

	response, err := h.commentService.Create(context, comment)
	if err != nil {
		handleError(err, context)
		return
	}

	responseBody := new(bytes.Buffer)
	h.marshaler.Marshal(responseBody, response)
	context.Response.Header.Add("Content-Type", "application/json; application/json; charset=utf-8")
	context.Write(responseBody.Bytes())
}

func (h *handler) updateComment(context *fasthttp.RequestCtx) {
	comment := &rpc.Comment{}

	if err := h.unmarshaler.Unmarshal(bytes.NewReader(context.Request.Body()), comment); err != nil {
		handleError(err, context)
		return
	}

	response, err := h.commentService.Update(context, comment)
	if err != nil {
		handleError(err, context)
		return
	}

	responseBody := new(bytes.Buffer)
	h.marshaler.Marshal(responseBody, response)
	context.Response.Header.Add("Content-Type", "application/json; application/json; charset=utf-8")
	context.Write(responseBody.Bytes())
}

func (h *handler) deleteComment(context *fasthttp.RequestCtx) {
	pathSegment := context.UserValue("id").(string)
	id, err := strconv.ParseInt(pathSegment, 10, 64)
	if err != nil {
		handleError(err, context)
		return
	}

	response, err := h.commentService.Delete(context, &rpc.Id{Id: id})
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
