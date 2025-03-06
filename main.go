package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/davidbyttow/govips/v2/vips"
)

// w - image output width
// h - image output height
// q - image output quality
// d - image rotation degrees
// ext - image output extension
// wmt - watermark text
// wmpos - watermark position
// wmclr- watermark color
// wmop - watermark opacity
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		widthQuery := query.Get("w")
		heightQuery := query.Get("h")
		qualityQuery := query.Get("q")
		degreesQuery := query.Get("d")
		// extension := query.Get("ext")
		quality := 100
		if qualityQuery != "" {
			q, err := strconv.Atoi(qualityQuery)
			if err == nil {
				quality = q
			}
		}

		pathToStorage := strings.ReplaceAll(r.URL.Path, "content/cluster", "mnt/gluster")
		img, err := vips.NewImageFromFile(fmt.Sprintf(".%s", pathToStorage))
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		img.Rotate(getImgRotateAngleFromQuery(degreesQuery))

		width, height := getImgWidthAndHeightFromQuery(img.Width(), img.Height(), widthQuery, heightQuery)
		// scaleFactor := math.Min(float64(width)/float64(img.Width()), float64(height)/float64(img.Height()))
		// img.Resize(scaleFactor, vips.KernelAuto)
		// img.ResizeWithVScale(float64(width)/float64(img.Width()), float64(height)/float64(img.Height()), vips.KernelAuto)
		if img.Width() > img.Height() {
			img.Thumbnail(width, height, vips.InterestingNone)
		} else {
			img.Thumbnail(height, width, vips.InterestingNone)
		}
		if img.Bands() < 3 {
			_ = img.ToColorSpace(vips.InterpretationSRGB)
		}
		if img.Bands() == 4 {
			img.ExtractBandToImage(3, 1)
			_ = img.ExtractBand(0, 3)
		}
		// vips.LabelParams{}
		img.Label(&vips.LabelParams{
			Text: "etagi",
			Font: "sans 12",
			Width: vips.Scalar{
				Value:    0.5,
				Relative: true,
			},
			Height: vips.Scalar{
				Value:    0.5,
				Relative: true,
			},
			OffsetX: vips.Scalar{
				Value:    0,
				Relative: true,
			},
			OffsetY: vips.Scalar{
				Value:    0,
				Relative: true,
			},
			Color: vips.Color{
				R: 0,
				G: 0,
				B: 0,
			},
			Alignment: vips.AlignLow,
		})
		b, _, _ := img.Export(&vips.ExportParams{
			Format:  vips.ImageTypeJPEG,
			Quality: quality,
		})
		w.WriteHeader(200)
		w.Write(b)
		// w.Write([]byte(fmt.Sprintf("path: %s, width: %s, height: %s, quality: %s, extension: %s", r.URL.Path, width, height, quality, extension)))
	})

	http.ListenAndServe(":8080", mux)

	// img, err := vips.NewImageFromFile("")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// img.Label(&vips.LabelParams{
	// 	Text:      "",
	// 	Font:      "",
	// 	Width:     vips.Scalar{},
	// 	Height:    vips.Scalar{},
	// 	OffsetX:   vips.Scalar{},
	// 	OffsetY:   vips.Scalar{},
	// 	Opacity:   0,
	// 	Color:     vips.Color{},
	// 	Alignment: 0,
	// })
}

func getImgRotateAngleFromQuery(degreesQuery string) vips.Angle {
	d := strings.TrimSpace(degreesQuery)
	switch d {
	case "90":
		return vips.Angle90
	case "180":
		return vips.Angle180
	case "270":
		return vips.Angle270
	default:
		return vips.Angle0
	}
}

func getImgWidthAndHeightFromQuery(width, height int, widthQuery, heightQuery string) (int, int) {
	if widthQuery != "" && heightQuery != "" {
		w, err := strconv.Atoi(widthQuery)
		if err != nil {
			return width, height
		}

		h, err := strconv.Atoi(heightQuery)
		if err != nil {
			return width, height
		}

		return w, h
	}

	return width, height
}
