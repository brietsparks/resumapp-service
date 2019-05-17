package app

import (
	"fmt"
	"github.com/brietsparks/resumapp-service/app/models"
	"github.com/brietsparks/resumapp-service/app/store"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RoutesParams struct {
	Router         *gin.Engine
	SysDomain      string
	AppName        string
	Logger         Logger
	FactsStore     *store.FactsStore
	ProfilesStore  *store.ProfilesStore
	ValidateToken  ValidateToken
}

func Routes(p RoutesParams) {
	p.Router.POST("/dummy", NewDummyHandler(p.ValidateToken))

	p.Router.GET("/handle-availability/:handle", NewGetHandleAvailabilityHandler(p.ProfilesStore, p.Logger))
	p.Router.GET("/profile/:handle", NewGetProfileByHandleHandler(p.ProfilesStore, p.Logger))
	p.Router.GET("/user/:user_id/profile", NewGetProfileByUserIdHandler(p.ProfilesStore, p.Logger))
	p.Router.POST("user/:user_id/profile", NewPostProfileHandler(p.ProfilesStore, p.Logger, p.ValidateToken))
	p.Router.GET("/user/:user_id/facts", NewGetFactsHandler(p.FactsStore, p.Logger))
	p.Router.POST("/user/:user_id/facts", NewPostFactsHandler(p.FactsStore, p.Logger, p.ValidateToken))
}

func handleError(c *gin.Context, log Logger, status int, err error, publicErr string) {
	log.Error(fmt.Sprintf("%s; %s", publicErr, err))
	c.JSON(status, gin.H{"message": publicErr})
}

func NewDummyHandler(validateToken ValidateToken) gin.HandlerFunc {
	DummyHandler := func(c *gin.Context) {
		token, claims, err := validateToken(c.Request)
		if err != nil {
			spew.Dump(err)

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": true,
				"token": token,
			})
			c.Abort()
			return
		}

		userId := claims.Subject

		c.JSON(http.StatusOK, gin.H{"user_id": userId})
	}

	return DummyHandler
}

func NewGetHandleAvailabilityHandler(profilesStore *store.ProfilesStore, log Logger) gin.HandlerFunc {
	GetHandleAvailabilityHandler := func(c *gin.Context) {
		handle := c.Param("handle")
		isAvailable, err := profilesStore.GetHandleAvailability(handle)

		if err != nil {
			handleError(c, log, http.StatusInternalServerError, err,
				fmt.Sprintf("error checking availability of handle %s", handle))
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": isAvailable})
	}

	return GetHandleAvailabilityHandler
}

func NewGetProfileByHandleHandler(profilesStore *store.ProfilesStore, log Logger) gin.HandlerFunc {
	GetProfileByHandleHandler := func(c *gin.Context) {
		handle := c.Param("handle")

		profile, err := profilesStore.GetProfileByHandle(handle)

		if err != nil {
			handleError(c, log, http.StatusInternalServerError, err,
				fmt.Sprintf("error retrieving profile by handle %s", handle))
			c.Abort()
			return
		}

		if profile == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("no profile found for handle %s", handle)})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": profile})
	}

	return GetProfileByHandleHandler
}

func NewGetProfileByUserIdHandler(profilesStore *store.ProfilesStore, log Logger) gin.HandlerFunc {
	GetProfileHandler := func(c *gin.Context) {
		userId := c.Param("user_id")

		profile, err := profilesStore.GetProfileByUserId(userId)

		if err != nil {
			handleError(c, log, http.StatusInternalServerError, err,
				fmt.Sprintf("error retrieving profile data for user %s", userId))
			c.Abort()
			return
		}

		if profile == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("no profile found for user %s", userId)})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": profile})
	}

	return GetProfileHandler
}

func NewPostProfileHandler(profilesStore *store.ProfilesStore, log Logger, validateToken ValidateToken) gin.HandlerFunc {
	PostProfileHandler := func(c *gin.Context) {
		// authorize via token
		_, claims, err := validateToken(c.Request)
		if err != nil {
			handleError(c, log, http.StatusUnauthorized, err,
				fmt.Sprintf("invalid token"))
			c.Abort()
			return
		}

		subject := claims.Subject
		userId := c.Param("user_id")

		if subject != userId {
			handleError(c, log, http.StatusBadRequest, err,
				fmt.Sprintf("token and url parameter are not the same userId"))
			c.Abort()
			return
		}

		// save the data
		var profile models.Profile
		c.BindJSON(&profile)

		profile.UserId = userId

		err = profilesStore.UpsertProfileByUserId(userId, profile)

		if err != nil {
			handleError(c, log, http.StatusInternalServerError, err,
				fmt.Sprintf("error saving profile data for user %s", userId))
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": "ok"})
	}

	return PostProfileHandler
}

func NewGetFactsHandler(factsStore *store.FactsStore, log Logger) gin.HandlerFunc {
	GetFactsHandler := func(c *gin.Context) {
		userId := c.Param("user_id")
		facts, err := factsStore.GetFactsByUserId(userId)

		if err != nil {
			handleError(c, log, http.StatusInternalServerError, err,
				fmt.Sprintf("error retrieving fact data for user %s", userId))
			return
		}

		if facts == "" {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("no fact data found for user %s", userId)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": facts})
	}

	return GetFactsHandler
}

func NewPostFactsHandler(factsStore *store.FactsStore, log Logger, validateToken ValidateToken) gin.HandlerFunc {
	type postFactsRequestPayload struct {
		Facts string
	}

	PostFactsHandler := func(c *gin.Context) {
		// authorize via token
		_, claims, err := validateToken(c.Request)
		if err != nil {
			handleError(c, log, http.StatusUnauthorized, err,
				fmt.Sprintf("invalid token"))
			c.Abort()
			return
		}

		subject := claims.Subject
		userId := c.Param("user_id")

		if subject != userId {
			handleError(c, log, http.StatusBadRequest, err,
				fmt.Sprintf("token and url parameter are not the same userId"))
			c.Abort()
			return
		}

		// save the data
		var p postFactsRequestPayload
		c.BindJSON(&p)

		err = factsStore.UpsertFactsByUserId(userId, p.Facts)

		if err != nil {
			handleError(c, log, http.StatusInternalServerError, err,
				fmt.Sprintf("error saving fact data for user %s", userId))
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": "ok"})
	}

	return PostFactsHandler
}
