package orders

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/matsapkov/wb_kafka_service/internal/models"
	"log"
	"time"
)

type OrdersRepository interface {
	GetOrder(ctx context.Context, Id string) (models.Order, error)
	ListRecent(ctx context.Context, limit int) ([]models.Order, error)
	SaveOrder(ctx context.Context, orderUID, trackNumber string, dateCreated, updatedAt time.Time, payload json.RawMessage) error
}

type PostgresOrders struct {
	db *sql.DB
}

func NewPostgresOrders(db *sql.DB) *PostgresOrders {
	return &PostgresOrders{db: db}
}

const getOrderQuery = `		SELECT order_uid, track_number, date_created, payload, updated_at
		FROM orders
		WHERE order_uid = $1`

const getRecentQuery = `SELECT order_uid, track_number, date_created, payload, updated_at
FROM orders
ORDER BY date_created DESC
LIMIT $1;`

func (p *PostgresOrders) GetOrder(ctx context.Context, orderUID string) (models.Order, error) {

	var order models.Order

	err := p.db.QueryRowContext(ctx, getOrderQuery, orderUID).
		Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.DateCreated,
			&order.Payload,
			&order.UpdatedAt,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Order{}, err
		}
		return models.Order{}, err
	}

	return order, nil
}

func (p *PostgresOrders) ListRecent(ctx context.Context, limit int) ([]models.Order, error) {
	rows, err := p.db.QueryContext(ctx, getRecentQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("repo: list recent: %w", err)
	}
	defer rows.Close()

	orders := make([]models.Order, 0, 128)

	for rows.Next() {
		var o models.Order
		if err := rows.Scan(
			&o.OrderUID,
			&o.TrackNumber,
			&o.DateCreated,
			&o.Payload,
			&o.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("repo: list recent scan: %w", err)
		}
		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo: list recent rows: %w", err)
	}

	return orders, nil
}

func (p *PostgresOrders) SaveOrder(ctx context.Context, orderUID, trackNumber string, dateCreated, updatedAt time.Time, payload json.RawMessage) error {
	query := `INSERT INTO orders (order_uid, track_number, date_created, payload, updated_at)
			  VALUES ($1, $2, $3, $4, $5)`

	_, err := p.db.ExecContext(ctx, query, orderUID, trackNumber, dateCreated, payload, updatedAt)
	if err != nil {
		log.Printf("Ошибка сохранения заказа: %v", err)
		return err
	}

	log.Printf("Заказ %s успешно сохранен", orderUID)
	return nil
}
