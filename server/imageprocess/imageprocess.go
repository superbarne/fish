package imageprocess

import (
	"log"

	"github.com/fogleman/gg"
)

// ProcessImage remove white background from image. TargetPath need a .png extension!
func ProcessImage(srcPath string, targetPath string) {
	src, err := gg.LoadImage(srcPath)
	if err != nil {
		log.Fatal(err)
	}

	w := src.Bounds().Size().X
	h := src.Bounds().Size().Y

	im := gg.NewContext(w, h)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := src.At(x, y)
			r, g, b, a := c.RGBA()
			if a == 0 {
				continue
			}

			whiteness := r + g + b
			if whiteness > 150000 {
				continue
			}

			im.SetColor(c)
			im.SetPixel(x, y)
		}
	}

	err = gg.SavePNG(targetPath, im.Image())
	if err != nil {
		log.Fatal(err)
	}
}
