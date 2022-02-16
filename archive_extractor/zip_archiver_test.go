package archive_extractor

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"strings"
	"testing"
)

func TestZipUnexpectedEofArchiver(t *testing.T) {
	za := &ZipArchiver{}
	funcParams := params()
	if err := za.ExtractArchive("./fixtures/test.deb", processingFunc, funcParams); err != nil {
		fmt.Print(err.Error() + "\n")
		assert.Equal(t, "No zip file found", strings.Trim(err.Error(), ""))
	}
}

func TestZipArchiver(t *testing.T) {
	za := &ZipArchiver{}
	funcParams := params()
	if err := za.ExtractArchive("./fixtures/test.zip", processingFunc, funcParams); err != nil {
		fmt.Print(err.Error())
		t.Fatal(err)
	}
	ad := funcParams["archiveData"].(*ArchiveData)
	assert.Equal(t, ad.Name, "test.txt")
	assert.Equal(t, ad.ModTime, int64(1534147868))
	assert.Equal(t, ad.IsFolder, false)
	assert.Equal(t, ad.Size, int64(0))
}

func TestZipArchiverReadAll(t *testing.T) {
	za := &ZipArchiver{}
	funcParams := params()
	err := za.ExtractArchive("./fixtures/test.zip", processingReadingFunc, funcParams)
	assert.NoError(t, err)
	assert.Zero(t, funcParams["read"])
}

func TestZipArchiverReadAllWithEntry(t *testing.T) {
	za := &ZipArchiver{MaxCompressRatio: 1}
	funcParams := params()
	err := za.ExtractArchive("./fixtures/testwithcontent.zip", processingReadingFunc, funcParams)
	assert.NoError(t, err)
	assert.Equal(t, int64(13), funcParams["read"])
}

func TestZipArchiverReadAllWithEntryMaxNumberOfEntriesOk(t *testing.T) {
	za := &ZipArchiver{MaxCompressRatio: 1, MaxNumberOfEntries: 100}
	funcParams := params()
	err := za.ExtractArchive("./fixtures/testwithmanyfiles.zip", processingReadingFunc, funcParams)
	assert.NoError(t, err)
}

func TestZipArchiverReadAllWithEntryMaxNumberOfEntriesTooHigh(t *testing.T) {
	za := &ZipArchiver{MaxCompressRatio: 1, MaxNumberOfEntries: 99}
	funcParams := params()
	err := za.ExtractArchive("./fixtures/testwithmanyfiles.zip", processingReadingFunc, funcParams)
	assert.EqualError(t, err, ErrTooManyEntries.Error())
}

func TestZipArchiverRatioAndMaxEntriesNotSet(t *testing.T) {
	za := &ZipArchiver{}
	funcParams := params()
	err := za.ExtractArchive("./fixtures/testwithcontent.zip", processingReadingFunc, funcParams)
	assert.NoError(t, err)
	assert.Equal(t, int64(13), funcParams["read"])
}

func TestZipArchiverRatioNotSet(t *testing.T) {
	za := &ZipArchiver{MaxNumberOfEntries: 1000}
	funcParams := params()
	err := za.ExtractArchive("./fixtures/testwithcontent.zip", processingReadingFunc, funcParams)
	assert.NoError(t, err)
	assert.Equal(t, int64(13), funcParams["read"])
}

func TestZipArchiverAggregationCauseError(t *testing.T) {
	za := &ZipArchiver{
		MaxCompressRatio: 1,
	}
	funcParams := params()
	err := za.ExtractArchive("./fixtures/testmanyfileswithcontent.zip", processingReadingFunc, funcParams)
	assert.True(t, IsErrCompressLimitReached(err))
}

func TestZipArchiverSingleFileRatioCauseError(t *testing.T) {
	za := &ZipArchiver{
		MaxCompressRatio: 1,
	}
	funcParams := params()
	err := za.ExtractArchive("./fixtures/testwithsinglelargefile.zip", processingReadingFunc, funcParams)
	assert.True(t, IsErrCompressLimitReached(err))
}

func TestZipArchiver_PrependedZip(t *testing.T) {
	prependedZipPath := "./fixtures/appendedZip"
	za := &ZipArchiver{}
	funcParams := params()
	_, err := zip.OpenReader(prependedZipPath)
	assert.True(t, errors.Is(err, zip.ErrFormat))
	err = za.ExtractArchive(prependedZipPath, processingFunc, funcParams)
	assert.NoError(t, err)
}

func TestZipArchiver_EmptyZip(t *testing.T) {
	appendedZipPath := "./fixtures/prepended.empty"
	za := &ZipArchiver{}
	funcParams := params()
	err := za.ExtractArchive(appendedZipPath, processingFunc, funcParams)
	assert.NoError(t, err)
}

func TestZipArchiver_PrependedEmptyZip(t *testing.T) {
	prependedEmptyZipPath := "./fixtures/prepended.empty"
	za := &ZipArchiver{}
	funcParams := params()
	err := za.ExtractArchive(prependedEmptyZipPath, processingFunc, funcParams)
	assert.NoError(t, err)
}

func TestZipArchiver_initZipReader_signatureAtBufStart(t *testing.T) {
	zipPath := "./fixtures/test.zip"
	f, err := os.Open(zipPath)
	require.NoError(t, err)
	defer f.Close()
	r := bufio.NewReader(f)
	fileBuf, err := io.ReadAll(r)
	require.NoError(t, err)
	for i := 0; i < 100; i++ {
		startBut := make([]byte, i)
		prependedBuf := append(startBut, fileBuf...)
		_, err := initZipReader(bytes.NewReader(prependedBuf), int64(len(startBut)+len(fileBuf)))
		assert.NoError(t, err)
	}
}

func TestZipArchiver_initZipReader_signatureAtBufEnd(t *testing.T) {
	zipPath := "./fixtures/test.zip"
	f, err := os.Open(zipPath)
	require.NoError(t, err)
	defer f.Close()
	r := bufio.NewReader(f)
	fileBuf, err := io.ReadAll(r)
	require.NoError(t, err)
	for i := 4092; i < 4096; i++ {
		startBut := make([]byte, i)
		prependedBuf := append(startBut, fileBuf...)
		_, err := initZipReader(bytes.NewReader(prependedBuf), int64(len(startBut)+len(fileBuf)))
		assert.NoError(t, err)
	}
}

func TestZipArchiver_initZipReader_dummySignaturesBeforeFile(t *testing.T) {
	zipPath := "./fixtures/test.zip"
	f, err := os.Open(zipPath)
	require.NoError(t, err)
	defer f.Close()
	r := bufio.NewReader(f)
	fileBuf, err := io.ReadAll(r)
	require.NoError(t, err)
	signatureBuf := []byte("ABPK\x03\x04CD")
	for i := 0; i < 10; i++ {
		startBut := make([]byte, 700)
		var prependedBuf []byte
		for j := 0; j < i; j++ {
			prependedBuf = append(prependedBuf, startBut...)
			prependedBuf = append(prependedBuf, signatureBuf...)
		}
		prependedBuf = append(prependedBuf, fileBuf...)
		_, err := initZipReader(bytes.NewReader(prependedBuf), int64(len(prependedBuf)))
		assert.NoError(t, err)
	}
}
