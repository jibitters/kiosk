// Copyright 2019 The Jibit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE.md file.

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
	Postgres    PostgresConfig    `json:"postgres"`
	GRPC        GRPCConfig        `json:"grpc"`
	WEB         WEBConfig         `json:"web"`
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

// PostgresConfig encapsulates postgres configuration properties.
type PostgresConfig struct {
	Host               string `json:"host"`
	Port               int    `json:"port"`
	Name               string `json:"name"`
	User               string `json:"user"`
	Password           string `json:"password"`
	ConnectionTimeout  int    `json:"connection_timeout"`
	MaxConnection      int    `json:"max_connection"`
	SSLMode            string `json:"ssl_mode"`
	MigrationDirectory string `json:"migration_directory"`
}

func (pc *PostgresConfig) validate() {
	if pc.Host == "" {
		pc.Host = "localhost"
	}

	if pc.Port <= 0 {
		pc.Port = 5432
	}

	if pc.Name == "" {
		pc.Name = "kiosk"
	}

	if pc.ConnectionTimeout <= 0 {
		pc.ConnectionTimeout = 10
	}

	if pc.MaxConnection <= 0 {
		pc.MaxConnection = 8
	}

	if pc.SSLMode == "" {
		pc.SSLMode = "disable"
	}
}

// GRPCConfig encapsulates gRPC server configuration properties.
type GRPCConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (gc *GRPCConfig) validate() {
	if gc.Host == "" {
		gc.Host = "localhost"
	}

	if gc.Port <= 0 {
		gc.Port = 9090
	}
}

// WEBConfig encapsulates web server configuration properties.
type WEBConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (wc *WEBConfig) validate() {
	if wc.Host == "" {
		wc.Host = "localhost"
	}

	if wc.Port <= 0 {
		wc.Port = 8080
	}
}

// Configure reads a configuration file from provided file path and populates an instance of Config struct.
func Configure(logger *logging.Logger, filePath string) (*Config, error) {
	logger.Info("loading configurations file from %s", filePath)

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.Unmarshal(content, config); err != nil {
		return nil, err
	}

	config.Application.validate()
	config.Logger.validate()
	config.Postgres.validate()
	config.GRPC.validate()
	config.WEB.validate()

	return config, nil
}
