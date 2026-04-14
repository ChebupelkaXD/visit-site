package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

var DB *sql.DB

func Init() error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	var err error

	DB, err = sql.Open("postgres", dsn)

	if err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		if err = DB.Ping(); err == nil {
			log.Println("✅ Подключились к PostgreSQL")
			return nil
		}
		log.Printf("Попытка %d/10 подключения к БД...", i+1)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("Ошибка подключения к БД: %w", err)
}
