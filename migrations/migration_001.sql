language plpgsql;

-- +migrate Up
CREATE TABLE
    users (
        id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
        name varchar(255) NOT NULL,
        email varchar(255) NOT NULL UNIQUE,
        email_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
        password_hash TEXT NOT NULL,
        is_blocked BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP
    );

CREATE TABLE
    roles (
        id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
        name varchar(255) NOT NULL UNIQUE,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP,
        permissions TEXT NOT NULL DEFAULT ''
    );

CREATE TABLE
    user_roles (
        user_id uuid REFERENCES users (id) ON DELETE CASCADE,
        role_id uuid REFERENCES roles (id) ON DELETE CASCADE,
        PRIMARY KEY (user_id, role_id)
    );

-- +migrate Down
DROP TABLE user_roles;

DROP TABLE roles;

DROP TABLE users;
