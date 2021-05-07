package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

var (
	ServerID = uuid.New().String()
)

func HandleNewConnection(id string, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\nReceived websocket connection with id %s", id)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	pool, ok := pools[id]
	log.Println("\npool, ok ", pool, ok)
	if !ok {
		pool = NewPool()
		go pool.Start()
		pools[id] = pool
		log.Println("\nAdded new pool with id ", id)
	}

	clientId := uuid.New().String()
	client := &Client{Id: clientId, pool: pool, conn: conn, send: make(chan SocketMessage, 256)}
	client.pool.register <- client

	go client.startWriting()
	go client.startReading()
}
