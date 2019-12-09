// Copyright 2019 The JIBit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE file.

package nats

import (
	"strings"

	"github.com/jibitters/kiosk/internal/app/kiosk/configuration"
	natsclient "github.com/nats-io/nats.go"
)

// ConnectToNats tries to make a connection to nats cluster.
func ConnectToNats(config *configuration.Config) (*natsclient.Conn, error) {
	return natsclient.Connect(strings.Join(config.Nats.Addresses, ","))
}
