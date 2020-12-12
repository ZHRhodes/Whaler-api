package websocket

import (
	"nhooyr.io/websocket"
)

var registry = make([]websocket.Conn, 0, 10)

func Register(conn websocket.Conn) error {
	registry = append(registry, conn)
	return nil
}

//Next up: look into othat warning ^
//https://eli.thegreenplace.net/2018/beware-of-copying-mutexes-in-go/
