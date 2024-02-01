CREATE TABLE IF NOT EXISTS gnr_notification (
    id BIGSERIAL PRIMARY KEY,
    user_ref BIGINT REFERENCES gnr_user(id) NOT NULL,
    type_c INT NOT NULL,
    date_c DATE NOT NULL
);