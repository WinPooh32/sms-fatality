-- +goose Up
-- +goose StatementBegin
ALTER TABLE message
    ALTER COLUMN phone SET NOT NULL,
    ALTER COLUMN body SET NOT NULL,
    ADD CONSTRAINT phonechk CHECK (char_length(phone) >= 1),
    ADD CONSTRAINT bodychk CHECK (char_length(body) >= 1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE message
    ALTER COLUMN phone DROP NOT NULL,
    ALTER COLUMN body DROP NOT NULL,
    DROP CONSTRAINT phonechk,
    DROP CONSTRAINT bodychk;
-- +goose StatementEnd
