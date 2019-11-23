// Copyright 2019 The Jibit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE.md file.

package services

import (
	"context"
	"testing"

	rpc "github.com/jibitters/kiosk/g/rpc/kiosk"
)

func TestEcho(t *testing.T) {
	service := NewEchoService()
	request := rpc.Message{Content: "echo"}

	response, err := service.Echo(context.Background(), &request)
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}

	if response.Content != request.Content {
		t.Logf("Actual : %v Expected: %v", response.Content, request.Content)
		t.FailNow()
	}
}
