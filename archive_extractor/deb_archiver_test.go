package archive_extractor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDebArchiver(t *testing.T) {
	za := &DebArchiver{}
	funcParams := params()
	if err := za.ExtractArchive("./fixtures/test.deb", processingFunc, funcParams); err != nil {
		fmt.Print(err.Error())
		t.Fatal(err)
	}
	ad := funcParams["archiveData"].(*ArchiveData)
	assert.Equal(t, ad.Name, "data.tar.xz")
	assert.Equal(t, ad.ModTime, int64(1485714631))
	assert.Equal(t, ad.IsFolder, false)
	assert.Equal(t, ad.Size, int64(42284))
}

func TestDebArchiverLimitNumberOfEntries(t *testing.T) {
	za := &DebArchiver{
		MaxNumberOfEntries: 1,
	}
	err := za.ExtractArchive("./fixtures/test.deb", processingReadingFunc, params())
	assert.EqualError(t, err, ErrTooManyEntries.Error())
}

func TestDebArchiverLimitNumberOfEntriesNotReached(t *testing.T) {
	za := &DebArchiver{
		MaxNumberOfEntries: 10,
	}
	err := za.ExtractArchive("./fixtures/test.deb", processingReadingFunc, params())
	assert.NoError(t, err)
}

func TestDebArchiverMaxRatioNotReached(t *testing.T) {
	za := &DebArchiver{
		MaxCompressRatio: 1,
	}
	err := za.ExtractArchive("./fixtures/test.deb", processingReadingFunc, params())
	assert.NoError(t, err)
}
