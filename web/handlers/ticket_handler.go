package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/jibitters/kiosk/errors"
	"github.com/jibitters/kiosk/models"
	"github.com/jibitters/kiosk/web/data"
	nc "github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// TicketHandler is the handler implementation of tickets related resource.
type TicketHandler struct {
	logger     *zap.SugaredLogger
	natsClient *nc.Conn
}

// NewTicketHandler returns back a newly created and ready to use TicketHandler.
func NewTicketHandler(logger *zap.SugaredLogger, natsClient *nc.Conn) *TicketHandler {
	return &TicketHandler{logger: logger, natsClient: natsClient}
}

// Create creates a new ticket with specified information.
func (h *TicketHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in, _ := ioutil.ReadAll(r.Body)

		response, e := h.natsClient.RequestWithContext(r.Context(), "kiosk.tickets.create", in)
		if e != nil {
			if e == nc.ErrTimeout {
				et := errors.RequestTimeout("")
				writeError(w, et)
			}

			et := errors.InternalServerError("unknown", "")
			h.logger.Error(et.FingerPrint, ": ", e.Error())
			writeError(w, et)
			return
		}

		if string(response.Data) != "" {
			et := &errors.Type{}
			_ = json.Unmarshal(response.Data, et)
			writeError(w, et)
			return
		}

		writeNoContent(w)
	}
}

// Filter filters tickets based on provided criteria values.
func (h *TicketHandler) Filter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		issuer := r.URL.Query().Get("issuer")
		owner := r.URL.Query().Get("owner")
		importanceLevel := r.URL.Query().Get("importanceLevel")
		status := r.URL.Query().Get("status")
		fromDate := r.URL.Query().Get("fromDate")
		toDate := r.URL.Query().Get("toDate")
		pageNumber, _ := strconv.Atoi(r.URL.Query().Get("pageNumber"))
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

		filterTicketsRequest := data.FilterTicketsRequest{Issuer: issuer, Owner: owner,
			ImportanceLevel: models.TicketImportanceLevel(importanceLevel), Status: models.TicketStatus(status),
			FromDate: fromDate, ToDate: toDate, PageNumber: pageNumber, PageSize: pageSize}

		in, _ := json.Marshal(filterTicketsRequest)
		response, e := h.natsClient.RequestWithContext(r.Context(), "kiosk.tickets.filter", in)
		if e != nil {
			if e == nc.ErrTimeout {
				et := errors.RequestTimeout("")
				writeError(w, et)
			}

			et := errors.InternalServerError("unknown", "")
			h.logger.Error(et.FingerPrint, ": ", e.Error())
			writeError(w, et)
			return
		}

		et := &errors.Type{}
		_ = json.Unmarshal(response.Data, et)
		if et.FingerPrint != "" {
			writeError(w, et)
			return
		}

		filterTicketsResponse := &data.FilterTicketsResponse{}
		_ = json.Unmarshal(response.Data, filterTicketsResponse)
		write(w, filterTicketsResponse)
	}
}
