package main

import (
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	// set up logging
	log.SetFormatter(&log.JSONFormatter{})

	file, err := os.OpenFile("log.txt", os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)

	// environment variables
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// variables from flags
	var appDomain, parentDomain, port, appName string
	flag.StringVar(&appDomain,"appDomain", "", "Service domain")
	flag.StringVar(&parentDomain,"parentDomain", "", "System domain")
	flag.StringVar(&port,"port", "", "Service port")
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

	if parentDomain == "" {
		parentDomain = os.Getenv("PARENT_DOMAIN")
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


	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     clientOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Credentials"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// routes
	setCookie := MakeSetCookie(parentDomain, isDev)
	r.GET("/hello", NewHelloHandler(setCookie, appName))

	// http Server
	port = fmt.Sprintf(":%s", port)
	if insecure {
		fmt.Println(fmt.Sprintf("Listening on port %s INSECURE", port))
		r.Run(port)
	} else {
		fmt.Println(fmt.Sprintf("Listening on port %s via TLS", port))
		r.RunTLS(port, certPath, secretKeyPath)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func NewHelloHandler(setCookie SetCookie, appName string) gin.HandlerFunc {
	HelloHandler := func(c *gin.Context) {
		setCookie(c, "dummy", "ok", 3600)
		c.JSON(200, gin.H{
			"ok": "ok",
			"appName": appName,
		})
	}

	return HelloHandler
}

type SetCookie = func(c *gin.Context, name string, value string, maxAge int)

func MakeSetCookie(domain string, insecure bool) SetCookie {
	return func(c *gin.Context, name string, value string, maxAge int) {
		fmt.Println(domain)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.SetCookie(name, value, maxAge, "/", domain, !insecure, true)
	}
}

func getDomain(su string) string {
	if su == "" {
		return ""
	}

	u, err := url.Parse(su)
	if err != nil {
		log.Fatal(err)
	}
	parts := strings.Split(u.Hostname(), ".")
	domain := parts[len(parts)-2] + "." + parts[len(parts)-1]
	return domain
}
