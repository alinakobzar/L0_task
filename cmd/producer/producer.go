package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"L0_task/internal/types"

	"github.com/segmentio/kafka-go"
)

// Функия-продьюсер создает тестовый заказ и отправляет его в KAfka в топик orders
func main() {
	brokerAddress := "localhost:9092"
	topic := "orders"

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{brokerAddress},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	t, _ := time.Parse(time.RFC3339, "2021-11-26T06:22:19Z")

	order := types.Order{
		OrderUID:    "test1",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: types.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: types.Payment{
			Transaction:  "b563feb7b2b84b6test",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDT:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
		},
		Items: []types.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:          "en",
		CustomerID:      "test",
		DeliveryService: "meest",
		ShardKey:        "9",
		SmID:            99,
		DateCreated:     t,
		OofShard:        "1",
	}

	data, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("Ошибка сериализации: %v", err)
	}

	fmt.Println("Отправка заказа в Kafka:\n", string(data))

	msg := kafka.Message{
		Key:   []byte(order.OrderUID),
		Value: data,
	}

	if err := writer.WriteMessages(context.Background(), msg); err != nil {
		log.Fatalf("Не удалось отправить сообщение: %v", err)
	}

	fmt.Println("Заказ отправлен успешно")
}
