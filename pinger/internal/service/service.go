package service

import (
	"context"
	"fmt"
	"github.com/Quizert/Docker-pinger/pinger/internal/model"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"log"
	"time"
)

type BackEndClient interface {
	SendPingResult(results []*model.ContainerInfo) error
}

type PingerService struct {
	BackEndClient BackEndClient
	DockerClient  *client.Client
}

func NewPingerService(backEndClient BackEndClient, dockerClient *client.Client) *PingerService {
	return &PingerService{BackEndClient: backEndClient, DockerClient: dockerClient}
}

func (p *PingerService) Ping(ctx context.Context) error {
	for {
		containers, err := getContainersInfo(p.DockerClient)
		if err != nil {
			log.Printf("Ошибка получения данных: %v", err)
			continue
		}

		if err := p.BackEndClient.SendPingResult(containers); err != nil {
			log.Printf("Ошибка отправки данных: %v", err)
		}
		fmt.Println(containers)
		time.Sleep(15 * time.Second) // Интервал опроса
	}
}

func getContainersInfo(cli *client.Client) ([]*model.ContainerInfo, error) {
	ctx := context.Background()

	// Получаем список контейнеров
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var result []*model.ContainerInfo

	for _, c := range containers {
		// Получаем детальную информацию для получения IP-адреса
		details, err := cli.ContainerInspect(ctx, c.ID)
		if err != nil {
			log.Printf("Ошибка инспекции контейнера %s: %v", c.ID, err)
			continue
		}

		// Извлекаем IP-адрес из сетевых настроек
		var ip string
		if details.NetworkSettings != nil {
			for _, network := range details.NetworkSettings.Networks {
				if network.IPAddress != "" {
					ip = network.IPAddress
					break
				}
			}
		}

		result = append(result, &model.ContainerInfo{
			Name:   c.Names[0][1:],
			IP:     ip,
			Status: c.Status,
		})
	}

	return result, nil
}
