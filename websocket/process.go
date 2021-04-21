package websocket

import (
	"encoding/json"
	"fmt"
)

func Process(bytes []byte) error {
	var message SocketMessage
	if err := json.Unmarshal(bytes, &message); err != nil {
		fmt.Println(err)
		return err
	}

	if message.Type == "docChange" {
		var change DocumentChange
		if err := json.Unmarshal(message.Data, &change); err != nil {
			fmt.Println(err)
			return err
		}

		ProcessDocumentChange(change)
	}

	return nil
}

func ProcessDocumentChange(change DocumentChange) {
	fmt.Print(change)
}

//map of connections grouped by docID
//add conn from map
//remove conn from map
//get con from map (if needed)
//send message to conn(ID)
