package rmads

import (
	"errors"
	"fmt"
	"os/exec"
	"path"

	"github.com/tartale/go/pkg/command"
	"github.com/tartale/remove-ads/pkg/config"
)

func makeGenerateStillFramesCmd(videoPath string) (*exec.Cmd, error) {

	fps := fmt.Sprintf("fps=1/%d", int(config.Values.StillFramesInterval.Seconds()))
	stillFramesDir, _, filePattern := getPreviewPaths(videoPath)
	stillFramesFilePattern := path.Join(stillFramesDir, filePattern)
	ffmpegCmd := exec.Command(config.Values.FFmpegFilePath, "-y", "-hwaccel", "auto", "-i", videoPath,
		"-vf", fps, stillFramesFilePattern)

	return ffmpegCmd, nil
}

func generateStillFrames(videoPath string) error {

	ffmpegCmd, err := makeGenerateStillFramesCmd(videoPath)
	if err != nil {
		return err
	}
	err = command.RunIf(ffmpegCmd)
	if err != nil && !errors.Is(err, command.ErrDryRun) {
		return err
	}

	return nil
}
