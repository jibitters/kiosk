// Copyright 2019 The JIBit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE file.

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics encapsulates all required counters, histograms and gauges to record some metrics.
type Metrics struct {
	HandledCounter   *prometheus.CounterVec
	HandledHistogram *prometheus.HistogramVec
}

// NewMetrics initializes all encapsulates metrics.
func NewMetrics() *Metrics {
	counter := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "kiosk_requests_count",
		Help: "Number of requests handled so far",
	}, []string{"type", "service", "method", "code", "message", "status"})

	histogram := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "kiosk_requests_seconds",
		Help: "The time took us to handle each request",
	}, []string{"type", "service", "method", "code", "message", "status"})

	return &Metrics{
		HandledCounter:   counter,
		HandledHistogram: histogram,
	}
}
