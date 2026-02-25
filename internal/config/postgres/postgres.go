package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

type PostgresConfig struct {
	dsn string
}

func NewPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		dsn: os.Getenv("DB_DSN"),
	}
}

func (p *PostgresConfig) PGconnect() (*sql.DB, error) {
	db, err := sql.Open("postgres", p.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
