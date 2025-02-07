package service

import "github.com/Quizert/Docker-pinger/backend/internal/model"

type Storage interface {
	SavePingResults(pingInfo model.ContainerInfo) error
	GetPingResults(limit int) ([]*model.ContainerInfo, error)
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) SavePingResults(pingInfo model.ContainerInfo) error {
	return s.storage.SavePingResults(pingInfo)
}

func (s *Service) GetPingResults(limit int) ([]*model.ContainerInfo, error) {
	return s.storage.GetPingResults(limit)
}
