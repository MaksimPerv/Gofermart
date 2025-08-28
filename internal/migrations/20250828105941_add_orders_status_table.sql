-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders
(
    id         SERIAL PRIMARY KEY,
    number     VARCHAR(255) NOT NULL UNIQUE,
    user_id    INTEGER      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ                                        DEFAULT CURRENT_TIMESTAMP,
    accrual    DECIMAL(10, 2)                                     DEFAULT 0,
    status_id  INTEGER      NOT NULL REFERENCES order_statuses(id) DEFAULT 1
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
