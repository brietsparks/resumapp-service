package main

import (
	"github.com/brietsparks/resumapp-service/app"
	"github.com/brietsparks/resumapp-service/app/store"
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

	if err != nil {
		log.Fatal(err)
	}

	db := app.NewDB(logger, config.DbDriver, config.DbUrl)
	factsStore := store.NewFactsStore(db)
	profilesStore := store.NewProfilesStore(db)
	validateToken := app.NewValidateToken(config.Auth0CertPath, config.Auth0Audience, config.Auth0Issuer)

	server := app.NewServer(app.ServerParams{
		Config: config,
		Log: logger,
		FactsStore: factsStore,
		ProfilesStore: profilesStore,
		ValidateToken: validateToken,
	})

	server.Run()

	if err != nil {
		log.Fatal(err)
	}
}
