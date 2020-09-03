package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/jibitters/kiosk/errors"
	"go.uber.org/zap"
)

func parse(logger *zap.SugaredLogger, w http.ResponseWriter, r *http.Request, t interface{}) (ok bool) {
	in, e := ioutil.ReadAll(r.Body)
	if e != nil {
		et := errors.InternalServerError("unknown", "")
		logger.Error(et.FingerPrint, ": ", e.Error())

		writeError(w, et)
		return false
	}

	if e = json.Unmarshal(in, t); e != nil {
		et := errors.InvalidRequestBody()
		logger.Warn(et.FingerPrint, ": Could not parse json: ", string(in))

		writeError(w, et)
		return false
	}

	return true
}

func write(w http.ResponseWriter, t interface{}) {
	out, _ := json.Marshal(t)
	_, _ = w.Write(out)
}

func writeNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func writeError(w http.ResponseWriter, e *errors.Type) {
	out, _ := json.Marshal(e)
	w.WriteHeader(e.HTTPStatusCode)
	_, _ = w.Write(out)
}
