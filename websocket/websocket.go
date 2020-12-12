package websocket

import (
	"log"
	"net/http"
)

func HandleNewConnection(id string, w http.ResponseWriter, r *http.Request) {
	log.Println("\nReceived websocket connection")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	pool, ok := pools[id]
	if !ok {
		pool = NewPool()
		pool.Start()
		pools[id] = pool
		log.Println("\nAdded new pool with id ", id)
	}
	client := &Client{pool: pool, conn: conn, send: make(chan Message, 256)}
	pool.register <- client

	go client.startWriting()
	go client.startReading()
}
