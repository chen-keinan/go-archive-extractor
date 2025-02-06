//go:build tests_group_all

package archive_extractor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"strings"
)

func TestTarUnexpectedEofArchiver(t *testing.T) {
	za := &TarArchiver{}
	funcParams := params()
	if err := za.ExtractArchive("./fixtures/test.deb", processingFunc, funcParams); err != nil {
		fmt.Print(err.Error() + "\n")
		assert.Equal(t, "archive/tar: invalid tar header", strings.Trim(err.Error(), ""))
	}
}

func TestTarArchiver(t *testing.T) {
	za := &TarArchiver{}
	funcParams := params()
	if err := za.ExtractArchive("./fixtures/test.tar.gz", processingFunc, funcParams); err != nil {
		fmt.Print(err.Error())
		t.Fatal(err)
	}
	ad := funcParams["archiveData"].(*ArchiveData)
	assert.Equal(t, ad.Name, "logRotator-1.0/log_rotator.go")
	assert.Equal(t, ad.ModTime, int64(1531307652))
	assert.Equal(t, ad.IsFolder, false)
	assert.Equal(t, ad.Size, int64(3685))
}

func TestTarArchiver_Lzma(t *testing.T) {
	za := &TarArchiver{}
	funcParams := params()
	if err := za.ExtractArchive("./fixtures/junit.tar.lzma", processingFunc, funcParams); err != nil {
		fmt.Print(err.Error())
		t.Fatal(err)
	}
	ad := funcParams["archiveData"].(*ArchiveData)
	assert.Equal(t, "junit-4.12.jar", ad.Name)
	assert.Equal(t, int64(1534397548), ad.ModTime)
	assert.False(t, ad.IsFolder)
	assert.Equal(t, int64(314932), ad.Size)
}

func TestTarArchiverMaxRatio(t *testing.T) {
	za := &TarArchiver{
		MaxCompressRatio: 2,
	}
	funcParams := params()
	err := za.ExtractArchive("./fixtures/testsinglelarge.tar.gz", processingReadingFunc, funcParams)
	assert.True(t, IsErrCompressLimitReached(err))
}

func TestTarArchiverMaxRatioNotReached(t *testing.T) {
	za := &TarArchiver{
		MaxCompressRatio: 100,
	}
	funcParams := params()
	err := za.ExtractArchive("./fixtures/testsinglelarge.tar.gz", processingReadingFunc, funcParams)
	assert.NoError(t, err)
}

func TestTarArchiverMaxEntriesReached(t *testing.T) {
	za := &TarArchiver{
		MaxNumberOfEntries: 12,
	}
	funcParams := params()
	err := za.ExtractArchive("./fixtures/testmanylarge.tar.gz", processingReadingFunc, funcParams)
	assert.EqualError(t, err, ErrTooManyEntries.Error())
}

func TestTarArchiverMaxEntriesNotReached(t *testing.T) {
	za := &TarArchiver{
		MaxNumberOfEntries: 20,
	}
	funcParams := params()
	err := za.ExtractArchive("./fixtures/testmanylarge.tar.gz", processingReadingFunc, funcParams)
	assert.NoError(t, err)
}

func TestTarArchiverAggregationCauseRatioLimitError(t *testing.T) {
	za := &TarArchiver{
		MaxCompressRatio: 4,
	}
	funcParams := params()
	err := za.ExtractArchive("./fixtures/testmanylarge.tar.gz", processingReadingFunc, funcParams)
	assert.True(t, IsErrCompressLimitReached(err))
}
