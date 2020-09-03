package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/jibitters/kiosk/errors"
	nc "github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// CommentHandler is the handler implementation of comments related resource.
type CommentHandler struct {
	logger     *zap.SugaredLogger
	natsClient *nc.Conn
}

// NewCommentHandler returns back a newly created and ready to use CommentHandler.
func NewCommentHandler(logger *zap.SugaredLogger, natsClient *nc.Conn) *CommentHandler {
	return &CommentHandler{logger: logger, natsClient: natsClient}
}

// Create creates a new comment with specified information.
func (h *CommentHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in, _ := ioutil.ReadAll(r.Body)

		response, e := h.natsClient.RequestWithContext(r.Context(), "kiosk.comments.create", in)
		if e != nil {
			if e == nc.ErrTimeout {
				et := errors.RequestTimeout("")
				writeError(w, et)
			} else {
				et := errors.InternalServerError("unknown", "")
				h.logger.Error(et.FingerPrint, ": ", e.Error())
				writeError(w, et)
			}

			return
		}

		et := &errors.Type{}
		_ = json.Unmarshal(response.Data, et)
		if et.FingerPrint != "" {
			writeError(w, et)
			return
		}

		writeNoContent(w)
	}
}
