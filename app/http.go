package app

import (
	"fmt"
	"github.com/brietsparks/resumapp-service/app/models"
	"github.com/brietsparks/resumapp-service/app/store"
	"github.com/gin-gonic/gin"
)

type RoutesParams struct {
	Router        *gin.Engine
	SysDomain     string
	AppName       string
	Logger        Logger
	FactsStore    *store.FactsStore
	ProfilesStore *store.ProfilesStore
}

func Routes(p RoutesParams) {
	p.Router.GET("/handle-availability/:handle", NewGetHandleAvailabilityHandler(p.ProfilesStore, p.Logger))
	p.Router.GET("/profile/:handle", NewGetProfileByHandleHandler(p.ProfilesStore, p.Logger))
	p.Router.GET("/user/:user_id/profile", NewGetProfileByUserIdHandler(p.ProfilesStore, p.Logger))
	p.Router.POST("user/:user_id/profile", NewPostProfileHandler(p.ProfilesStore, p.Logger))
	p.Router.GET("/user/:user_id/facts", NewGetFactsHandler(p.FactsStore, p.Logger))
	p.Router.POST("/user/:user_id/facts", NewPostFactsHandler(p.FactsStore, p.Logger))
}

func handleError(c *gin.Context, log Logger, status int, err error, publicErr string) {
	log.Error(fmt.Sprintf("%s; %s", publicErr, err))
	c.JSON(status, gin.H{"message": publicErr})
}

func NewGetHandleAvailabilityHandler(profilesStore *store.ProfilesStore, log Logger) gin.HandlerFunc {
	GetHandleAvailabilityHandler := func(c *gin.Context) {
		handle := c.Param("handle")
		isAvailable, err := profilesStore.GetHandleAvailability(handle)

		if err != nil {
			handleError(c, log, 500, err,
				fmt.Sprintf("error checking availability of handle %s", handle))
			return
		}

		c.JSON(200, gin.H{"result": isAvailable})
	}

	return GetHandleAvailabilityHandler
}

func NewGetProfileByHandleHandler(profilesStore *store.ProfilesStore, log Logger) gin.HandlerFunc {
	GetProfileByHandleHandler := func(c *gin.Context) {
		handle := c.Param("handle")

		profile, err := profilesStore.GetProfileByHandle(handle)

		if err != nil {
			handleError(c, log, 500, err,
				fmt.Sprintf("error retrieving profile by handle %s", handle))
			return
		}

		if profile == nil {
			c.JSON(404, gin.H{"message": fmt.Sprintf("no profile found for handle %s", handle)})
			return
		}

		c.JSON(200, gin.H{"result": profile})
	}

	return GetProfileByHandleHandler
}

func NewGetProfileByUserIdHandler(profilesStore *store.ProfilesStore, log Logger) gin.HandlerFunc {
	GetProfileHandler := func(c *gin.Context) {
		userId := c.Param("user_id")

		profile, err := profilesStore.GetProfileByUserId(userId)

		if err != nil {
			handleError(c, log, 500, err,
				fmt.Sprintf("error retrieving profile data for user %s", userId))
			return
		}

		if profile == nil {
			c.JSON(404, gin.H{"message": fmt.Sprintf("no profile found for user %s", userId)})
			return
		}

		c.JSON(200, gin.H{"result": profile})
	}

	return GetProfileHandler
}

func NewPostProfileHandler(profilesStore *store.ProfilesStore, log Logger) gin.HandlerFunc {
	PostProfileHandler := func(c *gin.Context) {
		userId := c.Param("user_id")

		var profile models.Profile
		c.BindJSON(&profile)

		profile.UserId = userId

		err := profilesStore.UpsertProfileByUserId(userId, profile)

		if err != nil {
			handleError(c, log, 500, err,
				fmt.Sprintf("error saving profile data for user %s", userId))
			return
		}

		c.JSON(200, gin.H{"result": "ok"})
	}

	return PostProfileHandler
}

func NewGetFactsHandler(factsStore *store.FactsStore, log Logger) gin.HandlerFunc {
	GetFactsHandler := func(c *gin.Context) {
		userId := c.Param("user_id")
		facts, err := factsStore.GetFactsByUserId(userId)

		if err != nil {
			handleError(c, log, 500, err,
				fmt.Sprintf("error retrieving fact data for user %s", userId))
			return
		}

		if facts == "" {
			c.JSON(404, gin.H{"message": fmt.Sprintf("no fact data found for user %s", userId)})
			return
		}

		c.JSON(200, gin.H{"result": facts})
	}

	return GetFactsHandler
}

func NewPostFactsHandler(factsStore *store.FactsStore, log Logger) gin.HandlerFunc {
	type postFactsRequestPayload struct {
		Facts string
	}

	PostFactsHandler := func(c *gin.Context) {
		//token := c.GetHeader("Authorization")
		//c.JSON(200, gin.H{"token": token})
		//return

		userId := c.Param("user_id")

		var p postFactsRequestPayload
		c.BindJSON(&p)

		err := factsStore.UpsertFactsByUserId(userId, p.Facts)

		if err != nil {
			handleError(c, log, 500, err,
				fmt.Sprintf("error saving fact data for user %s", userId))
			return
		}

		c.JSON(200, gin.H{"result": "ok"})
	}

	return PostFactsHandler
}
