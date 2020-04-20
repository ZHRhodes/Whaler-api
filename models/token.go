package models

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/heroku/whaler-api/utils"
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
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

const RefreshTokenValidTime = time.Hour * 72
const AuthTokenValidTime = time.Minute * 15

func CreateAccessToken(userID uint) string {
	exp := time.Now().Add(AuthTokenValidTime).Unix()
	tk := &AccessToken{UserID: userID, StandardClaims: jwt.StandardClaims{ExpiresAt: exp}}
	accessToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	accessTokenString, _ := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return accessTokenString
}

func CreateRefreshToken(userID uint) string {
	value := [256]byte{}
	_, err := rand.Read(value[:])
	if err != nil {
		fmt.Printf("Failed to generate a random refresh token %q", err)
		panic(err)
	}

	//base64 encode the random string
	valueEncoded := base64.StdEncoding.EncodeToString(value[:])

	//encrypt the random string
	valueEncrypted, _ := bcrypt.GenerateFromPassword(value[:], bcrypt.DefaultCost)

	//base64 encode the encrypted string (from step 1)
	valueEncryptedAndEncoded := base64.StdEncoding.EncodeToString(valueEncrypted)

	exp := time.Now().Add(RefreshTokenValidTime)

	//set the has to the encrypted and encoded value (from step 2)
	refreshToken := RefreshToken{UserID: userID, Hash: valueEncryptedAndEncoded, Exp: exp}

	refreshToken.StoreRefreshToken()

	//return the string from step 3
	return valueEncoded
}

func Retrieve(refreshTokenString string) (*RefreshToken, error) {
	refreshTokenStringBytes := []byte(refreshTokenString)
	refreshTokenEncrypted, _ := bcrypt.GenerateFromPassword(refreshTokenStringBytes, bcrypt.DefaultCost)
	hash := base64.StdEncoding.EncodeToString(refreshTokenEncrypted)

	refreshToken := &RefreshToken{}
	fmt.Printf("Fetching refresh token from DB")
	err := DB().Table("refresh_tokens").Where("hash = ?", hash).First(refreshToken).Error
	fmt.Printf("fetched refresh token with error: %q", err)
	if err != nil {
		fmt.Printf(err.Error())
	}

	return refreshToken, err
}

func (token RefreshToken) Validate(userId uint) bool {
	isExpired := token.Exp.Before(time.Now())
	if isExpired {
		fmt.Printf("Refresh token is expired")
		return false
	}

	idsMismatch := token.UserID != userId
	if idsMismatch {
		fmt.Printf("Refresh requested by the wrong user")
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

func Refresh(refreshTokenString string, userID uint) map[string]interface{} {
	refreshToken, err := Retrieve(refreshTokenString)

	if err != nil {
		fmt.Printf("Failed to retrieve RefreshToken from DB")
		return utils.Message(5001, "Unable to retrieve refresh token, connection error", true, map[string]interface{}{})
	}

	isTokenValid := refreshToken.Validate(userID)

	if !isTokenValid {
		fmt.Printf("Refresh token invalid")
		return utils.Message(4004, "Refresh token invalid", true, map[string]interface{}{})
	}

	accessTokenString := CreateAccessToken(userID)
	tokens := Tokens{AccessToken: accessTokenString, RefreshToken: refreshTokenString}
	data := map[string]interface{}{"tokens": tokens}
	resp := utils.Message(1000, "Token refreshed", false, data)
	return resp
}

func (token RefreshToken) Invalidate() {
	token.Exp = time.Now()
	err := DB().Save(&token).Error

	if err != nil {
		fmt.Printf("Failed to create refresh token in DB")
	}
}
