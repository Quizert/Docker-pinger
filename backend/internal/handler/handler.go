package handler

import (
	"encoding/json"
	"github.com/Quizert/Docker-pinger/backend/internal/model"
	"net/http"
	"strconv"
)

type Service interface {
	SavePingResults(pingInfo model.ContainerInfo) error
	GetPingResults(limit int) ([]*model.ContainerInfo, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) WritePingResults(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var containers []model.ContainerInfo
		err := json.NewDecoder(r.Body).Decode(&containers)
		if err != nil {
			http.Error(w, "Неверный JSON", http.StatusBadRequest)
			return
		}
		for i := range containers {
			err = h.service.SavePingResults(containers[i])
			if err != nil {
				http.Error(w, "Ошибка сохранения", http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusCreated)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) GetPingResults(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var limit int
		limitStr := r.URL.Query().Get("limit")
		if limitStr == "" {
			limit = 50
		} else {
			limit, _ = strconv.Atoi(limitStr)
		}
		results, err := h.service.GetPingResults(limit)
		if err != nil {
			http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}
