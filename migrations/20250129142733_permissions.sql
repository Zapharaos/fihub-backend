-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE "permissions"
(
    "id"          uuid PRIMARY KEY    NOT NULL DEFAULT (gen_random_uuid()),
    "value"       varchar(100) UNIQUE NOT NULL,
    "scope"       varchar(255)        NOT NULL,
    "description" varchar(255)        NOT NULL
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE "permissions";
