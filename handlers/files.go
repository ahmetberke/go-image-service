package handlers

import (
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

type filesHandler struct {
	store store
}

func NewFilesHandler(s store) *filesHandler {
	return &filesHandler{store: s}
}

func (f *filesHandler) UploadREST(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fn := vars["filename"]
	log.Printf("[POST] id: %s, filename: %s", id, fn)
	f.saveFile(id, fn, rw, r.Body)
}

func (f *filesHandler) saveFile(id, path string, rw http.ResponseWriter, r io.ReadCloser) {
	log.Printf("[Saving] id: %s, file path: %s", id, path)
	fp := filepath.Join(id, path)
	err := f.store.Save(fp, r)
	if err != nil {
		log.Printf("[Error] unable to save file, %s", err.Error())
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}
