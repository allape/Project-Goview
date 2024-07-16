package datasource

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/allape/goview/assets"
	"github.com/allape/goview/base"
	"github.com/allape/goview/datasource/driver"
	"github.com/allape/goview/datasource/driver/dufs"
	"github.com/allape/goview/datasource/driver/local"
	"github.com/allape/goview/env"
	"github.com/allape/goview/util"
	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"
	"gorm.io/gorm"
	"image"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

// region driver

type Type string

const (
	DUFS  Type = "dufs"
	LOCAL Type = "local"
)

var (
	locker  = make(chan struct{}, 1)
	drivers = map[Type]driver.Driver{}
)

func GetDriver(t Type) (driver.Driver, error) {
	locker <- struct{}{}
	defer func() {
		<-locker
	}()

	d, ok := drivers[t]
	if ok {
		return d, nil
	}

	switch t {
	case DUFS:
		d := &dufs.Driver{}

		caCertPool, err := env.TrustedCertsPoolFromEnv()
		if err != nil {
			return nil, err
		}
		tlsConfig := &tls.Config{
			RootCAs: caCertPool,
		}
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		}

		err = d.Setup(client)
		if err != nil {
			return nil, err
		}
		drivers[t] = d
	case LOCAL:
		drivers[t] = &local.Driver{}
	default:
		return nil, errors.New("type is not supported")
	}

	return drivers[t], nil
}

// endregion

type Preview struct {
	base.Model
	DatasourceID uint   `json:"datasourceId"`
	Key          string `json:"key" gorm:"uniqueIndex;type:varchar(255)"`
	Digest       string `json:"digest"`
	Cover        string `json:"cover"`
	FFProbeInfo  string `json:"ffprobeInfo"`
}

type Datasource struct {
	base.Model
	Name string `json:"name"`
	Type Type   `json:"type"`
	Cwd  string `json:"cwd"`
}

var GeneratePreviewLocker = make(chan struct{}, 1)

func BuildPreviewKey(datasource Datasource, file string) string {
	return fmt.Sprintf("goview://%d?file=%s", datasource.ID, url.QueryEscape(file))
}

func GeneratePreview(source driver.Driver, datasource Datasource, sourceFile, dstFolder string, finder func(digest string) (*Preview, error)) (*Preview, error) {
	GeneratePreviewLocker <- struct{}{}
	defer func() {
		<-GeneratePreviewLocker
	}()

	fullname := source.PathJoin(datasource.Cwd, sourceFile)

	f, err := source.Status(fullname)
	if err != nil {
		return nil, err
	} else if f.IsDir {
		return nil, errors.New("it is a directory")
	}

	tmpFile, err := os.CreateTemp(os.TempDir(), "goview_*_"+f.Name)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tmpFile.Close()
	}()

	err = source.Concatenate(fullname, tmpFile)
	if err != nil {
		return nil, err
	}

	digest, err := util.Sha256(tmpFile)
	if err != nil {
		return nil, err
	}

	key := BuildPreviewKey(datasource, sourceFile)

	found, err := finder(digest)
	if err == nil {
		found.ID = 0
		found.CreatedAt = time.Now()
		found.UpdatedAt = time.Now()
		found.DeletedAt = gorm.DeletedAt{}
		found.DatasourceID = datasource.ID
		found.Key = key
		return found, nil
	}

	prev := Preview{
		DatasourceID: datasource.ID,
		Key:          key,
		Digest:       digest,
	}

	fileType, err := filetype.MatchFile(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	dtsFile := digest + ".jpg"

	prev.FFProbeInfo, err = util.FFProbeInfo(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	switch fileType.MIME.Type {
	case "image":
		_, err = util.FFMpegScaleImage(path.Join(dstFolder, dtsFile), tmpFile.Name(), 0.1)
		if err != nil {
			return nil, err
		}
	case "video":
		_, err = util.FFMpegVideoSampleImage(tmpFile.Name(), path.Join(dstFolder, dtsFile), 4, image.Point{X: 10, Y: 10})
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("not supported")
	}

	prev.Cover = dtsFile

	return &prev, nil
}

func Setup(repo *gorm.DB, rout *gin.Engine, previewFolder string) error {
	err := repo.AutoMigrate(&Datasource{}, &Preview{})
	if err != nil {
		return err
	}

	ds := rout.Group("/datasource")

	ds.GET("/all", func(context *gin.Context) {
		var datasources []Datasource
		repo.Find(&datasources)
		context.JSON(http.StatusOK, base.R[[]Datasource]{
			Code: "200",
			Data: datasources,
		})
	})

	ds.POST("/save", func(context *gin.Context) {
		var datasource Datasource
		err := context.BindJSON(&datasource)
		if err != nil {
			context.JSON(http.StatusBadRequest, base.R[any]{
				Code:    "400",
				Message: err.Error(),
			})
			return
		}

		if datasource.Name == "" {
			context.JSON(http.StatusBadRequest, base.R[any]{
				Code:    "400",
				Message: "name is required",
			})
			return
		}

		if datasource.Type == "" {
			context.JSON(http.StatusBadRequest, base.R[any]{
				Code:    "400",
				Message: "type is required",
			})
			return
		}

		_, err = GetDriver(datasource.Type)
		if err != nil {
			context.JSON(http.StatusBadRequest, base.R[any]{
				Code:    "400",
				Message: err.Error(),
			})
			return
		}

		datasource.Name = strings.TrimSpace(datasource.Name)
		datasource.Cwd = strings.TrimSpace(datasource.Cwd)

		repo.Save(&datasource)
		context.JSON(http.StatusOK, base.R[Datasource]{
			Code: "200",
			Data: datasource,
		})
	})

	// /get/1/index.html
	ds.GET("/:func/:id/*file", func(context *gin.Context) {
		var datasource Datasource

		function := context.Param("func")

		repo.First(&datasource, context.Param("id"))

		if datasource.ID == 0 {
			context.JSON(http.StatusNotFound, base.R[any]{
				Code:    "404",
				Message: "not found",
			})
			return
		}

		file := context.Param("file")[1:]

		source, err := GetDriver(datasource.Type)
		if err != nil {
			context.JSON(http.StatusInternalServerError, base.R[any]{
				Code:    "500",
				Message: err.Error(),
			})
			return
		}

		file = source.PathJoin(datasource.Cwd, file)

		switch function {
		case "stat":
			stat, err := source.Status(file)
			if err != nil {
				context.JSON(http.StatusInternalServerError, base.R[any]{
					Code:    "500",
					Message: err.Error(),
				})
				return
			}
			context.JSON(http.StatusOK, base.R[driver.File]{
				Code: "200",
				Data: *stat,
			})
		case "ls":
			files, err := source.List(file)
			if err != nil {
				context.JSON(http.StatusInternalServerError, base.R[any]{
					Code:    "500",
					Message: err.Error(),
				})
				return
			}

			context.JSON(http.StatusOK, base.R[[]driver.File]{
				Code: "200",
				Data: files,
			})
		case "cat":
			context.Status(http.StatusOK)
			context.Writer.WriteHeaderNow()

			err := source.Concatenate(file, context.Writer)
			if err != nil {
				log.Println("failed to redirect data:", err)
				return
			}

			context.Writer.Flush()
		default:
			context.JSON(http.StatusBadRequest, base.R[any]{
				Code:    "400",
				Message: "function is not supported",
			})
		}
	})

	const NoPreviewRouter = "/preview/image/no-preview.jpg"

	preview := rout.Group("/preview")

	preview.GET("/gen/:dsid/*file", func(context *gin.Context) {
		var datasource Datasource
		repo.First(&datasource, context.Param("dsid"))
		if datasource.ID == 0 {
			context.JSON(http.StatusNotFound, base.R[any]{
				Code:    "404",
				Message: "not found",
			})
			return
		}

		file := context.Param("file")[1:]

		source, err := GetDriver(datasource.Type)
		if err != nil {
			context.JSON(http.StatusInternalServerError, base.R[any]{
				Code:    "500",
				Message: err.Error(),
			})
			return
		}

		key := BuildPreviewKey(datasource, file)

		var pre Preview
		repo.First(&pre, "`key` = ?", key)

		if pre.ID == 0 {
			prev, err := GeneratePreview(source, datasource, file, previewFolder, func(digest string) (*Preview, error) {
				var prev Preview
				err := repo.First(&prev, "`digest` = ?", digest).Error
				return &prev, err
			})
			if err != nil {
				context.JSON(http.StatusInternalServerError, base.R[any]{
					Code:    "500",
					Message: err.Error(),
				})
				return
			}
			pre = *prev
			repo.Save(&pre)
		}

		context.JSON(http.StatusOK, base.R[Preview]{
			Code: "200",
			Data: pre,
		})
	})

	preview.GET("/get/:dsid/*file", func(context *gin.Context) {
		var datasource Datasource
		err := repo.First(&datasource, context.Param("dsid")).Error
		if err != nil {
			log.Println(err)
			context.Data(http.StatusInternalServerError, "image/jpeg", assets.IV500)
			return
		}
		if datasource.ID == 0 {
			context.Data(http.StatusNotFound, "image/jpeg", assets.IV404)
			return
		}

		file := context.Param("file")[1:]
		key := BuildPreviewKey(datasource, file)

		var pre Preview
		repo.First(&pre, "`key` = ?", key)

		if pre.ID == 0 {
			context.Redirect(http.StatusFound, NoPreviewRouter)
			return
		}

		context.Header("X-FFProbe", pre.FFProbeInfo)
		context.File(path.Join(previewFolder, pre.Cover))
	})

	preview.GET("/info/:dsid/*file", func(context *gin.Context) {
		var datasource Datasource
		repo.First(&datasource, context.Param("dsid"))
		if datasource.ID == 0 {
			context.JSON(http.StatusNotFound, base.R[any]{
				Code:    "404",
				Message: "not found",
			})
			return
		}

		file := context.Param("file")[1:]
		key := BuildPreviewKey(datasource, file)

		var pre Preview
		repo.First(&pre, "`key` = ?", key)

		if pre.ID == 0 {
			context.JSON(http.StatusNotFound, base.R[any]{
				Code:    "404",
				Message: "not found",
			})
			return
		}

		context.JSON(http.StatusOK, base.R[Preview]{
			Code: "200",
			Data: pre,
		})
	})

	preview.GET("/static", func(context *gin.Context) {
		key := context.Query("key")

		if key == "" {
			context.Data(http.StatusBadGateway, "image/jpeg", assets.IV404)
			return
		}

		var pre Preview
		repo.First(&pre, "`key` = ?", key)
		if pre.ID == 0 {
			context.Redirect(http.StatusFound, NoPreviewRouter)
			return
		}

		context.Header("X-FFProbe", pre.FFProbeInfo)
		context.File(path.Join(previewFolder, pre.Cover))
	})

	rout.GET(NoPreviewRouter, func(context *gin.Context) {
		context.Data(http.StatusOK, "image/jpeg", assets.IVNoPreview)
	})

	return nil
}
