package websocket

import "encoding/json"

type SocketMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type DocumentDelta struct {
	value string `json:"value"`
}

// {"type": "docDelta", "data": {"documentID": "1", "value": "Hello World!"}}
