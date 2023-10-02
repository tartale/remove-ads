package rmads

import (
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/tartale/go/pkg/images"
	"golang.org/x/image/draw"
)

func createThumbnailImages(videoPath string) ([]images.DrawableImage, error) {

	stillFramesDir, _, _ := getPreviewPaths(videoPath)

	var thumbnails []images.DrawableImage
	err := filepath.WalkDir(stillFramesDir, func(imgFilename string, _ os.DirEntry, _ error) error {

		if !strings.HasSuffix(imgFilename, "png") {
			return nil
		}
		imgFile, err := os.OpenFile(imgFilename, os.O_RDWR, 0664)
		if err != nil {
			return err
		}
		defer imgFile.Close()

		// follows the technique described here:
		//   https://stackoverflow.com/a/67678654/1258206
		img, _, err := image.Decode(imgFile)
		if err != nil {
			return err
		}
		newWidth, newHeight := img.Bounds().Dx()/4, img.Bounds().Dy()/4
		resizedImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
		draw.NearestNeighbor.Scale(resizedImg, resizedImg.Rect, img, img.Bounds(), draw.Over, nil)
		thumbnails = append(thumbnails, resizedImg)

		return nil
	})

	return thumbnails, err
}
