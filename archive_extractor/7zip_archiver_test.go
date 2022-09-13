package archive_extractor

import (
	"fmt"
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
	assert.Equal(t, ad.ModTime, int64(-11644473600))
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
	assert.EqualError(t, err, ErrTooManyEntries.Error())
}

func Test7ZipAndRarArchiverLimitRatioAggregationCauseError(t *testing.T) {
	za := &SevenZipArchiver{
		MaxCompressRatio: 20,
	}
	funcParams := params()
	err := za.ExtractArchive("./fixtures/testwithmultiplelargeentries.7z", processingReadingFunc, funcParams)
	assert.True(t, IsErrCompressLimitReached(err))
}
