package router

import (
	"github.com/gorilla/mux"
	"github.com/pavva91/task-third-party/internal/handlers"
)

var (
	Router *mux.Router
)

func NewRouter() {
	Router = mux.NewRouter()

	initializeRoutes()
}

func initializeRoutes() {
	files := Router.PathPrefix("/files").Subrouter()
	files.HandleFunc("", handlers.FilesHandler.UploadBulkFiles).Methods("POST")
	files.HandleFunc("/", handlers.FilesHandler.UploadBulkFiles).Methods("POST")
	// file := Router.PathPrefix("/file").Subrouter()
	// file.HandleFunc("", handlers.FilesHandler.UploadFileOnLocalStorage).Methods("POST")
	// file.HandleFunc("/", handlers.FilesHandler.UploadFileOnLocalStorage).Methods("POST")
}
