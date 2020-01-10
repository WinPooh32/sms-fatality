-- +goose Up
CREATE TABLE message (
    id    bigserial NOT NULL PRIMARY KEY,
    phone varchar(11),
    body  varchar(2048)
);

INSERT INTO message VALUES
(0, '72888137328', 'Пробное сообщение бла бла бла'),
(1, '82888137328', 'Test message');

-- +goose Down
DROP TABLE message;

