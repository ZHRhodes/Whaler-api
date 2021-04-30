package websocket

import (
	"encoding/json"

	"github.com/heroku/whaler-api/OT/ot-master"
)

type SocketMessage struct {
	SenderId string          `json:"senderId"`
	Type     string          `json:"type"`
	Data     json.RawMessage `json:"data"`
}

type DocumentChange struct {
	ResourceId string   `json:"resourceId"`
	Rev        int      `json:"rev"`
	Ops        []*ot.Op `json:"ops"`
}

type ResourceConnection struct {
	ResourceId string `json:"resourceId"`
}

type ResourceConnectionConf struct {
	ResourceId   string `json:"resourceId,"`
	InitialState string `json:"initialState"`
}

// {"type": "docDelta", "data": {"documentID": "1", "value": "Hello World!"}}
