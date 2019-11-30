// Copyright 2019 The Jibit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE.md file.

package web

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	rpc "github.com/jibitters/kiosk/g/rpc/kiosk"
	"github.com/jibitters/kiosk/internal/app/kiosk/services"
)

// TicketController is the implementation of ticket related controller methods.
type TicketController struct {
	service *services.TicketService
	routes  []Route
}

// NewTicketController creates and returns a new ready to use ticket controller implementation.
func NewTicketController(service *services.TicketService) *TicketController {
	return &TicketController{
		service: service,
		routes: []Route{
			Route{http.MethodPost, "/tickets", func(w http.ResponseWriter, request *http.Request) {
				w.Header().Add("Content-Type", "application/json; charset=utf-8")
				createTicket(service, w, request)
			}},

			Route{http.MethodGet, "/tickets/{id}", func(w http.ResponseWriter, request *http.Request) {
				w.Header().Add("Content-Type", "application/json; charset=utf-8")
				readTicket(service, w, request)
			}},

			Route{http.MethodPut, "/tickets", func(w http.ResponseWriter, request *http.Request) {
				w.Header().Add("Content-Type", "application/json; charset=utf-8")
				updateTicket(service, w, request)
			}},

			Route{http.MethodDelete, "/tickets/{id}", func(w http.ResponseWriter, request *http.Request) {
				w.Header().Add("Content-Type", "application/json; charset=utf-8")
				deleteTicket(service, w, request)
			}},

			Route{http.MethodGet, "/tickets", func(w http.ResponseWriter, request *http.Request) {
				w.Header().Add("Content-Type", "application/json; charset=utf-8")
				filterTickets(service, w, request)
			}},
		},
	}
}

// Routes returns the registered routes for current controller.
func (controller *TicketController) Routes() []Route {
	return controller.routes
}

func createTicket(service *services.TicketService, w http.ResponseWriter, request *http.Request) {
	ticket := &rpc.Ticket{}

	bytesIn, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := unmarshaler.Unmarshal(bytes.NewReader(bytesIn), ticket); err != nil {
		handleError(w, err)
		return
	}

	response, err := service.Create(request.Context(), ticket)
	if err != nil {
		handleError(w, err)
		return
	}

	responseBody := new(bytes.Buffer)
	marshaler.Marshal(responseBody, response)
	w.Write(responseBody.Bytes())
}

func readTicket(service *services.TicketService, w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		handleError(w, err)
		return
	}

	response, err := service.Read(request.Context(), &rpc.Id{Id: id})
	if err != nil {
		handleError(w, err)
		return
	}

	responseBody := new(bytes.Buffer)
	marshaler.Marshal(responseBody, response)
	w.Write(responseBody.Bytes())
}

func updateTicket(service *services.TicketService, w http.ResponseWriter, request *http.Request) {
	ticket := &rpc.Ticket{}

	bytesIn, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := unmarshaler.Unmarshal(bytes.NewReader(bytesIn), ticket); err != nil {
		handleError(w, err)
		return
	}

	response, err := service.Update(request.Context(), ticket)
	if err != nil {
		handleError(w, err)
		return
	}

	responseBody := new(bytes.Buffer)
	marshaler.Marshal(responseBody, response)
	w.Write(responseBody.Bytes())
}

func deleteTicket(service *services.TicketService, w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		handleError(w, err)
		return
	}

	response, err := service.Delete(request.Context(), &rpc.Id{Id: id})
	if err != nil {
		handleError(w, err)
		return
	}

	responseBody := new(bytes.Buffer)
	marshaler.Marshal(responseBody, response)
	w.Write(responseBody.Bytes())
}

func filterTickets(service *services.TicketService, w http.ResponseWriter, request *http.Request) {
	issuer := request.FormValue("issuer")
	owner := request.FormValue("owner")
	ticketImportanceLevel := request.FormValue("ticket_importance_level")
	ticketStatus := request.FormValue("ticket_status")
	fromDate := request.FormValue("from_date")
	toDate := request.FormValue("to_data")
	pageNumber := request.FormValue("page_number")
	pageSize := request.FormValue("page_size")

	pn, err := strconv.ParseInt(pageNumber, 10, 64)
	if err != nil {
		handleError(w, err)
		return
	}

	ps, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil {
		handleError(w, err)
		return
	}

	response, err := service.Filter(request.Context(), &rpc.FilterTicketsRequest{
		Issuer:                issuer,
		Owner:                 owner,
		TicketImportanceLevel: rpc.TicketImportanceLevel(rpc.TicketImportanceLevel_value[ticketImportanceLevel]),
		TicketStatus:          rpc.TicketStatus(rpc.TicketImportanceLevel_value[ticketStatus]),
		FromDate:              fromDate,
		ToDate:                toDate,
		Page:                  &rpc.Page{Number: int32(pn), Size: int32(ps)},
	})
	if err != nil {
		handleError(w, err)
		return
	}

	responseBody := new(bytes.Buffer)
	marshaler.Marshal(responseBody, response)
	w.Write(responseBody.Bytes())
}
