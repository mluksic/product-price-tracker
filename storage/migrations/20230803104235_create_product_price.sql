-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table if not exists product_price(
    id bigserial not null primary key,
    name varchar(255) not null,
    price int not null,
    fetched_at timestamp(0) without time zone not null default CURRENT_TIMESTAMP
);

alter table product_price add column if not exists product_id int;
alter table product_price add constraint fk_product_id foreign key (product_id) references product(id);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

alter table product_price drop constraint fk_product_id;
drop table if exists product_price;
