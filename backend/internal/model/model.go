package model

import "time"

type ContainerInfo struct {
	Name      string    `json:"name"`
	IP        string    `json:"IP"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}
