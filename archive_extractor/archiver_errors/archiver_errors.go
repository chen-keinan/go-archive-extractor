package archiver_errors

import "fmt"

type ArchiverExtractorError struct {
	archiverError error
}

func New(e error) ArchiverExtractorError {
	return ArchiverExtractorError{archiverError: e}
}

func (aee ArchiverExtractorError) Error() string {
	return fmt.Sprintf("Failed to Open Archive,%s", aee.archiverError.Error())
}
