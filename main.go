package main

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const pongMessage = "PONG :tmi.twitch.tv\r\n"

func main() {
	fmt.Printf("[%s] Connecting to Twitch IRC...\n", time.Now())

	var (
		wg     = sync.WaitGroup{}
		dialer = websocket.Dialer{
			ReadBufferSize:   1024,
			WriteBufferSize:  1024,
			HandshakeTimeout: 30 * time.Second,
		}
	)

	conn, _, err := dialer.Dial("wss://irc-ws.chat.twitch.tv:443", nil)
	if err != nil {
		fmt.Printf("[%s] Cannot connect to Twitch IRC.\n", time.Now())
		return
	}

	fmt.Printf("[%s] Connected to Twitch IRC!\n", time.Now())

	conn.WriteMessage(1, []byte("PASS <password>\r\n"))
	conn.WriteMessage(1, []byte("NICK <name>\r\n"))
	conn.WriteMessage(1, []byte("JOIN <channel>\r\n"))

	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			_, messageBytes, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("[%s] Connection to Twitch IRC lost, err: %s\n", time.Now(), err)
				conn.Close()
				return
			}

			message := bytes.NewBuffer(messageBytes).String()
			fmt.Print(message)

			// Respond to heartbeat message and move on to the next message
			if message == "PING :tmi.twitch.tv\r\n" {
				conn.WriteMessage(1, []byte(pongMessage))
				fmt.Print(pongMessage)
				continue
			}
		}
	}()

	wg.Wait()
}
