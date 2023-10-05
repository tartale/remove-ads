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

	test.CheckFilesExist(t, test.TransportStreamPath)

	testTransportStreamFilename := path.Base(test.TransportStreamPath)
	expectedFfmpegCmd := fmt.Sprintf(`%s -y -hwaccel auto -i %s -vf fps=1/5 %s/%s-%%04d.png`,
		config.Values.FFmpegFilePath, test.TransportStreamPath, config.Values.TempDir, testTransportStreamFilename)
	ffmpegCmd, err := makeGenerateStillFramesCmd(test.TransportStreamPath)

	assert.Nil(t, err)
	assert.Equal(t, expectedFfmpegCmd, ffmpegCmd.String())
}
