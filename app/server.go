package app

import (
	"fmt"
	"github.com/brietsparks/resumapp-service/app/store"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Server struct {
	Config *Config
	Run func()
}

func NewServer(cfg *Config, log Logger, authClient AuthClient) *Server {
	s := &Server{Config: cfg}

	r := gin.Default()
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.ClientOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Credentials"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// todo move this to main
	db := NewDB(log, cfg.DbDriver, cfg.DbUrl)
	factsStore := store.NewFactsStore(db, log)

	Routes(RoutesParams{
		Router: r,
		SysDomain: cfg.SysDomain,
		Logger: log,
		FactsStore: factsStore,
		AuthClient: authClient,
	})

	port := fmt.Sprintf(":%s", cfg.Port)
	fmt.Println(fmt.Sprintf("Listening on port %s", port))
	err := r.Run(port)
	if err != nil {
		log.Fatal(err)
	}

	return s
}
