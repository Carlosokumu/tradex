package chat

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	//  Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	//holds messages from the clients
	messenger chan ChatDetails
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		messenger:  make(chan ChatDetails),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.messages)
			}
		case jsonMessage := <-h.messenger:
			for client := range h.clients {
				select {
				case client.messages <- jsonMessage:
				default:
					close(client.messages)
					delete(h.clients, client)
				}
			}

		}
	}
}
