-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

ALTER TABLE brokers
    ADD COLUMN image_id UUID DEFAULT NULL;

CREATE TABLE broker_image
(
    id         uuid PRIMARY KEY NOT NULL DEFAULT (gen_random_uuid()),
    broker_id  UUID             NOT NULL,
    name       VARCHAR(255)     NOT NULL,
    data       BYTEA            NOT NULL,
    FOREIGN KEY (broker_id) REFERENCES brokers (id) ON DELETE CASCADE
);

ALTER TABLE brokers
    ADD FOREIGN KEY (image_id) REFERENCES broker_image (id) ON DELETE SET NULL;

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

ALTER TABLE brokers
    DROP COLUMN image_id;

DROP TABLE broker_image;
