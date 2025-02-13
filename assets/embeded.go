package assets

import _ "embed"

var MIMEType = "image/jpeg"

//go:embed i_v_404.jpg
var IV404 []byte

//go:embed i_v_500.jpg
var IV500 []byte

//go:embed i_v_no_preview.jpg
var IVNoPreview []byte

//go:embed favicon.png
var Favicon []byte
