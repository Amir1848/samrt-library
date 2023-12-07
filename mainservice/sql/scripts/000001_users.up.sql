CREATE TABLE gnr_user(
    id BIGSERIAL PRIMARY KEY,
    student_code VARCHAR(50) NOT NULL,
    password_c VARCHAR(500) NOT NULL
);


CREATE TABLE gnr_user_role(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES gnr_user(id),
    role_c INT NOT NULL
);

CREATE INDEX IF NOT EXISTS gnr_user_role_user_id_idx ON gnr_user_role(user_id);