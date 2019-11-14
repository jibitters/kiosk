// Copyright 2019 The Jibit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE.md file.

package server

import (
	"fmt"
	"net"
	"net/http"

	rpc "github.com/jibitters/kiosk/g/rpc/kiosk"
	"github.com/jibitters/kiosk/internal/app/kiosk/configuration"
	"github.com/jibitters/kiosk/internal/app/kiosk/metrics"
	"github.com/jibitters/kiosk/internal/app/kiosk/services"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

// Listen creates a new gRPC server and listens on provided host and port.
func Listen(config *configuration.Config) (*grpc.Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.GRPC.Host, config.GRPC.Port))
	if err != nil {
		return nil, err
	}

	server := server(config)
	rpc.RegisterEchoServiceServer(server, services.NewEchoService())

	if config.Application.Metrics {
		go enableMetricsEndpoint(config.Application.MetricsHost, config.Application.MetricsPort)
	}

	go server.Serve(listener)
	return server, nil
}

// Creates an instance of gRPC server.
func server(config *configuration.Config) *grpc.Server {
	if config.Application.Metrics {
		return grpc.NewServer(grpc.UnaryInterceptor(metrics.UnaryInterceptor(metrics.NewMetrics())))
	}

	return grpc.NewServer()
}

// Enables metrics endpoint for prometheus server to scrape.
func enableMetricsEndpoint(host string, port int) {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
}
