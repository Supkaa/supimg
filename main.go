package main

import (
	"github.com/davidbyttow/govips/v2/vips"
	"log"
)

func main() {
	img, err := vips.NewImageFromFile("")
	if err != nil {
		log.Fatal(err)
	}

	img.Label(&vips.LabelParams{
		Text:      "",
		Font:      "",
		Width:     vips.Scalar{},
		Height:    vips.Scalar{},
		OffsetX:   vips.Scalar{},
		OffsetY:   vips.Scalar{},
		Opacity:   0,
		Color:     vips.Color{},
		Alignment: 0,
	})
}
