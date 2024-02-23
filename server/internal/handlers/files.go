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
	"github.com/pavva91/merkle-tree/libs/merkletree"
	"github.com/pavva91/merkle-tree/server/internal/errorhandlers"
)

type filesHandler struct{}

var (
	FilesHandler     = filesHandler{}
	MerkleTreeMatrix = [][]string{}
)

// TODO: Use configs (file and envvars)
const MAX_UPLOAD_SIZE = 2 * 1024 * 1024 // 2MB
const UPLOAD_FOLDER = "./uploads"

func (h filesHandler) BulkUpload(w http.ResponseWriter, r *http.Request) {
	// 32 MB is the default used by FormFile()
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		errorhandlers.BadRequestHandler(w, r, errors.New("The uploaded bulk files are too big."))
		return
	}

	// Get a reference to the fileHeaders.
	// They are accessible only after ParseMultipartForm is called
	files := r.MultipartForm.File["file"]

	err := os.RemoveAll(UPLOAD_FOLDER)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = os.MkdirAll(UPLOAD_FOLDER, os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var rFiles []*os.File
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
		fileName := fmt.Sprintf("%s/%d_%d_%s", UPLOAD_FOLDER, k+1, time.Now().UnixNano(), fileHeader.Filename)
		f, err := os.Create(fileName)
		// f, err := os.Create(fmt.Sprintf("./uploads/%d_%d_%s", k+1, time.Now().UnixNano(), fileHeader.Filename))
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

		file2, err := os.Open(fileName)
		if err != nil {
			fmt.Println(err)
			return
		}
		rFiles = append(rFiles, file2)
		defer file2.Close()
	}

	// TODO: Save MerkleTree
	// TODO: 1. Save in memory
	MerkleTreeMatrix, err = merkletree.ComputeMerkleTreeAsMatrix(rFiles...)
	if err != nil {
		errorhandlers.InternalServerErrorHandler(w, r)
		return
	}

	fmt.Println("merkle-root of created merkle tree", MerkleTreeMatrix[len(MerkleTreeMatrix)-1][0])
	fmt.Fprintf(w, "Upload successful")
}

func (h filesHandler) DownloadByName(w http.ResponseWriter, r *http.Request) {
	fileName := mux.Vars(r)["filename"]
	log.Println(fileName)

	filename := fileName

	dir, err := os.Open(UPLOAD_FOLDER)
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

	foundFilePath := UPLOAD_FOLDER + "/"

	for _, file := range files {
		ss := strings.SplitAfter(file.Name(), "_")
		if ss[len(ss)-1] == filename {
			fmt.Println("File found:", file.Name())
			foundFilePath = foundFilePath + file.Name()
		}
	}

	// fileBytes, err := os.ReadFile("./uploads/3_1708595975295766854_f3")
	fmt.Println(foundFilePath)
	fileBytes, err := os.ReadFile(foundFilePath)
	if err != nil {
		errorhandlers.NotFoundHandler(w, r, errors.New("file not found"))
		return
	}

	file, err := os.Open(foundFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// TODO: retrieve merkletree from DB
	merkleProofs := merkletree.ComputeMerkleProof(file, MerkleTreeMatrix)
	mps := fmt.Sprintf("%+q", merkleProofs)
	w.Header().Add("Merkle-Proof", mps)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
	// http.ServeFile(w, r, "./uploads/*_f1")
	// fmt.Fprintf(w, "Download successful")
	// TODO: return: file, merkle proofs
	// File: Content-Type: application/octet-stream
	// Merkle Proof: array on header
}
