package archive_extractor

import (
	"errors"
	"fmt"
	"io"
)

var ErrTooManyEntries = errors.New("too many entries in archive")

type ErrCompressLimitReached struct {
	SizeLimit int64
	CurrSize  int64
}

func newErrCompressLimitReached(sizeLimit, total int64) *ErrCompressLimitReached {
	return &ErrCompressLimitReached{SizeLimit: sizeLimit, CurrSize: total}
}

func IsErrCompressLimitReached(err error) bool {
	var errCompressLimitReached *ErrCompressLimitReached
	ok := errors.As(err, &errCompressLimitReached)
	return ok
}

func (ErrCompressLimit *ErrCompressLimitReached) Error() string {
	return fmt.Sprintf("total bytes limit reached with the following values: size limit: %d, total current size: %d", ErrCompressLimit.SizeLimit, ErrCompressLimit.CurrSize)
}

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
	if crc.Limit != 0 && *crc.Total > crc.Limit {
		return 0, newErrCompressLimitReached(crc.Limit, *crc.Total)
	}
	n, err := crc.Reader.Read(p)
	if err != nil && err != io.EOF {
		return n, err
	}
	*crc.Total += int64(n)
	if crc.Limit != 0 && *crc.Total > crc.Limit {
		return n, newErrCompressLimitReached(crc.Limit, *crc.Total)
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
