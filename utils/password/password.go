package password

import (
	"golang.org/x/crypto/bcrypt"
)

var SALT string

func InitPasswordVars(salt string) {
	SALT = salt
}

func PasswordToHash(password string) string {
	password += SALT
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func CompareHashPassword(hash, password string) bool {
	password += SALT
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
