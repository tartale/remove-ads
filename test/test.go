package test

import (
	"fmt"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tartale/go/pkg/errorz"
	"github.com/tartale/go/pkg/filez"
)

var (
	UtTivoMetadataPath  = path.Join(GetTestDir(), "data", string("unitTestTivoMetadata.json"))
	TivoMetadataPath    = path.Join(GetTestDir(), "data", string("intTestTivoMetadata.json"))
	TransportStreamPath = path.Join(GetTestDir(), "data", string("intTestTransportStream.ts"))
)

func GetTestDir() string {

	rootDir, found := filez.GetRootDirForCaller(1)
	if !found {
		panic(fmt.Errorf("%w: could not get test directory", errorz.ErrFatal))
	}

	return path.Join(rootDir, "test")
}

func CheckFilesExist(t *testing.T, paths ...string) {

	missingFiles := filez.Exist(paths...)
	if len(missingFiles) > 0 {
		assert.Fail(t, fmt.Sprintf("%s: test data file(s) '%s'", errorz.ErrNotFound.Error(), strings.Join(missingFiles, ",")))
	}
}
