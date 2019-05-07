package app

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Server struct {
	Config *Config
	Run func()
}

type NewServer func(cfg *Config) *Server


func MakeNewServer(log Logger) NewServer {
    NewServer := func(cfg *Config) *Server {
    	s := &Server{Config: cfg}

		r := gin.Default()

		r.Use(cors.New(cors.Config{
			AllowOrigins:     cfg.ClientOrigins,
			AllowMethods:     []string{http.MethodGet, http.MethodPost},
			AllowHeaders:     []string{"Origin"},
			ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Credentials"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))

		MakeRoutes(log)(r, cfg)

		port := fmt.Sprintf(":%s", cfg.Port)
		if cfg.Insecure {
			s.Run = func() {
				fmt.Println(fmt.Sprintf("Listening on port %s INSECURE", port))
				err := r.Run(port)
				if err != nil {
					log.Fatal(err)
				}
			}
		} else {
			s.Run = func() {
				fmt.Println(fmt.Sprintf("Listening on port %s via TLS", port))
				err := r.RunTLS(port, cfg.CertPath, cfg.SecretKeyPath)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		return s
	}

    return NewServer
}
