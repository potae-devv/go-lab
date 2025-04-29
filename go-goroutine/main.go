package main

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type Message struct {
	Data string `json:"data"`
}

type PubSub struct {
	subs []chan Message
	mu   sync.Mutex
}

func (ps *PubSub) Subscribe() chan Message {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ch := make(chan Message, 1)
	ps.subs = append(ps.subs, ch)
	return ch
}

func (ps *PubSub) Publish(msg Message) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	for _, ch := range ps.subs {
		ch <- msg
	}
}

func main() {
	app := fiber.New()

	// Middleware to handle panic
	app.Use(func(c *fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				c.Status(500).SendString("Internal Server Error")
			}
		}()
		return c.Next()
	})

	ps := &PubSub{}

	app.Post("/publisher", func(c *fiber.Ctx) error {
		message := new(Message)
		c.BodyParser(message)
		ps.Publish(*message)
		return c.JSON(&fiber.Map{
			"message": "add to subscribers",
		})
	})

	sub1 := ps.Subscribe()
	go func() {
		for msg := range sub1 {
			fmt.Println("Received message from subscriber 1:", msg)
		}
	}()

	sub2 := ps.Subscribe()
	go func() {
		for msg := range sub2 {
			fmt.Println("Received message from subscriber 2:", msg)
		}
	}()

	app.Listen(":8888")
}
