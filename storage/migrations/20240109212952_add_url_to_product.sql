-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
alter table product_price drop column url;
alter table product add column url text default null;
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
alter table product_price add column url text default null;
alter table product drop column url;
-- +goose StatementEnd
