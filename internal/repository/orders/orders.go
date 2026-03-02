package orders

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/matsapkov/wb_kafka_service/internal/models"
)

type OrdersRepository interface {
	GetOrder(ctx context.Context, id string) (models.Order, error)
	ListRecent(ctx context.Context, limit int) ([]models.Order, error)
	SaveOrder(ctx context.Context, orderUID, trackNumber string, dateCreated, updatedAt time.Time, payload json.RawMessage) error
}

type PostgresOrders struct {
	db *sql.DB
}

func NewPostgresOrders(db *sql.DB) *PostgresOrders {
	return &PostgresOrders{db: db}
}

const getOrderQuery = `
SELECT order_uid, track_number, date_created, payload, updated_at
FROM orders
WHERE order_uid = $1`

const getRecentQuery = `
SELECT order_uid, track_number, date_created, payload, updated_at
FROM orders
ORDER BY updated_at DESC
LIMIT $1`

const upsertOrderQuery = `
INSERT INTO orders (order_uid, track_number, date_created, payload, updated_at)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (order_uid) DO UPDATE
SET track_number = EXCLUDED.track_number,
    date_created = EXCLUDED.date_created,
    payload = EXCLUDED.payload,
    updated_at = EXCLUDED.updated_at`

func (p *PostgresOrders) GetOrder(ctx context.Context, orderUID string) (models.Order, error) {
	var order models.Order

	err := p.db.QueryRowContext(ctx, getOrderQuery, orderUID).Scan(
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
		return models.Order{}, fmt.Errorf("repo get order: %w", err)
	}

	return order, nil
}

func (p *PostgresOrders) ListRecent(ctx context.Context, limit int) ([]models.Order, error) {
	rows, err := p.db.QueryContext(ctx, getRecentQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("repo list recent: %w", err)
	}
	defer rows.Close()

	orders := make([]models.Order, 0, limit)
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.OrderUID, &o.TrackNumber, &o.DateCreated, &o.Payload, &o.UpdatedAt); err != nil {
			return nil, fmt.Errorf("repo list recent scan: %w", err)
		}
		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo list recent rows: %w", err)
	}

	return orders, nil
}

func (p *PostgresOrders) SaveOrder(ctx context.Context, orderUID, trackNumber string, dateCreated, updatedAt time.Time, payload json.RawMessage) error {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("repo begin tx: %w", err)
	}

	if _, err := tx.ExecContext(ctx, upsertOrderQuery, orderUID, trackNumber, dateCreated, payload, updatedAt); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("repo save order: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("repo commit tx: %w", err)
	}

	return nil
}
