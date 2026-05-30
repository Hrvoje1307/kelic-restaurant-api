CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS admins (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    email         TEXT        UNIQUE NOT NULL,
    password_hash TEXT        NOT NULL,
    full_name     TEXT        NOT NULL DEFAULT '',
    role          TEXT        NOT NULL DEFAULT 'admin' CHECK (role IN ('admin', 'superadmin')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    email         TEXT        UNIQUE NOT NULL,
    password_hash TEXT        NOT NULL,
    full_name     TEXT        NOT NULL DEFAULT '',
    phone         TEXT        NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Default superadmin — password: Admin@123456  (change after first login)
INSERT INTO admins (email, password_hash, full_name, role)
VALUES (
    'admin@elnigo.de',
    crypt('Admin@123456', gen_salt('bf', 10)),
    'Super Admin',
    'superadmin'
);
