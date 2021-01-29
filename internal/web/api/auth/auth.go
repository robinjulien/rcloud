package auth

import (
	"encoding/gob"
	"errors"
	"net/http"
	"os"

	"github.com/robinjulien/rcloud/pkg/enhancedmaps"
	"github.com/robinjulien/rcloud/pkg/sessions"
)

var (
	// SessionIDLength is default session id length in bytes
	SessionIDLength int = 128 // 1024 bits

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

	if len(os.Args) >= 3 {
		err := authstore.Users.ReadFile(os.Args[2]) // os.Args[2] is <database directory>

		if err != nil {
			if errors.Is(err, enhancedmaps.ErrorFileNotExist) {
				// File not exists or not having rights to read
				err2 := authstore.Users.WriteFile(os.Args[2])

				if err2 != nil {
					panic(err2)
				}
			} else {
				panic(err)
			}
		}
	}

}

// Handler returns auth module handler
func Handler() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/login", Login)
	router.HandleFunc("/logout", Logout)
	router.HandleFunc("/amiloggedin", AmILoggedIn)
	router.HandleFunc("/whoami", WhoAmI)
	return router
}

// GetUserByCookies returns a user given sessions cookies (by request)
// If an error occurs, it responds to request with corresponding http status and returns nil
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

	if session == nil {
		return nil
	}

	user := authstore.GetUserByID(session.UID)

	return user
}
