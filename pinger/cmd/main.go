package main

import (
	"context"
	"github.com/Quizert/Docker-pinger/pinger/clients"
	"github.com/Quizert/Docker-pinger/pinger/internal/service"
	"github.com/docker/docker/client"
	"log"
)

func main() {
	// Инициализация Docker клиента
	dockerClient, err := client.NewClientWithOpts(
		client.WithHost("unix:///var/run/docker.sock"),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		log.Fatal("Ошибка подключения к Docker API:", err)
	}

	// Инициализация клиента для отправки данных
	backendClient := clients.NewBackEndClient("http://backend-service:8080") // Используем имя контейнера
	svc := service.NewPingerService(backendClient, dockerClient)

	err = svc.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
