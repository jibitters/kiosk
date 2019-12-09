// Copyright 2019 The JIBit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE file.

package metrics

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryInterceptor intercepts each incoming request, record some metrics and exposes them to the prometheus endpoint.
func UnaryInterceptor(metrics *Metrics) func(context context.Context, request interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(context context.Context, request interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		response, err := handler(context, request)
		end := time.Now()
		recordMetrics(metrics, start, end, info, err)

		return response, err
	}
}

func recordMetrics(metrics *Metrics, start time.Time, end time.Time, info *grpc.UnaryServerInfo, err error) {
	service, method := extractServiceMethod(info.FullMethod)
	code, message := extractCodeMessage(err)
	responseStatus := extractStatus(err)

	metrics.HandledCounter.WithLabelValues("Unary", service, method, code, message, responseStatus).Inc()
	metrics.HandledHistogram.WithLabelValues("Unary", service, method, code, message, responseStatus).Observe(end.Sub(start).Seconds())
}

func extractServiceMethod(fullMethodName string) (string, string) {
	fullMethodName = strings.TrimPrefix(fullMethodName, "/")
	if i := strings.Index(fullMethodName, "/"); i >= 0 {
		return fullMethodName[:i], fullMethodName[i+1:]
	}

	return "", ""
}

func extractCodeMessage(err error) (string, string) {
	if err == nil {
		return codes.OK.String(), ""
	}

	grpcError, ok := status.FromError(err)
	if ok {
		return grpcError.Code().String(), grpcError.Message()
	}

	return "", ""
}

func extractStatus(err error) string {
	status := "Successful"
	if err != nil {
		status = "Failed"
	}

	return status
}
