package rmads

import (
	"context"
	"os/exec"
	"time"

	"github.com/elgs/gojq"
	"github.com/tartale/go/pkg/generics"
	"github.com/tartale/go/pkg/logz"
)

func GetShowDuration(ctx context.Context, inputFilePath string) (time.Duration, error) {

	var logger = logz.Logger()
	ffprobePath, err := exec.LookPath("ffprobe")
	if err != nil {
		return 0, err
	}

	// ffprobe -v error -show_entries format=duration -of json ./example.mp4 | jq -r '.format.duration'
	ffprobeCommand := exec.Command(ffprobePath, "-v", "error", "-show_entries", "format=duration", "-of", "json", inputFilePath)
	logger.Debugf("ffprobe command: %s\n", ffprobeCommand.String())

	output, err := ffprobeCommand.Output()
	if err != nil {
		return 0, err
	}
	parser, err := gojq.NewStringQuery(string(output))
	if err != nil {
		return 0, err
	}
	durationObj, err := parser.Query("format.duration")
	if err != nil {
		return 0, err
	}
	durationStr, err := generics.CastTo[string](durationObj)
	if err != nil {
		return 0, err
	}
	duration := *durationStr + "s"
	result, err := time.ParseDuration(duration)
	if err != nil {
		return 0, err
	}

	return result, nil
}
