package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

//Функция StartConsumer запускает потребителя для топика orders
func StartConsumer(ordersChan chan<- []byte) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"127.0.0.1:9092"},
		Topic:   "orders",
		GroupID: "orders-service-group",
	})
	defer reader.Close()

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		ordersChan <- m.Value

		fmt.Printf("Received message at offset %d: %s\n", m.Offset, string(m.Value))
	}
}