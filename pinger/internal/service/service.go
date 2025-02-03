package service

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type BackEndClient interface {
	sendContainerInfo() error
}

type PingerService struct {
	client BackEndClient
}

func NewPingerService(client BackEndClient) *PingerService {
	return &PingerService{client: client}
}

// getContainerInfo получает список всех контейнеров и их информацию
func getContainerInfo() ([]model.ContainerInfo, error) {
	var result []ContainerInfo

	// Получаем список всех контейнеров (их ID)
	cmd := exec.Command("docker", "ps", "-aq")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	containerIDs := strings.Fields(string(output)) // Разбиваем output по строкам
	if len(containerIDs) == 0 {
		return result, nil // Возвращаем пустой список, если контейнеров нет
	}

	// Для каждого контейнера получаем информацию
	for _, id := range containerIDs {
		info, err := inspectContainer(id)
		if err != nil {
			fmt.Println("Ошибка при получении информации о контейнере:", err)
			continue
		}
		result = append(result, info)
	}

	return result, nil
}

// inspectContainer получает имя, статус и IP контейнера
func inspectContainer(containerID string) (ContainerInfo, error) {
	cmd := exec.Command("docker", "inspect", "--format",
		"{{.Name}} {{.State.Status}} {{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}",
		containerID)

	var output bytes.Buffer
	cmd.Stdout = &output
	err := cmd.Run()
	if err != nil {
		return ContainerInfo{}, err
	}

	data := strings.Fields(output.String()) // Разбиваем строку
	if len(data) < 2 {
		return ContainerInfo{}, fmt.Errorf("неполные данные: %v", data)
	}

	// Убираем `/` перед именем контейнера
	name := strings.TrimPrefix(data[0], "/")
	status := data[1]
	ip := ""
	if len(data) > 2 {
		ip = data[2]
	}

	return ContainerInfo{Name: name, Status: status, IP: ip}, nil
}
