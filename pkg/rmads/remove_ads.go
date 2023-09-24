package rmads

import (
	"context"

	"github.com/tartale/remove-ads/pkg/config"
)

func RemoveAds(ctx context.Context) error {

	markers, err := ImportSkipFile(config.Values.SkipFilePath)
	if err != nil {
		return err
	}
	err = markers.Segments.Remove(ctx, config.Values.InputFilePath, config.Values.OutputFilePath)
	if err != nil {
		return err
	}

	return nil
}