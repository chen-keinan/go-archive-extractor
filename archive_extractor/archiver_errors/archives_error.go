package archiver_errors

import "fmt"

var (
	RarDecodeError      = fmt.Errorf("rardecode: RAR signature not found")
	SevenZipDecodeError = fmt.Errorf("sevenzip: not a valid 7-zip file")
)
