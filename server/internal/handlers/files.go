package handlers

import (
	"encoding/json"
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
	"github.com/pavva91/merkle-tree/server/config"
	"github.com/pavva91/merkle-tree/server/internal/dto"
	"github.com/pavva91/merkle-tree/server/internal/errorhandlers"
)

type filesHandler struct{}

var (
	FilesHandler = filesHandler{}
	// TODO: refactor: use models package
	MerkleTreeMatrix = [][]string{}
)

// TODO: Refactor code (handlers, services)

// Bulk Upload godoc
//
//	@Summary		Bulk Upload
//	@Description	Bulk Upload all files in a given folder and create merkle tree
//	@Tags			Files
//	@Accept			multipart/form-data
//	@Produce		text/plain
//	@Param			file	formData	[]file	true	"files to upload"
//	@Failure		400		{object}	string
//	@Failure		500		{object}	string
//	@Router			/files [post]
func (h filesHandler) BulkUpload(w http.ResponseWriter, r *http.Request) {
	// TODO: move fs interaction in "services" package
	if err := r.ParseMultipartForm(int64(config.Values.Server.MaxBulkUploadSize)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get a reference to the fileHeaders.
	// They are accessible only after ParseMultipartForm is called
	files := r.MultipartForm.File["file"]

	err := os.RemoveAll(config.Values.Server.UploadFolder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = os.MkdirAll(config.Values.Server.UploadFolder, os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var rFiles []*os.File
	for k, fileHeader := range files {
		if fileHeader.Size > int64(config.Values.Server.MaxUploadFileSize) {
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
		fileName := fmt.Sprintf("%s/%d_%d_%s", config.Values.Server.UploadFolder, k+1, time.Now().UnixNano(), fileHeader.Filename)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("merkle-root of created merkle tree", MerkleTreeMatrix[len(MerkleTreeMatrix)-1][0])
	fmt.Fprintf(w, "upload successful")
}

// Download godoc
//
//	@Summary		Download
//	@Description	Download By Name
//	@Tags			Files
//	@Accept			json
//	@Produce		json
//	@Param			filename	path		string	true	"File Name"	Format(string)
//	@Failure		400			{object}	string
//	@Failure		404			{object}	string
//	@Failure		500			{object}	string
//	@Router			/files/{filename} [get]
func (h filesHandler) DownloadByName(w http.ResponseWriter, r *http.Request) {
	// TODO: move fs interaction in "services" package

	if len(MerkleTreeMatrix) == 0 {
		err := errors.New("no merkle tree, upload files first")
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileName := mux.Vars(r)["filename"]
	log.Println(fileName)

	filename := fileName

	dir, err := os.Open(config.Values.Server.UploadFolder)
	if err != nil {
		fmt.Println("error opening directory:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println("error reading directory:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	foundFilePath := config.Values.Server.UploadFolder

	for k, f := range files {
		ss := strings.SplitAfter(f.Name(), "_")
		if ss[len(ss)-1] == filename {
			fmt.Printf("file %v found: %s\n", k+1, f.Name())
			// foundFilePath = foundFilePath + "/" + f.Name()
			foundFilePath = fmt.Sprintf("%s/%s", foundFilePath, f.Name())
		}
	}

	// fileBytes, err := os.ReadFile("./uploads/3_1708595975295766854_f3")
	fmt.Println(foundFilePath)
	fileBytes, err := os.ReadFile(foundFilePath)
	if err != nil {
		err = errors.New("file not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	file, err := os.Open(foundFilePath)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	merkleProofs := merkletree.ComputeMerkleProof(file, MerkleTreeMatrix)
	mps := fmt.Sprintf("%+q", merkleProofs)
	w.Header().Add("Merkle-Proof", mps)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(fileBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// List godoc
//
//	@Summary		List
//	@Description	List files
//	@Tags			Files
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.ListFilesResponse
//	@Failure		400	{object}	string
//	@Failure		404	{object}	string
//	@Failure		500	{object}	string
//	@Router			/files [get]
func (h filesHandler) ListNames(w http.ResponseWriter, r *http.Request) {
	if len(MerkleTreeMatrix) == 0 {
		err := errors.New("no merkle tree, upload files first")
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dir, err := os.Open(config.Values.Server.UploadFolder)
	if err != nil {
		fmt.Println("error opening directory:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println("error reading directory:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// foundFilePath := config.Values.Server.UploadFolder
	fileNames := []string{}

	for _, f := range files {
		ss := strings.SplitAfter(f.Name(), "_")
		fileNames = append(fileNames, ss[len(ss)-1])
	}

	// fileBytes, err := os.ReadFile("./uploads/3_1708595975295766854_f3")
	// fmt.Println(foundFilePath)

	var res dto.ListFilesResponse
	res.ToDTO(fileNames)

	js, err := json.Marshal(res)
	if err != nil {
		log.Println(err.Error())
		errorhandlers.InternalServerErrorHandler(w, r)
		return
	}

	_, err = w.Write(js)
	if err != nil {
		log.Println(err)
		errorhandlers.InternalServerErrorHandler(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}
