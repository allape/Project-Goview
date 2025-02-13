package controller

import (
	"github.com/allape/goview/assets"
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
		context.Data(http.StatusOK, assets.MIMEType, assets.Favicon)
	})

	return nil
}
