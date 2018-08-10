# go-archive-extractor

The archive-extractor is a library and set of tools
that can open archive types (tar , zip , rpm ,deb, 7zip) and invoke advance processing method
while iterating archive headers
This library encapsulate logic from 2 best licenses detection libraries :

Example :
 ```
type ArchiveData struct {
	ArchiveReader io.Reader
	IsFolder      bool
	Name          string
	ModTime       int64
	Size          int64
}
```
```
func advanceProcessingParams() map[string]interface{} {
	return map[string]interface{}{
		"archveData": &ArchiveData{},
	}
}
```
```
func advanceProcessing(header *ArchiveHeader, advanceProcessingParams map[string]interface{}) error {
	if len(advanceProcessingParams) == 0 {
		return errors.New("Advance processing params are missing")
	}
	var ok bool
	var archveData *ArchiveData
	if archveData, ok = advanceProcessingParams["archveData"].(*ArchiveData); !ok {
		return errors.New("Advance processing parentFile param is missing")
	}
	archveData.Name = header.Name
	archveData.ModTime=header.ModTime
	archveData.Size=header.Size
	archveData.IsFolder=header.IsFolder
	// create child files
	fmt.Print(archveData)
	return nil
}
```
```
func main() {
	za := &ZipArchvier{}
	err:=za.ExtractArchive("/User/Name/file.zip",advanceProcessing,advanceProcessingParams())
	if err != nil{
		fmt.Print(err)
	}
}
```
