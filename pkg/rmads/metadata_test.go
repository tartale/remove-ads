package rmads

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tartale/remove-ads/test"
)

func TestGetShowDuration(t *testing.T) {

	transportStreamPath := test.TransportStreamPath
	test.CheckFilesExist(t, transportStreamPath)

	ctx := context.Background()

	showDuration, err := GetShowDuration(ctx, transportStreamPath)

	assert.Nil(t, err)
	assert.Equal(t, "31m0.996611s", showDuration.String())
}
