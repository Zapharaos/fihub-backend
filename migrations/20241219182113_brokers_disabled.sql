-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

ALTER TABLE brokers
    ADD COLUMN disabled BOOLEAN DEFAULT FALSE;

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

ALTER TABLE brokers
    DROP COLUMN disabled;
