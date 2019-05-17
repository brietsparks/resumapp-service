package app

import (
	"flag"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

type Config struct {
	AppDomain     string
	SysDomain     string
	Port          string
	AppName       string
	IsDev         bool
	Insecure      bool
	SecretKeyPath string
	CertPath      string
	ClientOrigins []string
	DbUrl         string
	DbDriver      string
	Auth0CertPath string
	Auth0Audience string
	Auth0Issuer	  string
}

func SetConfigFromEnv(cfg *Config, log Logger) *Config {
	// environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// variables from flags
	var appDomain, sysDomain, port, appName, dbUrl, dbDriver, auth0CertPath, auth0Audience, auth0Issuer string

	flag.StringVar(&appDomain, "appDomain", "", "Service domain")
	flag.StringVar(&sysDomain, "sysDomain", "", "System domain")
	flag.StringVar(&port, "port", "", "Service port")
	flag.StringVar(&appName, "appName", "", "Unique string identifier of the service")
	flag.StringVar(&dbUrl, "dbUrl", "", "Database url")
	flag.StringVar(&dbDriver, "dbDriver", "", "Database driver")
	flag.StringVar(&auth0CertPath, "auth0CertPath", "", "Auth0 cert path")
	flag.StringVar(&auth0Audience, "auth0Audience", "", "Auth0 client id")
	flag.StringVar(&auth0Issuer, "auth0Issuer", "", "Auth0 issuer url")

	flag.Parse()

	// variables fallback to .env file
	isDev := os.Getenv("ENV") == "dev"
	insecure := os.Getenv("INSECURE") == "true"
	secretKeyPath := os.Getenv("SECRET_KEY_PATH")
	certPath := os.Getenv("CERT_PATH")

	if appDomain == "" {
		appDomain = os.Getenv("APP_DOMAIN")
	}

	if sysDomain == "" {
		sysDomain = os.Getenv("SYS_DOMAIN")
	}

	if port == "" {
		port = os.Getenv("PORT")
	}

	if auth0CertPath == "" {
		auth0CertPath = os.Getenv("AUTH0_CERT_PATH")
	}

	if auth0Audience == "" {
		auth0Audience = os.Getenv("AUTH0_AUDIENCE")
	}

	if auth0Issuer == "" {
		auth0Issuer = os.Getenv("AUTH0_ISSUER")
	}

	clientOrigins := strings.Split(os.Getenv("CLIENT_ORIGINS"), ",")

	if dbUrl == "" {
		dbUrl = os.Getenv("DB_URL")
	}

	if dbDriver == "" {
		dbDriver = os.Getenv("DB_DRIVER")
	}

	// dev variable defaults
	if isDev && secretKeyPath == "" {
		secretKeyPath = "./dev-server.key"
	}

	if isDev && certPath == "" {
		certPath = "./dev-server.crt"
	}

	// non-dev variable defaults
	if !isDev && port == "" {
		port = "443"
	}

	cfg.AppDomain = appDomain
	cfg.SysDomain = sysDomain
	cfg.Port = port
	cfg.AppName = appName
	cfg.IsDev = isDev
	cfg.Insecure = insecure
	cfg.SecretKeyPath = secretKeyPath
	cfg.CertPath = certPath
	cfg.ClientOrigins = clientOrigins
	cfg.DbDriver = dbDriver
	cfg.DbUrl = dbUrl
	cfg.Auth0CertPath = auth0CertPath
	cfg.Auth0Audience = auth0Audience
	cfg.Auth0Issuer = auth0Issuer

	return cfg
}
