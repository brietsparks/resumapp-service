package app

import (
	"fmt"
	"github.com/brietsparks/resumapp-service/app/store"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type Server struct {
	Config *Config
	Run func()
}

type ServerParams struct {
	Config *Config
	Log Logger
	FactsStore *store.FactsStore
	ProfilesStore *store.ProfilesStore
	ValidateToken ValidateToken
}

func NewServer(p ServerParams) *Server {
	s := &Server{Config: p.Config}

	r := gin.Default()
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     p.Config.ClientOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost},
		AllowHeaders:     []string{"Origin", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Credentials"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	Routes(RoutesParams{
		Router: r,
		SysDomain: p.Config.SysDomain,
		Logger: p.Log,
		FactsStore: p.FactsStore,
		ProfilesStore: p.ProfilesStore,
		ValidateToken: p.ValidateToken,
	})

	port := fmt.Sprintf(":%s", p.Config.Port)
	fmt.Println(fmt.Sprintf("Listening on port %s", port))
	err := r.Run(port)
	if err != nil {
		log.Fatal(err)
	}

	return s
}
