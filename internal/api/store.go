package api

import (
	"bytes"
	"time"

	"github.com/robinjulien/rcloud/pkg/enhancedmaps"
)

// AuthStore is the type that stores users and sessions
type AuthStore struct {
	Path     string
	Users    *enhancedmaps.Map
	Sessions []Session
}

// User is the user type as stored in an AuthStore
type User struct {
	ID      string
	PwdHash []byte
	Admin   bool
}

// Session is the session type as stored in an AuthStore
type Session struct {
	SID     []byte
	UID     string
	Expires time.Time
}

// GetUserByID returns a user given its Id in an authstore, or nil if it doesn't exists
func (as *AuthStore) GetUserByID(id string) *User {
	userint, err := as.Users.GetSafe(id)

	if err != nil {
		return nil
	}

	user, ok := userint.(User)

	if !ok {
		return nil
	}

	return &user
}

// GetSessionByID returns a session given its Id in an authstore, or nil if it doesn't exists
func (as *AuthStore) GetSessionByID(sid []byte) *Session {
	for _, s := range as.Sessions {
		if bytes.Equal(s.SID, sid) {
			return &s
		}
	}

	return nil
}

// AddSession adds a session in an authstore
func (as *AuthStore) AddSession(sid []byte, userid string, expires time.Time) {
	session := Session{
		SID:     sid,
		UID:     userid,
		Expires: expires,
	}
	as.Sessions = append(as.Sessions, session)
}

// RemoveSession removes a sesison in an authstore
func (as *AuthStore) RemoveSession(sid []byte) {
	var index int = -1

	for i, s := range as.Sessions {
		if bytes.Equal(s.SID, sid) {
			index = i
			break
		}
	}

	if index != -1 {
		as.Sessions = remove(as.Sessions, index)
	}
}

// remove the i-th element of a session slice
func remove(s []Session, i int) []Session {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
