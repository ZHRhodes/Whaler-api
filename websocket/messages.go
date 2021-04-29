package websocket

import (
	"encoding/json"

	"github.com/heroku/whaler-api/OT/ot-master"
)

type SocketMessage struct {
	SenderId string          `json:"id"`
	Type     string          `json:"type"`
	Data     json.RawMessage `json:"data"`
}

type DocumentChange struct {
	ResourceId string `json:"resourceId"`
	Rev        int    `json:"revision"`
	Ops        ot.Ops `json:"ops"`
}

type ResourceConnection struct {
	ResourceId string `json:"resourceId"`
}

type ResourceConnectionConf struct {
	ResourceId   string   `json:"resourceId,"`
	InitialState string   `json:"initialState"`
	Test         [][]rune `json:"rune"`
}

// {"type": "docDelta", "data": {"documentID": "1", "value": "Hello World!"}}
