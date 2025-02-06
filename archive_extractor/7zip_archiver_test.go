//go:build tests_group_all

package archive_extractor

import (
	"fmt"
	"github.com/jfrog/go-archive-extractor/archive_extractor/archiver_errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test7ZipAndRarArchiver(t *testing.T) {
	za := &SevenZipArchiver{}
	funcParams := params()
	if err := za.ExtractArchive("./fixtures/test.7z", processingFunc, funcParams); err != nil {
		fmt.Print(err.Error())
		t.Fatal(err)
	}
	ad := funcParams["archiveData"].(*ArchiveData)
	assert.Equal(t, ad.Name, "Interactive travel sample/.spxproperties")
	assert.Equal(t, ad.ModTime, int64(6802270473))
	assert.Equal(t, ad.IsFolder, false)
	assert.Equal(t, ad.Size, int64(44))
}

func Test7ZipAndRarArchiverReadAll(t *testing.T) {
	za := &SevenZipArchiver{}
	funcParams := params()
	err := za.ExtractArchive("./fixtures/testwithcontent.7z", processingReadingFunc, funcParams)
	assert.NoError(t, err)
	assert.Equal(t, int64(4410), funcParams["read"])
}

func TestRarArchiver_NoSevenZipFile(t *testing.T) {
	// zip file with .rar extension (changed manually)
	sz := &SevenZipArchiver{}
	funcParams := params()
	err := sz.ExtractArchive("./fixtures/notRarFile.rar", processingFunc, funcParams)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), archiver_errors.SevenZipDecodeError.Error())
}

func Test7ZipAndRarArchiverLimitRatio(t *testing.T) {
	za := &SevenZipArchiver{
		MaxCompressRatio: 3,
	}
	funcParams := params()
	err := za.ExtractArchive("./fixtures/testwithcontent.7z", processingReadingFunc, funcParams)
	assert.True(t, IsErrCompressLimitReached(err))
}

func Test7ZipAndRarArchiverLimitRatioHighEnough(t *testing.T) {
	za := &SevenZipArchiver{
		MaxCompressRatio: 4,
	}
	funcParams := params()
	err := za.ExtractArchive("./fixtures/testwithcontent.7z", processingReadingFunc, funcParams)
	assert.NoError(t, err)
}

func Test7ZipAndRarArchiverLimitNumberOfRecords(t *testing.T) {
	za := &SevenZipArchiver{
		MaxNumberOfEntries: 1,
	}
	funcParams := params()
	err := za.ExtractArchive("./fixtures/testwithmultipleentries.7z", processingReadingFunc, funcParams)
	assert.Contains(t, err.Error(), ErrTooManyEntries.Error())
}

func Test7ZipAndRarArchiverLimitRatioAggregationCauseError(t *testing.T) {
	za := &SevenZipArchiver{
		MaxCompressRatio: 20,
	}
	funcParams := params()
	err := za.ExtractArchive("./fixtures/testwithmultiplelargeentries.7z", processingReadingFunc, funcParams)
	assert.True(t, IsErrCompressLimitReached(err))
}
