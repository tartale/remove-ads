package rmads

import (
	"image"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tartale/go/pkg/images"
	"github.com/tartale/go/pkg/slicez"
)

var testImg = &image.RGBA{
	Rect: image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: 480, Y: 270},
	},
}

func TestGetSegmentOffsets(t *testing.T) {

	var drawableImage images.DrawableImage = testImg
	testImgs := slicez.Fill(drawableImage, 60)
	testSegment := Segment{
		StartOffset: 0,
		EndOffset:   10 * time.Second,
	}
	startIndex, startX, endIndex, endX := getSegmentOffsets(testSegment, testImgs)
	assert.Equal(t, 0, startIndex)
	assert.Equal(t, 0, startX)
	assert.Equal(t, 0, endIndex)
	assert.Equal(t, 80, endX)

	testSegment = Segment{
		StartOffset: 20 * time.Second,
		EndOffset:   50 * time.Second,
	}
	startIndex, startX, endIndex, endX = getSegmentOffsets(testSegment, testImgs)
	assert.Equal(t, 0, startIndex)
	assert.Equal(t, 160, startX)
	assert.Equal(t, 0, endIndex)
	assert.Equal(t, 400, endX)
}
