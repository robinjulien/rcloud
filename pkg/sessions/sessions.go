package sessions

import "crypto/rand"

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
