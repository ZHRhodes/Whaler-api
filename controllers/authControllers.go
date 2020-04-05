package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/heroku/whaler-api/models"
	"github.com/heroku/whaler-api/utils"
)

//CreateUser creates a user on the backend and returns it
var CreateUser = func(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid request"))
		return
	}

	resp := user.Create()
	utils.Respond(w, resp)
}

//Authenticate logs into the user
var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(user.Email, user.Password)
	utils.Respond(w, resp)
}
