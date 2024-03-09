package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pavva91/merkle-tree/server/internal/models"
	"github.com/pavva91/merkle-tree/server/internal/services"
	"github.com/pavva91/merkle-tree/server/internal/stubs"
)

func Test_filesHandler_BulkUpload(t *testing.T) {
	uploadFolder := "../../testfiles/3files"

	wrongRequest, err := http.NewRequest("POST", "/files", nil)
	if err != nil {
		t.Fatal(err)
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	files, err := os.ReadDir(uploadFolder)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, f := range files {
		if f.IsDir() {
			fmt.Println("use a folder with only the files you want to upload, without subfolders")
			return
		}
	}

	for _, f := range files {
		filePath := uploadFolder + "/" + f.Name()
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		part, err := writer.CreateFormFile("file", filepath.Base(filePath))
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = io.Copy(part, file)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	r, err := http.NewRequest("POST", "/files", payload)
	if err != nil {
		t.Fatal(err)
	}

	r.Header.Add("Content-Type", "multipart/form-data")
	r.Header.Set("Content-Type", writer.FormDataContentType())

	stubFileFail1 := stubs.FileService{}
	stubFileFail1.ResetUploadDirFn = func() error {
		return errors.New("stub error file 1")
	}
	stubFileFail2 := stubs.FileService{}
	stubFileFail2.ResetUploadDirFn = func() error {
		return nil
	}
	stubFileFail2.SaveBulkFn = func(_ []*multipart.FileHeader) error {
		return errors.New("stub error file 2")
	}

	stubFileFail3 := stubs.FileService{}
	stubFileFail3.ResetUploadDirFn = func() error {
		return nil
	}
	stubFileFail3.SaveBulkFn = func(_ []*multipart.FileHeader) error {
		err := fmt.Errorf("The uploaded file is too big: %s. Please use a file less than 2MB in size", "path/to/file/name")
		return err
	}

	stubFileOK := stubs.FileService{}
	stubFileOK.ResetUploadDirFn = func() error {
		return nil
	}
	stubFileOK.SaveBulkFn = func(_ []*multipart.FileHeader) error {
		return nil
	}

	stubMTfail := stubs.MerkleTreeService{}
	stubMTfail.CreateFn = func() (*models.MerkleTree, error) {
		return nil, errors.New("stub error mt")
	}

	merkleTreeStub := &models.MerkleTree{
		Matrix: [][]string{
			{"root-hash"},
		},
	}

	stubMTOK := stubs.MerkleTreeService{}
	stubMTOK.CreateFn = func() (*models.MerkleTree, error) {
		return merkleTreeStub, nil
	}

	type args struct {
		r *http.Request
	}
	tests := map[string]struct {
		args             args
		stub1            stubs.FileService
		stub2            stubs.MerkleTreeService
		wantErr          bool
		expectedHTTPCode int
		expectedResBody  string
	}{
		"wrong request type": {
			args{
				wrongRequest,
			},
			stubs.FileService{},
			stubs.MerkleTreeService{},
			true,
			400,
			"request Content-Type isn't multipart/form-data\n",
		},
		"request ok, failing reset upload dir": {
			args{
				r,
			},
			stubFileFail1,
			stubs.MerkleTreeService{},
			true,
			500,
			"stub error file 1\n",
		},
		"failing save bulk with generic error": {
			args{
				r,
			},
			stubFileFail2,
			stubs.MerkleTreeService{},
			true,
			500,
			"stub error file 2\n",
		},
		"failing save bulk with file too big error": {
			args{
				r,
			},
			stubFileFail3,
			stubs.MerkleTreeService{},
			true,
			400,
			"The uploaded file is too big: path/to/file/name. Please use a file less than 2MB in size\n",
		},
		"failing create merkle tree": {
			args{
				r,
			},
			stubFileOK,
			stubMTfail,
			true,
			500,
			"stub error mt\n",
		},
		"create merkle tree ok": {
			args{
				r,
			},
			stubFileOK,
			stubMTOK,
			false,
			200,
			"upload successful",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			services.File = tt.stub1
			services.MerkleTree = tt.stub2

			w := httptest.NewRecorder()

			handler := http.HandlerFunc(FilesHandler.BulkUpload)
			handler.ServeHTTP(w, tt.args.r)

			if (w.Code != 200) != tt.wantErr {
				t.Errorf("CreateTaskRequest.Validate() error = %v, wantErr %v", w.Code, tt.wantErr)
			}
			if w.Code != tt.expectedHTTPCode {
				t.Errorf("CreateTaskRequest.Validate() error = %v, wantErr %v", w.Code, tt.expectedHTTPCode)
			}
			if w.Body.String() != tt.expectedResBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					w.Body.String(), tt.expectedResBody)
			}
		})
	}
}

func Test_filesHandler_DownloadByName(t *testing.T) {
	stubMTfail1 := stubs.MerkleTreeService{}
	stubMTfail1.IsValidFn = func() bool {
		return false
	}

	stubMTfail2 := stubs.MerkleTreeService{}
	stubMTfail2.IsValidFn = func() bool {
		return true
	}
	stubMTfail2.CreateMerkleProofFn = func(_ string) ([]string, error) {
		return []string{}, errors.New("stub merkle tree error 1")
	}

	stubMTOK := stubs.MerkleTreeService{}
	stubMTOK.IsValidFn = func() bool {
		return true
	}
	stubMTOK.CreateMerkleProofFn = func(_ string) ([]string, error) {
		return []string{}, nil
	}

	stubFileFail1 := stubs.FileService{}
	stubFileFail1.GetByNameFn = func(_ string) ([]byte, string, error) {
		return []byte{}, "", errors.New("stub error file 1")
	}

	stubFileFail2 := stubs.FileService{}
	stubFileFail2.GetByNameFn = func(_ string) ([]byte, string, error) {
		return []byte{}, "", errors.New("file not found")
	}

	stubFileOK := stubs.FileService{}
	stubFileOK.GetByNameFn = func(_ string) ([]byte, string, error) {
		return []byte{}, "", nil
	}

	type args struct {
		vars map[string]string
	}
	tests := map[string]struct {
		args             args
		stub1            stubs.MerkleTreeService
		stub2            stubs.FileService
		wantErr          bool
		expectedHTTPCode int
		expectedResBody  string
	}{
		"not initialized merkle tree": {
			args{
				vars: map[string]string{
					"filename": "",
				},
			},
			stubMTfail1,
			stubs.FileService{},
			true,
			400,
			"no merkle tree, upload files first\n",
		},
		"get file fails, generic error": {
			args{
				vars: map[string]string{
					"filename": "",
				},
			},
			stubMTfail2,
			stubFileFail1,
			true,
			500,
			"stub error file 1\n",
		},
		"get file fails, file not found error": {
			args{
				vars: map[string]string{
					"filename": "",
				},
			},
			stubMTfail2,
			stubFileFail2,
			true,
			404,
			"file not found\n",
		},
		"create merkle proof fails": {
			args{
				vars: map[string]string{
					"filename": "",
				},
			},
			stubMTfail2,
			stubFileOK,
			true,
			500,
			"stub merkle tree error 1\n",
		},
		"create merkle proof ok": {
			args{
				vars: map[string]string{
					"filename": "",
				},
			},
			stubMTOK,
			stubFileOK,
			false,
			200,
			"",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			services.MerkleTree = tt.stub1
			services.File = tt.stub2

			r, err := http.NewRequest("GET", "/files", nil)
			if err != nil {
				t.Fatal(err)
			}

			r = mux.SetURLVars(r, tt.args.vars)

			w := httptest.NewRecorder()

			handler := http.HandlerFunc(FilesHandler.DownloadByName)
			handler.ServeHTTP(w, r)

			if (w.Code != 200) != tt.wantErr {
				t.Errorf("CreateTaskRequest.Validate() error = %v, wantErr %v", w.Code, tt.wantErr)
			}
			if w.Code != tt.expectedHTTPCode {
				t.Errorf("CreateTaskRequest.Validate() error = %v, wantErr %v", w.Code, tt.expectedHTTPCode)
			}
			if w.Body.String() != tt.expectedResBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					w.Body.String(), tt.expectedResBody)
			}
		})
	}
}

func Test_filesHandler_ListNames(t *testing.T) {
	stubMTfail1 := stubs.MerkleTreeService{}
	stubMTfail1.IsValidFn = func() bool {
		return false
	}

	stubMTOK := stubs.MerkleTreeService{}
	stubMTOK.IsValidFn = func() bool {
		return true
	}

	stubFileFail1 := stubs.FileService{}
	stubFileFail1.ListFn = func() ([]string, error) {
		return []string{}, errors.New("stub error get file")
	}

	stubFileEmpty := stubs.FileService{}
	stubFileEmpty.ListFn = func() ([]string, error) {
		return []string{}, nil
	}

	stubFileOK := stubs.FileService{}
	stubFileOK.ListFn = func() ([]string, error) {
		return []string{"f1", "f2", "f3"}, nil
	}

	tests := map[string]struct {
		stub1            stubs.MerkleTreeService
		stub2            stubs.FileService
		wantErr          bool
		expectedHTTPCode int
		expectedResBody  string
	}{
		"not initialized merkle tree": {
			stubMTfail1,
			stubs.FileService{},
			true,
			400,
			"no merkle tree, upload files first\n",
		},
		"list files fails, generic error": {
			stubMTOK,
			stubFileFail1,
			true,
			500,
			"stub error get file\n",
		},
		"list files ok, empty list": {
			stubMTOK,
			stubFileEmpty,
			false,
			200,
			"{\"filenames\":[]}",
		},
		"list files ok": {
			stubMTOK,
			stubFileOK,
			false,
			200,
			"{\"filenames\":[\"f1\",\"f2\",\"f3\"]}",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			services.MerkleTree = tt.stub1
			services.File = tt.stub2

			r, err := http.NewRequest("GET", "/files", nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()

			handler := http.HandlerFunc(FilesHandler.ListNames)
			handler.ServeHTTP(w, r)

			if (w.Code != 200) != tt.wantErr {
				t.Errorf("CreateTaskRequest.Validate() error = %v, wantErr %v", w.Code, tt.wantErr)
			}
			if w.Code != tt.expectedHTTPCode {
				t.Errorf("CreateTaskRequest.Validate() error = %v, wantErr %v", w.Code, tt.expectedHTTPCode)
			}
			if w.Body.String() != tt.expectedResBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					w.Body.String(), tt.expectedResBody)
			}
		})
	}
}
