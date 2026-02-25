-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
                                      order_uid     text PRIMARY KEY,
                                      track_number  text NOT NULL,
                                      date_created  timestamptz NOT NULL,
                                      payload       jsonb NOT NULL,
                                      updated_at    timestamptz NOT NULL DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS orders_date_created_idx ON orders(date_created);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS orders_date_created_idx;
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd