package websocket

import (
	"encoding/json"

	"github.com/google/uuid"
)

type ChangeConsumer struct{}

func (c ChangeConsumer) ModelChanged(id string, senderId *string) {
	resourceUpdate := ResourceUpdate{ResourceId: id, SenderId: senderId}
	bytes, _ := json.Marshal(resourceUpdate)
	message := SocketMessage{SenderId: ServerID, MessageId: uuid.New().String(), Type: "resourceUpdate", Data: bytes}
	if pool := pools[id]; pool != nil {
		pool.broadcast <- message
	}
}
