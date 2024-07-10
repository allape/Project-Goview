package main

import (
	"database/sql"
	"github.com/allape/goenv"
	"github.com/allape/goview/datasource"
	"github.com/allape/goview/env"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const TAG = "[main]"

var rout *gin.Engine
var repo *gorm.DB

func SetupRepo() (*gorm.DB, error) {
	log.Println(TAG, "opening database")
	db, err := sql.Open(
		"mysql",
		goenv.Getenv(
			env.DatabaseURL,
			"root:Root_123456@(127.0.0.1:3306)/goview?charset=utf8mb4&parseTime=true",
		),
	)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	log.Println(TAG, "pinging database")
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	rep, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}))
	if err != nil {
		return nil, err
	}

	return rep, nil
}

func main() {
	var err error

	repo, err = SetupRepo()
	if err != nil {
		log.Fatal(TAG, " setup repo: ", err)
	}

	rout = gin.New()
	rout.Use(cors.Default())

	indexHTML := goenv.Getenv(env.UIIndexHTML, "ui/dist/index.html")
	rout.StaticFile("/", indexHTML)
	rout.StaticFile("/index", indexHTML)
	rout.StaticFile("/index.htm", indexHTML)
	rout.StaticFile("/index.html", indexHTML)

	// region logic

	err = datasource.Setup(repo, rout, goenv.Getenv(env.PreviewFolder, "preview"))
	if err != nil {
		log.Fatal(TAG, " setup datasource: ", err)
	}

	// endregion

	go func() {
		err = rout.Run(goenv.Getenv(env.HttpBinding, ":8080"))
		if err != nil {
			log.Fatal(TAG, " run: ", err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	log.Println(TAG, "started")
	sig := <-sigs
	log.Println(TAG, "exiting with", sig)
}
