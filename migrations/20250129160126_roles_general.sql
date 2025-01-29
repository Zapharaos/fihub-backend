-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE "user_roles"
(
    "user_id" uuid NOT NULL,
    "role_id" uuid NOT NULL,
    PRIMARY KEY ("user_id", "role_id"),
    FOREIGN KEY ("user_id") REFERENCES "users" (id) ON DELETE CASCADE,
    FOREIGN KEY ("role_id") REFERENCES "roles" (id) ON DELETE CASCADE
);

CREATE TABLE "role_permissions"
(
    "role_id"       uuid NOT NULL,
    "permission_id" uuid NOT NULL,
    PRIMARY KEY ("role_id", "permission_id"),
    FOREIGN KEY ("role_id") REFERENCES "roles" (id) ON DELETE CASCADE,
    FOREIGN KEY ("permission_id") REFERENCES "permissions" (id) ON DELETE CASCADE
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE "user_roles";
DROP TABLE "role_permissions";