-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE password_reset_tokens
(
    id         uuid PRIMARY KEY NOT NULL DEFAULT (gen_random_uuid()),
    user_id    uuid             NOT NULL,
    token      VARCHAR(255)     NOT NULL,
    expires_at timestamptz      NOT NULL,
    created_at timestamptz      NOT NULL DEFAULT (NOW()),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE password_reset_tokens;
