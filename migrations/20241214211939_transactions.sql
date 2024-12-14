-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE "transactions"
(
    "id"               uuid PRIMARY KEY NOT NULL DEFAULT (gen_random_uuid()),
    "user_id"          uuid             NOT NULL,
    "broker_id"        uuid             NOT NULL,
    "date"             timestamp        NOT NULL,
    "transaction_type" varchar(100)     NOT NULL,
    "asset"            varchar(100)     NOT NULL,
    "quantity"         numeric          NOT NULL,
    "price"            numeric          NOT NULL,
    "price_unit"       numeric          NOT NULL,
    "fee"              numeric          NOT NULL,

    FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("broker_id") REFERENCES "brokers" ("id") ON DELETE CASCADE
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop table if exists transactions;