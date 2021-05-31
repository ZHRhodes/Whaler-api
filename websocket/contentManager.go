package websocket

import (
	"fmt"

	"github.com/heroku/whaler-api/OT/ot-master"
	"github.com/heroku/whaler-api/models"
)

type ContentFetcher interface {
	FetchDocument(resourceId string) (string, error)
	SaveDocument(resourceId string, content string)
}

type ActiveServerDoc struct {
	ServerDoc *ot.ServerDoc
	Clients   []*Client
}

type ContentManager struct {
	ActiveServerDocs  map[string]ActiveServerDoc
	ClientResourceIDs map[*Client]string
}

var Fetcher ContentFetcher

func (manager *ContentManager) clients(resourceId string) []*Client {
	return manager.ActiveServerDocs[resourceId].Clients
}

func (manager *ContentManager) serverDoc(resourceId string) *ot.ServerDoc {
	return manager.ActiveServerDocs[resourceId].ServerDoc
}

func (manager *ContentManager) registerClient(client *Client, resourceId string) error {
	if existingServerDoc, ok := manager.ActiveServerDocs[resourceId]; ok {
		fmt.Printf("\n Adding client with id %s to existing doc for resourceId %s", client.Id, resourceId)
		existingServerDoc.Clients = append(existingServerDoc.Clients, client)
		manager.ActiveServerDocs[resourceId] = existingServerDoc
	} else {
		fmt.Printf("\n Creating new server doc for resourceId %s", resourceId)
		note, err := models.FetchNote("", resourceId)

		if err != nil {
			fmt.Printf("\nFailed to fetch note with id %s", resourceId)
			return err
		}

		newDoc := ot.NewDocFromStr(note.Content)
		newServerDoc := &ot.ServerDoc{Doc: newDoc, History: []ot.Ops{}}
		manager.ActiveServerDocs[resourceId] = ActiveServerDoc{ServerDoc: newServerDoc, Clients: []*Client{client}}
	}

	manager.ClientResourceIDs[client] = resourceId
	return nil
}

func (manager *ContentManager) unregisterClient(client *Client) {
	resourceId := manager.ClientResourceIDs[client]
	if existingServerDoc, ok := manager.ActiveServerDocs[resourceId]; ok {
		fmt.Printf("\nUnregistering client with id %s from existing doc for resourceId %s", client.Id, resourceId)
		for i, existingClient := range existingServerDoc.Clients {
			if client.Id == existingClient.Id {
				existingServerDoc.Clients = append(existingServerDoc.Clients[:i], existingServerDoc.Clients[i+1:]...)
				if len(existingServerDoc.Clients) > 0 {
					manager.ActiveServerDocs[resourceId] = existingServerDoc
				} else {
					docString := existingServerDoc.ServerDoc.Doc.String()
					models.SaveNoteContent(resourceId, docString)
					delete(manager.ActiveServerDocs, resourceId)
				}
			}
		}
	}
}
