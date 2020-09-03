package handlers

import (
	"net/http"

	"github.com/jibitters/kiosk/web/data"
	"go.uber.org/zap"
)

// EchoHandler is the handler implementation for health checking.
type EchoHandler struct {
	logger *zap.SugaredLogger
}

// NewEchoHandler returns back a newly created and ready to use EchoHandler.
func NewEchoHandler(logger *zap.SugaredLogger) *EchoHandler {
	return &EchoHandler{logger: logger}
}

// Echo returns back the same message that receives.
func (h *EchoHandler) Echo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		echoRequest := data.EchoRequest{}
		if ok := parse(h.logger, w, r, &echoRequest); !ok {
			return
		}

		write(w, echoRequest)
	}
}
