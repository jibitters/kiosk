// Copyright 2019 The Jibit Team. All rights reserved.
// Use of this source code is governed by a Apache-style license that can be found in the LICENSE.md file.

//go:generate protoc -I ../../api/protobuf-spec --go_out=plugins=grpc:../../ models.proto
//go:generate protoc -I ../../api/protobuf-spec --go_out=plugins=grpc:../../ echo_service.proto
//go:generate protoc -I ../../api/protobuf-spec --go_out=plugins=grpc:../../ ticket_service.proto
//go:generate protoc -I ../../api/protobuf-spec --go_out=plugins=grpc:../../ message_service.proto
package main

import (
	"flag"

	"github.com/jibitters/kiosk/internal/app/kiosk/configuration"
	"github.com/jibitters/kiosk/internal/pkg/logging"
)

// Command line options to parse.
var (
	config = flag.String("config", "./configs/kiosk.json", "JSON configuration file path.")
)

// The kiosk application definition.
type kiosk struct {
	config *configuration.Config
	logger *logging.Logger
}

func main() {
	flag.Parse()

	kiosk := &kiosk{config: &configuration.Config{}, logger: logging.New(logging.InfoLevel)}
	kiosk.configure()
}

// Configures kiosk application instance based on provided configuration properties.
func (k *kiosk) configure() {
	k.config = configuration.Configure(k.logger, *config)
	k.logger = logging.NewWithLevel(k.config.Logger.Level)
}
