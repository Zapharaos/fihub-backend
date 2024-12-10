-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE "brokers"
(
    "id"   uuid PRIMARY KEY    NOT NULL DEFAULT (gen_random_uuid()),
    "name" varchar(100) UNIQUE NOT NULL
);

CREATE TABLE "user_brokers"
(
    "user_id"   uuid NOT NULL,
    "broker_id" uuid NOT NULL,
    PRIMARY KEY ("user_id", "broker_id"),
    FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("broker_id") REFERENCES "brokers" ("id") ON DELETE CASCADE
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop table if exists user_brokers;

drop table if exists brokers;
