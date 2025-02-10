package controller

import (
	"github.com/allape/gocrud"
	"github.com/allape/goview/assets"
	"github.com/allape/goview/env"
	"github.com/allape/goview/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path"
)

func SetupPreviewController(group *gin.RouterGroup, db *gorm.DB) error {
	group.PUT("/:datasource/*filename", func(context *gin.Context) {
		datasourceId := context.Param("datasource")
		filename := context.Param("filename")

		var datasource model.Datasource
		if err := db.Model(&datasource).First(&datasource, datasourceId).Error; err != nil {
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.NotFound(), err.Error())
			return
		}

		preview, err := model.GeneratePreview(datasource, filename, env.PreviewFolder, func(digest string) (*model.Preview, error) {
			var pre model.Preview
			err := db.Model(&pre).First(&pre, "`digest` = ?", digest).Error
			return &pre, err
		})
		if err != nil {
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.InternalServerError(), err.Error())
			return
		}

		context.JSON(http.StatusOK, gocrud.R[model.Preview]{
			Code: gocrud.RestCoder.OK(),
			Data: *preview,
		})
	})

	group.GET("/:datasource/*filename", func(context *gin.Context) {
		datasourceId := context.Param("datasource")
		filename := context.Param("filename")

		var datasource model.Datasource
		if err := db.Model(&datasource).First(&datasource, datasourceId).Error; err != nil {
			context.Data(http.StatusNotFound, assets.MIMEType, assets.IV404)
			return
		}

		key := model.BuildPreviewKey(datasource, filename)

		var preview model.Preview
		if err := db.Model(&preview).First(&preview, "`key` =? ", key).Error; err != nil {
			context.Data(http.StatusNotFound, assets.MIMEType, assets.IV404)
			return
		}

		cover := path.Join(env.PreviewFolder, preview.Cover)

		stat, err := os.Stat(cover)
		if err != nil {
			context.Data(http.StatusNotFound, assets.MIMEType, assets.IV404)
			return
		} else if stat.IsDir() {
			context.Data(http.StatusMethodNotAllowed, assets.MIMEType, assets.IVNoPreview)
			return
		}

		context.File(cover)
	})

	return nil
}
