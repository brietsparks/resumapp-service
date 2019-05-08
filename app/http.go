package app

import (
	"github.com/brietsparks/resumapp-service/app/store"
	"github.com/gin-gonic/gin"
)

type RoutesParams struct {
	Router *gin.Engine
	SysDomain string
	AppName string
	Logger Logger
	FactsStore *store.FactsStore
}

func Routes(p RoutesParams) {
	p.Router.GET("/health", NewHealthCheckHandler(p.AppName))
	p.Router.GET("/user/:user_id/facts", NewGetFactsHandler(p.FactsStore, p.Logger))
}

func NewHealthCheckHandler(appName string) gin.HandlerFunc {
	HelloHandler := func(c *gin.Context) {
		c.JSON(200, gin.H{
			"health": "ok",
			"appName": appName,
		})
	}

	return HelloHandler
}

func NewGetFactsHandler(factsStore *store.FactsStore, log Logger) gin.HandlerFunc {
	GetFactsHandler := func(c *gin.Context) {
		userId := c.Param("user_id")
		facts, err := factsStore.GetFactsByUserId(userId)

		if err != nil {
			c.JSON(500, gin.H{"message": "Error retrieving the data"})
			return
		}

		if facts == "" {
			c.JSON(404, gin.H{"message": "No data found for user"})
			return
		}

		c.JSON(200, gin.H{
			"result": facts,
		})
	}

	return GetFactsHandler
}
