package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pavva91/merkle-tree/server/config"
	"github.com/pavva91/merkle-tree/server/internal/dto"
	"github.com/pavva91/merkle-tree/server/internal/services"
)

type filesHandler struct{}

var (
	FilesHandler = filesHandler{}
)

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
	if err := r.ParseMultipartForm(int64(config.Values.Server.MaxBulkUploadSize)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["file"]

	err := services.File.ResetUploadDir()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = services.File.SaveBulk(files)
	if err != nil {
		if strings.Contains(err.Error(), "uploaded file is too big") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	merkleTree, err := services.MerkleTree.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("merkle-root of created merkle tree", merkleTree.Matrix[len(merkleTree.Matrix)-1][0])
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

	if !services.MerkleTree.IsValid() {
		err := errors.New("no merkle tree, upload files first")
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileName := mux.Vars(r)["filename"]

	fileBytes, foundFilePath, err := services.File.GetByName(fileName)
	if err != nil {
		if strings.Contains(err.Error(), "file not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	merkleProofs, err := services.MerkleTree.CreateMerkleProof(foundFilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	if !services.MerkleTree.IsValid() {
		err := errors.New("no merkle tree, upload files first")
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileNames, err := services.File.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var res dto.ListFilesResponse
	res.ToDTO(fileNames)

	js, err := json.Marshal(res)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(js)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}
