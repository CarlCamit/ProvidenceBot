package websocket

import "github.com/gorilla/websocket"

var dialer = websocket.Dialer{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client represents the websocket connection
type Client struct {
	conn *websocket.Conn
}

// NewConn produces a websocket connection to Twitch
func NewConn() (*Client, error) {
	conn, _, err := dialer.Dial("wss://irc-ws.chat.twitch.tv:443", nil)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn}, nil
}

// ReadMessage reads the next message in the websocket connection
func (c *Client) ReadMessage() ([]byte, error) {
	_, messageBytes, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return messageBytes, nil
}

// WriteText sends a text data message to the websocket connection
func (c *Client) WriteText(message string) {
	c.conn.WriteMessage(1, []byte(message))
}

// CloseConn terminates the websocket connection
func (c *Client) CloseConn() {
	c.conn.Close()
}
