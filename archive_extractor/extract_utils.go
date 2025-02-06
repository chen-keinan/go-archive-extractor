package archive_extractor

import (
	"context"
	"github.com/jfrog/go-archive-extractor/archive_extractor/archiver_errors"
	"github.com/jfrog/go-archive-extractor/utils"
	"github.com/mholt/archives"
	"io"
)

type processingArchiveFunc func(*ArchiveHeader, map[string]interface{}) error

func extract(ctx context.Context, ex archives.Extractor, arcReader io.Reader, MaxNumberOfEntries int, provider LimitAggregatingReadCloserProvider, processingFunc processingArchiveFunc, params map[string]any) error {
	entriesCount := 0
	var multiErrors *archiver_errors.MultiError
	err := ex.Extract(ctx, arcReader, func(ctx context.Context, fileInfo archives.FileInfo) error {
		if MaxNumberOfEntries != 0 && entriesCount >= MaxNumberOfEntries {
			return ErrTooManyEntries
		}
		entriesCount++
		file, err := fileInfo.Open()
		defer func() {
			if file != nil {
				_ = file.Close()
			}
		}()
		if err != nil {
			multiErrors = archiver_errors.Append(multiErrors, archiver_errors.NewArchiverExtractorError(fileInfo.NameInArchive, err))
		} else if !fileInfo.IsDir() && !utils.PlaceHolderFolder(fileInfo.Name()) {
			countingReadCloser := provider.CreateLimitAggregatingReadCloser(file)
			archiveHeader := NewArchiveHeader(countingReadCloser, fileInfo.NameInArchive, fileInfo.ModTime().Unix(), fileInfo.Size())
			processingError := processingFunc(archiveHeader, params)
			if processingError != nil {
				return processingError
			}
		}
		return nil
	})
	//multi error can be skipped or not skipped by caller, therefore we distinguish between err and multiErrors
	if err == nil && multiErrors != nil {
		return multiErrors
	}
	return err
}
