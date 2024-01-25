ALTER TABLE gnr_library ADD COLUMN IF NOT EXISTS is_online bool;


CREATE TABLE IF NOT EXISTS gnr_library_item (
    id bigserial primary key,
    title varchar(500),
    library_id bigint references gnr_library(id),
    status int
);