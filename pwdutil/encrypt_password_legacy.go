package pwdutil

import (
	"golang.org/x/crypto/scrypt"
)

func EncryptPasswordLegacy(password, salt string) (string, error) {
	dk, err := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, 32)
	if err != nil {
		return "", err
	}
	return string(dk), nil
}
