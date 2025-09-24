package imagemapper

import (
	"bytes"
	"image/color"
	"log"
	"math/rand"
	"path/filepath"

	"github.com/fogleman/gg"
)

func RenderMessageIntoImage(message string) (*bytes.Buffer, error) {
	fontPath := filepath.Join("fonts", "OpenSans-ExtraBold.ttf")
	dc := gg.NewContext(1000, 1000) // canvas 1000px by 1000px
	if err := dc.LoadFontFace(fontPath, 80); err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	padding := 10.0
	usableWidth := float64(dc.Width()) - 2*padding
	usableHeight := float64(dc.Height()) - 2*padding
	cx := padding + usableWidth/2
	cy := padding + usableHeight/2

	dc.SetColor(color.White)
	dc.DrawRectangle(0, 0, 1000, 1000)
	dc.Fill()
	dc.SetColor(color.Black)

	rotation := rand.Float64() * 3.1415
	dc.RotateAbout(rotation, cx, cy)
	dc.DrawStringWrapped(message, cx, cy, 0.5, 0.5, usableWidth, 1.5, gg.AlignCenter)

	buf := new(bytes.Buffer)
	if err := dc.EncodePNG(buf); err != nil {
		return nil, err
	}

	return buf, nil
}

func RenderMessageIntoImageWithBackgroundImage(message, strFilePath string) (*bytes.Buffer, error) {
	fontPath := filepath.Join("fonts", "OpenSans-Regular.ttf")
	dc := gg.NewContext(1000, 1000) // canvas 1000px by 1000px
	if err := dc.LoadFontFace(fontPath, 80); err != nil {
		log.Printf(err.Error())
		return nil, err
	}

	padding := 10.0
	usableWidth := float64(dc.Width()) - 2*padding
	usableHeight := float64(dc.Height()) - 2*padding

	cx := padding + usableWidth/2
	cy := padding + usableHeight/2
	im, err := gg.LoadPNG(strFilePath) //backgrounds/wp12782247.png
	if err != nil {
		return nil, err
	}
	dc.SetColor(color.White)
	dc.DrawImageAnchored(im, im.Bounds().Dx(), im.Bounds().Dy(), 1, 1)
	dc.RotateAbout(3.1415, cx, cy)
	dc.DrawStringWrapped(message, cx, cy, 0.5, 0.5, usableWidth, 1.5, gg.AlignCenter)
	buf := new(bytes.Buffer)
	if err := dc.EncodePNG(buf); err != nil {
		return nil, err
	}

	return buf, nil
}
