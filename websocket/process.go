package websocket

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/heroku/whaler-api/OT/ot-master"
	"github.com/heroku/whaler-api/models"
)

func Process(message SocketMessage, client *Client) error {
	fmt.Println("Processing socket message...")

	if message.Type == "docChange" {
		return processDocChange(message, client)
	} else if message.Type == "resourceConnection" {
		return processResourceConnection(message, client)
	}

	return nil
}

func processDocChange(message SocketMessage, client *Client) error {
	var change DocumentChange
	if err := json.Unmarshal(message.Data, &change); err != nil {
		fmt.Println(err)
		if err, ok := err.(*json.UnmarshalTypeError); ok {
			line, character, lcErr := lineAndCharacter(j, int(jsonError.Offset))
			fmt.Fprintf(os.Stderr, "test %d failed with error: The JSON type '%v' cannot be converted into the Go '%v' type on struct '%s', field '%v'. See input file line %d, character %d\n", i+1, jsonError.Value, jsonError.Type.Name(), jsonError.Struct, jsonError.Field, line, character)
			if lcErr != nil {
				fmt.Fprintf(os.Stderr, "test %d failed with error: Couldn't find the line and character position of the error due to error %v\n", i+1, lcErr)
			}
			continue
		}
		return err
	}

	return nil
	// doc := ot.ServerDocs[change.ResourceId]
	// ops, err := doc.Recv(change.Rev, change.Ops)
	// if err != nil {
	// 	fmt.Printf("\nFailed sending changes to doc. %s", err)
	// 	return err
	// }
	// err = returnOps(client, change.ResourceId, ops)
	// return err
}

func returnOps(client *Client, resourceId string, ops ot.Ops) error {
	message := struct {
		ResoureceId string `json:"resourceId"`
		Ops         ot.Ops `json:"ops"`
	}{ResoureceId: resourceId, Ops: ops}

	bytes, err := json.Marshal(message)

	if err != nil {
		fmt.Println("\nFailed to marshal return ops into bytes")
		return err
	}

	sendMessage(bytes, ServerID, "docChangeReturnOps", client)
	return nil
}

func processResourceConnection(message SocketMessage, client *Client) error {
	var request ResourceConnection
	if err := json.Unmarshal(message.Data, &request); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(request)

	//Resource connections are for notes only for now
	note, err := models.FetchNote("", request.ResourceId)

	if err != nil {
		fmt.Printf("\nFailed to fetch note with id %s", request.ResourceId)
		return err
	}

	doc := ot.NewDocFromStr(note.Content)
	serverDoc := ot.ServerDoc{Doc: doc, History: []ot.Ops{}}
	ot.ServerDocs[request.ResourceId] = serverDoc
	sendResourceConnectionConfirmation(request.ResourceId, note.Content, client)
	return nil
}

func sendResourceConnectionConfirmation(resourceId string, initialState string, client *Client) {
	conf := ResourceConnectionConf{ResourceId: resourceId, InitialState: initialState}
	bytes, err := json.Marshal(conf)
	sendMessage(bytes, ServerID, "resourceConnectionConf", client)

	if err != nil {
		fmt.Println("\nFailed to marshal conf message into bytes")
		return
	}
}

func sendMessage(bytes []byte, senderId string, messageType string, client *Client) {
	socketMessage := SocketMessage{SenderId: senderId, Type: messageType, Data: bytes}
	select {
	case client.send <- socketMessage:
	default:
		fmt.Println("Client send channel is closed.")
	}
}

func lineAndCharacter(input string, offset int) (line int, character int, err error) {
	lf := rune(0x0A)

	if offset > len(input) || offset < 0 {
		return 0, 0, fmt.Errorf("Couldn't find offset %d within the input.", offset)
	}

	// Humans tend to count from 1.
	line = 1

	for i, b := range input {
		if b == lf {
			line++
			character = 0
		}
		character++
		if i == offset {
			break
		}
	}

	return line, character, nil
}

type data struct {
	Key string `json:"key"`
}
