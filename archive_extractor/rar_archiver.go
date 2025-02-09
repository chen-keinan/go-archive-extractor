package archive_extractor

import (
	"context"
	"github.com/jfrog/go-archive-extractor/archive_extractor/archiver_errors"
	"github.com/mholt/archives"
	"os"
	"strings"
)

type RarArchiver struct {
	MaxCompressRatio   int64
	MaxNumberOfEntries int
}

func (ra RarArchiver) ExtractArchive(path string,
	processingFunc func(*ArchiveHeader, map[string]interface{}) error, params map[string]interface{}) error {
	ctx := context.Background()
	maxBytesLimit, err := maxBytesLimit(path, ra.MaxCompressRatio)
	if err != nil {
		return archiver_errors.New(err)
	}
	provider := LimitAggregatingReadCloserProvider{
		Limit: maxBytesLimit,
	}
	format := archives.Rar{}
	rarFile, err := os.Open(path)
	if err != nil {
		return archiver_errors.NewOpenError(path, err)
	}
	defer func() {
		_ = rarFile.Close()
	}()
	err = extract(ctx, format, rarFile, ra.MaxNumberOfEntries, provider, processingFunc, params)
	if err != nil && strings.Contains(err.Error(), archiver_errors.RarDecodeError.Error()) {
		return archiver_errors.NewOpenError(path, err)
	}
	return err
}
