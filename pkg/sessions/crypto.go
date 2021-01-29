package sessions

import (
	"crypto/hmac"
	"crypto/sha512"
	"hash"
)

// Defaults
var (
	// HashType is the hash type used by HMAC
	HashType func() hash.Hash = sha512.New
)

// ValidMAC reports whether messageMAC is a valid HMAC tag for message.
// This function is copied from the official Go documentation. This is why there is no error checking.
func ValidMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(HashType, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

// GetMAC return corresponding HMAC tag for message
// This function follows the same syntax as the ValidMAC function which is copied from official Go documentation. This is why there is no error checking.
func GetMAC(message, key []byte) []byte {
	mac := hmac.New(HashType, key)
	mac.Write(message)
	return mac.Sum(nil)
}
