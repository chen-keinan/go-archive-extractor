//go:build tests_group_all

package archive_extractor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGzMetadataArchiver(t *testing.T) {
	ga := GzMetadataArchiver{}
	funcParams := params()
	err := ga.ExtractArchive("./fixtures/test.gz", processingReadingFunc, funcParams)
	assert.NoError(t, err)
	assert.Equal(t, int64(6), funcParams["read"])
}

func TestGzMetadataArchiverWithRatioOk(t *testing.T) {
	ga := GzMetadataArchiver{MaxCompressRatio: 1}
	funcParams := params()
	err := ga.ExtractArchive("./fixtures/test.gz", processingReadingFunc, funcParams)
	assert.NoError(t, err)
	assert.Equal(t, int64(6), funcParams["read"])
}

func TestGzMetadataArchiverWithRatio(t *testing.T) {
	ga := GzMetadataArchiver{MaxCompressRatio: 2}
	funcParams := params()
	err := ga.ExtractArchive("./fixtures/testwithcontent.gz", processingReadingFunc, funcParams)
	assert.True(t, IsErrCompressLimitReached(err))
}

func TestGzMetadataArchiverWithReasonableRatio(t *testing.T) {
	ga := GzMetadataArchiver{MaxCompressRatio: 3}
	funcParams := params()
	err := ga.ExtractArchive("./fixtures/testwithcontent.gz", processingReadingFunc, funcParams)
	assert.True(t, IsErrCompressLimitReached(err))
}
