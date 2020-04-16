package models

import (
	"fmt"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/heroku/whaler-api/utils"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	DBModel
	Email          string `json:"email" gorm:"unique, not null"`
	Password       string `json:"password"`
	Token          string `json:"token" sql:"-"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	IsAdmin        bool   `json:"isAdmin"`
	OrganizationID uint   `json:"organizationId"`
}

func (user *User) validate() map[string]interface{} {
	if !strings.Contains(user.Email, "@") {
		return utils.Message(4001, "Email address is required", true, map[string]interface{}{})
	}

	if len(user.Password) < 6 {
		return utils.Message(4001, "Password is required", true, map[string]interface{}{})
	}

	temp := &User{}

	err := DB().Table("users").Where("email = ?", user.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		errMessage := fmt.Sprintf("Connection error. Please retry. temp: %q", err)
		return utils.Message(4001, errMessage, true, map[string]interface{}{})
	}
	if temp.Email != "" {
		return utils.Message(4001, "Email address already in use by another user.", true, map[string]interface{}{})
	}

	return utils.Message(4001, "Requirement passed", false, map[string]interface{}{})
}

func (user *User) Create() map[string]interface{} {
	resp := user.validate()
	if resp["hasError"] == true {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	DB().Create(user)

	if user.ID <= 0 {
		fmt.Print(fmt.Sprint("the user id was less than zero"))
		return utils.Message(5001, "Failed to create user, connection error.", true, map[string]interface{}{})
	}

	tk := &AccessToken{UserID: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	user.Token = tokenString

	user.Password = ""

	data := map[string]interface{}{"user": user}
	response := utils.Message(2000, "User has been created", false, data)
	return response
}

func Login(email, password string) map[string]interface{} {
	user := &User{}
	err := DB().Table("users").Where("email = ?", email).First(user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.Message(4001, "Email address not found", true, map[string]interface{}{})
		}
		return utils.Message(5001, "Connection error", true, map[string]interface{}{})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return utils.Message(4001, "Invalid login credentials", true, map[string]interface{}{})
	}
	user.Password = ""

	user.Password = ""

	tk := &AccessToken{UserID: user.ID}
	accessToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	accessTokenString, _ := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	user.Token = accessTokenString

	refreshTokenString := CreateRefreshToken(user.ID)
	data := map[string]interface{}{"user": user, "refreshToken": refreshTokenString}

	resp := utils.Message(1000, "Logged in", false, data)
	return resp
}

func FetchUser(userID uint) *User {
	acc := &User{}
	DB().Table("users").Where("id = ?", userID).First(acc)
	if acc.Email == "" {
		return nil
	}

	acc.Password = ""
	return acc

}
