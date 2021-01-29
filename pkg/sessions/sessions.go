package sessions

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// GenerateSessionID generates cryptographically secure pseudorandom session id and stores it in a byte slice of given size.
// size argument should be big enough (eg >= 256 bits, ie size >= 32)
func GenerateSessionID(size int) []byte {
	var sid []byte = make([]byte, size, size)

	n, err := rand.Read(sid)

	if n != len(sid) || err != nil {
		panic(err)
	}

	return sid
}

// GeneratePwdHash returns the hash of the password using bcrypt, and a boolean telling if there are an error
func GeneratePwdHash(password []byte) ([]byte, bool) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	return hash, err == nil
}

// CheckPassword returns if the password stored in the database with a hash match the given password at login
func CheckPassword(hash, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}

// ToBase64 converts raw bytes to base64 encoding to fit in cookies
func ToBase64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

// FromBase64 converts base64 string to raw bytes
func FromBase64(s string) []byte {
	res, _ := base64.StdEncoding.DecodeString(s)

	return res
}
