package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/heroku/whaler-api/middleware"
	"github.com/heroku/whaler-api/models"
	"github.com/heroku/whaler-api/utils"
	"nhooyr.io/websocket"
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
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{"echo"},
	})
	if err != nil {
		fmt.Print("Failed to accept web socket")
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	if c.Subprotocol() != "echo" {
		c.Close(websocket.StatusPolicyViolation, "the client must speak the echo subprotocol")
	}

	for {
		err = echo(r.Context(), c)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			fmt.Print(fmt.Sprintf("failed to echo with %v: %v", r.RemoteAddr, err))
			return
		}
	}
}

func echo(ctx context.Context, conn *websocket.Conn) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	messageType, ioReader, err := conn.Reader(ctx)
	if err != nil {
		return err
	}

	writeCloser, err := conn.Writer(ctx, messageType)
	if err != nil {
		return err
	}

	_, err = io.Copy(writeCloser, ioReader)
	if err != nil {
		return fmt.Errorf("failed to io.Copy: %w", err)
	}

	err = writeCloser.Close()
	return err
}
