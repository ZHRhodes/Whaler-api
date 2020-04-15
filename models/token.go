package models

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AccessToken struct {
	jwt.StandardClaims
	UserID uint
}

type RefreshToken struct {
	UserID uint
	Value  string
	Exp    time.Time
}

const RefreshTokenValidTime = time.Hour * 72
const AuthTokenValidTime = time.Minute * 15

func CreateRefreshToken(userID uint) RefreshToken {
	value := [256]byte{}
	_, err := rand.Read(value[:])
	if err != nil {
		fmt.Printf("Failed to generate a random refresh token %q", err)
		panic(err)
	}
	valueEncoded := base64.StdEncoding.EncodeToString(value[:])
	exp := time.Now().Add(AuthTokenValidTime)
	return RefreshToken{UserID: userID, Value: valueEncoded, Exp: exp}
}

func StoreRefreshToken(refreshToken RefreshToken) {
	DB().Create(refreshToken)

	if len(refreshToken.Value) == 0 {
		fmt.Printf("")
	}

}
