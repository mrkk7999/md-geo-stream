package main

import (
	"fmt"
	"log"
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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
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

	// Connect to Database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	repo := repository.New(db)

	svc := implementation.New(repo)

	controller := controller.New(svc)

	handler := httpTransport.SetUpRouter(controller)

	kafka.StartConsumer(brokers, topic, groupID, svc)

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Println("Server is running " + httpAddr)

	go func() {
		server := &http.Server{
			Addr:    httpAddr,
			Handler: handler,
		}
		errs <- server.ListenAndServe()
	}()

	log.Println("exit", <-errs)
}
