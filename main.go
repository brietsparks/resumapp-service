package main

import (
	"github.com/brietsparks/resumapp-service/app"
	"github.com/keratin/authn-go/authn"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
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
	logger.SetOutput(file)

	config := app.SetConfigFromEnv(&app.Config{}, logger)
	authClient, err := authn.NewClient(authn.Config{
		Issuer:         config.AuthUrl,
		Audience:       config.SysDomain,
		Username:       config.AuthUsername,
		Password:       config.AuthPassword,
	})

	if err != nil {
		log.Fatal(err)
	}

	server := app.NewServer(config, logger, authClient)

	server.Run()

	if err != nil {
		log.Fatal(err)
	}
}
