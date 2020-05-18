package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/heroku/whaler-api/utils"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	DBModel
	Email          string      `json:"email" gorm:"unique, not null"`
	Password       string      `json:"password"`
	FirstName      string      `json:"firstName"`
	LastName       string      `json:"lastName"`
	IsAdmin        bool        `json:"isAdmin"`
	OrganizationID int         `json:"organizationId"`
	Workspaces     []Workspace `json:"workspaces" gorm:"many2many:workspace_user"`
}

//DEPRECATED -- REST
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

	user.Password = ""

	accessTokenString := CreateAccessToken(user.ID)
	refreshTokenString := CreateRefreshToken(user.ID)
	tokens := Tokens{AccessToken: accessTokenString, RefreshToken: refreshTokenString}

	response := utils.MessageWithTokens(2000, "User has been created", false, user, tokens)
	return response
}

func CreateUser(email string, password string) (*User, error) {
	validatationErr := validateUserCreds(email, password)
	if validatationErr != nil {
		return nil, *validatationErr
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &User{Email: email, Password: string(hashedPassword)}

	err := DB().Create(user).Error

	if user.ID <= 0 {
		return nil, err
	}

	user.Password = ""

	return user, nil
}

func LogIn(email, password string) map[string]interface{} {
	user := &User{}
	err := DB().Table("users").Where("email = ?", email).Preload("Workspaces").First(user).Error

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

	accessTokenString := CreateAccessToken(user.ID)
	refreshTokenString := CreateRefreshToken(user.ID)
	tokens := Tokens{AccessToken: accessTokenString, RefreshToken: refreshTokenString}

	resp := utils.MessageWithTokens(1000, "Logged in", false, user, tokens)
	return resp
}

//DEPRECATED -- REST
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

func validateUserCreds(email string, password string) *error {
	if !strings.Contains(email, "@") {
		err := errors.New("Email address is required")
		return &err
	}

	if len(password) < 6 {
		err := errors.New("Password is required")
		return &err
	}

	temp := &User{}

	err := DB().Table("users").Where("email = ?", email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		err := fmt.Errorf("connection error. Please retry. temp: %q", err)
		return &err
	}
	if temp.Email != "" {
		err := errors.New("email address already in use by another user")
		return &err
	}

	return nil
}

func fetchUser(userID int) *User {
	acc := &User{}
	DB().Table("users").Where("id = ?", userID).First(acc)
	if acc.Email == "" {
		return nil
	}

	acc.Password = ""
	return acc

}
