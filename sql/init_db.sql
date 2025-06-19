CREATE TABLE logins (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    login_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,                      -- Unique identifier for each role (auto-incrementing)
    name TEXT UNIQUE NOT NULL,                  -- Role name (e.g., 'admin', 'editor', 'viewer')
    permissions TEXT[] NOT NULL,                -- Array of permissions associated with the role
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW() -- Timestamp when the role was created
);

ALTER TABLE roles ADD CONSTRAINT uni_roles_name UNIQUE (name);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,                      -- Unique identifier (auto-incrementing)
    first_name VARCHAR(100) NOT NULL,           -- User's first name (required)
    last_name VARCHAR(100) NOT NULL,            -- User's last name (required)
    email VARCHAR(255) UNIQUE NOT NULL,         -- Unique email (required)
    password_hash TEXT NOT NULL,                -- Hashed password for security
    role_id INTEGER NOT NULL,                   -- Foreign key to roles table
    is_active BOOLEAN DEFAULT TRUE,             -- Status (active/inactive)
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, -- Last update timestamp
    CONSTRAINT fk_users_role
        FOREIGN KEY (role_id) REFERENCES roles(id)
);
