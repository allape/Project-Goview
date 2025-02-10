package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupUIController(engine *gin.Engine, folder string) error {
	grp := engine.Group("/ui")
	grp.Static("/", folder)

	engine.GET("/", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "/ui")
	})
	engine.GET("/favicon.ico", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "/ui/favicon.png")
	})

	return nil
}
