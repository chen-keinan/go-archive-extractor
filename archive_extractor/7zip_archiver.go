package archive_extractor

import (
	"context"
	"github.com/jfrog/go-archive-extractor/archive_extractor/archiver_errors"
	"github.com/mholt/archives"
	"os"
	"strings"
)

type SevenZipArchiver struct {
	MaxCompressRatio   int64
	MaxNumberOfEntries int
}

func (sa SevenZipArchiver) ExtractArchive(path string,
	processingFunc func(*ArchiveHeader, map[string]interface{}) error, params map[string]interface{}) error {
	ctx := context.Background()
	maxBytesLimit, err := maxBytesLimit(path, sa.MaxCompressRatio)
	if err != nil {
		return err
	}
	provider := LimitAggregatingReadCloserProvider{
		Limit: maxBytesLimit,
	}
	format := archives.SevenZip{}
	archFile, err := os.Open(path)
	if err != nil {
		return archiver_errors.NewOpenError(path, err)
	}
	defer func() {
		_ = archFile.Close()
	}()

	err = extract(ctx, format, archFile, sa.MaxNumberOfEntries, provider, processingFunc, params)
	if err != nil && strings.Contains(err.Error(), archiver_errors.SevenZipDecodeError.Error()) {
		return archiver_errors.NewOpenError(path, err)
	}
	return err
}
