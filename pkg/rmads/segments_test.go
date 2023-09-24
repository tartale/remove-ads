package rmads

import (
	"context"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tartale/go/pkg/filez"
	"github.com/tartale/remove-ads/test"
)

func TestInvertSegments(t *testing.T) {

	inputSegments := Segments{
		Segment{Description: "Segment 1", StartOffset: 5 * time.Second, EndOffset: 10 * time.Second},
		Segment{Description: "Segment 2", StartOffset: 15 * time.Second, EndOffset: 20 * time.Second},
		Segment{Description: "Segment 3", StartOffset: 30 * time.Second, EndOffset: 40 * time.Second},
	}
	expectedSegments := Segments{
		Segment{StartOffset: 0 * time.Second, EndOffset: 5 * time.Second},
		Segment{StartOffset: 10 * time.Second, EndOffset: 15 * time.Second},
		Segment{StartOffset: 20 * time.Second, EndOffset: 30 * time.Second},
		Segment{StartOffset: 40 * time.Second, EndOffset: 60 * time.Second},
	}

	segments := inputSegments.Invert(60 * time.Second)

	assert.Equal(t, expectedSegments, segments)
}

func TestInvertSegments_NoInputSegments(t *testing.T) {

	inputSegments := Segments{}
	expectedSegments := Segments{Segment{StartOffset: 0, EndOffset: 60 * time.Second}}

	segments := inputSegments.Invert(60 * time.Second)

	assert.Equal(t, expectedSegments, segments)
}

func TestInvertSegments_InputStartsAtZero(t *testing.T) {

	inputSegments := Segments{
		Segment{Description: "Segment 1", StartOffset: 0, EndOffset: 10 * time.Second},
		Segment{Description: "Segment 2", StartOffset: 15 * time.Second, EndOffset: 20 * time.Second},
		Segment{Description: "Segment 3", StartOffset: 30 * time.Second, EndOffset: 40 * time.Second},
	}
	expectedSegments := Segments{
		Segment{StartOffset: 10 * time.Second, EndOffset: 15 * time.Second},
		Segment{StartOffset: 20 * time.Second, EndOffset: 30 * time.Second},
		Segment{StartOffset: 40 * time.Second, EndOffset: 60 * time.Second},
	}

	segments := inputSegments.Invert(60 * time.Second)

	assert.Equal(t, expectedSegments, segments)
}

func TestInvertSegments_InputEndsAtEndtime(t *testing.T) {

	inputSegments := Segments{
		Segment{Description: "Segment 1", StartOffset: 5 * time.Second, EndOffset: 10 * time.Second},
		Segment{Description: "Segment 2", StartOffset: 15 * time.Second, EndOffset: 20 * time.Second},
		Segment{Description: "Segment 3", StartOffset: 30 * time.Second, EndOffset: 60 * time.Second},
	}
	expectedSegments := Segments{
		Segment{StartOffset: 0 * time.Second, EndOffset: 5 * time.Second},
		Segment{StartOffset: 10 * time.Second, EndOffset: 15 * time.Second},
		Segment{StartOffset: 20 * time.Second, EndOffset: 30 * time.Second},
	}

	segments := inputSegments.Invert(60 * time.Second)

	assert.Equal(t, expectedSegments, segments)
}

func TestInvertSegments_StartAtZeroEndAtEndtime(t *testing.T) {

	inputSegments := Segments{
		Segment{Description: "Segment 1", StartOffset: 0, EndOffset: 10 * time.Second},
		Segment{Description: "Segment 2", StartOffset: 15 * time.Second, EndOffset: 20 * time.Second},
		Segment{Description: "Segment 3", StartOffset: 30 * time.Second, EndOffset: 60 * time.Second},
	}
	expectedSegments := Segments{
		Segment{StartOffset: 10 * time.Second, EndOffset: 15 * time.Second},
		Segment{StartOffset: 20 * time.Second, EndOffset: 30 * time.Second},
	}

	segments := inputSegments.Invert(60 * time.Second)

	assert.Equal(t, expectedSegments, segments)
}

func TestMakeRemoveSegmentsCmd(t *testing.T) {

	testMetadataPath, _, testTransportStreamPath := test.GetTestFiles()
	test.CheckFilesExist(t, testMetadataPath, testTransportStreamPath)

	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		assert.Fail(t, "failing test; ffmpeg required in path")
	}

	metadataInput := filez.MustReadAll(testMetadataPath)
	markers, err := ImportTivoClipMetadata(metadataInput)
	assert.Nil(t, err)

	expectedFfmpegCmd := fmt.Sprintf(`%s -y -i %s -vf select='between(t\,1\,20)+between(t\,30\,40)+between(t\,2\,19)+between(t\,23\,39),setpts=N/FRAME_RATE/TB' -af aselect='between(t\,1\,20)+between(t\,30\,40)+between(t\,2\,19)+between(t\,23\,39),asetpts=N/SR/TB' -`, ffmpegPath, testTransportStreamPath)
	ffmpegCmd, err := markers.Segments.makeRemoveCommand(testTransportStreamPath, "")

	assert.Nil(t, err)
	assert.Equal(t, expectedFfmpegCmd, ffmpegCmd.String())
}

func TestRemoveSegments(t *testing.T) {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		assert.Fail(t, "failing test; ffmpeg required in path")
	}

	ctx := context.Background()
	_, metadataPath, transportStreamPath := test.GetTestFiles()
	metadataInput := filez.MustReadAll(metadataPath)
	markers, err := ImportTivoClipMetadata(metadataInput)
	assert.Nil(t, err)

	markers.Segments.Remove(ctx, transportStreamPath, "")
}