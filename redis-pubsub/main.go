package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Subscriber struct {
	Redis *redis.Client
	Topic string
}

type Publisher struct {
	Redis *redis.Client
}

func NewSubscriber(rdb *redis.Client, topic string) *Subscriber {
	return &Subscriber{
		Redis: rdb,
		Topic: topic,
	}
}

func NewPublisher(rdb *redis.Client) *Publisher {
	return &Publisher{
		Redis: rdb,
	}
}

func (s *Subscriber) Listen(ctx context.Context) {
	fmt.Println("Listening for messages...")
	pubSub := s.Redis.Subscribe(ctx, s.Topic)
	defer pubSub.Close()

	ch := pubSub.Channel()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Subscriber shutting down...")
			return
		case msg, ok := <-ch:
			if !ok {
				fmt.Println("Channel closed")
				return
			}

			if msg.Payload == "" {
				fmt.Println("Empty message received")
				continue
			}

			var data ProductMessage
			err := json.Unmarshal([]byte(msg.Payload), &data)
			if err != nil {
				fmt.Println("Failed to unmarshal message:", err)
				continue
			}

			fmt.Printf("Received - Product ID: %d, Name: %s, Action: %s\n",
				data.Product.ID, data.Product.Name, data.Action)
		}
	}
}

func (p *Publisher) Publish(ctx context.Context, topic string, message string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second) // Set timeout for publishing
	defer cancel()

	err := p.Redis.Publish(ctx, topic, message).Err()
	if err != nil {
		log.Println("Failed to publish message:", err)
	}
	return err
}

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProductMessage struct {
	Product *Product `json:"product"`
	Action  string   `json:"action"`
}

func NewProduct(id int, name string) *Product {
	return &Product{ID: id, Name: name}
}

func NewProductMessage(product *Product, action string) *ProductMessage {
	return &ProductMessage{Product: product, Action: action}
}

func (p *ProductMessage) ToBytes() ([]byte, error) {
	return json.Marshal(p)
}

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password
		DB:       0,  // Default DB
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return
	}
	fmt.Println("Connected to Redis")

	productSub := NewSubscriber(rdb, "product")
	productPub := NewPublisher(rdb)

	ctx, cancel := context.WithCancel(ctx)
	go productSub.Listen(ctx)

	// Publish a test message
	time.Sleep(1 * time.Second) // Give some time for subscriber to start

	product := NewProduct(1, "Laptop")
	productBytes, err := NewProductMessage(product, "create").ToBytes()
	if err != nil {
		fmt.Println("Failed to marshal product message:", err)
		return
	}

	err = productPub.Publish(ctx, "product", string(productBytes))
	if err != nil {
		fmt.Println("Failed to publish message:", err)
		return
	}
	fmt.Println("Message published")

	productTwo := NewProduct(2, "Laptop A")
	productTwoBytes, err := NewProductMessage(productTwo, "update").ToBytes()
	if err != nil {
		fmt.Println("Failed to marshal product message:", err)
		return
	}

	err = productPub.Publish(ctx, "product", string(productTwoBytes))
	if err != nil {
		fmt.Println("Failed to publish message:", err)
		return
	}
	fmt.Println("Message published")

	time.Sleep(10 * time.Second)
	cancel()
	fmt.Println("Shutting down...")
}
