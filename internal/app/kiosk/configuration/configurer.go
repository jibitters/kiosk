// Copyright 2019 The JIBit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE file.

package configuration

import (
	"encoding/json"
	"io/ioutil"

	"github.com/jibitters/kiosk/internal/pkg/logging"
)

// Config encapsulates toplevel configuration properties.
type Config struct {
	Application   ApplicationConfig   `json:"application"`
	Logger        LoggerConfig        `json:"logger"`
	Postgres      PostgresConfig      `json:"postgres"`
	Nats          NatsConfig          `json:"nats"`
	Notifier      NotifierConfig      `json:"notifier"`
	Notifications NotificationsConfig `json:"notifications"`
	GRPC          GRPCConfig          `json:"grpc"`
	WEB           WEBConfig           `json:"web"`
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

// NatsConfig encapsulates the nats cluster configuration properties.
type NatsConfig struct {
	Addresses []string `json:"addresses"`
}

func (nc *NatsConfig) validate() {
	if nc.Addresses == nil || len(nc.Addresses) == 0 {
		nc.Addresses = append(nc.Addresses, "nats://localhost:4222")
	}
}

// NotifierConfig encapsulates the notifier project configuration properties.
type NotifierConfig struct {
	Subject string `json:"subject"`
}

// NotificationsConfig encapsulates the notifications configuration properties required on different action types.
type NotificationsConfig struct {
	Ticket  TicketNotificationConfig  `json:"ticket"`
	Comment CommentNotificationConfig `json:"comment"`
}

// TicketNotificationConfig encapsulates the notifications configuration properties related to ticket services.
type TicketNotificationConfig struct {
	New NewTicketNotificationConfig `json:"new"`
}

// NewTicketNotificationConfig encapsulates the notifications configuration properties related to new ticket created.
type NewTicketNotificationConfig struct {
	Low      NotificationConfig `json:"low"`
	Medium   NotificationConfig `json:"medium"`
	High     NotificationConfig `json:"high"`
	Critical NotificationConfig `json:"critical"`
}

// CommentNotificationConfig encapsulates the notifications configuration properties related to comment services.
type CommentNotificationConfig struct {
	New NotificationConfig `json:"new"`
}

// NotificationConfig encapsulates a single notification configuration properties.
type NotificationConfig struct {
	Type       string   `json:"type"`
	Recipients []string `json:"recipients"`
	Sender     string   `json:"sender"`
	CC         []string `json:"cc"`
	BCC        []string `json:"bcc"`
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
	config.Nats.validate()
	config.GRPC.validate()
	config.WEB.validate()

	return config, nil
}
