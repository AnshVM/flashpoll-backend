package ws

type Hub struct {
	clients    map[*Client]bool
	Broadcast  chan UpdatePollResponse
	register   chan *Client
	unregister chan *Client
	rooms      map[uint][]*Client // map of poll ID and the clients connected to it
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan UpdatePollResponse),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:      make(map[uint][]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.rooms[client.pollID] = append(h.rooms[client.pollID], client)

		case client := <-h.unregister:

			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.Broadcast:
			for _, client := range h.rooms[message.ID] {
				select {
				case client.send <- message:
				//incase client send channel is full or closed, assume client is dead or stuck and close the channel and delete the client
				default:
					close(client.send)
					delete(h.clients, client)
				}

			}
		}
	}
}
