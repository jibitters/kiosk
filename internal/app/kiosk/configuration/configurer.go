// Copyright 2019 The Jibit Team. All rights reserved.
// Use of this source code is governed by a Apache-style license that can be found in the LICENSE.md file.

package configuration

import (
	"encoding/json"
	"io/ioutil"

	"github.com/jibitters/kiosk/internal/pkg/logging"
)

// Config encapsulates toplevel configuration properties.
type Config struct {
	Application ApplicationConfig `json:"application"`
	Logger      LoggerConfig      `json:"logger"`
}

// ApplicationConfig encapsulates default application properties.
type ApplicationConfig struct {
	Metrics     bool   `json:"metrics"`
	MetricsHost string `json:"metrics_host"`
	MetricsPort int    `json:"metrics_port"`
}

func (ac *ApplicationConfig) validate() {
	if ac.MetricsHost == "" {
		ac.MetricsHost = "localhost"
	}

	if ac.MetricsPort <= 0 {
		ac.MetricsPort = 9091
	}
}

// LoggerConfig encapsulates logger configuration properties.
type LoggerConfig struct {
	Level string `json:"level"`
}

func (lc *LoggerConfig) validate() {
	if lc.Level == "" {
		lc.Level = "info"
	}
}

// Configure reads a configuration file from provided file path and populates an instance of Config struct.
func Configure(logger *logging.Logger, filePath string) *Config {
	logger.Info("loading configurations file from %s", filePath)

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Fatal("failed to load configurations file: %v", err)
	}

	config := &Config{}
	if err := json.Unmarshal(content, config); err != nil {
		logger.Fatal("failed to parse configurations file: %v", err)
	}

	config.Application.validate()
	config.Logger.validate()

	return config
}
