package websocket

var pools = make(map[string]*Pool)

type Pool struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
}

func NewPool() *Pool {
	return &Pool{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message),
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
			for client := range p.clients {
				select {
				case client.send <- message:
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
