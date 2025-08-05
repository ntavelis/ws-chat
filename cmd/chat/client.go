package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"

	"github.com/coder/websocket"
)

type Client struct {
	ID       string
	Username string
	Send     chan []byte
	Conn     *websocket.Conn
}

var funnyAdjectives = []string{
	"Sneaky", "Bouncy", "Giggly", "Dizzy", "Fluffy", "Wobbly", "Sparkly", "Quirky",
	"Zany", "Peppy", "Bubbly", "Goofy", "Jolly", "Spunky", "Zippy", "Wacky",
}

var funnyNouns = []string{
	"Penguin", "Banana", "Llama", "Pickle", "Waffle", "Unicorn", "Taco", "Hamster",
	"Muffin", "Ninja", "Pirate", "Dragon", "Robot", "Wizard", "Pancake", "Narwhal",
}

func generateFunnyUsername() string {
	adjIdx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(funnyAdjectives))))
	nounIdx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(funnyNouns))))
	num, _ := rand.Int(rand.Reader, big.NewInt(100))

	return fmt.Sprintf("%s%s%d",
		funnyAdjectives[adjIdx.Int64()],
		funnyNouns[nounIdx.Int64()],
		num.Int64())
}

func NewClient(id, username string, conn *websocket.Conn) *Client {
	if username == "" {
		username = generateFunnyUsername()
	}

	client := &Client{
		ID:       id,
		Username: username,
		Send:     make(chan []byte, 256),
		Conn:     conn,
	}
	go client.writePump()
	return client
}

func (c *Client) SendWelcomeMessage() {
	welcomeMsg := fmt.Sprintf(`User %s connected.
You can pass your username by providing it as a query parameter ws://host@port?user=xxx
Type /exit to exit`, c.Username)
	if err := c.Conn.Write(context.Background(), websocket.MessageText, []byte(welcomeMsg)); err != nil {
		log.Printf("Error sending welcome message to client %s: %v", c.ID, err)
	}
}

func (c *Client) InformHowManyUsersAreInTheRoom(onlineClients int) {
	message := fmt.Sprintf("There are currently %d users in the room.", onlineClients)
	if err := c.Conn.Write(context.Background(), websocket.MessageText, []byte(message)); err != nil {
		log.Printf("Error sending user count to client %s: %v", c.ID, err)
	}
}

func (c *Client) writePump() {
	for msg := range c.Send {
		if err := c.Conn.Write(context.Background(), websocket.MessageText, msg); err != nil {
			log.Printf("Error writing message to client %s: %v", c.ID, err)
			return
		}
	}
}

