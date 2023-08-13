package api

import (
	"log"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type SignedClaims struct {
	Email     string
	FirstName string
	LastName  string
	UserType  string
	Uid       string
	jwt.RegisteredClaims
}

var SECRET_KEY = "secret"

func HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(hash)
}

func VerifyPassword(user_password, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(user_password))
	return err == nil
}

func GenerateTokens(email, first_name, last_name, user_id string) (token, refresh_token string) {
	var err error
	claim := &SignedClaims{
		Email:     email,
		FirstName: first_name,
		LastName:  last_name,
		Uid:       user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	refresh_claim := &SignedClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Fatal(err)
		return
	}
	refresh_token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refresh_claim).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Fatal(err)
		return
	}
	return token, refresh_token
}
