-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
alter table product_price drop column name;
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
alter table product_price add column name varchar(255) not null default '';
-- +goose StatementEnd
