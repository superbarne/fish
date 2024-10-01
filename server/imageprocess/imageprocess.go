package imageprocess

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
)

// ProcessImage remove white background from image. TargetPath need a .png extension!
func ProcessImage(srcPath string, targetPath string, log *slog.Logger) error {
	src, err := gg.LoadImage(srcPath)
	if err != nil {
		log.Error("Failed to load image", slog.String("error", err.Error()))
		return err
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

	// save image
	filePath := filepath.Dir(targetPath)
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		log.Error("Failed to create folder", slog.String("error", err.Error()))
		return err
	}

	err = gg.SavePNG(targetPath, im.Image())
	if err != nil {
		log.Error("Failed to save image", slog.String("error", err.Error()))
		return err
	}

	return nil
}
