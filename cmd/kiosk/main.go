// Copyright 2019 The JIBit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE file.

//go:generate protoc -I ../../api/protobuf-spec --go_out=plugins=grpc:../../ models.proto
//go:generate protoc -I ../../api/protobuf-spec --go_out=plugins=grpc:../../ echo_service.proto
//go:generate protoc -I ../../api/protobuf-spec --go_out=plugins=grpc:../../ ticket_service.proto
//go:generate protoc -I ../../api/protobuf-spec --go_out=plugins=grpc:../../ comment_service.proto
//go:generate protoc -I ../../api/protobuf-spec/notifier --go_out=plugins=grpc:../../ models.proto
package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jibitters/kiosk/internal/app/kiosk/configuration"
	"github.com/jibitters/kiosk/internal/app/kiosk/database"
	"github.com/jibitters/kiosk/internal/app/kiosk/nats"
	"github.com/jibitters/kiosk/internal/app/kiosk/server"
	"github.com/jibitters/kiosk/internal/app/kiosk/web"
	"github.com/jibitters/kiosk/internal/pkg/logging"
	_ "github.com/lib/pq"
	natsclient "github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

// Command line options to parse.
var (
	config = flag.String("config", "./configs/kiosk.json", "JSON configuration file path.")
	pguser = flag.String("pguser", "", "Postgres database user. If provided rewrites the value of postgres.user in configuration file.")
	pgpass = flag.String("pgpass", "", "Postgres database password. If provided rewrites the value of postgres.password in configuration file.")
)

// The kiosk application definition.
type kiosk struct {
	config *configuration.Config
	logger *logging.Logger
	db     *pgxpool.Pool
	nats   *natsclient.Conn
	grpc   *grpc.Server
	web    *http.Server
}

func main() {
	flag.Parse()

	kiosk := &kiosk{config: &configuration.Config{}, logger: logging.New(logging.InfoLevel)}
	kiosk.configure()
	kiosk.migrate()
	kiosk.connectToDatabase()
	kiosk.connectToNats()
	kiosk.listen()
	kiosk.listenWeb()
	kiosk.addInterruptHook()
}

// Configures kiosk application instance based on provided configuration properties.
func (k *kiosk) configure() {
	config, err := configuration.Configure(k.logger, *config)
	if err != nil {
		k.logger.Fatal("failed to load configurations file: %v", err)
	}

	k.config = config
	k.logger = logging.NewWithLevel(k.config.Logger.Level)
}

// Tries to connect to a postgres instance and then runs the database migration.
func (k *kiosk) migrate() {
	if *pguser != "" {
		k.config.Postgres.User = *pguser
	}

	if *pgpass != "" {
		k.config.Postgres.Password = *pgpass
	}

	if err := database.Migrate(k.config); err != nil {
		k.stop()
		k.logger.Fatal("failed to run database migration: %v", err)
	}

	k.logger.Debug("successfully executed database migration")
}

// Tries to setup a connection to postgres instance.
func (k *kiosk) connectToDatabase() {
	if *pguser != "" {
		k.config.Postgres.User = *pguser
	}

	if *pgpass != "" {
		k.config.Postgres.Password = *pgpass
	}

	db, err := database.ConnectToDatabase(k.config)
	if err != nil {
		k.stop()
		k.logger.Fatal("failed to connect to postgres instance: %v", err)
	}

	k.db = db
}

// Tries to setup a connection to nats cluster.
func (k *kiosk) connectToNats() {
	nats, err := nats.ConnectToNats(k.config)
	if err != nil {
		k.stop()
		k.logger.Fatal("failed to connect to nats cluster: %v", err)
	}

	k.nats = nats
}

// Listens on provided host and port to provide a series of gRPC services.
func (k *kiosk) listen() {
	server, err := server.Listen(k.config, k.logger, k.db, k.nats)
	if err != nil {
		k.stop()
		k.logger.Fatal("failed to start gRPC server: %v", err)
	}

	k.grpc = server
	k.logger.Info("successfully started gRPC server and listening on %s:%d", k.config.GRPC.Host, k.config.GRPC.Port)
}

// Listens on provided host and port to provide a series of REST apis.
func (k *kiosk) listenWeb() {
	k.web = web.ListenWeb(k.config, k.logger, k.db, k.nats)

	k.logger.Info("successfully started web server and listening on %s:%d", k.config.WEB.Host, k.config.WEB.Port)
}

// Adds interrupt hook for application to be called on os terminate signal.
func (k *kiosk) addInterruptHook() {
	signalReceiver := make(chan os.Signal, 1)
	signal.Notify(signalReceiver, os.Interrupt)

	<-signalReceiver
	k.stop()
}

// Gracefully stops all components.
func (k *kiosk) stop() {
	// First we should stop gRPC to deny incoming calls.
	if k.web != nil {
		k.logger.Debug("stopping web server ...")
		if err := k.web.Close(); err != nil {
			k.logger.Error("could not stop web server successfully")
		}
	}

	if k.grpc != nil {
		k.logger.Debug("stopping gRPC server ...")
		k.grpc.GracefulStop()
	}

	if k.nats != nil {
		k.logger.Debug("stopping nats client ...")
		k.nats.Close()
	}

	if k.db != nil {
		k.logger.Debug("closing database connection ...")
		k.db.Close()
	}
}
