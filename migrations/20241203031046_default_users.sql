-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE "users"
(
    "id"         uuid PRIMARY KEY NOT NULL DEFAULT (gen_random_uuid()),
    "email"      varchar(100)     NOT NULL,
    "password"   varchar(100)     NOT NULL,
    "first_name" varchar(100)     NOT NULL,
    "last_name"  varchar(100)     NOT NULL,
    "updated_at" timestamptz      NOT NULL,
    "created_at" timestamptz      NOT NULL DEFAULT (NOW())
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
