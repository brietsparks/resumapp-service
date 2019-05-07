package app

import (
	"github.com/gin-gonic/gin"
)

type Routes func(r *gin.Engine, cfg *Config)

func MakeRoutes(log Logger) Routes {
	Routes := func(r *gin.Engine, cfg *Config) {
		setCookie := MakeSetCookie(cfg.SysDomain, cfg.IsDev)
		r.GET("/hello", NewHelloHandler(setCookie, cfg.AppName))
	}

	return Routes
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
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.SetCookie(name, value, maxAge, "/", domain, !insecure, true)
	}
}
