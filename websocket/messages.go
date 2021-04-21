package websocket

import "encoding/json"

type SocketMessage struct {
	Id   string          `json:"id"`
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type DocumentChange struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	Range string `json:"range"`
}

// {"type": "docDelta", "data": {"documentID": "1", "value": "Hello World!"}}
