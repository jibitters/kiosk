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
