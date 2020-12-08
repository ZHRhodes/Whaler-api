package sockets

import (
	"encoding/json"
	"fmt"
)

func Retrieve(bytes []byte) error {
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
