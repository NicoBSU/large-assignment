package middleware

import (
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

func ValidateIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, exists := vars["id"]

		if !exists {
			http.Error(w, "ID is missing in the request URL", http.StatusBadRequest)
			return
		}

		if len(id) > 32 {
			http.Error(w, "ID exceeds the maximum allowed length of 32", http.StatusBadRequest)
			return
		}

		match, err := regexp.MatchString("^[a-zA-Z0-9]+$", id)
		if !match || err != nil {
			http.Error(w, "ID contains non-alphanumeric characters", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
