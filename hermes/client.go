package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan interface{}
}

type ClientSet map[*Client]bool
type ClientChannels map[*Client]map[string]bool

func (c *Client) message(msg interface{})  {
	c.send <- msg
}

func (c *Client) reportError(e error) {
	errorDescription := map[string]string{errorMessageType: e.Error()}
	c.message(errorDescription)
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		// Only acceptable messages from user are subscribe/unsubscribe messages
		var subMsg SubscribeMessage

		// Validates SubscribeMessage
		if err := c.conn.ReadJSON(&subMsg); err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				if websocket.IsUnexpectedCloseError(
						err,
						websocket.CloseNormalClosure,
						websocket.CloseNoStatusReceived,
						websocket.CloseGoingAway,
						websocket.CloseAbnormalClosure,
					) {
					log.Println("Client - Error reading message", err)
				}
				break
			}

			if !definedErrorType(err) {
				err = GenericError
			}
			c.reportError(err)
		}

		log.Println("Client - Read message")
		subRequest := SubscribeRequest{Client: c, SubMsg: subMsg}
		c.hub.subscribe <- subRequest
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				log.Println("Client - Hub closed channel")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			log.Println("Client - Writing message to client")
			if err := c.conn.WriteJSON(message); err != nil {
				log.Println("Client - Error writing message to client", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan interface{}, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
