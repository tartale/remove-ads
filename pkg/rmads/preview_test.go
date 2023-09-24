package rmads

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tartale/remove-ads/pkg/config"
	"github.com/tartale/remove-ads/test"
)

// ffmpeg -i input.mp4 -vf fps=1/4 %04d.png

func TestMakeCreatePreviewCmd(t *testing.T) {

	_, _, testTransportStreamPath := test.GetTestFiles()
	test.CheckFilesExist(t, testTransportStreamPath)

	expectedFfmpegCmd := fmt.Sprintf(`%s -y -i %s -vf fps=1/4 %s/%%04d.png`,
		config.Values.FFmpegFilePath, testTransportStreamPath, config.Values.TempDir)
	ffmpegCmd, err := makeCreatePreviewCommand(testTransportStreamPath)

	assert.Nil(t, err)
	assert.Equal(t, expectedFfmpegCmd, ffmpegCmd.String())
}
