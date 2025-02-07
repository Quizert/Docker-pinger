package postgresql

import (
	"context"
	"fmt"
	"github.com/Quizert/Docker-pinger/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) SavePingResults(pingInfo model.ContainerInfo) error {
	query := `
		INSERT INTO ping_data (name, ip, status, last_ping)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (name) DO UPDATE
		SET ip = EXCLUDED.ip, status = EXCLUDED.status, last_ping = NOW();
	`
	_, err := r.db.Exec(context.Background(), query, pingInfo.Name, pingInfo.IP, pingInfo.Status)
	if err != nil {
		return fmt.Errorf("ошибка в выполнении запроса : %w", err)
	}
	return nil
}

func (r *Repository) GetPingResults(limit int) ([]*model.ContainerInfo, error) {
	query := `
		SELECT name, ip, status, last_ping
		from ping_data order by last_ping desc
		LIMIT $1
	`
	fmt.Println(query, limit)
	rows, err := r.db.Query(context.Background(), query, limit)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("ошибка в запросе: %w", err)
	}
	defer rows.Close()

	var containersInfo []*model.ContainerInfo
	for rows.Next() {
		var pingInfo model.ContainerInfo
		err = rows.Scan(&pingInfo.Name, &pingInfo.IP, &pingInfo.Status, &pingInfo.Timestamp)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("ошибка в парсинге данных: %w", err)
		}
		containersInfo = append(containersInfo, &pingInfo)
	}
	return containersInfo, nil
}
