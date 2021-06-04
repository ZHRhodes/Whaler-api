package websocket

import (
	"encoding/json"

	"github.com/google/uuid"
)

type ChangeConsumer struct{}

func (c ChangeConsumer) ModelChanged(id string) {
	resourceUpdate := ResourceUpdate{ResourceId: id}
	bytes, _ := json.Marshal(resourceUpdate)
	message := SocketMessage{SenderId: ServerID, MessageId: uuid.New().String(), Type: "resourceUpdate", Data: bytes}
	pools[id].broadcast <- message
}
