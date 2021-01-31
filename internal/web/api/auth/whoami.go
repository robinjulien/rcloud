package auth

import (
	"net/http"

	"github.com/robinjulien/rcloud/internal/web/api/common"
)

// PublicUser is user's public informations to be displayed
type PublicUser struct {
	ID    string `json:"id"`
	Admin bool   `json:"admin"`
}

type responseWhoAmI struct {
	LoggedIn bool `json:"loggedIn"`
	PublicUser
}

// WhoAmI /auth/amiloggedin tells who you are if you're connected
func WhoAmI(w http.ResponseWriter, r *http.Request) {
	if !common.CheckMethod(w, r, "GET") {
		return
	}

	u := GetUserByCookies(r)
	if u == nil {
		common.RespondJSON(w, responseWhoAmI{LoggedIn: false, PublicUser: PublicUser{}})
		return
	}

	common.RespondJSON(w, responseWhoAmI{LoggedIn: true, PublicUser: PublicUser{
		ID:    u.ID,
		Admin: u.Admin,
	}})
}
