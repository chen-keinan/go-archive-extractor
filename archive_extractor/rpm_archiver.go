package archive_extractor

import (
	"fmt"
	"github.com/chen-keinan/go-archive-extractor/archive_extractor/archiver_errors"
	"github.com/chen-keinan/go-archive-extractor/compression"
	"github.com/chen-keinan/go-rpm"
	cpio "github.com/chen-keinan/gocpio"
	"io"
	"os"
	"path/filepath"
)

type RpmArchvier struct {
}

//Extract extract rpm archive
//accept rpm file path
//return file header metadata
func (za RpmArchvier) Extract(path string) ([]*ArchiveHeader, error) {
	headers := make([]*ArchiveHeader, 0)
	rpm, err := rpm.OpenPackageFile(filepath.Clean(path))
	if err != nil {
		return nil, archiver_errors.New(err)
	}
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	//Read content of cpio archive which starts after headers
	headerEnd := rpm.Headers[1].End
	archiveHead := make([]byte, 6)

	_, err = file.ReadAt(archiveHead, int64(headerEnd))
	if err != nil {
		return nil, err
	}

	//rewind to start of the file
	file.Seek(int64(headerEnd), 0)
	fileReader, err := compression.CreateCompressionFromBytes(archiveHead).GetReader(file)
	if err != nil || fileReader == nil {
		return nil, archiver_errors.New(err)
	}
	defer func() {
		err := fileReader.Close()
		if err != nil {
			fmt.Print(err.Error())
		}
	}()
	rc := cpio.NewReader(fileReader)
	var count = 0
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
					return nil, err
				}
				archiveHeader.PkgMeta = RpmMeta(rpm)
				headers = append(headers, archiveHeader)
			}
		}
	}
	return headers, nil
}
func RpmMeta(data *rpm.PackageFile) map[string]interface{} {
	return map[string]interface{}{
		"Name":     data.Name(),
		"Version":  data.Version(),
		"Release":  data.Release(),
		"Licenses": data.License(),
	}
}
