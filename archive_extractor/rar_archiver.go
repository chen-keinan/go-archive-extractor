package archive_extractor

import (
	"github.com/jfrog/go-archive-extractor/archive_extractor/archiver_errors"
	"github.com/jfrog/go-archive-extractor/utils"
	"github.com/mholt/archiver/v3"
	"io"
)

type RarArchiver struct {
	MaxCompressRatio   int64
	MaxNumberOfEntries int
}

func (ra RarArchiver) ExtractArchive(path string,
	processingFunc func(*ArchiveHeader, map[string]interface{}) error, params map[string]interface{}) error {
	maxBytesLimit, err := maxBytesLimit(path, ra.MaxCompressRatio)
	if err != nil {
		return archiver_errors.New(err)
	}
	provider := LimitAggregatingReadCloserProvider{
		Limit: maxBytesLimit,
	}
	format := archiver.Rar{}
	err = format.OpenFile(path)
	defer format.Close()
	if err != nil {
		return archiver_errors.New(err)
	}
	entriesCount := 0
	for {
		if ra.MaxNumberOfEntries != 0 && entriesCount > ra.MaxNumberOfEntries {
			return ErrTooManyEntries
		}
		entriesCount++
		file, err := format.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return archiver_errors.New(err)
		}
		if !file.FileInfo.IsDir() && !utils.PlaceHolderFolder(file.FileInfo.Name()) {
			countingReadCloser := provider.CreateLimitAggregatingReadCloser(file.ReadCloser)
			archiveHeader := NewArchiveHeader(countingReadCloser, file.FileInfo.Name(), file.FileInfo.ModTime().Unix(), file.FileInfo.Size())
			err = processingFunc(archiveHeader, params)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
