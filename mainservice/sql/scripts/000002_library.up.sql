CREATE TABLE gnr_library (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    user_id BIGINT NOT NUL REFERENCES gnr_user(id),
    token VARCHAR(50) NOT NULL
);