package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/heroku/whaler-api/middleware"
	"github.com/heroku/whaler-api/models"
	"github.com/heroku/whaler-api/utils"
	"nhooyr.io/websocket"
)

type DocumentDelta struct {
	Type       string `json:"type"`
	DocumentID string `json:"documentID"`
	Value      string `json:"value"`
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
	fmt.Print("Request received at /socket\n")
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{},
	})
	if err != nil {
		fmt.Print("Failed to accept web socket: ", err, "\n")
		return
	}
	fmt.Print("-->socket accepted!")
	defer conn.Close(websocket.StatusInternalError, "the sky is falling")

	// // if c.Subprotocol() != "echo" {
	// // 	c.Close(websocket.StatusPolicyViolation, "the client must speak the echo subprotocol")
	// // }

	for {
		err = echo(r.Context(), conn)
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
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	_, ioReader, err := conn.Reader(ctx)
	if err != nil {
		return err
	}

	//Write it back
	//
	// writer, err := conn.Writer(ctx, messageType)
	// if err != nil {
	// 	return err
	// }
	//
	// _, err = io.Copy(writer, ioReader)
	//
	// if err != nil {
	// 	return fmt.Errorf("failed to io.Copy: %w", err)
	//}
	// err = writer.Close()

	//parse message as json
	bytes, err := ioutil.ReadAll(ioReader)
	if err != nil {
		return err
	}

	var delta DocumentDelta
	if err := json.Unmarshal(bytes, &delta); err != nil {
		fmt.Print(err)
		return err
	}

	fmt.Print(delta)

	return err
}
