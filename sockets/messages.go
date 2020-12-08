package sockets

import "encoding/json"

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type DocumentDelta struct {
	DocumentID string `json:"documentID"`
	Value      string `json:"value"`
}

// {"type": "docDelta", "data": {"documentID": "1", "value": "Hello World!"}}
