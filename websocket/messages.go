package websocket

import (
	"encoding/json"
)

type SocketMessage struct {
	SenderId string          `json:"senderId"`
	Type     string          `json:"type"`
	Data     json.RawMessage `json:"data"`
}

type DocumentChange struct {
	ResourceId string   `json:"resourceId"`
	Rev        int      `json:"rev"`
	N          []int    `json:"n"`
	S          []string `json:"s"`
}

type DocumentChangeReturn struct {
	ResoureceId string   `json:"resourceId"`
	N           []int    `json:"n"`
	S           []string `json:"s"`
}

type ResourceConnection struct {
	ResourceId string `json:"resourceId"`
}

type ResourceConnectionConf struct {
	ResourceId   string `json:"resourceId,"`
	InitialState string `json:"initialState"`
}

// {"type": "docDelta", "data": {"documentID": "1", "value": "Hello World!"}}
