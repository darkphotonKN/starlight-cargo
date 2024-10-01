package fileservice

import "fmt"

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

// TODO: requires implementation
func (fs *FileService) UploadFile(file []byte) error {
	fmt.Println("Upload file")
	return nil
}
