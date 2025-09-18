package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/avast/retry-go"
	_ "github.com/lib/pq"
)


//Функция Connect создаёт подключение к PostgreSQL с повторными попытками подключения
// При неудаче после нескольких попыток завершает программу
func Connect() *sql.DB {
	connStr := "postgres://l0_user:alina789589@localhost:5432/l0_db?sslmode=disable"

	var db *sql.DB

	err := retry.Do(
		func() error {
			var err error
			db, err = sql.Open("postgres", connStr)
			if err != nil {
				return err
			}
			return db.Ping()
		},
		retry.Attempts(5),          
		retry.Delay(2*time.Second), 
		retry.DelayType(retry.FixedDelay),
	)

	if err != nil {
		log.Fatalf("Не удалось подключиться к БД после нескольких попыток: %v", err)
	}

	fmt.Println("Успешное подключение к Postgres")
	return db
}
