package websocket

import (
	"encoding/json"
)

type SocketMessage struct {
	SenderId  string          `json:"senderId"`
	MessageId string          `json:"messageId"`
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
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
	ResourceId   string `json:"resourceId"`
	InitialState string `json:"initialState"`
	Revision     int    `json:"revision"`
}

type ResourceUpdate struct {
	ResourceId string  `json:"resourceId"`
	SenderId   *string `json:"senderId"`
}
