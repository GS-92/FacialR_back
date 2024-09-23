package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(hash, password string) bool {
	plain := []byte(password)
	err := bcrypt.CompareHashAndPassword([]byte(hash), plain)
	if err != nil{
		fmt.Println("No match")
	}
	return err == nil
}
