// Copyright 2019 The Jibit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE.md file.

package web

import (
	"bytes"
	"io/ioutil"
	"net/http"

	rpc "github.com/jibitters/kiosk/g/rpc/kiosk"
	"github.com/jibitters/kiosk/internal/app/kiosk/services"
)

// EchoController is the implementation of echo related controller methods.
type EchoController struct {
	service *services.EchoService
	routes  []Route
}

// NewEchoController creates and returns a new ready to use echo controller implementation.
func NewEchoController(service *services.EchoService) *EchoController {
	return &EchoController{
		service: service,
		routes: []Route{
			Route{http.MethodPost, "/echo", func(w http.ResponseWriter, request *http.Request) {
				w.Header().Add("Content-Type", "application/json; charset=utf-8")
				echo(service, w, request)
			}},
		},
	}
}

// Routes returns the registered routes for current controller.
func (controller *EchoController) Routes() []Route {
	return controller.routes
}

func echo(service *services.EchoService, w http.ResponseWriter, request *http.Request) {
	message := &rpc.Message{}

	bytesIn, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := unmarshaler.Unmarshal(bytes.NewReader(bytesIn), message); err != nil {
		handleError(w, err)
		return
	}

	response, err := service.Echo(request.Context(), message)
	if err != nil {
		handleError(w, err)
		return
	}

	responseBody := new(bytes.Buffer)
	marshaler.Marshal(responseBody, response)
	w.Write(responseBody.Bytes())
}
