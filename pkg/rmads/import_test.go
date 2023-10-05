package rmads

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tartale/go/pkg/filez"
	"github.com/tartale/remove-ads/test"
)

func TestImportTivoSkip(t *testing.T) {

	testMetadataPath := test.TivoMetadataPath
	test.CheckFilesExist(t, testMetadataPath)
	testMetadataInput := filez.MustReadAll(testMetadataPath)

	markers, err := ImportTivoClipMetadata(testMetadataInput, time.Duration(0))
	assert.Nil(t, err)
	assert.NotNil(t, markers)
	assert.Len(t, markers.Segments, 4)
	assert.Equal(t, time.Duration(1)*time.Second, markers.Segments[0].StartOffset)
	assert.Equal(t, time.Duration(20)*time.Second, markers.Segments[0].EndOffset)
	assert.Equal(t, time.Duration(30)*time.Second, markers.Segments[1].StartOffset)
	assert.Equal(t, time.Duration(40)*time.Second, markers.Segments[1].EndOffset)
	assert.Equal(t, time.Duration(2)*time.Second, markers.Segments[2].StartOffset)
	assert.Equal(t, time.Duration(19)*time.Second, markers.Segments[2].EndOffset)
	assert.Equal(t, time.Duration(23)*time.Second, markers.Segments[3].StartOffset)
	assert.Equal(t, time.Duration(39)*time.Second, markers.Segments[3].EndOffset)
}
