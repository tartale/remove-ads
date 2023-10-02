package rmads

import (
	"fmt"
	"time"

	"github.com/tartale/go/pkg/errorz"
	"github.com/tartale/go/pkg/filez"
	"github.com/tartale/go/pkg/jsonx"
	"github.com/tartale/go/pkg/primitives"
)

func ImportSkipFile(skipFilePath string) (*Markers, error) {

	skipData := filez.MustReadAll(skipFilePath)
	if markers, err := ImportTivoClipMetadata(skipData); err == nil {
		return markers, nil
	}
	if markers, err := ImportVideoRedoV3Skip(skipData); err == nil {
		return markers, nil
	}

	return nil, fmt.Errorf("%w: Unable to import skip data", errorz.ErrFatal)
}

func ImportTivoClipMetadata(tivoClipMetadata []byte) (*Markers, error) {

	var markers Markers

	clipMetadata, err := jsonx.QueryToType[[]any]("clipMetadata", string(tivoClipMetadata))
	if err != nil {
		return nil, err
	}
	for _, cm := range *clipMetadata {
		segments, err := jsonx.QueryObjToType[[]any]("segment", cm)
		if err != nil {
			return nil, err
		}
		for _, segmentObj := range *segments {
			var (
				offsetStr *string
				offset    int
				segment   Segment
			)

			offsetStr, err = jsonx.QueryObjToType[string]("startOffset", segmentObj)
			if err != nil {
				return nil, fmt.Errorf("%w: field 'startOffset' is required", errorz.ErrInvalidArgument)
			}
			err = primitives.Parse(*offsetStr, &offset)
			if err != nil {
				return nil, err
			}
			segment.StartOffset = time.Duration(offset) * time.Millisecond

			offsetStr, err = jsonx.QueryObjToType[string]("endOffset", segmentObj)
			if err != nil {
				return nil, fmt.Errorf("%w: field 'endOffset' is required", errorz.ErrInvalidArgument)
			}
			err = primitives.Parse(*offsetStr, &offset)
			if err != nil {
				return nil, err
			}
			segment.EndOffset = time.Duration(offset) * time.Millisecond

			markers.Segments = append(markers.Segments, segment)
		}

		syncMarks, err := jsonx.QueryObjToType[[]any]("syncMark", cm)
		if err != nil {
			return nil, err
		}
		for _, syncMarkObj := range *syncMarks {
			var (
				timestampStr *string
				timstampVal  int
				timestamp    Timestamp
			)

			timestampStr, err = jsonx.QueryObjToType[string]("timestamp", syncMarkObj)
			if err != nil {
				return nil, fmt.Errorf("%w: field 'timestamp' is required", errorz.ErrInvalidArgument)
			}
			err = primitives.Parse(*timestampStr, &timstampVal)
			if err != nil {
				return nil, err
			}
			timestamp.Timestamp = time.Duration(timstampVal) * time.Millisecond
			markers.Timestamps = append(markers.Timestamps, timestamp)
		}
	}

	return &markers, nil
}

func ImportVideoRedoV3Skip(videoRedoData []byte) (*Markers, error) {

	return nil, errorz.ErrFatal
}
