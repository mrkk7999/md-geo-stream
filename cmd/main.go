package main

import (
	"fmt"
	"md-geo-track/controller"
	"md-geo-track/implementation"
	"md-geo-track/kafka"
	"md-geo-track/repository"
	httpTransport "md-geo-track/transport/http"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Initialize Logger
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	err := godotenv.Load("../.env")
	if err != nil {
		log.WithError(err).Fatal("Error loading .env file")
	}

	var (
		httpAddr = os.Getenv("HTTP_ADDR")
		// DB ENV
		dbHost     = os.Getenv("DB_HOST")
		dbUser     = os.Getenv("DB_USER")
		dbPassword = os.Getenv("DB_PASSWORD")
		dbName     = os.Getenv("DB_NAME")
		dbPort     = os.Getenv("DB_PORT")
		dbSSLMode  = os.Getenv("DB_SSLMODE")
		dbTimeZone = os.Getenv("DB_TIMEZONE")
		// KAFKA ENV
		brokers = []string{os.Getenv("BROKERS")}
		topic   = os.Getenv("TOPIC")
		groupID = os.Getenv("GROUP_ID")
	)

	// PostgreSQL DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode, dbTimeZone)

	log.Info("Connecting to database...")

	// Connect to Database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to database")
	}
	log.Info("Successfully connected to the database.")

	repo := repository.New(db)

	svc := implementation.New(repo, log)

	controller := controller.New(svc)

	handler := httpTransport.SetUpRouter(controller)

	kafka.StartConsumer(brokers, topic, groupID, svc, log)

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	log.WithField("address", httpAddr).Info("Starting HTTP server...")

	go func() {
		server := &http.Server{
			Addr:    httpAddr,
			Handler: handler,
		}
		errs <- server.ListenAndServe()
	}()

	log.Error("exit", <-errs)
}
