package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal Server Error :%s path: %s error: %s", r.Method, r.URL.Path, err)
	_ = writeJSONError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad Request Error :%s path: %s error: %s", r.Method, r.URL.Path, err)
	_ = writeJSONError(w, http.StatusBadRequest, "Bad Request Error")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Conflict Error :%s path: %s error: %s", r.Method, r.URL.Path, err)
	_ = writeJSONError(w, http.StatusConflict, "Conflict Error")
}

func (app *application) unprocessableEntityResponse(w http.ResponseWriter, r *http.Request, err any) {
	log.Printf("Unprocessable Entity Error :%s path: %s error: %s", r.Method, r.URL.Path, err)
	_ = writeValidationJSONError(w, http.StatusUnprocessableEntity, "Bad Request", err)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Not Found Error :%s path: %s error: %s", r.Method, r.URL.Path, err)
	_ = writeJSONError(w, http.StatusNotFound, "Not Found Error")
}
