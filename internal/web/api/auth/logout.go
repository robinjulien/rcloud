package auth

import (
	"errors"
	"net/http"

	"github.com/robinjulien/rcloud/internal/web/api/common"
	"github.com/robinjulien/rcloud/pkg/sessions"
)

// Logout /auth/logout endpoint to logout user
func Logout(w http.ResponseWriter, r *http.Request) {
	if !common.CheckMethod(w, r, "POST") {
		// Logout uniquement avec POST
		return
	}

	sessionIDCookie, errSID := r.Cookie("sessionid")
	_, errMAC := r.Cookie("signature")

	if !errors.Is(errSID, http.ErrNoCookie) {
		sessionID := sessions.FromBase64(sessionIDCookie.Value)
		authstore.RemoveSession(sessionID)

		http.SetCookie(w, &http.Cookie{
			Name:   "sessionid",
			MaxAge: -1,
			Path:   "/",
		})
	}

	if !errors.Is(errMAC, http.ErrNoCookie) {
		http.SetCookie(w, &http.Cookie{
			Name:   "signature",
			MaxAge: -1,
			Path:   "/",
		})
	}

	w.WriteHeader(http.StatusOK)
}
