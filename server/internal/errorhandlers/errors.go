package errorhandlers

import (
	"log"
	"net/http"
)

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte("500 Internal Server Error"))
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func BadRequestHandler(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusNotFound)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		log.Println(err.Error())
		return
	}
}
