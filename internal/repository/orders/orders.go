package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/matsapkov/wb_kafka_service/internal/models"
)

type OrdersRepository interface {
	GetOrder(ctx context.Context, Id string) (models.Order, error)
	ListRecent(ctx context.Context, limit int) ([]models.Order, error)
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
