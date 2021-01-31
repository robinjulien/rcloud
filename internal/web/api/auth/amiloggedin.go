package auth

import (
	"net/http"

	"github.com/robinjulien/rcloud/internal/web/api/common"
)

type responseAmILoggedIn struct {
	LoggedIn bool `json:"loggedIn"`
}

// AmILoggedIn /auth/amiloggedin tells wheras the given sessionid is valid or not
func AmILoggedIn(w http.ResponseWriter, r *http.Request) {
	if !common.CheckMethod(w, r, "GET") {
		return
	}

	u := GetUserByCookies(r)

	if u == nil {
		common.RespondJSON(w, responseAmILoggedIn{
			LoggedIn: false,
		})
		return
	}

	common.RespondJSON(w, responseAmILoggedIn{
		LoggedIn: true,
	})
}
