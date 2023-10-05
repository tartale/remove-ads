package rmads

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tartale/go/pkg/errorz"
	"github.com/tartale/go/pkg/filez"
	"github.com/tartale/go/pkg/mathx"
	"github.com/tartale/go/pkg/primitives"
)

func ImportSkipFile(skipFilePath string, shift time.Duration) (*Markers, error) {

	skipData := filez.MustReadAll(skipFilePath)
	if markers, err := ImportTivoClipMetadata(skipData, shift); err == nil {
		return markers, nil
	}
	if markers, err := ImportVideoRedoV3Skip(skipData); err == nil {
		return markers, nil
	}

	return nil, fmt.Errorf("%w: Unable to import skip data", errorz.ErrFatal)
}

func ImportTivoClipMetadata(tivoClipMetadata []byte, shift time.Duration) (*Markers, error) {

	var (
		markers Markers
		tcm     TivoClipMetadata
	)

	err := json.Unmarshal(tivoClipMetadata, &tcm)
	if err != nil {
		return nil, err
	}
	if len(tcm.ClipMetadata) == 0 {
		return nil, fmt.Errorf("%w: at least one clip metadata object expected", errorz.ErrInvalidArgument)
	}
	// Just use the first clipMetadata until we can figure out how the others factor in
	clipMetadata := tcm.ClipMetadata[1]

	for _, segment := range clipMetadata.Segment {
		startOffsetMs := primitives.MustParseTo[int](segment.StartOffset)
		endOffsetMs := primitives.MustParseTo[int](segment.EndOffset)
		segment := Segment{
			StartOffset: time.Duration(startOffsetMs)*time.Millisecond + shift,
			EndOffset:   time.Duration(endOffsetMs)*time.Millisecond + shift,
		}
		segment.StartOffset = mathx.Max(segment.StartOffset, time.Duration(0))
		segment.EndOffset = mathx.Max(segment.EndOffset, time.Duration(0))
		markers.Segments = append(markers.Segments, segment)
	}

	return &markers, nil
}

func ImportVideoRedoV3Skip(videoRedoData []byte) (*Markers, error) {

	return nil, errorz.ErrFatal
}
