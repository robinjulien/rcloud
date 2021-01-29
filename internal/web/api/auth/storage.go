package auth

import (
	"bytes"
	"time"

	"github.com/robinjulien/rcloud/pkg/enhancedmaps"
)

type AuthStore struct {
	Users    *enhancedmaps.Map
	Sessions []Session
}

type User struct {
	ID      string
	PwdHash []byte
	Admin   bool
}

type Session struct {
	SID     []byte
	UID     string
	Expires time.Time
}

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

func (as *AuthStore) GetSessionByID(sid []byte) *Session {
	for _, s := range as.Sessions {
		if bytes.Equal(s.SID, sid) {
			return &s
		}
	}

	return nil
}

func (as *AuthStore) AddSession(sid []byte, userid string, expires time.Time) {
	session := Session{
		SID:     sid,
		UID:     userid,
		Expires: expires,
	}
	as.Sessions = append(as.Sessions, session)
}

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

func remove(s []Session, i int) []Session {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (as *AuthStore) CreateUser(u User) {
	as.Users.Set(u.ID, u)
}

func (as *AuthStore) RemoveUserByID(uid string) {
	as.Users.Remove(uid)
}

func (as *AuthStore) EditUser(uid string, u User) {
	if uid != u.ID {
		return
	}

	as.Users.Set(u.ID, u)
}
