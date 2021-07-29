package archive_extractor

import (
	"archive/zip"
	"fmt"
	"github.com/jfrog/go-archive-extractor/archive_extractor/archiver_errors"
)

type ZipArchiver struct {
	MaxCompressRatio   int64
	MaxNumberOfEntries int
}

func (za ZipArchiver) ExtractArchive(path string,
	processingFunc func(*ArchiveHeader, map[string]interface{}) error, params map[string]interface{}) error {
	maxBytesLimit, err := maxBytesLimit(path, za.MaxCompressRatio)
	if err != nil {
		return err
	}
	rcProvider := LimitAggregatingReadCloserProvider{Limit: maxBytesLimit}
	r, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer r.Close()
	var multiArchiveErr error
	if za.MaxNumberOfEntries > 0 && len(r.File) > za.MaxNumberOfEntries {
		return ErrTooManyEntries
	}
	for _, archiveEntry := range r.File {
		rc, err := archiveEntry.Open()
		if err != nil {
			if rc != nil {
				rc.Close()
			}
			multiArchiveErr = archiver_errors.Append(multiArchiveErr, fmt.Errorf("failed to open %s: %v", path, err))
			continue
		}
		countingReadCloser := rcProvider.CreateLimitAggregatingReadCloser(rc)
		archiveHeader := NewArchiveHeader(countingReadCloser, archiveEntry.Name, archiveEntry.ModTime().Unix(), archiveEntry.FileInfo().Size())
		err = processingFunc(archiveHeader, params)
		if err != nil {
			if rc != nil {
				rc.Close()
			}
			return err
		}
		rc.Close()
	}
	if multiArchiveErr != nil {
		return archiver_errors.New(multiArchiveErr)
	}
	return nil
}
