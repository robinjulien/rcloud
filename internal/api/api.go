package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/robinjulien/rcloud/pkg/enhancedmaps"
	"github.com/robinjulien/rcloud/pkg/sessions"
)

// BaseResponse is the base response resturned by the API
type BaseResponse struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage"`
}

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
	router.Handle("/users/", http.StripPrefix("/users", AdminMiddleware(UsersHandler())))

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

// GetUserByCookies return a user or nil given an http request
func GetUserByCookies(r *http.Request) *User {
	sidcookie, errCSID := r.Cookie("sessionid")
	maccookie, errCMAC := r.Cookie("signature")

	if errCSID != nil || errCMAC != nil {
		return nil
	}

	sid := sessions.FromBase64(sidcookie.Value)
	mac := sessions.FromBase64(maccookie.Value)

	if !sessions.ValidMAC(sid, mac, ServerKey) {
		return nil
	}

	session := authstore.GetSessionByID(sid)

	if session == nil || session.Expires.Before(time.Now()) {
		return nil
	}

	u := authstore.GetUserByID(session.UID)

	return u // Can be nil
}

// AuthMiddleware is used as a middleware to know if a user is authenticated
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := GetUserByCookies(r)

		if u == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUser, *u)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminMiddleware is used as a middleware to know if a user is authenticated as admin
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := GetUserByCookies(r)

		if u == nil || !u.Admin {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUser, *u)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserFromContext returns the user contained in the context values from AuthMiddleware
func UserFromContext(ctx context.Context) User {
	return ctx.Value(ContextKeyUser).(User)
}
