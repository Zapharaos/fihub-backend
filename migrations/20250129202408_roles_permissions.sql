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

    -- Admin overview
    ('10781004-bfa0-4a3e-b07c-156f868481da', 'front.admin.overview', 'front', 'Admin overview view'),

    -- Permissions
    ('f2ac4ddf-5698-43e9-be56-712877d34d3b', 'admin.permissions.create', 'admin', 'Create permission'),
    ('846dc244-c2cc-4e0a-b18a-1a2d07ff1b9d', 'admin.permissions.read', 'admin', 'Read permission'),
    ('0b1de3c4-f161-4313-bd55-905bf71ea8b6', 'admin.permissions.update', 'admin', 'Update permission'),
    ('764c4961-82d4-41fe-aa8c-d6b05986f58c', 'admin.permissions.delete', 'admin', 'Delete permission'),
    ('ae2cf96d-5bdb-4ed0-aa0a-402224e3b94b', 'admin.permissions.list', 'admin', 'List permissions'),

    -- Roles
    ('a1eb52c5-d517-4b39-9657-c2ff8786a34a', 'front.admin.roles', 'front', 'Admin roles view'),
    ('bb1f32a9-7de5-4ebe-9993-083906183701', 'admin.roles.create', 'admin', 'Create role'),
    ('89f7d88f-f89d-4a9c-9e84-a67ae2ba6307', 'admin.roles.read', 'admin', 'Read role'),
    ('a0bee99f-36b7-4ebc-b780-87efbceabd18', 'admin.roles.update', 'admin', 'Update role'),
    ('d5c0d0e9-3c11-4857-b28c-c09d4aa6063c', 'admin.roles.delete', 'admin', 'Delete role'),
    ('28987a68-dfb7-4e2d-b07b-9e9fa1efd608', 'admin.roles.list', 'admin', 'List roles'),

    -- Role permissions
    ('b8017630-d702-48fc-84b6-7c0d17b59a56', 'admin.roles.permissions.list', 'admin', 'List role permission'),
    ('1302a210-63fd-411c-9bd4-022e3f4b84d4', 'admin.roles.permissions.update', 'admin', 'Update role permission'),

    -- User
    ('c6ea4ccc-5d10-43ea-995c-6edf21fe5ad9', 'front.admin.users', 'front', 'Admin users view'),
    ('341ce3de-136a-419e-b837-3dc3d5aa8f2d', 'admin.users.read', 'admin', 'Read user'),
    ('ca0e7cc2-55f4-4700-a84b-a5e5694f697f', 'admin.users.update', 'admin', 'Update user'),
    ('55363c9b-f8db-4d37-8fb7-df1f3674897e', 'admin.users.list', 'admin', 'List users'),

    -- User roles
    ('2fa33825-e636-44c5-9193-e39d4909b826', 'admin.users.roles.list', 'admin', 'List user role'),
    ('8cc7abbf-d640-4cd2-a872-6537663782b8', 'admin.users.roles.update', 'admin', 'Update user role'),

    -- Brokers
    ('3bbeaf54-c92c-405f-a403-b01b9dc2ba7f', 'front.admin.brokers', 'front', 'Admin brokers view'),
    ('a3a92759-8d51-4564-b536-9c092b600471', 'admin.brokers.create', 'admin', 'Create broker'),
    ('bc2343a0-f99c-4a05-8172-271cfad4f422', 'admin.brokers.read', 'admin', 'Read broker'),
    ('97e63f39-9e19-4e7d-84ce-6fd95ed95e26', 'admin.brokers.update', 'admin', 'Update broker'),
    ('d01adea7-4f46-4e35-a286-190099ac9299', 'admin.brokers.delete', 'admin', 'Delete broker'),
    ('3b0ea488-de95-49e6-a084-73a80eb8b4fc', 'admin.brokers.list', 'admin', 'List brokers')
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
