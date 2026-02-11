-- +goose Up
-- +goose StatementBegin
ALTER TABLE contract
    DROP COLUMN agent_token,
    DROP COLUMN patient_agent_token,
    DROP COLUMN doctor_agent_token;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE contract
    ADD COLUMN agent_token varchar(255) not null default '',
    ADD COLUMN patient_agent_token varchar(255) not null default '',
    ADD COLUMN doctor_agent_token varchar(255) not null default '';
-- +goose StatementEnd
