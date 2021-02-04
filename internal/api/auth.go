package api

import (
	"encoding/gob"
	"errors"
	"net/http"
	"time"

	"github.com/robinjulien/rcloud/pkg/enhancedmaps"
	"github.com/robinjulien/rcloud/pkg/sessions"
)

const (
	// SessionIDLength is default session id length in bytes
	SessionIDLength int = 128 // 1024 bits
)

// ContextKey is type of keys used n context.WithValue
type ContextKey int

const (
	// ContextKeyUser is the key of a user
	ContextKeyUser ContextKey = iota
)

var (
	// ServerKey is the key used to get the mac
	ServerKey []byte

	authstore AuthStore
)

func init() {
	ServerKey = sessions.GenerateSessionID(32)

	authstore = AuthStore{
		Users:    enhancedmaps.New(),
		Sessions: make([]Session, 0),
	}

	// Default user
	hash, _ := sessions.GeneratePwdHash([]byte("admin"))
	u := User{ID: "admin", PwdHash: hash, Admin: true}

	authstore.Users = enhancedmaps.New()
	authstore.Users.Set("admin", u)
	gob.Register(User{})
}

// AuthHandler returns auth module handler
func AuthHandler() http.Handler {
	router := http.NewServeMux()
	router.Handle("/login", MethodMiddleware("POST", http.HandlerFunc(Login)))
	router.Handle("/logout", MethodMiddleware("POST", http.HandlerFunc(Logout)))
	router.Handle("/amiloggedin", MethodMiddleware("GET", http.HandlerFunc(AmILoggedIn)))
	router.Handle("/whoami", MethodMiddleware("GET", http.HandlerFunc(WhoAmI)))
	return router
}

type responseLogin struct {
	Success bool `json:"success"`
}

// Login endpoint /auth/login
func Login(w http.ResponseWriter, r *http.Request) {
	// Everytime one try to log in, all expired sessions are removed from the store
	for i, s := range authstore.Sessions {
		if s.Expires.Before(time.Now()) {
			//authstore.RemoveSession(s.SID)
			authstore.Sessions = remove(authstore.Sessions, i)
		}
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
		RespondJSON(w, responseLogin{
			Success: false,
		})
		return
	}

	if !sessions.CheckPassword(user.PwdHash, []byte(password)) {
		// Mot de passe incorrect
		RespondJSON(w, responseLogin{
			Success: false,
		})
		return
	}

	generatedSID := sessions.GenerateSessionID(SessionIDLength) // 128 bytes long session ids, 1024 bits
	expiration := time.Now().Add(2 * time.Hour)                 // Validite de 2 heures par default

	authstore.AddSession(generatedSID, user.ID, expiration)
	setSessionCookies(w, generatedSID, expiration)

	RespondJSON(w, responseLogin{
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

// Logout /auth/logout endpoint to logout user
func Logout(w http.ResponseWriter, r *http.Request) {
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

type responseAmILoggedIn struct {
	LoggedIn bool `json:"loggedIn"`
}

// AmILoggedIn /auth/amiloggedin tells wheras the given sessionid is valid or not
func AmILoggedIn(w http.ResponseWriter, r *http.Request) {
	u := GetUserByCookies(r)

	if u == nil {
		RespondJSON(w, responseAmILoggedIn{
			LoggedIn: false,
		})
		return
	}

	RespondJSON(w, responseAmILoggedIn{
		LoggedIn: true,
	})
}

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
	u := GetUserByCookies(r)
	if u == nil {
		RespondJSON(w, responseWhoAmI{LoggedIn: false, PublicUser: PublicUser{}})
		return
	}

	RespondJSON(w, responseWhoAmI{LoggedIn: true, PublicUser: PublicUser{
		ID:    u.ID,
		Admin: u.Admin,
	}})
}
