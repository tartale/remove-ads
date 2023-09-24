package rmads

import (
	"context"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tartale/remove-ads/test"
)

func TestGetShowDuration(t *testing.T) {

	_, _, transportStreamPath := test.GetTestFiles()
	test.CheckFilesExist(t, transportStreamPath)
	_, err := exec.LookPath("ffprobe")
	if err != nil {
		assert.Fail(t, "failing test; ffprobe required in path")
	}

	ctx := context.Background()
	showDuration, err := GetShowDuration(ctx, transportStreamPath)

	assert.Nil(t, err)
	assert.Equal(t, "31m0.996611s", showDuration.String())
}
