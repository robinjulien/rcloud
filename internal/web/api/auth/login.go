package auth

import (
	"net/http"
	"time"

	"github.com/robinjulien/rcloud/internal/web/api/common"
	"github.com/robinjulien/rcloud/pkg/sessions"
)

type responseLogin struct {
	Success bool `json:"success"`
}

// Login endpoint /auth/login
func Login(w http.ResponseWriter, r *http.Request) {
	if !common.CheckMethod(w, r, "POST") {
		return
	}

	id := r.PostFormValue("id")
	password := r.PostFormValue("password")

	if id == "" || password == "" {
		// Request doesn't contains required informations
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := authstore.GetUserByID(id)

	if user == nil {
		// Utilisateur non trouve
		common.RespondJSON(w, responseLogin{
			Success: false,
		})
		return
	}

	if !sessions.CheckPassword(user.PwdHash, []byte(password)) {
		// Mot de passe incorrect
		common.RespondJSON(w, responseLogin{
			Success: false,
		})
		return
	}

	generatedSID := sessions.GenerateSessionID(SessionIDLength) // 128 bytes long session ids, 1024 bits
	expiration := time.Now().Add(2 * time.Hour)                 // Validite de 2 heures par default

	authstore.AddSession(generatedSID, user.ID, expiration)
	setSessionCookies(w, generatedSID, expiration)

	common.RespondJSON(w, responseLogin{
		Success: true,
	})
}

func createCookie(expiration time.Time) *http.Cookie {
	var cookie *http.Cookie = new(http.Cookie)
	cookie.HttpOnly = true                    // XSS Attack
	cookie.Secure = false                     // HTTPS only // TODO change this
	cookie.SameSite = http.SameSiteStrictMode // CRSF Attack SameSiteLaxMode peut être plus permissif
	cookie.Expires = expiration
	cookie.Path = "/"
	return cookie
}

// TODO CHANGE Secure to true

func setSessionCookies(w http.ResponseWriter, sid []byte, expiration time.Time) {
	MACSignature := sessions.GetMAC(sid, ServerKey)

	cookiesid := createCookie(expiration)
	cookiesid.Name = "sessionid"
	cookiesid.Value = sessions.ToBase64(sid)

	cookiesignature := createCookie(expiration)
	cookiesignature.Name = "signature"
	cookiesignature.Value = sessions.ToBase64(MACSignature)

	http.SetCookie(w, cookiesid)
	http.SetCookie(w, cookiesignature)
}
