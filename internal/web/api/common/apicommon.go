package common

import (
	"encoding/json"
	"net/http"
)

// RespondJSON set headers, marshall json and send it to he client
func RespondJSON(w http.ResponseWriter, JSON interface{}) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK) // Implicit with w.Write

	json.NewEncoder(w).Encode(JSON)
}

func hasCorrectMethod(r *http.Request, method string) bool {
	switch method {
	case "GET":
		return r.Method == http.MethodGet
	case "POST":
		return r.Method == http.MethodPost
	case "PUT":
		return r.Method == http.MethodPut
	case "DELETE":
		return r.Method == http.MethodDelete
	default:
		return false
	}
}

// CheckMethod return if the current request has the right method, if not it send StatusMethodNotAllowed
func CheckMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if hasCorrectMethod(r, method) {
		return true
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	return false
}
