CREATE TABLE IF NOT EXISTS gnr_student_history(
    id bigserial primary key,
    user_ref bigint references gnr_user(id) not null,
    date_c date not null,
    library_item_ref bigint references gnr_library_item(id)
);