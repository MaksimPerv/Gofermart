-- +goose Up
-- +goose StatementBegin
CREATE TABLE history_withdrawal
(
    user_id      INTEGER        NOT NULL REFERENCES users(id),
    number VARCHAR(255) NOT NULL UNIQUE ,
    sum DECIMAL(10, 2) NOT NULL,
    processed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, number)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS history_withdrawal;
-- +goose StatementEnd
