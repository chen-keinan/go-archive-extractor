package archive_extractor

import (
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLimitingReaderErrorWhenLimitReached(t *testing.T) {
	provider := LimitAggregatingReadCloserProvider{
		Total: 0,
		Limit: 100,
	}
	reader := provider.CreateLimitAggregatingReadCloser(rand.Reader)
	b := make([]byte, 150)
	_, err := reader.Read(b)
	assert.EqualError(t, err, ErrCompressLimitReached.Error())
}

func TestLimitingReaderLimitNotReached(t *testing.T) {
	provider := LimitAggregatingReadCloserProvider{
		Total: 0,
		Limit: 100,
	}
	reader := provider.CreateLimitAggregatingReadCloser(rand.Reader)
	b := make([]byte, 10)
	for i := 0; i < 10; i++ {
		_, err := reader.Read(b)
		assert.NoError(t, err)
	}
}

func TestLimitingReaderErrorWhenAggregatingMultipleReadersFromSameProvider(t *testing.T) {
	provider := LimitAggregatingReadCloserProvider{
		Total: 0,
		Limit: 100,
	}
	b := make([]byte, 10)
	for i := 0; i < 10; i++ {
		reader := provider.CreateLimitAggregatingReadCloser(rand.Reader)
		_, err := reader.Read(b)
		assert.NoError(t, err)
	}
	reader := provider.CreateLimitAggregatingReadCloser(rand.Reader)
	_, err := reader.Read(b)
	assert.EqualError(t, err, ErrCompressLimitReached.Error())
}
