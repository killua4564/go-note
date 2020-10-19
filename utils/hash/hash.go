package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/pbkdf2"

	"github.com/killua4564/go-note/config"
)

func HS256(password string) string {
	salt := []byte(config.Salt())
	h := hmac.New(sha256.New, salt)
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func PBKDF2(password string) string {
	pwd := []byte(password)
	salt := []byte(config.Salt())
	keyLen := 32
	iter := 2000
	digest := sha256.New

	dk := pbkdf2.Key(pwd, salt, iter, keyLen, digest)
	return hex.EncodeToString(dk)
}
