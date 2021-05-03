package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/heroku/whaler-api/OT/ot-master"
)

var ServerDocClients = map[*ot.ServerDoc][]*Client{}

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
		return err
	}

	doc := ot.ServerDocs[change.ResourceId]
	ops := []ot.Op{}
	for i, n := range change.N {
		ops = append(ops, ot.Op{N: n, S: change.S[i]})
	}
	ops, err := doc.Recv(change.Rev, ops)
	if err != nil {
		fmt.Printf("\nFailed sending changes to doc. %s", err)
		return err
	}
	err = returnOps(client, doc, message.MessageId, change.ResourceId, ops)
	return err
}

func returnOps(client *Client, serverDoc *ot.ServerDoc, messageId string, resourceId string, ops ot.Ops) error {
	n := []int{}
	s := []string{}

	for _, op := range ops {
		n = append(n, op.N)
		s = append(s, op.S)
	}

	message := DocumentChangeReturn{ResoureceId: resourceId, N: n, S: s}

	bytes, err := json.Marshal(message)

	if err != nil {
		fmt.Println("\nFailed to marshal return ops into bytes")
		return err
	}

	for _, client := range ServerDocClients[serverDoc] {
		sendMessage(bytes, messageId, ServerID, "docChangeReturnOps", client)
	}
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
	// note, err := models.FetchNote("", request.ResourceId)

	// if err != nil {
	// 	fmt.Printf("\nFailed to fetch note with id %s", request.ResourceId)
	// 	return err
	// }

	serverDoc := ot.ServerDocs[request.ResourceId]
	if serverDoc == nil {
		doc := ot.NewDocFromStr("")
		serverDoc := ot.ServerDoc{Doc: doc, History: []ot.Ops{}}
		ot.ServerDocs[request.ResourceId] = &serverDoc
	} else {
		fmt.Println("Connected to existing server doc")
	}
	ServerDocClients[serverDoc] = append(ServerDocClients[serverDoc], client) //TODO: clear this map out somewhere?
	//TODO in general, the unregister flow hasn't really been looked at yet
	sendResourceConnectionConfirmation(message.MessageId, request.ResourceId, serverDoc.Doc.String(), client)
	return nil
}

func sendResourceConnectionConfirmation(messageId string, resourceId string, initialState string, client *Client) {
	conf := ResourceConnectionConf{ResourceId: resourceId, InitialState: initialState}
	bytes, err := json.Marshal(conf)
	sendMessage(bytes, messageId, ServerID, "resourceConnectionConf", client)

	if err != nil {
		fmt.Println("\nFailed to marshal conf message into bytes")
		return
	}
}

func sendMessage(bytes []byte, messageId string, senderId string, messageType string, client *Client) {
	fmt.Printf("\nSending messageId %s, senderId %s, messageType %s, to clientId %s", messageId, senderId, messageType, client.id)
	socketMessage := SocketMessage{SenderId: senderId, MessageId: messageId, Type: messageType, Data: bytes}
	select {
	case client.send <- socketMessage:
	default:
		fmt.Println("Client send channel is closed.")
	}
}
