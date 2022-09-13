package archive_extractor

import (
	"archive/tar"
	"github.com/jfrog/go-archive-extractor/archive_extractor/archiver_errors"
	"github.com/jfrog/go-archive-extractor/compression"
	"github.com/jfrog/go-archive-extractor/utils"
	"io"
)

type TarArchiver struct {
	MaxCompressRatio   int64
	MaxNumberOfEntries int
}

func (ta TarArchiver) ExtractArchive(path string,
	processingFunc func(*ArchiveHeader, map[string]interface{}) error, params map[string]interface{}) error {
	maxBytesLimit, err := maxBytesLimit(path, ta.MaxCompressRatio)
	if err != nil {
		return err
	}
	provider := LimitAggregatingReadCloserProvider{
		Limit: maxBytesLimit,
	}
	cReader, _, err := compression.NewReader(path)
	if compression.IsGetReaderError(err) {
		return archiver_errors.New(err)
	}
	if err != nil {
		return err
	}
	limitingReader := provider.CreateLimitAggregatingReadCloser(cReader)
	defer limitingReader.Close()
	rc := tar.NewReader(limitingReader)
	entriesCount := 0
	for {
		if ta.MaxNumberOfEntries != 0 && entriesCount > ta.MaxNumberOfEntries {
			return ErrTooManyEntries
		}
		entriesCount++
		archiveEntry, err := rc.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if !archiveEntry.FileInfo().IsDir() && !utils.PlaceHolderFolder(archiveEntry.FileInfo().Name()) {
			archiveHeader := NewArchiveHeader(rc, archiveEntry.Name, archiveEntry.ModTime.Unix(), archiveEntry.FileInfo().Size())
			err = processingFunc(archiveHeader, params)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
