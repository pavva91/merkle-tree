package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/pavva91/task-third-party/internal/errorhandlers"
)

type filesHandler struct{}

var (
	FilesHandler = filesHandler{}
)

const MAX_UPLOAD_SIZE = 2 * 1024 * 1024 // 2MB

func (h filesHandler) BulkUpload(w http.ResponseWriter, r *http.Request) {
	// 32 MB is the default used by FormFile()
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		errorhandlers.BadRequestHandler(w, r, errors.New("The uploaded bulk files are too big."))
		return
	}

	// Get a reference to the fileHeaders.
	// They are accessible only after ParseMultipartForm is called
	files := r.MultipartForm.File["file"]

	err := os.RemoveAll("./uploads")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for k, fileHeader := range files {
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
		// TODO: add client identifier (to handle multiple clients)
		f, err := os.Create(fmt.Sprintf("./uploads/%d_%d_%s", k+1, time.Now().UnixNano(), fileHeader.Filename))
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
		// TODO: Create Merkle Tree
	}

	fmt.Fprintf(w, "Upload successful")
}

func (h filesHandler) DownloadByName(w http.ResponseWriter, r *http.Request) {
	fileName := mux.Vars(r)["filename"]
	log.Println(fileName)

	dirname := "./uploads"
	filename := fileName

	dir, err := os.Open(dirname)
	if err != nil {
		fmt.Println("Error opening directory:", err)
		return
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	foundFilePath := dirname + "/"

	for _, file := range files {
		ss := strings.SplitAfter(file.Name(), "_")
		if ss[len(ss)-1] == filename {
			fmt.Println("File found:", file.Name())
			foundFilePath = foundFilePath + file.Name()
		}
	}

	arr := []string{"foo", "bar", "baz"}
	mp := fmt.Sprintf("%+q", arr)
	// NOTE: To reconstruct string[] from mp:
	// result1 := mp[1 : len(mp)-2]
	result1 := strings.Replace(mp, "[", "", -1)
	result2 := strings.Replace(result1, "]", "", -1)
	result3 := strings.Replace(result2, "\"", "", -1)
	result := strings.SplitAfter(result3, " ")
	log.Println(result[0])
	log.Println(result[1])
	log.Println(result[2])

	w.Header().Add("Merkle-Proof", mp)

	// fileBytes, err := os.ReadFile("./uploads/3_1708595975295766854_f3")
	fmt.Println(foundFilePath)
	fileBytes, err := os.ReadFile(foundFilePath)
	if err != nil {
		errorhandlers.NotFoundHandler(w, r, errors.New("file not found"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
	// http.ServeFile(w, r, "./uploads/*_f1")
	// fmt.Fprintf(w, "Download successful")
	// TODO: return: file, merkle proofs
	// File: Content-Type: application/octet-stream
	// Merkle Proof: array on header
}
