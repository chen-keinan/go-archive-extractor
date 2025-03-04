package archive_extractor

import (
	"context"
	"errors"

	"github.com/jfrog/go-archive-extractor/archive_extractor/archiver_errors"
	"github.com/mholt/archives"
)

type TarArchiver struct {
	MaxCompressRatio   int64
	MaxNumberOfEntries int
}

var ErrUnsupportedArchiveType = errors.New("unsupported archive type")

func (ta TarArchiver) ExtractArchive(path string, processingFunc func(*ArchiveHeader, map[string]interface{}) error, params map[string]interface{}) error {
	ctx := context.Background()
	maxBytesLimit, err := maxBytesLimit(path, ta.MaxCompressRatio)
	if err != nil {
		return archiver_errors.New(err)
	}
	provider := LimitAggregatingReadCloserProvider{
		Limit: maxBytesLimit,
	}
	format, _, err := archives.Identify(ctx, path, nil)
	if err != nil {
		return archiver_errors.New(err)
	}
	extractor, ok := format.(archives.Extractor)
	if !ok {
		return archiver_errors.New(ErrUnsupportedArchiveType)
	}
	return extractWithSymlinks(ctx, extractor, path, ta.MaxNumberOfEntries, provider, processingFunc, params)
}
