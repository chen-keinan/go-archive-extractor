package archive_extractor

import (
	"errors"
	"io"
)

var ErrTooManyEntries = errors.New("too many entries in archive")
var ErrCompressLimitReached = errors.New("total bytes limit reached")

type LimitAggregatingReadCloserProvider struct {
	Total int64
	Limit int64
}

func (provider *LimitAggregatingReadCloserProvider) CreateLimitAggregatingReadCloser(rc io.Reader) LimitAggregatingReadCloser {
	return &limitAggregatingReadCloser{
		Reader: rc,
		Total:  &provider.Total,
		Limit:  provider.Limit,
	}
}

type LimitAggregatingReadCloser interface {
	Read(p []byte) (int, error)
	Close() error
}

type limitAggregatingReadCloser struct {
	Reader io.Reader
	Total  *int64
	Limit  int64
}

func (crc *limitAggregatingReadCloser) Read(p []byte) (int, error) {
	n, err := crc.Reader.Read(p)
	if err != nil && err != io.EOF {
		return n, err
	}
	*crc.Total += int64(n)
	if crc.Limit != 0 && *crc.Total > crc.Limit {
		return n, ErrCompressLimitReached
	}
	return n, err
}

func (crc *limitAggregatingReadCloser) Close() error {
	closer, ok := crc.Reader.(io.Closer)
	if ok {
		return closer.Close()
	}
	return nil
}
