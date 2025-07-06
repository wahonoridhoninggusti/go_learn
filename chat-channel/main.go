package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrClientDisconnected   = errors.New("client disconnected")
)

type Client struct {
	Username   string
	Disconnect bool
	ch         chan string
	Time       time.Time
	mu         sync.Mutex
}

type ChatServer struct {
	Clients map[string]*Client
	mu      sync.Mutex
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		Clients: make(map[string]*Client),
	}
}

func (s *ChatServer) Connect(username string) (*Client, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exist := s.Clients[username]; exist {
		return nil, ErrUsernameAlreadyTaken
	}

	client := &Client{
		Username: username,
		ch:       make(chan string, 10),
		Time:     time.Now(),
	}

	s.Clients[username] = client

	return client, nil
}

func (s *ChatServer) Disconnect(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if exist, ok := s.Clients[client.Username]; ok && exist == client {
		client.mu.Lock()
		client.Disconnect = true
		close(client.ch)
		client.mu.Unlock()
		delete(s.Clients, client.Username)
	}
}

func (c *Client) Send(message string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Disconnect {
		return
	}

	select {
	case c.ch <- message:
	default:
		fmt.Printf("Message to %s dropped: %s\n", c.Username, message)
	}

}

func (c *Client) Receive() string {
	msg, ok := <-c.ch
	if !ok {
		return ""
	}
	return msg
}

func (s *ChatServer) PrivateMessage(sender *Client, recipient string, message string) error {
	sender.mu.Lock()
	if sender.Disconnect {
		sender.mu.Unlock()
		return ErrClientDisconnected
	}
	sender.mu.Unlock()

	s.mu.Lock()
	target, ok := s.Clients[recipient]
	s.mu.Unlock()
	if !ok {
		return ErrRecipientNotFound
	}

	target.mu.Lock()
	disconnected := target.Disconnect
	target.mu.Unlock()

	if disconnected {
		return ErrClientDisconnected
	}

	formatted := fmt.Sprintf("[private from %s] %s", sender.Username, message)
	target.Send(formatted)
	return nil
}

func (s *ChatServer) Broadcast(sender *Client, message string) {
	formatted := fmt.Sprintf("[broadcast from %s] %s", sender.Username, message)
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, target := range s.Clients {
		if target.Username != sender.Username {
			target.Send(formatted)
		}
	}
}

func main() {
	server := NewChatServer()
	client1, err := server.Connect("wahono")
	if err != nil {
		fmt.Println(err)
		return
	}

	client2, err := server.Connect("dimas")
	if err != nil {
		fmt.Println(err)
		return
	}

	client3, err := server.Connect("wahyu")
	if err != nil {
		fmt.Println(err)
		return
	}

	// server.Disconnect(client1)
	server.PrivateMessage(client2, "wahono", "halo")
	// msg1 := client1.Receive()

	server.Broadcast(client1, "halo from client1")
	msg2 := client2.Receive()
	msg3 := client3.Receive()
	fmt.Println("receive: ", msg2, msg3)
	// select {}
}
