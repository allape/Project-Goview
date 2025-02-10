package controller

import (
	"github.com/allape/gocrud"
	"github.com/allape/goview/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

	return nil
}
