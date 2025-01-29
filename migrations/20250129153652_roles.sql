-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE "roles"
(
    "id"   uuid PRIMARY KEY    NOT NULL DEFAULT (gen_random_uuid()),
    "name" varchar(100) UNIQUE NOT NULL
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE "roles";