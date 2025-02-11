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
	"strings"
)

func servePreviewByKey(context *gin.Context, db *gorm.DB, key model.FileKey) {
	key = model.FileKey(strings.TrimSpace(string(key)))
	if key == "" {
		context.Header("Cache-Control", "no-cache")
		context.Data(http.StatusNotFound, assets.MIMEType, assets.IV404)
		return
	}

	key = model.FileKey(strings.TrimPrefix(string(key), "/"))

	var preview model.Preview
	if err := db.Model(&preview).First(&preview, "`key` = ?", key).Error; err != nil {
		context.Header("Cache-Control", "no-cache")
		context.Data(http.StatusNotFound, assets.MIMEType, assets.IV404)
		return
	}

	cover := path.Join(env.PreviewFolder, preview.Cover)

	stat, err := os.Stat(cover)
	if err != nil {
		context.Header("Cache-Control", "no-cache")
		context.Data(http.StatusNotFound, assets.MIMEType, assets.IV404)
		return
	} else if stat.IsDir() {
		context.Header("Cache-Control", "no-cache")
		context.Data(http.StatusMethodNotAllowed, assets.MIMEType, assets.IVNoPreview)
		return
	}

	context.File(cover)
}

func SetupPreviewController(group *gin.RouterGroup, db *gorm.DB) error {
	err := gocrud.New(group, db, gocrud.CRUD[model.Preview]{
		SearchHandlers: map[string]gocrud.SearchHandler{
			"datasourceId": gocrud.KeywordEqual("datasource_id", nil),
			"mime":         gocrud.KeywordLike("mime", nil),
			"key":          gocrud.KeywordLike("key", nil),
			"ffprobeInfo":  gocrud.KeywordLike("ff_probe_info", nil),
			"digest":       gocrud.KeywordEqual("digest", nil),
			"deleted":      gocrud.NewSoftDeleteSearchHandler(""),
		},
		OnDelete: gocrud.NewSoftDeleteHandler[model.Preview](gocrud.RestCoder),
	})

	if err != nil {
		return err
	}

	group.PUT("/from-ds/:datasource/*filename", func(context *gin.Context) {
		datasourceId := context.Param("datasource")
		filename := context.Param("filename")

		var datasource model.Datasource
		if err := db.Model(&datasource).First(&datasource, datasourceId).Error; err != nil {
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.NotFound(), err.Error())
			return
		}

		key := model.BuildPreviewKey(datasource, filename)

		var found model.Preview
		if err := db.Model(&found).First(&found, "`key` = ?", key).Error; err == nil {
			context.JSON(http.StatusOK, gocrud.R[model.Preview]{
				Code: gocrud.RestCoder.OK(),
				Data: found,
			})
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

		if err := db.Save(preview).Error; err != nil {
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.InternalServerError(), err.Error())
			return
		}

		context.JSON(http.StatusOK, gocrud.R[model.Preview]{
			Code: gocrud.RestCoder.OK(),
			Data: *preview,
		})
	})

	group.GET("/by-ds/:datasource/*filename", func(context *gin.Context) {
		datasourceId := context.Param("datasource")
		filename := context.Param("filename")

		var datasource model.Datasource
		if err := db.Model(&datasource).First(&datasource, datasourceId).Error; err != nil {
			context.Header("Cache-Control", "no-cache")
			context.Data(http.StatusNotFound, assets.MIMEType, assets.IV404)
			return
		}

		key := model.BuildPreviewKey(datasource, filename)

		servePreviewByKey(context, db, key)
	})

	group.GET("/by-key/*key", func(context *gin.Context) {
		key := context.Param("key")
		servePreviewByKey(context, db, model.FileKey(key))
	})

	group.GET("/no-preview", func(context *gin.Context) {
		context.Data(http.StatusOK, assets.MIMEType, assets.IVNoPreview)
	})

	return nil
}
