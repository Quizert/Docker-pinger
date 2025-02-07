package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Quizert/Docker-pinger/pinger/internal/model"
	"net/http"
)

type BackEndClient struct {
	baseURL string
	client  *http.Client
}

func NewBackEndClient(baseURL string) *BackEndClient {
	return &BackEndClient{baseURL: baseURL, client: http.DefaultClient}
}

func (client *BackEndClient) SendPingResult(results []*model.ContainerInfo) error {
	jsonData, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("ошибка сериализации JSON: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, client.baseURL+"/ping", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("ошибка от Backend: %w", resp.Status)
	}
	return nil
}
