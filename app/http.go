package app

import (
	"fmt"
	"github.com/brietsparks/resumapp-service/app/store"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

type RoutesParams struct {
	Router *gin.Engine
	SysDomain string
	AppName string
	Logger Logger
	FactsStore *store.FactsStore
	AuthClient AuthClient
}

func Routes(p RoutesParams) {
	p.Router.GET("/health", NewHealthCheckHandler(p.AppName))
	p.Router.GET("/user/:user_id/facts", NewGetFactsHandler(p.FactsStore, p.Logger))
	p.Router.POST("user/:user_id/facts", NewPostFactsHandler(p.FactsStore, p.Logger, p.AuthClient))
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
			c.JSON(500, gin.H{"message": fmt.Sprintf("Error retrieving fact data for user %s", userId)})
			return
		}

		if facts == "" {
			c.JSON(404, gin.H{"message": fmt.Sprintf("No data found for user %s", userId)})
			return
		}

		c.JSON(200, gin.H{
			"result": facts,
		})
	}

	return GetFactsHandler
}

type PostFactsRequestPayload struct {
	Facts string
}

func NewPostFactsHandler(factsStore *store.FactsStore, log Logger, authClient AuthClient) gin.HandlerFunc {
	PostFactsHandler := func(c *gin.Context) {
		cookie, err := c.Request.Cookie("authn")

		token := cookie.Value

		subject, err := authClient.SubjectFrom(token)

		if err != nil {
			spew.Dump(err)
			log.Error(err)
		}

		c.JSON(200, gin.H{
			"token": token,
			"err": err,
			"subject": subject,
		})
		return

		userId := c.Param("user_id")

		var p PostFactsRequestPayload
		c.BindJSON(&p)

		err = factsStore.UpsertFactsByUserId(userId, p.Facts)

		if err != nil {
			c.JSON(500, gin.H{"message": fmt.Sprintf("Error saving data for user %s", userId)})
			return
		}

		c.JSON(200, gin.H{"result": "ok"})
	}

	return PostFactsHandler
}
