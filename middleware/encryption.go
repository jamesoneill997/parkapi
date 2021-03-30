package middleware

import "golang.org/x/crypto/bcrypt"

//HashPass will receive a password as a string and hash it
func HashPass(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	return string(bytes), err
}

//CheckPassword checks a password vs the hash in the DB, used for user auth
func CheckPassword(pass string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}
