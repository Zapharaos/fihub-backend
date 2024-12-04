-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

ALTER TABLE "users" DROP COLUMN "first_name";
ALTER TABLE "users" DROP COLUMN "last_name";

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

ALTER TABLE "users" ADD COLUMN "first_name" varchar(100) NOT NULL;
ALTER TABLE "users" ADD COLUMN "last_name" varchar(100) NOT NULL;
