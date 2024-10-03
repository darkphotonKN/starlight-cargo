package fileservice

import (
	"bytes"
	"fmt"
)

/**
* A service to group the methods responsible for handling files.
*
* These methods are abstracted from the transport to maintain separation of concerns.
**/
type FileService struct {
}

func NewFileService() *FileService {
	return &FileService{}
}

// TODO: requires implementation
func (fs *FileService) DownloadFile(id string) ([]byte, error) {

	return []byte{}, nil
}

/**
* Handles uploading files.
* Pre-defined structure of the file payload is:
* A slice of bytes encoded from the file size + an empty space + file name
**/
func (fs *FileService) UploadFile(file []byte) error {
	byteParts := bytes.SplitN(file, []byte(" "), 2)
	fileName := string(byteParts[0])
	fileSize := string(byteParts[1])

	fmt.Println("fileName:", fileName)
	fmt.Println("fileSize:", fileSize)

	return nil
}
