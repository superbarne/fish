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

	// create map
	heightMap := make([][]bool, h)
	for y := 0; y < h; y++ {
		heightMap[y] = make([]bool, w)

		for x := 0; x < w; x++ {
			c := src.At(x, y)
			heightMap[y][x] = isWhite(c.RGBA())
		}
	}

	pixels := findConnectedPixels(heightMap, 1, 1)

	im := gg.NewContext(w, h)

	// draw all non-white pixel
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if pixels[y][x] {
				continue
			}

			c := src.At(x, y)

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

func isWhite(r, g, b, a uint32) bool {
	if a == 0 {
		return true
	}

	whiteness := r + g + b
	return whiteness > 150000
}

// dfs performing the depth search
func dfs(matrix [][]bool, visited [][]bool, x, y int) {
	directions := [8][2]int{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1},
	}

	visited[x][y] = true

	for _, dir := range directions {
		newX, newY := x+dir[0], y+dir[1]

		if newX >= 0 && newX < len(matrix) && newY >= 0 && newY < len(matrix[0]) && !visited[newX][newY] && matrix[newX][newY] {
			dfs(matrix, visited, newX, newY)
		}
	}
}

func findConnectedPixels(matrix [][]bool, startX, startY int) [][]bool {
	rows := len(matrix)
	cols := len(matrix[0])

	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	if matrix[startX][startY] {
		dfs(matrix, visited, startX, startY)
	}

	return visited
}
