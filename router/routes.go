package router

import (
	"large-assignment/handlers"
	"large-assignment/middleware"

	"github.com/gorilla/mux"
)

func InitRoutes(h handlers.Handler) *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.ValidateIDMiddleware)
	r.HandleFunc("/object/{id}", h.GetObjectHandler).Methods("GET")
	r.HandleFunc("/object/{id}", h.PutObjectHandler).Methods("PUT")

	return r
}
