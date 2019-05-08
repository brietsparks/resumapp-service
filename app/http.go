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
	setCookie := MakeSetCookie(p.SysDomain)
	p.Router.GET("/health", NewHelloHandler(setCookie, p.AppName))
	//p.Router.POST("/")
}

func NewHelloHandler(setCookie SetCookie, appName string) gin.HandlerFunc {
	HelloHandler := func(c *gin.Context) {
		setCookie(c, "dummy", "ok", 3600)
		c.JSON(200, gin.H{
			"health": "ok",
			"appName": appName,
		})
	}

	return HelloHandler
}


type SetCookie = func(c *gin.Context, name string, value string, maxAge int)

func MakeSetCookie(domain string) SetCookie {
	return func(c *gin.Context, name string, value string, maxAge int) {
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.SetCookie(name, value, maxAge, "/", domain, true, true)
	}
}
