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
	DBModel
	UserID uint
	Hash   string
	Exp    time.Time
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
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

func Retrieve(refreshTokenString string) (*RefreshToken, error) {
	refreshToken := &RefreshToken{}
	err := DB().Table("refresh_tokens").Where("hash = ?", refreshTokenString).First(refreshToken).Error

	if err != nil {
		fmt.Printf(err.Error())
	}

	return refreshToken, err
}

func (token RefreshToken) Validate(userId uint) bool {
	isExpired := token.Exp.Before(time.Now())
	if isExpired {
		return false
	}

	idsMismatch := token.UserID != userId
	if idsMismatch {
		return false
	}

	return true
}

func (token RefreshToken) StoreRefreshToken() {
	err := DB().Create(token).Error

	if err != nil {
		fmt.Printf("Failed to create refresh token in DB")
	}
}

func Refresh() {

}

func (token RefreshToken) Invalidate() {
	token.Exp = time.Now()
	err := DB().Save(&token).Error

	if err != nil {
		fmt.Printf("Failed to create refresh token in DB")
	}
}
