# go-archive-extractor

The archive-extractor is a library and set of tools
that can open archive types (tar , zip , rpm ,deb, 7zip) and invoke advance processing method
while iterating archive headers
This library encapsulate logic from 2 best licenses detection libraries :

- Define advance params to be uses in advance processing method :
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
func params() map[string]interface{} {
	return map[string]interface{}{
		"archveData": &ArchiveData{},
	}
}
```
- Define advance processing method to be invoke during archive extraction :
```
func advanceProcessing(header *ArchiveHeader, params map[string]interface{}) error {
	if len(advanceProcessingParams) == 0 {
		return errors.New("Advance processing params are missing")
	}
	var ok bool
	var archiveData *ArchiveData
	if archiveData, ok = params["archiveData"].(*ArchiveData); !ok {
		return errors.New("Advance processing archveData param is missing")
	}
	archiveData.Name = header.Name
	archiveData.ModTime=header.ModTime
	archiveData.Size=header.Size
	archiveData.IsFolder=header.IsFolder
 	fmt.Print(archveData)
	return nil
}
```
- create archive extractor type and pass advance processing method and params :
```
func main() {
	za := &ZipArchvier{}
	err:=za.ExtractArchive("/User/Name/file.zip",advanceProcessing,params())
	if err != nil{
		fmt.Print(err)
	}
}
```
