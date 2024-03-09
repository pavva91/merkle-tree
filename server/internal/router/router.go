package router

import (
	"github.com/gorilla/mux"
	"github.com/pavva91/merkle-tree/server/internal/handlers"
)

var Router *mux.Router

func NewRouter() {
	Router = mux.NewRouter()

	initializeRoutes()
}

func initializeRoutes() {
	files := Router.PathPrefix("/files").Subrouter()
	files.HandleFunc("", handlers.FilesHandler.BulkUpload).Methods("POST")
	files.HandleFunc("/", handlers.FilesHandler.BulkUpload).Methods("POST")
	files.HandleFunc("/{filename:[a-z]+[0-9]+}", handlers.FilesHandler.DownloadByName).Methods("GET")
	files.HandleFunc("", handlers.FilesHandler.ListNames).Methods("GET")
	files.HandleFunc("/", handlers.FilesHandler.ListNames).Methods("GET")
}
