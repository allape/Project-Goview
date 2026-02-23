package main

import (
	"time"

	"github.com/allape/gocrud"
	"github.com/allape/gogger"
	"github.com/allape/goview/controller"
	"github.com/allape/goview/env"
	"github.com/allape/goview/model"
	"github.com/allape/goview/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var l = gogger.New("main")

func main() {
	err := gogger.InitFromEnv()
	if err != nil {
		l.Error().Fatalln(err)
	}

	db, err := gorm.Open(mysql.Open(env.DatabaseDSN), &gorm.Config{
		Logger: logger.New(gogger.New("db").Debug(), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Info,
			Colorful:      true,
		}),
	})
	if err != nil {
		l.Error().Fatalln(err)
	}

	err = db.AutoMigrate(&model.Datasource{}, &model.Preview{})
	if err != nil {
		l.Error().Fatalf("Failed to auto migrate database: %v", err)
	}

	engine := gin.Default()

	if env.EnableCors {
		config := cors.DefaultConfig()
		config.AddAllowHeaders("Authorization")
		config.AllowAllOrigins = true
		engine.Use(cors.New(config))
	}

	err = gocrud.NewSingleHTMLServe(engine.Group("/ui"), env.UIFolder, &gocrud.SingleHTMLServeConfig{
		AllowReplace: true,
	})
	if err != nil {
		l.Error().Fatalf("Failed to setup ui controller: %v", err)
	}

	apiGroup := engine.Group("api")

	err = controller.SetupDatasourceController(apiGroup.Group("datasource"), db)
	if err != nil {
		l.Error().Fatalf("Failed to setup datasource controller: %v", err)
	}

	err = controller.SetupPreviewController(apiGroup.Group("preview"), db)
	if err != nil {
		l.Error().Fatalf("Failed to setup preview controller: %v", err)
	}

	go func() {
		err := engine.Run(env.BindAddr)
		if err != nil {
			l.Error().Fatalf("Failed to start http server: %v", err)
		}
	}()

	util.Wait4CtrlC()
}
