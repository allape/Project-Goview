package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	URIUI = "/ui"
)

func SetupUIController(engine *gin.Engine, folder string) error {
	grp := engine.Group(URIUI)
	grp.Static("/", folder)

	engine.GET("/", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, URIUI)
	})
	engine.GET("/favicon.ico", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "/ui/favicon.png")
	})

	return nil
}
