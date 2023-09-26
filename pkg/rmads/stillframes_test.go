package rmads

import (
	"fmt"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tartale/remove-ads/pkg/config"
	"github.com/tartale/remove-ads/test"
)

// ffmpeg -i input.mp4 -vf fps=1/4 %04d.png
func TestMakeThumbnailsCmd(t *testing.T) {

	_, _, testTransportStreamPath := test.GetTestFiles()
	test.CheckFilesExist(t, testTransportStreamPath)

	testTransportStreamFilename := path.Base(testTransportStreamPath)
	expectedFfmpegCmd := fmt.Sprintf(`%s -y -hwaccel auto -i %s -vf fps=1/5 %s/%s-%%04d.png`,
		config.Values.FFmpegFilePath, testTransportStreamPath, config.Values.TempDir, testTransportStreamFilename)
	ffmpegCmd, err := makeGenerateStillFramesCmd(testTransportStreamPath)

	assert.Nil(t, err)
	assert.Equal(t, expectedFfmpegCmd, ffmpegCmd.String())
}
