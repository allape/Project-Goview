package env

import "github.com/allape/goenv"

const (
	trustedCerts  = "GOVIEW_TRUSTED_CERTS"
	uiIndexHTML   = "GOVIEW_UI_INDEX_HTML"
	previewFolder = "GOVIEW_PREVIEW_FOLDER"
	bindAddr      = "GOVIEW_BIND_ADDR"
	enbaleCors    = "GOVIEW_ENABLE_CORS"
	databaseDSN   = "GOVIEW_DATABASE_DSN"
)

var (
	TrustedCerts  = goenv.Getenv(trustedCerts, "")
	UIIndexHTML   = goenv.Getenv(uiIndexHTML, "./ui/dist/index.html")
	PreviewFolder = goenv.Getenv(previewFolder, "./preview")
	BindAddr      = goenv.Getenv(bindAddr, ":8080")
	EnableCors    = goenv.Getenv(enbaleCors, true)
	DatabaseDSN   = goenv.Getenv(databaseDSN, "root:Root_123456@tcp(localhost:3306)/goview?charset=utf8mb4&parseTime=True&loc=Local")
)
