package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"L0_task/internal/db"
	kafka_consumer "L0_task/internal/kafka"
	"L0_task/internal/service"
	"L0_task/internal/types"
)


var (
	// кэш заказов в памяти
	ordersCache = make(map[string]types.Order)
	mu          sync.RWMutex
)

//Функция loadOrdersToCache загружает все заказы из бд и помещает их в кэш
func loadOrdersToCache(database *sql.DB) {
	allOrders, err := db.LoadAllOrders(database)
	if err != nil {
		log.Fatalf("failed to load orders: %v", err)
	}
	mu.Lock()
	defer mu.Unlock()
	for _, o := range allOrders {
		ordersCache[o.OrderUID] = o
	}
	log.Printf("Loaded %d orders into cache\n", len(allOrders))
}

//Функция handleKafkaMessages - обработчик сообщений из Kafka, помещает их в кэш
func handleKafkaMessages(kafkaChan <-chan []byte, database *sql.DB) {
	for msg := range kafkaChan {
		var o types.Order
		if err := json.Unmarshal(msg, &o); err != nil {
			log.Printf("Invalid order message: %v", err)
			continue
		}

		mu.Lock()
		ordersCache[o.OrderUID] = o
		mu.Unlock()

		if err := db.SaveOrder(database, o); err != nil {
			log.Printf("Failed to save order from Kafka: %v", err)
		} else {
			log.Printf("Order %s saved from Kafka", o.OrderUID)
		}
	}
}

// Фунция setupHTTPHandlers настраивает ручки для сервера
func setupHTTPHandlers(database *sql.DB) {
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		postOrderHandler(w, r, database)
	})
	http.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		getOrderHandler(w, r, database)
	})
}

// Функция pingHandler - вспомогательная для проверки, что сервер работает
func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "pong")
}

// Функция postOrderHandler - обработчик для POST-ручки, создание заказа, сохраняет данные в кэш и бд
func postOrderHandler(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var o types.Order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := service.ValidateOrder(o); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := service.ValidateDelivery(o.Delivery); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := service.ValidatePayment(o.Payment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, item := range o.Items {
		if err := service.ValidateItem(item); err != nil {
			http.Error(w, fmt.Sprintf("item %d: %s", i, err.Error()), http.StatusBadRequest)
			return
		}
	}

	mu.Lock()
	ordersCache[o.OrderUID] = o
	mu.Unlock()

	if err := db.SaveOrder(database, o); err != nil {
		log.Printf("SaveOrder error: %v", err)
		http.Error(w, "failed to save order in db", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(o)
}

// Функция getOrderHandler - обработчик для GET-ручки для получения заказа, ищет заказ в кэше, если не нашел, ищет в бд
func getOrderHandler(w http.ResponseWriter, r *http.Request, database *sql.DB) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	id := parts[2]

	mu.RLock()
	o, ok := ordersCache[id]
	mu.RUnlock()

	if !ok {
		var err error
		o, err = db.LoadSingleOrder(database, id)
		if err != nil {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}
		mu.Lock()
		ordersCache[id] = o
		mu.Unlock()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(o)
}

func main() {
	database := db.Connect()
	defer database.Close()

	loadOrdersToCache(database)

	kafkaChan := make(chan []byte)
	go kafka_consumer.StartConsumer(kafkaChan)
	go handleKafkaMessages(kafkaChan, database)

	setupHTTPHandlers(database)

	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
