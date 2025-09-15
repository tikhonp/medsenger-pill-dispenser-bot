-- +goose Up
-- +goose StatementBegin
INSERT INTO hardware_type (id, name, description)
VALUES ('HW_4X7_V2', '28 ячеек (REV2)',
        'V2 is for REV1.3 pill dispensers with 28 compartments and on board buttons.');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM hardware_type WHERE id = 'HW_4X7_V2';
-- +goose StatementEnd
