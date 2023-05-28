// Password Manager

package users

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}
	return bytes, nil

}
