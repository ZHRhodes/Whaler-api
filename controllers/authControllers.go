package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/heroku/whaler-api/middleware"
	"github.com/heroku/whaler-api/models"
	"github.com/heroku/whaler-api/websocket"
	"github.com/heroku/whaler-api/utils"
)

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

	userID := middleware.UserIDFromContext(r.Context())

	resp := models.Refresh(tokens.RefreshToken, userID)
	utils.Respond(w, resp)
}

var LogOut = func(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	go models.InvalidateTokens(userID)
	utils.Respond(w, utils.Message(2000, "success", false, map[string]interface{}{}))
}

var Socket = func(w http.ResponseWriter, r *http.Request) {
	pool := websocket.NewPool()
	go pool.Start()
	websocket.HandleNewConnection(pool, w, r)
}