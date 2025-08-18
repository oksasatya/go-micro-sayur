CREATE TABLE if not exists "user_role" (
    id SERIAL PRIMARY KEY,
    roles_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (roles_id) REFERENCES roles(id) ON DELETE CASCADE ON UPDATE CASCADE
);