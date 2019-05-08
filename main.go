package main

import (
	"github.com/brietsparks/resumapp-service/app"
	"github.com/sirupsen/logrus"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	file, err := os.OpenFile("log.txt", os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)

	config := app.SetConfigFromEnv(&app.Config{}, logger)
	server := app.NewServer(config, logger)

	server.Run()

	if err != nil {
		log.Fatal(err)
	}
}
