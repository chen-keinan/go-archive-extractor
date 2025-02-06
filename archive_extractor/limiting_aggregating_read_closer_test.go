//go:build tests_group_all

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
	expectedErr := newErrCompressLimitReached(100, 150)
	reader := provider.CreateLimitAggregatingReadCloser(rand.Reader)
	b := make([]byte, 150)
	_, err := reader.Read(b)
	assert.Error(t, err)
	actualErr, ok := err.(*ErrCompressLimitReached)
	assert.True(t, ok)
	assert.Equal(t, *actualErr, *expectedErr)
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
	assert.True(t, IsErrCompressLimitReached(err))
}
