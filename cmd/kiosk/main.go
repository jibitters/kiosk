// Copyright 2019 The Jibit Team. All rights reserved.
// Use of this source code is governed by a Apache-style license that can be found in the LICENSE.md file.

//go:generate protoc -I ../../api/protobuf-spec --go_out=plugins=grpc:../../ models.proto
//go:generate protoc -I ../../api/protobuf-spec --go_out=plugins=grpc:../../ echo_service.proto
//go:generate protoc -I ../../api/protobuf-spec --go_out=plugins=grpc:../../ ticket_service.proto
//go:generate protoc -I ../../api/protobuf-spec --go_out=plugins=grpc:../../ message_service.proto
package main

import (
	"flag"
)

// Command line options to parse.
var (
	config = flag.String("config", "./configs/orbital.json", "JSON configuration file path.")
)

func main() {
	flag.Parse()
}
