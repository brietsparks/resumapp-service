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
}

type SetConfigFromEnv func(cfg *Config) *Config

func MakeSetConfigFromEnv(log Logger) SetConfigFromEnv {
	SetConfigFromEnv := func(cfg *Config) *Config {
		// environment variables
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		// variables from flags
		var appDomain, sysDomain, port, appName string
		flag.StringVar(&appDomain, "appDomain", "", "Service domain")
		flag.StringVar(&sysDomain, "sysDomain", "", "System domain")
		flag.StringVar(&port, "port", "", "Service port")
		flag.StringVar(&appName, "appName", "", "Unique string identifier of the service")
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

		return cfg
	}

	return SetConfigFromEnv
}
