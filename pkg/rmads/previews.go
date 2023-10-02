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
	"github.com/tartale/go/pkg/images"
	"github.com/tartale/go/pkg/logz"
	"github.com/tartale/go/pkg/mathx"
	"github.com/tartale/remove-ads/pkg/config"
	"golang.org/x/image/colornames"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const previewBlockDuration = 1 * time.Minute

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
		err = createImageFile(previewsDir, previewFilename, previewImage)
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

func createImageFile(dir, filename string, img images.DrawableImage) error {

	path := filepath.Join(dir, filename)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}

func createPreviewImages(videoPath string) ([]images.DrawableImage, error) {

	thumbnailImages, err := createThumbnailImages(videoPath)
	if err != nil {
		return nil, err
	}
	addTimestampLabels(thumbnailImages)

	// break the preview images into 1-minute blocks
	oneMinuteMills := int((previewBlockDuration).Milliseconds())
	thumbnailIntervalMillis := int(config.Values.StillFramesInterval.Milliseconds())
	imagesPerBlock := oneMinuteMills / thumbnailIntervalMillis
	thumbnailCount := len(thumbnailImages)
	var previewImages []images.DrawableImage
	for i := 0; i < thumbnailCount; {
		start := i
		end := mathx.Min(thumbnailCount-1, start+imagesPerBlock)
		previewImages = append(previewImages, stitchImages(thumbnailImages[start:end]))
		i += imagesPerBlock
	}

	markers, err := ImportSkipFile(config.Values.SkipFilePath)
	if err != nil {
		return nil, err
	}
	addMarkers(previewImages, markers)

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

func addTimestampLabels(thumbnailImgs []images.DrawableImage) {

	timeOffset := time.Duration(0)

	for _, img := range thumbnailImgs {
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

func stitchImages(imgs []images.DrawableImage) images.DrawableImage {

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

func addMarkers(previewImgs []images.DrawableImage, markers *Markers) {

	var (
		lineThickness = 20
		lineColor     = colornames.Red
	)
	logger := logz.Logger()
	for _, segment := range markers.Segments {
		logger.Debugf("segment times: %s to %s\n", segment.StartOffset.String(), segment.EndOffset.String())

		startImgIndex, startX, endImgIndex, endX := getSegmentOffsets(segment, previewImgs)
		if startImgIndex == endImgIndex {
			img := previewImgs[startImgIndex]
			height := img.Bounds().Max.Y
			rect := image.Rect(startX, 0, endX, height)
			images.DrawRectangle(img, rect, lineThickness, lineColor)
			// createImageFile(config.Values.TempDir, "current-image.png", img)
			continue
		}

		for i := startImgIndex; i <= endImgIndex; i++ {
			img := previewImgs[i]
			width := img.Bounds().Max.X
			height := img.Bounds().Max.Y
			if i == startImgIndex {
				images.DrawFullVerticalLine(img, startX, lineThickness, lineColor)
				images.DrawHorizontalLineBetween(img, startX, width, 0, lineThickness, lineColor)
				images.DrawHorizontalLineBetween(img, startX, width, height, lineThickness, lineColor)
			} else if i == endImgIndex {
				images.DrawFullVerticalLine(img, endX, lineThickness, lineColor)
				images.DrawHorizontalLineBetween(img, 0, endX, 0, lineThickness, lineColor)
				images.DrawHorizontalLineBetween(img, 0, endX, height, lineThickness, lineColor)
			} else {
				images.DrawFullHorizontalLine(img, 0, lineThickness, lineColor)
				images.DrawFullHorizontalLine(img, height, lineThickness, lineColor)
			}
			// createImageFile(config.Values.TempDir, "current-image.png", img)
		}

		logger.Debugf("added markers: startIndex: %d; startX: %d; endIndex: %d; endX: %d\n", startImgIndex, startX, endImgIndex, endX)
	}

	// for _, ts := range markers.Timestamps {
	// 	logger.Debugf("timestamp: %s\n", ts.Timestamp.String())

	// 	imgIndex, x := getOffsets(ts.Timestamp, previewImgs)
	// 	img := previewImgs[imgIndex]
	// 	images.DrawFullVerticalLine(img, mathx.Floor(x), lineThickness, colornames.Green)
	// }
}

func getSegmentOffsets(segment Segment, imgs []images.DrawableImage) (startImgIndex, startX, endImgIndex, endX int) {

	var startFloat, endFloat float64
	startImgIndex, startFloat = getOffsets(segment.StartOffset, imgs)
	endImgIndex, endFloat = getOffsets(segment.EndOffset, imgs)
	startX = mathx.Ceil(startFloat)
	endX = mathx.Floor(endFloat)

	return
}

func getOffsets(offset time.Duration, imgs []images.DrawableImage) (imgIndex int, x float64) {

	offsetTime := time.Time{}.Add(offset)
	offsetSeconds := offsetTime.Second()
	offsetMinutes := offsetTime.Minute()
	imgIndex = mathx.Min(offsetMinutes, len(imgs)-1)
	img := imgs[imgIndex]
	imgWidth := img.Bounds().Dx()
	offsetPercent := float64(offsetSeconds) / 60.0
	x = float64(imgWidth) * offsetPercent

	return
}
