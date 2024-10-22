package handlers

import (
	"errors"
	"fmt"
	"io"
	"large-assignment/models"
	"large-assignment/services"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	clientManager services.MinioServiceManager
	bucketName    string
}

func NewHandler(clientManager services.MinioServiceManager, bucketName string) *Handler {
	return &Handler{
		clientManager: clientManager,
		bucketName:    bucketName,
	}
}

func (h *Handler) PutObjectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	client := h.clientManager.GetMinioService(id)
	if client == nil {
		http.Error(w, "No suitable MinIO client found", http.StatusInternalServerError)
		return
	}

	err = client.PutObject(r.Context(), h.bucketName, id, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to upload object: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Object uploaded successfully"))
}

func (h *Handler) GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	client := h.clientManager.GetMinioService(id)
	if client == nil {
		http.Error(w, "No suitable MinIO client found", http.StatusInternalServerError)
		return
	}

	data, err := client.GetObject(r.Context(), h.bucketName, id)
	if err != nil {
		if errors.Is(err, models.ERR_NOT_FOUND) {
			http.Error(w, "Object not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve object", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", id))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
