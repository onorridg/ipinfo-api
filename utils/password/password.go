package password

import (
	"golang.org/x/crypto/bcrypt"
)

const SALT = `t^BjFhr3jbtwDNm+azddi=ndp7sZc,4pM%-]aQyQsWpfVu]wE6>sn5cMHhq_DE0?`

func PasswordToHash(password string)string{
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