package models

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/heroku/whaler-api/utils"
	"github.com/jinzhu/gorm"
)

type Token struct {
	UserID uint
	jwt.StandardClaims
}

type User struct {
	gorm.Model
	Email     string `gorm:"unique, not null"`
	Password  string
	Token     string `sql:"-"`
	FirstName string
	LastName  string
}

func (user *User) Validate() (map[string]interface{}, bool) {
	if !strings.Contains(user.Email, "@") {
		return utils.Message(false, "Email address is required"), false
	}

	if len(user.Password) < 6 {
		return utils.Message(false, "Password is required"), false
	}

	temp := &User{}

	err := DB().Table("users").Where("email = ?", user.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		errMessage := fmt.Sprintf("Connection error. Please retry. temp: %q", err)
		return utils.Message(false, errMessage), false
	}
	if temp.Email != "" {
		return utils.Message(false, "Email address already in use by another user."), false
	}

	return utils.Message(false, "Requirement passed"), true
}

func (user *User) Create() map[string]interface{} {
	if resp, ok := user.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	DB().Create(user)

	if user.ID <= 0 {
		fmt.Print(fmt.Sprint("the user id was less than zero"))
		return utils.Message(false, "Failed to create user, connection error.")
	}

	tk := &Token{UserID: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	user.Token = tokenString

	user.Password = ""

	response := utils.Message(true, "User has been created")
	response["user"] = user
	return response
}

func Login(email, password string) map[string]interface{} {
	user := &User{}
	err := DB().Table("users").Where("email = ?", email).First(user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.Message(false, "Email address not found")
		}
		return utils.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return utils.Message(false, "Invalid login credentials. Please try again")
	}
	user.Password = ""

	user.Password = ""

	tk := &Token{UserID: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	user.Token = tokenString

	resp := utils.Message(true, "Logged in")
	resp["user"] = user
	return resp
}

func FindUser(userID uint) *User {
	acc := &User{}
	DB().Table("users").Where("id = ?", userID).First(acc)
	if acc.Email == "" {
		return nil
	}

	acc.Password = ""
	return acc

}
