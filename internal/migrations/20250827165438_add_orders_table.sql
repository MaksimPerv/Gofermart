-- +goose Up
-- +goose StatementBegin

CREATE TABLE order_statuses
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(50) NOT NULL UNIQUE,
    description TEXT
);

INSERT INTO order_statuses (name, description)
VALUES ('NEW', 'Заказ загружен в систему, но не попал в обработку'),
       ('PROCESSING', 'Вознаграждение за заказ рассчитывается'),
       ('INVALID', 'Система расчёта вознаграждений отказала в расчёте'),
       ('PROCESSED', 'Данные по заказу проверены и информация о расчёте успешно получена');


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_statuses;
-- +goose StatementEnd
