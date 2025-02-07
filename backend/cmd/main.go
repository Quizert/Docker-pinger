package main

import (
	"context"
	"github.com/Quizert/Docker-pinger/backend/internal/handler"
	"log"
	"net/http"

	"github.com/Quizert/Docker-pinger/backend/internal/service"
	"github.com/Quizert/Docker-pinger/backend/internal/storage/postgresql"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Подключаемся к PostgreSQL
	dbURL := "postgres://postgres:postgres@db:5432/postgres?sslmode=disable"

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Ошибка подключения к PostgreSQL:", err)
	}
	defer db.Close() // Закрываем соединение при завершении программы

	// Создаём репозиторий с подключением к БД
	store := postgresql.NewRepository(db)

	svc := service.NewService(store)
	pingHandler := handler.NewHandler(svc)
	http.HandleFunc("/ping", pingHandler.WritePingResults)
	http.HandleFunc("/get", pingHandler.GetPingResults)
	log.Println("Server is running on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
