package archive_extractor

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"github.com/jfrog/go-archive-extractor/archive_extractor/archiver_errors"
	"io"
	"os"
)

const fileHeaderSignatureString = "PK\x03\x04"

type ZipArchiver struct {
	MaxCompressRatio   int64
	MaxNumberOfEntries int
}

type ZipReadCloser struct {
	*zip.Reader
	io.Closer
}

func (za ZipArchiver) ExtractArchive(path string,
	processingFunc func(*ArchiveHeader, map[string]interface{}) error, params map[string]interface{}) error {
	maxBytesLimit, err := maxBytesLimit(path, za.MaxCompressRatio)
	if err != nil {
		return err
	}
	rcProvider := LimitAggregatingReadCloserProvider{Limit: maxBytesLimit}
	r, err := openZipReader(path)
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

func openZipReader(name string) (*ZipReadCloser, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	r, err := initZipReader(f, fi.Size())
	if err != nil {
		f.Close()
		return nil, err
	}
	zr := &ZipReadCloser{
		Reader: r,
		Closer: f,
	}
	return zr, nil
}

// This method was added in order to support prepended zip files.
// In case we failed get a zip.Reader from r because of zip.ErrFormat error,
// it searches for begin-of-file signature 0x04034b50 (in little-endian the value of the signature is PK0304).
// Once the signature was found, tries to read the file starting at that offset using zip.NewReader.
func initZipReader(r io.ReaderAt, size int64) (*zip.Reader, error) {
	zr, err := zip.NewReader(r, size)
	if err == nil || !errors.Is(err, zip.ErrFormat) {
		return zr, err
	}
	const BUFSIZE = 4096
	var buf [BUFSIZE + 4]byte
	for i := int64(0); (i-1)*BUFSIZE < size; i++ {
		len, err := r.ReadAt(buf[:], i*BUFSIZE)
		if err != nil && err != io.EOF {
			break
		}
		n := 0
		for {
			m := bytes.Index(buf[n:len], []byte(fileHeaderSignatureString))
			if m == -1 {
				break
			}
			off := i*BUFSIZE + int64(n+m)
			zipSize := size - int64(off)
			sr := io.NewSectionReader(r, int64(off), zipSize)
			if zr, ze := zip.NewReader(sr, zipSize+1); ze == nil {
				return zr, nil
			}
			n += m + 1
		}
		if err == io.EOF {
			break
		}
	}
	return nil, errors.New("No zip file found")
}
