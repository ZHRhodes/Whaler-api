package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/heroku/whaler-api/models"
	"github.com/heroku/whaler-api/utils"
)

//CreateUser creates a user in the database and returns it
var CreateUser = func(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		utils.Respond(w, utils.Message(4000, "Invalid request - malformed user", true, map[string]interface{}{}))
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
		utils.Respond(w, utils.Message(4000, "Invalid request - malformed user", true, map[string]interface{}{}))
		return
	}

	resp := models.Login(user.Email, user.Password)
	utils.Respond(w, resp)
}

//CreateOrg creates an org in the database and returns it
var CreateOrg = func(w http.ResponseWriter, r *http.Request) {
	org := &models.Organization{}
	err := json.NewDecoder(r.Body).Decode(org)
	if err != nil {
		utils.Respond(w, utils.Message(4000, "Invalid request - malformed organziation", true, map[string]interface{}{}))
		return
	}

	resp := org.Create()
	utils.Respond(w, resp)
}

var FetchOrg = func(w http.ResponseWriter, r *http.Request) {
	orgID := r.URL.Query().Get("id")
	resp := models.FetchOrg(orgID)
	utils.Respond(w, resp)
}
