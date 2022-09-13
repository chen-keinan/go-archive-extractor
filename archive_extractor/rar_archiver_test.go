package archive_extractor

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRarArchiver(t *testing.T) {
	ra := &RarArchiver{}
	funcParams := params()
	err := ra.ExtractArchive("./fixtures/test.rar", processingFunc, funcParams)
	require.NoError(t, err)
	ad, ok := funcParams["archiveData"].(*ArchiveData)
	assert.True(t, ok)
	assert.Equal(t, ad.Name, "Interactive travel sample/media/Maldives.jpg")
	assert.Equal(t, ad.ModTime, int64(1454061977))
	assert.Equal(t, ad.IsFolder, false)
	assert.Equal(t, ad.Size, int64(695028))
}

func TestRarArchiver_NoRarFile(t *testing.T) {
	// zip file with .rar extension (changed manually)
	ra := &RarArchiver{}
	funcParams := params()
	err := ra.ExtractArchive("./fixtures/notRarFile.rar", processingFunc, funcParams)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "rardecode: RAR signature not found")
}

func TestRarArchiver_ExtractArchive(t *testing.T) {
	ra := &RarArchiver{}
	funcParams := params()
	err := ra.ExtractArchive("./fixtures/testwithcontent.rar", processingReadingFunc, funcParams)
	assert.NoError(t, err)
	assert.Equal(t, int64(4410), funcParams["read"])
}

func TestRarArchiver_LimitRatioReached(t *testing.T) {
	ra := &RarArchiver{
		MaxCompressRatio: 3,
	}
	funcParams := params()
	err := ra.ExtractArchive("./fixtures/testwithcontent.rar", processingReadingFunc, funcParams)
	assert.True(t, IsErrCompressLimitReached(err))
}

func TestRarArchiver_LimitRatioNotReached(t *testing.T) {
	ra := &RarArchiver{
		MaxCompressRatio: 4,
	}
	funcParams := params()
	err := ra.ExtractArchive("./fixtures/testwithcontent.rar", processingReadingFunc, funcParams)
	assert.NoError(t, err)
}

func TestRarArchiver_MaxNumberOfEntriesNotReached(t *testing.T) {
	ra := &RarArchiver{MaxCompressRatio: 1, MaxNumberOfEntries: 100}
	funcParams := params()
	err := ra.ExtractArchive("./fixtures/testwithmanyfiles.rar", processingReadingFunc, funcParams)
	assert.NoError(t, err)
}

func TestRarArchiver_MaxNumberOfEntriesReached(t *testing.T) {
	ra := &RarArchiver{MaxCompressRatio: 1, MaxNumberOfEntries: 99}
	funcParams := params()
	err := ra.ExtractArchive("./fixtures/testwithmanyfiles.rar", processingReadingFunc, funcParams)
	assert.EqualError(t, err, ErrTooManyEntries.Error())
}

func TestRarArchiver_AggregationCauseRatioLimitError(t *testing.T) {
	ra := &RarArchiver{
		MaxCompressRatio: 2,
	}
	funcParams := params()
	err := ra.ExtractArchive("./fixtures/testmanylarge.rar", processingReadingFunc, funcParams)
	assert.True(t, IsErrCompressLimitReached(err))
}
