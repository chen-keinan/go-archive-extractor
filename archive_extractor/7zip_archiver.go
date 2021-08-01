package archive_extractor

import (
	"errors"
	"fmt"
	archive "github.com/gen2brain/go-unarr"
	"github.com/jfrog/go-archive-extractor/utils"
	"io"
)

type SevenZipArchiver struct {
	MaxCompressRatio   int64
	MaxNumberOfEntries int
}

func (sa SevenZipArchiver) ExtractArchive(path string,
	processingFunc func(*ArchiveHeader, map[string]interface{}) error, params map[string]interface{}) error {
	maxBytesLimit, err := maxBytesLimit(path, sa.MaxCompressRatio)
	if err != nil {
		return err
	}
	provider := LimitAggregatingReadCloserProvider{
		Limit: maxBytesLimit,
	}
	r, err := archive.NewArchive(path)
	if err != nil {
		return err
	}
	allFiles, err := r.List()
	if err != nil {
		return err
	}
	defer r.Close()

	if sa.MaxNumberOfEntries > 0 && len(allFiles) > sa.MaxNumberOfEntries {
		return ErrTooManyEntries
	}
	for _, archiveEntry := range allFiles {
		err := r.EntryFor(archiveEntry)
		if err != nil {
			return err
		}
		if !utils.IsFolder(archiveEntry) {
			rc := &SevenZipReader{Archive: r, Size: r.Size()}
			countingReadCloser := provider.CreateLimitAggregatingReadCloser(rc)
			archiveHeader := NewArchiveHeader(countingReadCloser, r.Name(), r.ModTime().Unix(), int64(r.Size()))
			err = processingFunc(archiveHeader, params)
			rc.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type SevenZipReader struct {
	Archive *archive.Archive
	Size    int
}

func (a *SevenZipReader) Read(p []byte) (n int, err error) {
	if a.Size <= 0 {
		return 0, io.EOF
	}
	size := len(p)
	if len(p) > a.Size {
		size = a.Size
	}
	b := make([]byte, size)
	n, err = a.Archive.Read(b)
	if err != nil && err != io.EOF {
		return 0, err
	}
	copied := copy(p, b)
	if copied != n {
		return 0, errors.New(fmt.Sprintf("copy arrays failed, copied only %v from %v bytes", copied, n))
	}
	a.Size -= n
	return n, nil
}

func (a *SevenZipReader) Close() error {
	return nil
}
