package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jibitters/kiosk/db/postgres"
	"github.com/lireza/lib/configuring"
	"go.uber.org/zap"
)

var config = flag.String("config", "./configs/kiosk.json", "configuration file")

// Kiosk is the main program encapsulation that holds all required components.
type Kiosk struct {
	logger *zap.SugaredLogger
	config *configuring.Config
	db     *pgxpool.Pool
}

func main() {
	kiosk := setup()

	kiosk.configure()
	kiosk.connectToDatabase()
	kiosk.migrateDatabase()
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

func (k *Kiosk) awaitTermination() {
	receiver := make(chan os.Signal)
	signal.Notify(receiver, os.Interrupt)

	<-receiver
	k.logger.Debug("Received interrupt signal!")

	k.stop()
}

func (k *Kiosk) stop() {
	k.logger.Info("Stopping the process ...")

	if k.db != nil {
		k.db.Close()
	}

	_ = k.logger.Sync()
}
