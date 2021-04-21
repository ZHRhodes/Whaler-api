package websocket

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func HandleNewConnection(id string, w http.ResponseWriter, r *http.Request) {
	log.Println("\nReceived websocket connection")
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
	client := &Client{id: clientId, pool: pool, conn: conn, send: make(chan SocketMessage, 256)}
	client.pool.register <- client

	go client.startWriting()
	go client.startReading()
}
