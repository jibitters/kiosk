// Copyright 2019 The JIBit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	rpc "github.com/jibitters/kiosk/g/rpc/kiosk"
	"github.com/jibitters/kiosk/internal/app/kiosk/configuration"
	"github.com/jibitters/kiosk/internal/app/kiosk/metrics"
	"github.com/jibitters/kiosk/internal/app/kiosk/services"
	"github.com/jibitters/kiosk/internal/pkg/logging"
	natsclient "github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

// Listen creates a new gRPC server and listens on provided host and port.
func Listen(logger *logging.Logger, config *configuration.Config, db *pgxpool.Pool, nats *natsclient.Conn) (*grpc.Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.GRPC.Host, config.GRPC.Port))
	if err != nil {
		return nil, err
	}

	server := server(config)
	rpc.RegisterEchoServiceServer(server, services.NewEchoService())
	rpc.RegisterTicketServiceServer(server, services.NewTicketService(logger, config, db, nats))
	rpc.RegisterCommentServiceServer(server, services.NewCommentService(logger, config, db, nats))

	if config.Application.Metrics {
		go enableMetricsEndpoint(config.Application.MetricsHost, config.Application.MetricsPort)
	}

	go server.Serve(listener)
	return server, nil
}

func server(config *configuration.Config) *grpc.Server {
	if config.Application.Metrics {
		return grpc.NewServer(grpc.UnaryInterceptor(metrics.UnaryInterceptor(metrics.NewMetrics())))
	}

	return grpc.NewServer()
}

func enableMetricsEndpoint(host string, port int) error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
}
