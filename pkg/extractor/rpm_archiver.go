package extractor

import (
	"fmt"
	compression2 "github.com/chen-keinan/go-archive-extractor/pkg/compression"
	aerrors2 "github.com/chen-keinan/go-archive-extractor/pkg/extractor/aerrors"
	"github.com/chen-keinan/go-rpm"
	cpio "github.com/chen-keinan/gocpio"
	"io"
	"os"
	"path/filepath"
)

//RpmArchvier object
type RpmArchvier struct {
}

//Extract extract rpm archive
//accept rpm file path
//return file header metadata
func (za RpmArchvier) Extract(path string) ([]*ArchiveHeader, error) {
	headers := make([]*ArchiveHeader, 0)
	rpm, err := rpm.OpenPackageFile(filepath.Clean(path))
	if err != nil {
		return nil, aerrors2.New(err)
	}
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = file.Close(); err != nil {
			fmt.Print(err.Error())
		}
	}()
	headerEnd := rpm.Headers[1].End
	archiveHead := make([]byte, 6)

	_, err = file.ReadAt(archiveHead, int64(headerEnd))
	if err != nil {
		return nil, err
	}

	//rewind to start of the file
	_, err = file.Seek(int64(headerEnd), 0)
	if err != nil {
		return nil, err
	}
	fileReader, err := compression2.CreateCompressionFromBytes(archiveHead).GetReader(file)
	if err != nil || fileReader == nil {
		return nil, aerrors2.New(err)
	}
	defer func() {
		err = fileReader.Close()
		if err != nil {
			fmt.Print(err.Error())
		}
	}()
	rc := cpio.NewReader(fileReader)
	var count = 0
	headers, archiveHeaders, err := za.extractHeaders(rc, count, rpm, headers)
	if err != nil {
		return archiveHeaders, err
	}
	return headers, nil
}

func (za RpmArchvier) extractHeaders(rc *cpio.Reader, count int, rpm *rpm.PackageFile, headers []*ArchiveHeader) ([]*ArchiveHeader, []*ArchiveHeader, error) {
	for {
		archiveEntry, err := rc.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		count++
		//skip trailer
		if archiveEntry.IsTrailer() {
			break
		}
		if archiveEntry.Mode != cpio.TYPE_DIR {
			if archiveEntry != nil {
				archiveHeader, err := NewArchiveHeader(rc, archiveEntry.Name, archiveEntry.Mtime, archiveEntry.Size)
				if err != nil {
					return nil, nil, err
				}
				archiveHeader.PkgMeta = RpmMeta(rpm)
				headers = append(headers, archiveHeader)
			}
		}
	}
	return headers, nil, nil
}

//RpmMeta return rpm metadata as key/value
// accept rpm headers
func RpmMeta(data *rpm.PackageFile) map[string]interface{} {
	return map[string]interface{}{
		"Name":     data.Name(),
		"Version":  data.Version(),
		"Release":  data.Release(),
		"Licenses": data.License(),
	}
}
