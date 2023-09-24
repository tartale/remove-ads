package rmads

import (
	"os/exec"
	"path"

	"github.com/tartale/go/pkg/filez"
	"github.com/tartale/go/pkg/logz"
	"github.com/tartale/remove-ads/pkg/config"
)

func CreatePreview(inputFilePath string) error {

	var logger = logz.Logger()
	ffmpegCmd, err := makeCreatePreviewCommand(inputFilePath)
	if err != nil {
		return err
	}
	logger.Debugf("ffmpeg command: %s\n", ffmpegCmd.String())

	filez.MustMkdirAll(config.Values.TempDir)

	return nil
}

func makeCreatePreviewCommand(inputFilePath string) (*exec.Cmd, error) {

	if inputFilePath == "" {
		inputFilePath = "-"
	}
	outputFilePattern := path.Join(config.Values.TempDir, "%04d.png")
	ffmpegCmd := exec.Command(config.Values.FFmpegFilePath, "-y", "-i", inputFilePath,
		"-vf", "fps=1/4", outputFilePattern)

	return ffmpegCmd, nil
}
