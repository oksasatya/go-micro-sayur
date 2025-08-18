CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(17) NULL,
    photo VARCHAR(255) NULL,
    address text NULL,
    lat varchar(50) NULL,
    lng varchar(50) NULL,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at timestamp null,
    deleted_at timestamp null

);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_is_verified ON users(is_verified);