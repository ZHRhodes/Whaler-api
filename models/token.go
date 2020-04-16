package models

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AccessToken struct {
	jwt.StandardClaims
	UserID uint
}

type RefreshToken struct {
	UserID uint
	Hash   string
	Exp    time.Time
}

const RefreshTokenValidTime = time.Hour * 72
const AuthTokenValidTime = time.Minute * 15

func CreateRefreshToken(userID uint) string {
	value := [256]byte{}
	_, err := rand.Read(value[:])
	if err != nil {
		fmt.Printf("Failed to generate a random refresh token %q", err)
		panic(err)
	}
	valueEncrypted, _ := bcrypt.GenerateFromPassword(value[:], bcrypt.DefaultCost)
	valueEncoded := base64.StdEncoding.EncodeToString(value[:])
	valueEncryptedAndEncoded := base64.StdEncoding.EncodeToString(valueEncrypted)
	exp := time.Now().Add(AuthTokenValidTime)
	refreshToken := RefreshToken{UserID: userID, Hash: valueEncryptedAndEncoded, Exp: exp}
	refreshToken.StoreRefreshToken()
	return valueEncoded
}

func (token RefreshToken) StoreRefreshToken() {
	DB().Create(token)

	// if len(refreshToken.Value) == 0 {
	// 	fmt.Printf("")
	// }

}
