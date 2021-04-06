package models

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/heroku/whaler-api/utils"
)

type AccessToken struct {
	jwt.StandardClaims
	UserID       string                 `json:"userID"`
	HasuraClaims map[string]interface{} `json:"https://hasura.io/jwt/claims"`
}

type RefreshToken struct {
	DBModel
	UserID string
	Hash   string
	Exp    time.Time
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

const RefreshTokenValidTime = time.Hour * 24 * 90
const AuthTokenValidTime = time.Minute * 1

func CreateAccessToken(userID string) string {
	exp := time.Now().Add(AuthTokenValidTime).Unix()
	allowedRoles := [3]string{"editor", "user", "mod"}
	hasuraClaims := map[string]interface{}{
		"x-hasura-allowed-roles": allowedRoles,
		"x-hasura-default-role":  "user",
		"x-hasura-user-id":       userID,
	}
	standardClaims := jwt.StandardClaims{ExpiresAt: exp, Audience: "hasura", Issuer: "getwhaler-auth"}
	tk := &AccessToken{UserID: userID, HasuraClaims: hasuraClaims, StandardClaims: standardClaims}
	accessToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	accessTokenString, _ := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return accessTokenString
}

func CreateRefreshToken(userID string) string {
	value := [256]byte{}
	_, err := rand.Read(value[:])
	if err != nil {
		fmt.Printf("Failed to generate a random refresh token %q\n", err)
		panic(err)
	}

	tokenValue := base64.StdEncoding.EncodeToString(value[:])
	exp := time.Now().Add(RefreshTokenValidTime)
	refreshToken := RefreshToken{UserID: userID, Hash: tokenValue, Exp: exp}

	refreshToken.store()

	return tokenValue
}

func Retrieve(refreshTokenString string) (*RefreshToken, error) {
	refreshToken := &RefreshToken{}
	err := DB().Table("refresh_tokens").Where("hash = ?", refreshTokenString).First(refreshToken).Error

	if err != nil {
		fmt.Printf(err.Error())
	}

	return refreshToken, err
}

func (token RefreshToken) Validate(userId string) bool {
	isExpired := token.Exp.Before(time.Now())
	if isExpired {
		fmt.Printf("Refresh token is expired\n")
		return false
	}

	idsMismatch := token.UserID != userId
	if idsMismatch {
		fmt.Printf("Refresh requested by the wrong user\n")
		return false
	}

	return true
}

func (token *RefreshToken) store() {
	err := DB().Create(token).Error

	if err != nil {
		fmt.Printf("Failed to create refresh token in DB -- %q\n", err)
	}
}

func Refresh(refreshTokenString string, userID string) map[string]interface{} {
	refreshToken, err := Retrieve(refreshTokenString)

	if err != nil {
		fmt.Printf("Failed to retrieve RefreshToken from DB\n")
		return utils.Message(5001, "Unable to retrieve refresh token, connection error", true, map[string]interface{}{})
	}

	isTokenValid := refreshToken.Validate(userID)

	if !isTokenValid {
		fmt.Printf("Refresh token invalid\n")
		return utils.Message(4004, "Refresh token invalid", true, map[string]interface{}{})
	}

	accessTokenString := CreateAccessToken(userID)
	tokens := Tokens{AccessToken: accessTokenString, RefreshToken: refreshTokenString}
	user := FetchUser(userID)
	resp := utils.MessageWithTokens(1000, "Tokens refreshed", false, user, tokens)
	fmt.Println(fmt.Sprintf("New Token %s", tokens.AccessToken))
	return resp
}

func (token *RefreshToken) Invalidate() {
	token.Exp = time.Now()
	err := DB().Save(token).Error

	if err != nil {
		fmt.Printf("Failed to incalidate refresh token in DB\n")
	}
}

func InvalidateTokens(userID int) {
	//this will re-invalidate all tokens a user has ever had
	//to negate the perf impact, should either remove revoked tokens from db
	//or invalidate the specific refresh token only
	DB().Table("refresh_tokens").Where("user_id == ?", userID).Updates(map[string]interface{}{"Exp": time.Now()})
}
