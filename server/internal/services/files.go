package services

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/pavva91/merkle-tree/server/config"
)

var File Filer = file{}

type Filer interface {
	ResetUploadDir() error
	SaveBulk(files []*multipart.FileHeader) error
	GetByName(fileName string) ([]byte, string, error)
	List() ([]string, error)
}

type file struct{}

func (s file) ResetUploadDir() error {
	err := os.RemoveAll(config.Values.Server.UploadFolder)
	if err != nil {
		return err
	}
	err = os.MkdirAll(config.Values.Server.UploadFolder, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (s file) SaveBulk(files []*multipart.FileHeader) error {
	for k, fileHeader := range files {
		if fileHeader.Size > int64(config.Values.Server.MaxUploadFileSize) {
			sizeMB := (config.Values.Server.MaxUploadFileSize / 1024) / 1024
			err := fmt.Errorf("the uploaded file is too big: %s. please use a file less than %vmb in size", fileHeader.Filename, sizeMB)
			return err
		}

		file, err := fileHeader.Open()
		if err != nil {
			return err
		}

		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			return err
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}

		fileName := fmt.Sprintf("%s/%d_%d_%s", config.Values.Server.UploadFolder, k+1, time.Now().UnixNano(), fileHeader.Filename)
		f, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s file) GetByName(fileName string) ([]byte, string, error) {
	dir, err := os.Open(config.Values.Server.UploadFolder)
	if err != nil {
		log.Println("error opening directory:", err)
		return nil, "", err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		log.Println("error reading directory:", err)
		return nil, "", err
	}

	foundFilePath := config.Values.Server.UploadFolder

	for k, f := range files {
		ss := strings.SplitAfter(f.Name(), "_")
		if ss[len(ss)-1] == fileName {
			fmt.Printf("file %v found: %s\n", k+1, f.Name())
			foundFilePath = fmt.Sprintf("%s/%s", foundFilePath, f.Name())
		}
	}

	fileBytes, err := os.ReadFile(foundFilePath)
	if err != nil {
		err = errors.New("file not found")
		return nil, "", err
	}

	return fileBytes, foundFilePath, nil
}

func (s file) List() ([]string, error) {
	dir, err := os.Open(config.Values.Server.UploadFolder)
	if err != nil {
		log.Println("error opening directory:", err)
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		log.Println("error reading directory:", err)
		return nil, err
	}

	fileNames := []string{}

	for _, f := range files {
		ss := strings.SplitAfter(f.Name(), "_")
		fileNames = append(fileNames, ss[len(ss)-1])
	}
	return fileNames, nil
}
