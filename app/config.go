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
}

func SetConfigFromEnv(cfg *Config, log Logger) *Config {
	// environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// variables from flags
	var appDomain, sysDomain, port, appName, dbUrl, dbDriver string
	flag.StringVar(&appDomain, "appDomain", "", "Service domain")
	flag.StringVar(&sysDomain, "sysDomain", "", "System domain")
	flag.StringVar(&port, "port", "", "Service port")
	flag.StringVar(&appName, "appName", "", "Unique string identifier of the service")
	flag.StringVar(&dbUrl, "dbUrl", "", "Database url")
	flag.StringVar(&dbDriver, "dbDriver", "", "Database driver")
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

	return cfg
}
