package rmads

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"github.com/tartale/go/pkg/filez"
	"github.com/tartale/go/pkg/mathx"
	"github.com/tartale/remove-ads/pkg/config"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func CreatePreviews(videoPath string) error {

	err := generateStillFrames(videoPath)
	if err != nil {
		return err
	}
	previewImages, err := createPreviewImages(videoPath)
	if err != nil {
		return err
	}
	_, previewsDir, filenamePattern := getPreviewPaths(videoPath)
	for i, previewImage := range previewImages {
		previewFilename := fmt.Sprintf(filenamePattern, i+1)
		previewPath := filepath.Join(previewsDir, previewFilename)
		previewFile, err := os.Create(previewPath)
		if err != nil {
			return err
		}
		defer previewFile.Close()
		err = png.Encode(previewFile, previewImage)
		if err != nil {
			return err
		}
	}

	return nil
}

func DisplayPreview(videoPath string) {

	previewApp := app.New()
	previewWindow := previewApp.NewWindow("Preview")

	_, previewsDir, _ := getPreviewPaths(videoPath)
	img := canvas.NewImageFromFile(filepath.Join(previewsDir, "intTestTransportStream-000001.png"))

	content := container.NewVBox(img)
	previewWindow.SetContent(content)
	previewWindow.ShowAndRun()
}

func createPreviewImages(videoPath string) ([]image.Image, error) {

	thumbnailImages, err := createThumbnailImages(videoPath)
	if err != nil {
		return nil, err
	}
	addTimestampLabels(thumbnailImages)

	// break the preview images into 1-minute blocks
	oneMinuteMills := int((1 * time.Minute).Milliseconds())
	previewIntervalMillis := int(config.Values.StillFramesInterval.Milliseconds())
	imagesPerBlock := oneMinuteMills / previewIntervalMillis
	thumbnailCount := len(thumbnailImages)
	var previewImages []image.Image
	for i := 0; i < thumbnailCount; {
		start := i
		end := mathx.Min(thumbnailCount-1, start+imagesPerBlock)
		previewImages = append(previewImages, stitchImages(thumbnailImages[start:end]))
		i += imagesPerBlock
	}

	return previewImages, nil
}

func getPreviewPaths(videoPath string) (stillFramesDir, previewsDir, filePattern string) {

	base := filez.NameWithoutExtension(videoPath)
	stillFramesDir = path.Join(config.Values.TempDir, base, "stillframes")
	previewsDir = path.Join(config.Values.TempDir, base, "previews")
	filePattern = fmt.Sprintf("%s-%%06d.png", base)

	filez.MustMkdirAll(stillFramesDir)
	filez.MustMkdirAll(previewsDir)

	return
}

func addTimestampLabels(imgs []*image.RGBA) {

	timeOffset := time.Duration(0)

	for _, img := range imgs {
		imgWidth, imgHeight := img.Bounds().Dx(), img.Bounds().Dy()
		labelX := imgWidth - (imgWidth / 5)
		labelY := imgHeight - (imgHeight / 8)

		col := color.White
		point := fixed.Point26_6{X: fixed.I(labelX), Y: fixed.I(labelY)}
		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(col),
			Face: basicfont.Face7x13,
			Dot:  point,
		}
		formattedTimeOffset := time.Unix(0, 0).UTC().Add(time.Duration(timeOffset)).
			Format(time.TimeOnly)
		d.DrawString(formattedTimeOffset)
		timeOffset += config.Values.StillFramesInterval
	}
}

func stitchImages(imgs []*image.RGBA) *image.RGBA {

	if len(imgs) == 0 {
		return nil
	}

	// get the width/height of the first image
	imgWidth, imgHeight := imgs[0].Bounds().Dx(), imgs[0].Bounds().Dy()
	// create the stitched image's background
	stitchedImg := image.NewRGBA(image.Rect(0, 0, imgWidth*len(imgs), imgHeight))
	// set the background color
	draw.Draw(stitchedImg, stitchedImg.Bounds(), &image.Uniform{color.Black}, image.Point{0, 0}, draw.Src)

	for i, img := range imgs {
		//set image offset
		offset := image.Pt(i*imgWidth, 0)
		//combine the image
		draw.Draw(stitchedImg, img.Bounds().Add(offset), img, image.Point{0, 0}, draw.Over)
	}

	return stitchedImg
}
