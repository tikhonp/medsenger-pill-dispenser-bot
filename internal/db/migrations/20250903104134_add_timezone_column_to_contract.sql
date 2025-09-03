-- +goose Up
-- +goose StatementBegin

-- add timezone string column to contract table
ALTER TABLE contract ADD COLUMN timezone VARCHAR(255);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE contract DROP COLUMN timezone;
-- +goose StatementEnd
