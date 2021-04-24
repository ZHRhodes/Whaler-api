package websocket

import "fmt"

var pools = make(map[string]*Pool)

type Pool struct {
	clients    map[*Client]bool
	broadcast  chan SocketMessage
	register   chan *Client
	unregister chan *Client
}

func NewPool() *Pool {
	return &Pool{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		broadcast:  make(chan SocketMessage),
	}
}

func (p *Pool) Start() {
	for {
		select {
		case client := <-p.register:
			p.clients[client] = true
		case client := <-p.unregister:
			if _, ok := p.clients[client]; ok {
				delete(p.clients, client)
				close(client.send)
			}
		case message := <-p.broadcast:
			fmt.Printf("\nBroadcasting message of type %s", message.Type)
			for client := range p.clients {
				// if client.id == message.Id {
				// 	fmt.Printf("\nNot sending message to sender - client with id %s", client.id)
				// 	return
				// }
				select {
				case client.send <- message:
					fmt.Printf("\nSending message with id %s to client with id %s", message.Id, client.id)
				default:
					close(client.send)
					delete(p.clients, client)
				}
			}
		}
	}
}

func AddPool(id string, pool *Pool) {
	pools[id] = pool
}
