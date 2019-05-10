package main

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/carlcamit/ProvidenceBot/websocket"
)

const (
	pingMessage = "PING :tmi.twitch.tv\r\n"
	pongMessage = "PONG :tmi.twitch.tv\r\n"
)

func main() {
	fmt.Printf("[%s] Connecting to Twitch IRC...\n", time.Now())

	ws, err := websocket.NewConn()
	if err != nil {
		fmt.Printf("[%s] Cannot connect to Twitch IRC.\n", time.Now())
		return
	}

	fmt.Printf("[%s] Connected to Twitch IRC!\n", time.Now())

	ws.WriteText("PASS <password>\r\n")
	ws.WriteText("NICK <name>\r\n")
	ws.WriteText("JOIN <channel>\r\n")

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			messageBytes, err := ws.ReadMessage()
			if err != nil {
				fmt.Printf("[%s] Connection to Twitch IRC lost, err: %s\n", time.Now(), err)
				ws.CloseConn()
				return
			}

			message := bytes.NewBuffer(messageBytes).String()
			fmt.Print(message)

			// Respond to heartbeat message and move on to the next message
			if message == pingMessage {
				ws.WriteText(pongMessage)
				fmt.Print(pongMessage)
				continue
			}
		}
	}()

	wg.Wait()
}
