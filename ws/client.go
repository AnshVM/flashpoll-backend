package ws

import (
	"fmt"
	"net/http"
	"os"

	"github.com/AnshVM/flashpoll-backend/types"
	"github.com/gorilla/websocket"
)

type UpdatePollResponse = types.UpdatePollResponse

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan types.UpdatePollResponse
	pollID uint
}

// upgrades http to webscocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// write sends messages from hub to the websocket connection
func (c *Client) write() {
	defer c.conn.Close()
	for {
		message, ok := <-c.send
		if !ok {
			// if hub closed the channel
			// send close message
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		err := c.conn.WriteJSON(message)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, pollID uint) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return r.Header["Origin"][0] == os.Getenv("CLIENT_URL")
	}
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan UpdatePollResponse), pollID: pollID}
	hub.register <- client
	go client.write()

}
