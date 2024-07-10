package env

const (
	// TrustedCerts
	// Cert files separated by comma
	TrustedCerts = "GOVIEW_TRUSTED_CERTS"

	// UIIndexHTML
	// The index.html rendered in path of "/"
	UIIndexHTML = "GOVIEW_UI_INDEX_HTML"

	// PreviewFolder
	// Folder to store preview files
	PreviewFolder = "GOVIEW_PREVIEW_FOLDER"

	// HttpBinding
	// Port to run the server
	HttpBinding = "GOVIEW_HTTP_BINDING"

	// DatabaseURL
	// Database URL in format like "user:password@tcp(host:port)/database"
	DatabaseURL = "GOVIEW_DATABASE_URL"
)
