-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table if not exists product(
    id bigserial primary key,
    name varchar(255) not null,
    is_tracked bool not null default true,
    created_at timestamp(0) without time zone not null default CURRENT_TIMESTAMP,
    updated_at timestamp(0) without time zone not null default CURRENT_TIMESTAMP
);

create index product_name_idx on product(name);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop index if exists product_name_idx;
drop table if exists product;
