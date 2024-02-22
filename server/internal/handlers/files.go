package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pavva91/task-third-party/internal/errorhandlers"
)

type filesHandler struct{}

var (
	FilesHandler = filesHandler{}
)

const MAX_UPLOAD_SIZE = 2 * 1024 * 1024 // 2MB

// Upload File godoc
//
//	@Summary		Upload File
//	@Description	Upload a File
//	@Tags			File
//	@Accept			json
//	@Produce		json
//
// TODO: swag for file upload (Content-Type multipart/form-data)
//
//	@Param			request	body		dto.UploadFileRequest	true	"query params"
//	@Success		200		{object}	dto.UploadFileResponse
//	@Failure		400		{object}	string
//	@Failure		500		{object}	string
//	@Router			/file [post]
//
func (h filesHandler) UploadFileOnLocalStorage(w http.ResponseWriter, r *http.Request) {

	// func (h *FilesHandler) UploadFileOnLocalStorage(w http.ResponseWriter, r *http.Request) {
	// Parse request body as multipart form data with 32MB max memory
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Println(err)
	}

	// Get file uploaded via Form
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		errorhandlers.InternalServerErrorHandler(w, r)
		return
	}
	defer file.Close()

	// Create file locally
	filePath := fmt.Sprintf("tmp/%d_%s", time.Now().UnixNano(), handler.Filename)
	localFile, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		errorhandlers.InternalServerErrorHandler(w, r)
		return
	}
	defer localFile.Close()

	// Copy the uploaded file data to the newly created file on the filesystem
	if _, err := io.Copy(localFile, file); err != nil {
		log.Println(err)
		errorhandlers.InternalServerErrorHandler(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h filesHandler) UploadBulkFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 32 MB is the default used by FormFile()
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "The uploaded bulk files are too big.", http.StatusBadRequest)
		return
	}

	// Get a reference to the fileHeaders.
	// They are accessible only after ParseMultipartForm is called
	files := r.MultipartForm.File["file"]

	for _, fileHeader := range files {
		// Restrict the size of each uploaded file to 1MB.
		// To prevent the aggregate size from exceeding
		// a specified value, use the http.MaxBytesReader() method
		// before calling ParseMultipartForm()
		if fileHeader.Size > MAX_UPLOAD_SIZE {
			http.Error(w, fmt.Sprintf("The uploaded image is too big: %s. Please use an image less than 1MB in size", fileHeader.Filename), http.StatusBadRequest)
			return
		}

		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = os.MkdirAll("./uploads", os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// f, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
		f, err := os.Create(fmt.Sprintf("./uploads/%d_%s", time.Now().UnixNano(), fileHeader.Filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	fmt.Fprintf(w, "Upload successful")
}
