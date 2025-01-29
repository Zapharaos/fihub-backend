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

CREATE TABLE "roles"
(
    "id"   uuid PRIMARY KEY    NOT NULL DEFAULT (gen_random_uuid()),
    "name" varchar(100) UNIQUE NOT NULL
);

CREATE TABLE "role_permissions"
(
    "role_id"       uuid NOT NULL,
    "permission_id" uuid NOT NULL,
    PRIMARY KEY ("role_id", "permission_id"),
    FOREIGN KEY ("role_id") REFERENCES "roles" (id) ON DELETE CASCADE,
    FOREIGN KEY ("permission_id") REFERENCES "permissions" (id) ON DELETE CASCADE
);

CREATE TABLE "user_roles"
(
    "user_id" uuid NOT NULL,
    "role_id" uuid NOT NULL,
    PRIMARY KEY ("user_id", "role_id"),
    FOREIGN KEY ("user_id") REFERENCES "users" (id) ON DELETE CASCADE,
    FOREIGN KEY ("role_id") REFERENCES "roles" (id) ON DELETE CASCADE
);

INSERT INTO permissions (id, value, scope, description)
VALUES
    -- Superadmin wildcard
    ('28b1a126-f340-476d-a369-0d44d7af9f3f', '*', 'all', 'Superadmin wildcard'),

    -- Roles
    ('bb1f32a9-7de5-4ebe-9993-083906183701', 'admin.roles.create', 'admin', 'Create role'),
    ('89f7d88f-f89d-4a9c-9e84-a67ae2ba6307', 'admin.roles.read', 'admin', 'Read role'),
    ('a0bee99f-36b7-4ebc-b780-87efbceabd18', 'admin.roles.update', 'admin', 'Update role'),
    ('d5c0d0e9-3c11-4857-b28c-c09d4aa6063c', 'admin.roles.delete', 'admin', 'Delete role'),
    ('28987a68-dfb7-4e2d-b07b-9e9fa1efd608', 'admin.roles.list', 'admin', 'List roles'),

    -- Permissions
    ('f2ac4ddf-5698-43e9-be56-712877d34d3b', 'admin.permissions.create', 'admin', 'Create permission'),
    ('846dc244-c2cc-4e0a-b18a-1a2d07ff1b9d', 'admin.permissions.read', 'admin', 'Read permission'),
    ('0b1de3c4-f161-4313-bd55-905bf71ea8b6', 'admin.permissions.update', 'admin', 'Update permission'),
    ('764c4961-82d4-41fe-aa8c-d6b05986f58c', 'admin.permissions.delete', 'admin', 'Delete permission'),
    ('ae2cf96d-5bdb-4ed0-aa0a-402224e3b94b', 'admin.permissions.list', 'admin', 'List permissions'),

    -- User roles
    ('2fa33825-e636-44c5-9193-e39d4909b826', 'admin.users.roles.list', 'admin', 'List user role'),
    ('8cc7abbf-d640-4cd2-a872-6537663782b8', 'admin.users.roles.update', 'admin', 'Update user role'),

    -- Role permissions
    ('b8017630-d702-48fc-84b6-7c0d17b59a56', 'admin.roles.permissions.list', 'admin', 'List role permission'),
    ('1302a210-63fd-411c-9bd4-022e3f4b84d4', 'admin.roles.permissions.update', 'admin', 'Update role permission')
;

-- Insert super admin role
INSERT INTO roles (id, name)
VALUES ('d2b1a126-f340-476d-a369-0d44d7af9f3f', 'superadmin');

-- Assign super admin wildcard permission to super admin role
INSERT INTO role_permissions (role_id, permission_id)
VALUES ('d2b1a126-f340-476d-a369-0d44d7af9f3f', '28b1a126-f340-476d-a369-0d44d7af9f3f');

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE "user_roles";
DROP TABLE "role_permissions";
DROP TABLE "roles";
DROP TABLE "permissions";
