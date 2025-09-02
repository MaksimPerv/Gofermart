-- +goose Up
-- +goose StatementBegin
CREATE TABLE balance
(
    user_id   INTEGER PRIMARY KEY NOT NULL REFERENCES users (id),
    current   DECIMAL(10, 2) DEFAULT 0,
    withdrawn DECIMAL(10, 2) DEFAULT 0
);

INSERT INTO balance (user_id, current, withdrawn)
SELECT id, 0, 0
FROM users;

CREATE OR REPLACE FUNCTION create_balance_for_new_user()
    RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO balance (user_id, current, withdrawn)
    VALUES (NEW.id, 0, 0)
    ON CONFLICT (user_id) DO NOTHING;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_create_balance
    AFTER INSERT
    ON users
    FOR EACH ROW
EXECUTE FUNCTION create_balance_for_new_user();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_create_balance ON users;

DROP FUNCTION IF EXISTS create_balance_for_new_user();

DROP TABLE IF EXISTS balcance;
-- +goose StatementEnd
