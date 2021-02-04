package api

import (
	"net/http"

	"github.com/robinjulien/rcloud/pkg/sessions"
)

// UsersHandler is the function that returns handler for users management
func UsersHandler() http.Handler {
	router := http.NewServeMux()
	router.Handle("/list", MethodMiddleware("GET", http.HandlerFunc(ListUsers)))
	router.Handle("/add", MethodMiddleware("POST", http.HandlerFunc(AddUser)))
	router.Handle("/del", MethodMiddleware("POST", http.HandlerFunc(DelUser)))
	router.Handle("/edit", MethodMiddleware("POST", http.HandlerFunc(EditUser)))
	return router
}

type responseListUsers struct {
	BaseResponse
	Users []PublicUser `json:"users"`
}

// ListUsers /users/list lists all users
func ListUsers(w http.ResponseWriter, r *http.Request) {
	res := responseListUsers{}
	values := authstore.Users.Values()

	for _, v := range values {
		u, ok := v.(User)

		if !ok {
			RespondJSON(w, BaseResponse{Success: false})
			return
		}

		res.Users = append(res.Users, PublicUser{ID: u.ID, Admin: u.Admin})
	}

	res.Success = true
	RespondJSON(w, res)
}

// AddUser /users/add adds a user
func AddUser(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("id")
	pwdRaw := r.PostFormValue("pwd")
	admin := r.PostFormValue("admin")

	if id == "" || pwdRaw == "" {
		RespondJSON(w, BaseResponse{Success: false})
		return
	}

	if exists, err := authstore.Users.Exists(id); exists || err != nil { // User already exists or there is an error
		RespondJSON(w, BaseResponse{Success: false})
		return
	}

	pwdHash, ok := sessions.GeneratePwdHash([]byte(pwdRaw))

	if !ok {
		RespondJSON(w, BaseResponse{Success: false})
		return
	}

	err := authstore.Users.Set(id, User{
		ID:      id,
		PwdHash: pwdHash,
		Admin:   admin == "true",
	})

	if err != nil {
		RespondJSON(w, BaseResponse{Success: false})
		return
	}

	err = authstore.Users.WriteFile(authstore.Path)

	if err != nil {
		RespondJSON(w, BaseResponse{Success: false})
		return
	}

	RespondJSON(w, BaseResponse{Success: true})
}

// DelUser /users/del delete a user
func DelUser(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("id")

	if id == "" {
		RespondJSON(w, BaseResponse{Success: false})
		return
	}

	err := authstore.Users.Remove(id)

	if err != nil {
		RespondJSON(w, BaseResponse{Success: false})
		return
	}

	err = authstore.Users.WriteFile(authstore.Path)

	if err != nil {
		RespondJSON(w, BaseResponse{Success: false})
		return
	}

	RespondJSON(w, BaseResponse{Success: true})
}

// EditUser /users/edit edit an existing user
func EditUser(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("id")
	pwdRaw := r.PostFormValue("pwd")
	admin := r.PostFormValue("admin")

	if id == "" {
		RespondJSON(w, BaseResponse{Success: false})
		return
	}

	if exists, err := authstore.Users.Exists(id); !exists || err != nil { // User doesn't exists or there is an error
		RespondJSON(w, BaseResponse{Success: false})
		return
	}

	if pwdRaw != "" {
		pwdHash, ok := sessions.GeneratePwdHash([]byte(pwdRaw))

		if !ok {
			RespondJSON(w, BaseResponse{Success: false})
			return
		}

		uRaw, err := authstore.Users.GetSafe(id)

		u, ok := uRaw.(User)

		if !ok || err != nil {
			RespondJSON(w, BaseResponse{Success: false})
			return
		}

		u.PwdHash = pwdHash

		err = authstore.Users.Set(id, u)

		if err != nil {
			RespondJSON(w, BaseResponse{Success: false})
			return
		}
	}

	if admin != "" {
		uRaw, err := authstore.Users.GetSafe(id)

		u, ok := uRaw.(User)

		if !ok || err != nil {
			RespondJSON(w, BaseResponse{Success: false})
			return
		}

		u.Admin = admin == "true"

		err = authstore.Users.Set(id, u)

		if err != nil {
			RespondJSON(w, BaseResponse{Success: false})
			return
		}
	}

	err := authstore.Users.WriteFile(authstore.Path)

	if err != nil {
		RespondJSON(w, BaseResponse{Success: false})
		return
	}

	RespondJSON(w, BaseResponse{Success: true})
}
