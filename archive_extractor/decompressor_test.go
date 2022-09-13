package archive_extractor

import (
	"fmt"
	"github.com/jfrog/go-archive-extractor/archive_extractor/archiver_errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDecompressor_ExtractArchive_CompressedFile(t *testing.T) {
	dc := &Decompressor{}
	funcParams := params()
	var testCases = []struct {
		Name             string
		FilePath         string
		ExpectedName     string
		ExpectedModTime  int64
		ExpectedIsFolder bool
		ExpectedSize     int64
	}{
		{
			Name:             "Test xz compression",
			FilePath:         "./fixtures/test.txt.xz",
			ExpectedName:     "test.txt",
			ExpectedModTime:  1661433804,
			ExpectedIsFolder: false,
			ExpectedSize:     64,
		},
		{
			Name:             "Test bzip2 compression",
			FilePath:         "./fixtures/test.txt.bz2",
			ExpectedName:     "test.txt",
			ExpectedModTime:  1661837894,
			ExpectedIsFolder: false,
			ExpectedSize:     43,
		},
		{
			Name:             "Test gzip compression",
			FilePath:         "./fixtures/test.txt.gz",
			ExpectedName:     "test.txt",
			ExpectedModTime:  1661837894,
			ExpectedIsFolder: false,
			ExpectedSize:     36,
		},
		{
			Name:             "Test lzma compression",
			FilePath:         "./fixtures/test.txt.lzma",
			ExpectedName:     "test.txt",
			ExpectedModTime:  1661837894,
			ExpectedIsFolder: false,
			ExpectedSize:     30,
		},
		{
			Name:             "Test lzw compression",
			FilePath:         "./fixtures/test.txt.Z",
			ExpectedName:     "test.txt",
			ExpectedModTime:  1661434675,
			ExpectedIsFolder: false,
			ExpectedSize:     11,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := dc.ExtractArchive(tc.FilePath, processingFunc, funcParams)
			require.NoError(t, err)
			ad, ok := funcParams["archiveData"].(*ArchiveData)
			assert.True(t, ok)
			assert.Equal(t, tc.ExpectedName, ad.Name)
			assert.Equal(t, tc.ExpectedModTime, ad.ModTime)
			assert.Equal(t, tc.ExpectedIsFolder, ad.IsFolder)
			assert.Equal(t, tc.ExpectedSize, ad.Size)
		})
	}
}

func TestDecompressor_ExtractArchive_NotCompressedFile(t *testing.T) {
	dc := &Decompressor{}
	funcParams := params()
	filePath := "./fixtures/test.txt"
	expectedErr := archiver_errors.New(fmt.Errorf(NotCompressedOrNotSupportedError, filePath))
	err := dc.ExtractArchive(filePath, processingFunc, funcParams)
	assert.EqualError(t, err, expectedErr.Error())
}

func TestExtractArchive_MaxRatioReached_ShouldReturnError(t *testing.T) {
	dc := &Decompressor{
		MaxCompressRatio: 2,
	}
	funcParams := params()
	err := dc.ExtractArchive("./fixtures/testsinglelarge.txt.xz", processingReadingFunc, funcParams)
	assert.True(t, IsErrCompressLimitReached(err))
}

func TestExtractArchive_MaxRatioNotReached(t *testing.T) {
	dc := &Decompressor{
		MaxCompressRatio: 100,
	}
	funcParams := params()
	err := dc.ExtractArchive("./fixtures/testsinglelarge.txt.xz", processingReadingFunc, funcParams)
	assert.NoError(t, err)
}
