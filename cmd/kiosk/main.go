package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jibitters/kiosk/db/postgres"
	"github.com/lireza/lib/configuring"
	nc "github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

var config = flag.String("config", "./configs/kiosk.json", "configuration file")

// Kiosk is the main program encapsulation that holds all required components.
type Kiosk struct {
	logger     *zap.SugaredLogger
	config     *configuring.Config
	db         *pgxpool.Pool
	natsClient *nc.Conn
}

func main() {
	kiosk := setup()

	kiosk.configure()
	kiosk.connectToDatabase()
	kiosk.migrateDatabase()
	kiosk.prepareNatsClient()
	kiosk.awaitTermination()
}

func setup() *Kiosk {
	flag.Parse()
	logger, _ := zap.NewDevelopment()
	config := configuring.New()

	return &Kiosk{logger: logger.Sugar(), config: config}
}

func (k *Kiosk) configure() {
	k.logger.Info("Loading configuration file from ", *config)
	if _, e := k.config.LoadJSON(*config); e != nil {
		k.logger.Fatal(e.Error())
	}

	environment := k.config.Get("logger.environment").StringOrElse("DEVELOPMENT")
	k.logger.Debug("logger.environment -> ", environment)

	if environment == "PRODUCTION" {
		logger, _ := zap.NewProduction()
		k.logger = logger.Sugar()
	}
}

func (k *Kiosk) connectToDatabase() {
	db, e := postgres.Connect(k.logger, k.config)
	if e != nil {
		k.stop()
		k.logger.Fatal(e.Error())
	}

	k.db = db
}

func (k *Kiosk) migrateDatabase() {
	if e := postgres.Migrate(k.logger, k.config); e != nil {
		k.stop()
		k.logger.Fatal(e.Error())
	}
}

func (k *Kiosk) prepareNatsClient() {
	addresses := k.config.Get("nats.addresses").SliceOfStringOrElse([]string{"nats://localhost:4222"})
	k.logger.Debug("nats.addresses -> ", addresses)

	client, e := nc.Connect(strings.Join(addresses, ","), nc.Name("Kiosk"))
	if e != nil {
		k.stop()
		k.logger.Fatal(e.Error())
	}

	k.natsClient = client
}

func (k *Kiosk) awaitTermination() {
	receiver := make(chan os.Signal)
	signal.Notify(receiver, os.Interrupt)

	<-receiver
	k.logger.Debug("Received interrupt signal!")

	k.stop()
}

func (k *Kiosk) stop() {
	k.logger.Info("Stopping the process ...")

	if k.natsClient != nil {
		k.natsClient.Close()
	}

	if k.db != nil {
		k.db.Close()
	}

	_ = k.logger.Sync()
}
