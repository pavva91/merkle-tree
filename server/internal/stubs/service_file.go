package stubs

import (
	"mime/multipart"
)

type FileService struct {
	ResetUploadDirFn func() error
	SaveBulkFn       func(files []*multipart.FileHeader) error
	GetByNameFn      func(fileName string) ([]byte, string, error)
	ListFn           func() ([]string, error)
}

func (stub FileService) ResetUploadDir() error {
	return stub.ResetUploadDirFn()
}

func (stub FileService) SaveBulk(files []*multipart.FileHeader) error {
	return stub.SaveBulkFn(files)
}

func (stub FileService) GetByName(fileName string) ([]byte, string, error) {
	return stub.GetByNameFn(fileName)
}

func (stub FileService) List() ([]string, error) {
	return stub.ListFn()
}
