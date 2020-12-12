package websocket

import (
	"encoding/json"
	"fmt"
)

func Process(bytes []byte) error {
	var message Message
	if err := json.Unmarshal(bytes, &message); err != nil {
		fmt.Println(err)
		return err
	}

	if message.Type == "docDelta" {
		var delta DocumentDelta
		if err := json.Unmarshal(message.Data, &delta); err != nil {
			fmt.Println(err)
			return err
		}

		ProcessDocumentDelta(delta)
	}

	return nil
}

func ProcessDocumentDelta(delta DocumentDelta) {
	fmt.Print(delta)
}

//map of connections grouped by docID
//add conn from map
//remove conn from map
//get con from map (if needed)
//send message to conn(ID)

//first iteration.. everyone is sharing the same note!
