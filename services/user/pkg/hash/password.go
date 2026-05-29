package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) string {
	salt := make([]byte, 16)
	rand.Read(salt)

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// guarda salt + hash juntos
	return fmt.Sprintf("%x:%x", salt, hash)
}

func VerifyPassword(password, hashed string) bool {
	parts := strings.Split(hashed, ":")
	salt, _ := hex.DecodeString(parts[0])
	expected, _ := hex.DecodeString(parts[1])

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return subtle.ConstantTimeCompare(hash, expected) == 1
}
