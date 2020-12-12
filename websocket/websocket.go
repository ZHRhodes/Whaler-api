package websocket

import (
	"log"
	"net/http"
)

func HandleNewConnection(id string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	pool, ok := pools[id]
	if !ok {
		pool = NewPool()
		pools[id] = pool
	}
	client := &Client{pool: pool, conn: conn, send: make(chan Message, 256)}
	pool.register <- client

	go client.startWriting()
	go client.startReading()
}
