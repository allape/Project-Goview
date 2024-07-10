package assets

import (
	_ "embed"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"image/color"
	"image/jpeg"
	"os"
)

//go:embed Roboto-Regular.ttf
var FontBytes []byte

var Font *truetype.Font

func CreateImage(
	width, height int,
	filename, text string,
	fontSize float64,
) error {
	dc := gg.NewContext(width, height)
	dc.SetColor(color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 255})
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Fill()

	if Font == nil {
		var err error
		Font, err = truetype.Parse(FontBytes)
		if err != nil {
			return err
		}
	}
	dc.SetFontFace(truetype.NewFace(Font, &truetype.Options{Size: fontSize}))

	dc.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	dc.DrawStringAnchored(text, float64(width/2), float64(height/2), 0.5, 0.5)

	//if timestamp {
	//	nowStr := time.Now().Format(time.DateTime)
	//	dc.SetFontFace(truetype.NewFace(Font, &truetype.Options{Size: 32}))
	//	dc.DrawStringAnchored(nowStr, float64(width-50), float64(height-50), 1, 0)
	//}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	img := dc.Image()
	err = jpeg.Encode(file, img, &jpeg.Options{Quality: 80})
	if err != nil {
		return err
	}

	return nil
}
