// Copyright 2019 The Jibit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE.md file.

package services

import (
	"context"

	rpc "github.com/jibitters/kiosk/g/rpc/kiosk"
)

// EchoService is the implementation of echo rpc methods.
type EchoService struct {
}

// NewEchoService creates and returns a new ready to use echo service implementation.
func NewEchoService() *EchoService {
	return &EchoService{}
}

// Echo is an echo rpc.
func (service *EchoService) Echo(context context.Context, request *rpc.Message) (*rpc.Message, error) {
	return &rpc.Message{Content: request.Content}, nil
}
