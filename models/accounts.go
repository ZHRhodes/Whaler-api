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

type BaseModel struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	DeletedAt string `json:"deletedAt"`
}

type Token struct {
	UserID uint
	jwt.StandardClaims
}

type User struct {
	BaseModel
	Email     string `json:"email" gorm:"unique, not null"`
	Password  string `json:"password"`
	Token     string `json:"token" sql:"-"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Account struct {
	BaseModel
	Name        string
	Industry    string
	Description string
}

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (user *User) Validate() map[string]interface{} {
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
	resp := user.Validate()
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

	tk := &Token{UserID: user.ID}
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

	tk := &Token{UserID: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	user.Token = tokenString

	data := map[string]interface{}{"user": user}

	resp := utils.Message(1000, "Logged in", false, data)
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
