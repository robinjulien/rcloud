package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/robinjulien/rcloud/pkg/enhancedmaps"
)

// SetUp sets up all the ressources needed for the use of the api
func SetUp(directorypath string, databasepath string) {
	os.Chdir(directorypath)

	if err := authstore.Users.ReadFile(databasepath); err != nil { // os.Args[2] is <database directory>
		if errors.Is(err, enhancedmaps.ErrorFileNotExist) { // File not exists or not having rights to read
			if err2 := authstore.Users.WriteFile(os.Args[2]); err2 != nil {
				panic(err2)
			}
		} else {
			panic(err)
		}
	}
}

// Handler returns the API endpoints handler
func Handler() http.Handler {
	router := http.NewServeMux()

	router.Handle("/auth/", http.StripPrefix("/auth", AuthHandler()))
	router.Handle("/fm/", http.StripPrefix("/fm", AuthMiddleware(FmHandler())))

	return router
}

// RespondJSON set headers, marshall json and send it to he client
func RespondJSON(w http.ResponseWriter, JSON interface{}) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK) // Implicit with w.Write

	json.NewEncoder(w).Encode(JSON)
}

func hasCorrectMethod(method string, r *http.Request) bool {
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

// MethodMiddleware is a middleware used to check the method of the handler
func MethodMiddleware(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !hasCorrectMethod(method, r) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}
