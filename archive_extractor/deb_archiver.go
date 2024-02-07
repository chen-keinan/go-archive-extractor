package archive_extractor

import (
	"errors"
	"fmt"
	"github.com/blakesmith/ar"
	"io"
	"os"

	"github.com/jfrog/go-archive-extractor/utils"
)

type DebArchiver struct {
	MaxCompressRatio   int64
	MaxNumberOfEntries int
}

const DebArchiverSkipFoldersCheckParamsKey = "DebArchiverSkipFoldersCheckParamsKey"

func (da DebArchiver) ExtractArchive(path string,
	processingFunc func(*ArchiveHeader, map[string]interface{}) error, params map[string]interface{}) error {
	maxBytesLimit, err := maxBytesLimit(path, da.MaxCompressRatio)
	if err != nil {
		return err
	}
	provider := LimitAggregatingReadCloserProvider{
		Limit: maxBytesLimit,
	}
	debFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer debFile.Close()
	rc := ar.NewReader(debFile)
	if rc == nil {
		return errors.New(fmt.Sprintf("Failed to open deb file : %s", path))
	}

	entriesCount := 0
	for {
		if da.MaxNumberOfEntries != 0 && entriesCount > da.MaxNumberOfEntries {
			return ErrTooManyEntries
		}
		entriesCount++
		archiveEntry, err := rc.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if archiveEntry == nil {
			return errors.New(fmt.Sprintf("Failed to open file : %s", path))
		}
		if skipFolderCheck(params) || !utils.IsFolder(archiveEntry.Name) {
			limitingReader := provider.CreateLimitAggregatingReadCloser(rc)
			archiveHeader := NewArchiveHeader(limitingReader, archiveEntry.Name, archiveEntry.ModTime.Unix(), archiveEntry.Size)
			err = processingFunc(archiveHeader, params)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func skipFolderCheck(params map[string]interface{}) bool {
	value, found := params[DebArchiverSkipFoldersCheckParamsKey]
	if !found {
		return false
	}
	boolValue, ok := value.(bool)
	return ok && boolValue
}
