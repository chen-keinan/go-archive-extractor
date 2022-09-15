package archive_extractor

import (
	"fmt"
	"github.com/jfrog/go-archive-extractor/archive_extractor/archiver_errors"
	"github.com/jfrog/go-archive-extractor/compression"
	"os"
	"path/filepath"
	"strings"
)

type Decompressor struct {
	MaxCompressRatio int64
}

const (
	NotCompressedOrNotSupportedError = "file %v is not compressed or the compression method is not supported"
)

func (dc Decompressor) ExtractArchive(path string,
	processingFunc func(*ArchiveHeader, map[string]interface{}) error, params map[string]interface{}) error {
	maxBytesLimit, err := maxBytesLimit(path, dc.MaxCompressRatio)
	if err != nil {
		return archiver_errors.New(err)
	}
	provider := LimitAggregatingReadCloserProvider{
		Limit: maxBytesLimit,
	}
	cReader, isCompressed, err := compression.NewReader(path)
	if err != nil {
		return archiver_errors.New(err)
	}
	defer cReader.Close()
	if !isCompressed {
		return archiver_errors.New(fmt.Errorf(NotCompressedOrNotSupportedError, path))
	}
	limitingReader := provider.CreateLimitAggregatingReadCloser(cReader)
	defer limitingReader.Close()
	f, err := os.Open(path)
	if err != nil {
		return archiver_errors.New(err)
	}
	defer f.Close()
	fInfo, err := f.Stat()
	if err != nil {
		return archiver_errors.New(err)
	}
	// removing the compression extension since now we have a decompressed file
	name := strings.TrimSuffix(fInfo.Name(), filepath.Ext(fInfo.Name()))
	archiveHeader := NewArchiveHeader(limitingReader, name, fInfo.ModTime().Unix(), fInfo.Size())
	err = processingFunc(archiveHeader, params)
	if err != nil {
		return err
	}
	return nil
}
