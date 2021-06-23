package archiver_errors

import "fmt"

//ArchiverExtractorError archive error object
type ArchiverExtractorError struct {
	archiverError error
}

//New instantiate new archive error object
func New(e error) ArchiverExtractorError {
	return ArchiverExtractorError{archiverError: e}
}

//Error return archive error
func (aee ArchiverExtractorError) Error() string {
	return fmt.Sprintf("Archive extractor error,%s", aee.archiverError.Error())
}
