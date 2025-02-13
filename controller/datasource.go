package controller

import (
	"github.com/allape/gocrud"
	"github.com/allape/goview/assets"
	"github.com/allape/goview/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io/fs"
	"net/http"
	"path"
	"strconv"
	"sync"
	"time"
)

type FileInfo struct {
	Name       string    `json:"name"`
	IsDir      bool      `json:"isDir"`
	Size       int64     `json:"size"`
	MTime      time.Time `json:"mtime"`
	HasPreview bool      `json:"hasPreview"`
}

func SetupDatasourceController(group *gin.RouterGroup, db *gorm.DB) error {
	err := gocrud.New(group, db, gocrud.CRUD[model.Datasource]{
		EnableGetAll:  true,
		DisableCount:  true,
		DisablePage:   true,
		DisableGetOne: true,
		SearchHandlers: map[string]gocrud.SearchHandler{
			"deleted": gocrud.NewSoftDeleteSearchHandler(""),
		},
		OnDelete: gocrud.NewSoftDeleteHandler[model.Datasource](gocrud.RestCoder),
	})

	if err != nil {
		return err
	}

	group.GET("/readdir/:datasource/*wd", func(context *gin.Context) {
		datasourceId := context.Param("datasource")
		wd := context.Param("wd")

		var datasource model.Datasource
		if err := db.First(&datasource, datasourceId).Error; err != nil {
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.NotFound(), err)
			return
		}

		dfs, err := model.GetFS(datasource)
		if err != nil {
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.InternalServerError(), err)
			return
		}

		entries, err := dfs.ReadDir(wd)
		if err != nil {
			gocrud.MakeErrorResponse(context, gocrud.RestCoder.InternalServerError(), err)
			return
		}

		var waitGroup sync.WaitGroup

		files := make([]FileInfo, len(entries))
		for i, entry := range entries {
			waitGroup.Add(1)
			go func(i int, entry fs.DirEntry) {
				defer waitGroup.Done()

				info, err := entry.Info()
				if err != nil {
					return
				}

				hasPreview := false

				key := model.BuildPreviewKey(datasource, path.Join(wd, info.Name()))
				var preview model.Preview
				if err := db.First(&preview, "`key` = ?", key).Error; err == nil {
					hasPreview = true
				}

				files[i] = FileInfo{
					Name:       info.Name(),
					IsDir:      info.IsDir(),
					Size:       info.Size(),
					MTime:      info.ModTime(),
					HasPreview: hasPreview,
				}
			}(i, entry)
		}

		context.JSON(http.StatusOK, gocrud.R[[]FileInfo]{
			Code: gocrud.RestCoder.OK(),
			Data: files,
		})
	})

	group.GET("/fetch/:datasource/*wd", func(context *gin.Context) {
		datasourceId := context.Param("datasource")
		wd := context.Param("wd")

		var datasource model.Datasource
		if err := db.First(&datasource, datasourceId).Error; err != nil {
			context.Header("Cache-Control", "no-cache")
			context.Data(http.StatusNotFound, assets.MIMEType, assets.IV404)
			return
		}

		dfs, err := model.GetFS(datasource)
		if err != nil {
			l.Error().Printf("Failed to get fs for datasource %d: %v", datasource.ID, err)
			context.Header("Cache-Control", "no-cache")
			context.Data(http.StatusInternalServerError, assets.MIMEType, assets.IV500)
			return
		}

		file, err := dfs.Open(wd)
		if err != nil {
			l.Error().Printf("Failed to open file %s: %v", wd, err)
			context.Header("Cache-Control", "no-cache")
			context.Data(http.StatusInternalServerError, assets.MIMEType, assets.IV500)
			return
		}

		stat, err := file.Stat()
		if err != nil {
			l.Error().Printf("Failed to stat file %s: %v", wd, err)
			context.Header("Cache-Control", "no-cache")
			context.Data(http.StatusInternalServerError, assets.MIMEType, assets.IV500)
			return
		} else if stat.IsDir() {
			context.Header("Cache-Control", "no-cache")
			context.Data(http.StatusMethodNotAllowed, assets.MIMEType, assets.IVNoPreview)
			return
		}

		contentType := "stream/octet"

		key := model.BuildPreviewKey(datasource, wd)
		var preview model.Preview
		if err := db.First(&preview, "`key` = ?", key).Error; err == nil {
			contentType = preview.MIME
		}

		context.Header("Content-Type", contentType)
		context.Header("Content-Length", strconv.FormatInt(stat.Size(), 10))
		context.Header("Last-Modified", stat.ModTime().Format(http.TimeFormat))
		context.Writer.WriteHeaderNow()
		context.Writer.Flush()

		_, err = file.WriteTo(context.Writer)
		if err != nil {
			l.Error().Printf("Failed to write file %s: %v", wd, err)
			return
		}
		context.Writer.Flush()
	})

	return nil
}
