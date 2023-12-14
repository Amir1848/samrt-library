CREATE TABLE gnr_library (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    user_id BIGINT NOT NULL REFERENCES gnr_user(id) UNIQUE,
    token VARCHAR(50) NOT NULL
);