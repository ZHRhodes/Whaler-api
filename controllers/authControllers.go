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

	resp := models.LogIn(user.Email, user.Password)
	utils.Respond(w, resp)
}

var Refresh = func(w http.ResponseWriter, r *http.Request) {
	tokens := &models.Tokens{}
	err := json.NewDecoder(r.Body).Decode(tokens)
	if err != nil {
		utils.Respond(w, utils.Message(4000, "Invalid request - malformed refresh token", true, map[string]interface{}{}))
		return
	}

	userID := r.Context().Value("userID").(uint)

	resp := models.Refresh(tokens.RefreshToken, userID)
	utils.Respond(w, resp)
}

var LogOut = func(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	go models.InvalidateTokens(userID)
	utils.Respond(w, utils.Message(2000, "success", false, map[string]interface{}{}))
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

var CreateWorkspace = func(w http.ResponseWriter, r *http.Request) {
	workspace := &models.Workspace{}
	err := json.NewDecoder(r.Body).Decode(workspace)
	if err != nil {
		utils.Respond(w, utils.Message(4000, "Invalid request - malformed workspace", true, map[string]interface{}{}))
		return
	}

	resp := workspace.Create()
	utils.Respond(w, resp)
}

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		utils.Respond(w, utils.Message(4000, "Invalid request - malformed workspace", true, map[string]interface{}{}))
		return
	}

	resp := account.Create()
	utils.Respond(w, resp)
}

var CreateContact = func(w http.ResponseWriter, r *http.Request) {
	contact := &models.Contact{}
	err := json.NewDecoder(r.Body).Decode(contact)
	if err != nil {
		utils.Respond(w, utils.Message(4000, "Invalid request - malformed contact", true, map[string]interface{}{}))
		return
	}

	resp := contact.Create()
	utils.Respond(w, resp)
}

var FetchOrg = func(w http.ResponseWriter, r *http.Request) {
	orgID := r.URL.Query().Get("id")
	resp := models.FetchOrg(orgID)
	utils.Respond(w, resp)
}

var FetchWorkspace = func(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.URL.Query().Get("workspaceID")
	resp := models.FetchAccounts(workspaceID)
	utils.Respond(w, resp)
}
