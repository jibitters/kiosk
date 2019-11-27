// Copyright 2019 The Jibit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE.md file.

package configuration

import (
	"io/ioutil"
	"testing"

	"github.com/jibitters/kiosk/internal/pkg/logging"
)

func TestValidateApplicationConfig(t *testing.T) {
	config := ApplicationConfig{MetricsHost: ""}
	config.validate()
	if config.MetricsHost != "localhost" {
		t.Logf("Actual value: %v Expected value: localhost", config.MetricsHost)
		t.FailNow()
	}

	config = ApplicationConfig{MetricsHost: "", MetricsPort: 0}
	config.validate()
	if config.MetricsPort != 9091 {
		t.Logf("Actual value: %v Expected value: 9091", config.MetricsPort)
		t.FailNow()
	}

	config = ApplicationConfig{MetricsHost: "", MetricsPort: -1}
	config.validate()
	if config.MetricsPort != 9091 {
		t.Logf("Actual value: %v Expected value: 9091", config.MetricsPort)
		t.FailNow()
	}
}

func TestValidateLoggerConfig(t *testing.T) {
	config := LoggerConfig{Level: ""}
	config.validate()
	if config.Level != "info" {
		t.Logf("Actual value: %v Expected value: info", config.Level)
		t.FailNow()
	}
}

func TestValidatePostgresConfig(t *testing.T) {
	config := PostgresConfig{Host: ""}
	config.validate()
	if config.Host != "localhost" {
		t.Logf("Actual value: %v Expected value: localhost", config.Host)
		t.FailNow()
	}

	config = PostgresConfig{Port: 0}
	config.validate()
	if config.Port != 5432 {
		t.Logf("Actual value: %v Expected value: 5432", config.Port)
		t.FailNow()
	}

	config = PostgresConfig{Port: -1}
	config.validate()
	if config.Port != 5432 {
		t.Logf("Actual value: %v Expected value: 5432", config.Port)
		t.FailNow()
	}

	config = PostgresConfig{Name: ""}
	config.validate()
	if config.Name != "kiosk" {
		t.Logf("Actual value: %v Expected value: kiosk", config.Name)
		t.FailNow()
	}

	config = PostgresConfig{ConnectionTimeout: 0}
	config.validate()
	if config.ConnectionTimeout != 10 {
		t.Logf("Actual value: %v Expected value: 10", config.ConnectionTimeout)
		t.FailNow()
	}

	config = PostgresConfig{ConnectionTimeout: -1}
	config.validate()
	if config.ConnectionTimeout != 10 {
		t.Logf("Actual value: %v Expected value: 10", config.ConnectionTimeout)
		t.FailNow()
	}

	config = PostgresConfig{MaxConnection: 0}
	config.validate()
	if config.MaxConnection != 8 {
		t.Logf("Actual value: %v Expected value: 8", config.MaxConnection)
		t.FailNow()
	}

	config = PostgresConfig{MaxConnection: -1}
	config.validate()
	if config.MaxConnection != 8 {
		t.Logf("Actual value: %v Expected value: 8", config.MaxConnection)
		t.FailNow()
	}

	config = PostgresConfig{SSLMode: "disable"}
	config.validate()
	if config.SSLMode != "disable" {
		t.Logf("Actual value: %v Expected value: disable", config.MaxConnection)
		t.FailNow()
	}
}

func TestValidateNatsConfig(t *testing.T) {
	config := NatsConfig{Addresses: nil}
	config.validate()
	if len(config.Addresses) != 1 {
		t.Logf("Actual value: %v Expected value: 1", len(config.Addresses))
		t.FailNow()
	}

	if config.Addresses[0] != "nats://localhost:4222" {
		t.Logf("Actual value: %v Expected value: nats://localhost:4222", config.Addresses[0])
		t.FailNow()
	}

	config = NatsConfig{Addresses: []string{}}
	config.validate()
	if len(config.Addresses) != 1 {
		t.Logf("Actual value: %v Expected value: 1", len(config.Addresses))
		t.FailNow()
	}

	if config.Addresses[0] != "nats://localhost:4222" {
		t.Logf("Actual value: %v Expected value: nats://localhost:4222", config.Addresses[0])
		t.FailNow()
	}
}

func TestValidateGRPCConfig(t *testing.T) {
	config := GRPCConfig{Host: ""}
	config.validate()
	if config.Host != "localhost" {
		t.Logf("Actual value: %v Expected value: localhost", config.Host)
		t.FailNow()
	}

	config = GRPCConfig{Port: 0}
	config.validate()
	if config.Port != 9090 {
		t.Logf("Actual value: %v Expected value: 9090", config.Port)
		t.FailNow()
	}

	config = GRPCConfig{Port: -1}
	config.validate()
	if config.Port != 9090 {
		t.Logf("Actual value: %v Expected value: 9090", config.Port)
		t.FailNow()
	}
}

func TestValidateWebConfig(t *testing.T) {
	config := WEBConfig{Host: ""}
	config.validate()
	if config.Host != "localhost" {
		t.Logf("Actual value: %v Expected value: localhost", config.Host)
		t.FailNow()
	}

	config = WEBConfig{Port: 0}
	config.validate()
	if config.Port != 8080 {
		t.Logf("Actual value: %v Expected value: 8080", config.Port)
		t.FailNow()
	}

	config = WEBConfig{Port: -1}
	config.validate()
	if config.Port != 8080 {
		t.Logf("Actual value: %v Expected value: 8080", config.Port)
		t.FailNow()
	}
}

func TestConfigure(t *testing.T) {
	file, err := ioutil.TempFile("", "kiosk*.json")
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer file.Close()

	file.WriteString(`
	{
		"application": {
			"metrics" : true,
			"metrics_host" : "localhost",
			"metrics_port" : 9091
		},
	
		"logger": {
			"level": "debug"
		},
	
		"postgres": {
			"host": "localhost",
			"port": 5432,
			"name": "kiosk",
			"user": "",
			"password": "",
			"connection_timeout": 10,
			"max_connection" : 8,
			"ssl_mode": "disable",
			"migration_directory": "file://migration/postgres"
		},
	
		"nats": {
			"addresses": ["nats://localhost:4222"]
		},
	
		"notifier": {
			"notify_by_sms_subject": "notifier.notifications.sms",
			"notify_by_call_subject": "notifier.notifications.call",
			"notify_by_email_subject": "notifier.notifications.email"
		},
	
		"notifications": {
			"ticket": {
				"new": {
					"low": {
						"type": "EMAIL",
						"recipients": ["support@example.com"]
					},
	
					"medium": {
						"type": "EMAIL",
						"recipients": ["support@example.com"]
					},
	
					"high": {
						"type": "EMAIL",
						"recipients": ["support@example.com"]
					},
	
					"critical": {
						"type": "SMS",
						"recipients": ["09120000000"]
					}
				}
			},
			
			"comment": {
				"new": {
					"type": "EMAIL",
					"recipients": ["support@example.com"],
					"ignore_owners": [""]
				}
			}
		},
	
		"grpc": {
			"host": "localhost",
			"port": 9090
		},
	
		"web": {
			"host": "localhost",
			"port": 8080
		}
	}
	`)

	config, err := Configure(logging.NewWithLevel("info"), file.Name())
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}

	if config.Application.Metrics != true {
		t.Logf("Actual value: %v Expected value: true", config.Application.Metrics)
		t.FailNow()
	}

	if config.Application.MetricsHost != "localhost" {
		t.Logf("Actual value: %v Expected value: localhost", config.Application.MetricsHost)
		t.FailNow()
	}

	if config.Application.MetricsPort != 9091 {
		t.Logf("Actual value: %v Expected value: 9091", config.Application.MetricsPort)
		t.FailNow()
	}

	if config.Logger.Level != "debug" {
		t.Logf("Actual value: %v Expected value: debug", config.Logger.Level)
		t.FailNow()
	}

	if config.Postgres.Host != "localhost" {
		t.Logf("Actual value: %v Expected value: localhost", config.Postgres.Host)
		t.FailNow()
	}

	if config.Postgres.Port != 5432 {
		t.Logf("Actual value: %v Expected value: 5432", config.Postgres.Port)
		t.FailNow()
	}

	if config.Postgres.Name != "kiosk" {
		t.Logf("Actual value: %v Expected value: kiosk", config.Postgres.Name)
		t.FailNow()
	}

	if config.Postgres.User != "" {
		t.Logf("Actual value: %v Expected value: ", config.Postgres.User)
		t.FailNow()
	}

	if config.Postgres.Password != "" {
		t.Logf("Actual value: %v Expected value: ", config.Postgres.Password)
		t.FailNow()
	}

	if config.Postgres.ConnectionTimeout != 10 {
		t.Logf("Actual value: %v Expected value: 10", config.Postgres.ConnectionTimeout)
		t.FailNow()
	}

	if config.Postgres.MaxConnection != 8 {
		t.Logf("Actual value: %v Expected value: 8", config.Postgres.MaxConnection)
		t.FailNow()
	}

	if config.Postgres.SSLMode != "disable" {
		t.Logf("Actual value: %v Expected value: disable", config.Postgres.SSLMode)
		t.FailNow()
	}

	if config.Postgres.MigrationDirectory != "file://migration/postgres" {
		t.Logf("Actual value: %v Expected value: file://migration/postgres", config.Postgres.MigrationDirectory)
		t.FailNow()
	}

	if config.Nats.Addresses[0] != "nats://localhost:4222" {
		t.Logf("Actual value: %v Expected value: nats://localhost:4222", config.Nats.Addresses[0])
		t.FailNow()
	}

	if config.Notifier.NotifyBySMSSubject != "notifier.notifications.sms" {
		t.Logf("Actual value: %v Expected value: notifier.notifications.sms", config.Notifier.NotifyBySMSSubject)
		t.FailNow()
	}

	if config.Notifier.NotifyByCallSubject != "notifier.notifications.call" {
		t.Logf("Actual value: %v Expected value: notifier.notifications.call", config.Notifier.NotifyByCallSubject)
		t.FailNow()
	}

	if config.Notifier.NotifyByEmailSubject != "notifier.notifications.email" {
		t.Logf("Actual value: %v Expected value: notifier.notifications.email", config.Notifier.NotifyByEmailSubject)
		t.FailNow()
	}

	if config.Notifications.Ticket.New.Low.Type != "EMAIL" {
		t.Logf("Actual value: %v Expected value: EMAIL", config.Notifications.Ticket.New.Low.Type)
		t.FailNow()
	}

	if config.Notifications.Ticket.New.Low.Recipients[0] != "support@example.com" {
		t.Logf("Actual value: %v Expected value: support@example.com", config.Notifications.Ticket.New.Low.Recipients[0])
		t.FailNow()
	}

	if config.Notifications.Ticket.New.Medium.Type != "EMAIL" {
		t.Logf("Actual value: %v Expected value: EMAIL", config.Notifications.Ticket.New.Medium.Type)
		t.FailNow()
	}

	if config.Notifications.Ticket.New.Medium.Recipients[0] != "support@example.com" {
		t.Logf("Actual value: %v Expected value: support@example.com", config.Notifications.Ticket.New.Medium.Recipients[0])
		t.FailNow()
	}

	if config.Notifications.Ticket.New.High.Type != "EMAIL" {
		t.Logf("Actual value: %v Expected value: EMAIL", config.Notifications.Ticket.New.High.Type)
		t.FailNow()
	}

	if config.Notifications.Ticket.New.High.Recipients[0] != "support@example.com" {
		t.Logf("Actual value: %v Expected value: support@example.com", config.Notifications.Ticket.New.High.Recipients[0])
		t.FailNow()
	}

	if config.Notifications.Ticket.New.Critical.Type != "SMS" {
		t.Logf("Actual value: %v Expected value: SMS", config.Notifications.Ticket.New.Critical.Type)
		t.FailNow()
	}

	if config.Notifications.Ticket.New.Critical.Recipients[0] != "09120000000" {
		t.Logf("Actual value: %v Expected value: 09120000000", config.Notifications.Ticket.New.Critical.Recipients[0])
		t.FailNow()
	}

	if config.Notifications.Comment.New.Type != "EMAIL" {
		t.Logf("Actual value: %v Expected value: EMAIL", config.Notifications.Comment.New.Type)
		t.FailNow()
	}

	if config.Notifications.Comment.New.Recipients[0] != "support@example.com" {
		t.Logf("Actual value: %v Expected value: support@example.com", config.Notifications.Comment.New.Recipients[0])
		t.FailNow()
	}

	if config.GRPC.Host != "localhost" {
		t.Logf("Actual value: %v Expected value: localhost", config.GRPC.Host)
		t.FailNow()
	}

	if config.GRPC.Port != 9090 {
		t.Logf("Actual value: %v Expected value: 9090", config.GRPC.Port)
		t.FailNow()
	}

	if config.WEB.Host != "localhost" {
		t.Logf("Actual value: %v Expected value: localhost", config.WEB.Host)
		t.FailNow()
	}

	if config.WEB.Port != 8080 {
		t.Logf("Actual value: %v Expected value: 8080", config.WEB.Port)
		t.FailNow()
	}
}

func TestConfigure_FileNotFound(t *testing.T) {
	_, err := Configure(logging.NewWithLevel("info"), "")
	if err == nil {
		t.Logf("Expected error here!")
		t.FailNow()
	}
}

func TestConfigure_InvalidJsonFormat(t *testing.T) {
	file, err := ioutil.TempFile("", "kiosk*.json")
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer file.Close()

	file.WriteString(`
		"application": {
			"metrics" : true,
			"metrics_host" : "localhost",
			"metrics_port" : 9091
		},
	
		"logger": {
			"level": "debug"
		},
	
		"postgres": {
			"host": "localhost",
			"port": 5432,
			"name": "kiosk",
			"user": "",
			"password": "",
			"connection_timeout": 10,
			"max_connection" : 8,
			"ssl_mode": "disable",
			"migration_directory": "file://migration/postgres"
		},
	
		"grpc": {
			"host": "localhost",
			"port": 9090
		},
	
		"web": {
			"host": "localhost",
			"port": 8080
		}
	}
	`)

	if _, err := Configure(logging.NewWithLevel("info"), file.Name()); err == nil {
		t.Logf("Expected error here!")
		t.FailNow()
	}
}
