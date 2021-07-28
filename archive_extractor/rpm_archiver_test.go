package archive_extractor

import (
	"fmt"
	"github.com/jfrog/go-archive-extractor/archive_extractor/archiver_errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRpmArchiver(t *testing.T) {
	za := &RpmArchiver{}
	funcParams := params()
	if err := za.ExtractArchive("./fixtures/test.rpm", processingFunc, funcParams); err != nil {
		fmt.Print(err.Error())
		t.Fatal(err)
	}
	ad := funcParams["archiveData"].(*ArchiveData)
	assert.Equal(t, ad.Name, "./usr/share/doc/php-zstd-devel/tests/info.phpt")
	assert.Equal(t, ad.ModTime, int64(1517299253))
	assert.Equal(t, ad.IsFolder, false)
	assert.Equal(t, ad.Size, int64(183))
	rpmPkg := funcParams["rpmPkg"].(*RpmPkg)
	assert.Equal(t, rpmPkg.Release, "1.fc24.remi.7.0")
	assert.Equal(t, rpmPkg.Version, "0.4.11")
	assert.Equal(t, rpmPkg.Name, "php-zstd-devel")
	assert.Equal(t, rpmPkg.ModularityLabel, "")
}

func TestRpmArchiverTooManyEntries(t *testing.T) {
	za := &RpmArchiver{
		MaxNumberOfEntries: 1,
	}
	err := za.ExtractArchive("./fixtures/test.rpm", processingFunc, params())
	assert.EqualError(t, err, archiver_errors.New(ErrTooManyEntries).Error())
}

func TestRpmArchiverTooManyEntriesNotReached(t *testing.T) {
	za := &RpmArchiver{
		MaxNumberOfEntries: 100,
	}
	err := za.ExtractArchive("./fixtures/test.rpm", processingFunc, params())
	assert.NoError(t, err)
}

func TestRpmArchiverRatioOk(t *testing.T) {
	za := &RpmArchiver{
		MaxCompressRatio: 1,
	}
	err := za.ExtractArchive("./fixtures/test.rpm", processingFunc, params())
	assert.NoError(t, err)
}
