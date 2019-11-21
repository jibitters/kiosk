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
