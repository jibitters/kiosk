// Copyright 2019 The JIBit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE file.

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

// CommentController is the implementation of comment related controller methods.
type CommentController struct {
	service *services.CommentService
	routes  []Route
}

// NewCommentController creates and returns a new ready to use comment controller implementation.
func NewCommentController(service *services.CommentService) *CommentController {
	return &CommentController{
		service: service,
		routes: []Route{
			Route{http.MethodPost, "/comments", func(w http.ResponseWriter, request *http.Request) {
				w.Header().Add("Content-Type", "application/json; charset=utf-8")
				createComment(service, w, request)
			}},

			Route{http.MethodPut, "/comments", func(w http.ResponseWriter, request *http.Request) {
				w.Header().Add("Content-Type", "application/json; charset=utf-8")
				updateComment(service, w, request)
			}},

			Route{http.MethodDelete, "/comments/{id}", func(w http.ResponseWriter, request *http.Request) {
				w.Header().Add("Content-Type", "application/json; charset=utf-8")
				deleteComment(service, w, request)
			}},
		},
	}
}

// Routes returns the registered routes for current controller.
func (controller *CommentController) Routes() []Route {
	return controller.routes
}

func createComment(service *services.CommentService, w http.ResponseWriter, request *http.Request) {
	comment := &rpc.Comment{}

	bytesIn, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := unmarshaler.Unmarshal(bytes.NewReader(bytesIn), comment); err != nil {
		handleError(w, err)
		return
	}

	response, err := service.Create(request.Context(), comment)
	if err != nil {
		handleError(w, err)
		return
	}

	responseBody := new(bytes.Buffer)
	marshaler.Marshal(responseBody, response)
	w.Write(responseBody.Bytes())
}

func updateComment(service *services.CommentService, w http.ResponseWriter, request *http.Request) {
	comment := &rpc.Comment{}

	bytesIn, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := unmarshaler.Unmarshal(bytes.NewReader(bytesIn), comment); err != nil {
		handleError(w, err)
		return
	}

	response, err := service.Update(request.Context(), comment)
	if err != nil {
		handleError(w, err)
		return
	}

	responseBody := new(bytes.Buffer)
	marshaler.Marshal(responseBody, response)
	w.Write(responseBody.Bytes())
}

func deleteComment(service *services.CommentService, w http.ResponseWriter, request *http.Request) {
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
