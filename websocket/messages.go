package websocket

import "encoding/json"

type SocketMessage struct {
	SenderId string          `json:"id"`
	Type     string          `json:"type"`
	Data     json.RawMessage `json:"data"`
}

type DocumentChange struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	Range string `json:"range"`
}

type ResourceConnection struct {
	ResourceId string `json:"resourceId"`
}

type ResourceConnectionConf struct {
	ResourceId   string `json:"resourceId,"`
	InitialState string `json:"initialState"`
}

// {"type": "docDelta", "data": {"documentID": "1", "value": "Hello World!"}}
