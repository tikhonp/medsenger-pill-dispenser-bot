-- +goose Up
-- +goose StatementBegin

-- Add description column to hardware_type table with default empty string
ALTER TABLE hardware_type
ADD COLUMN description TEXT NOT NULL DEFAULT '';

-- Add description column to pill_dispenser table with default empty string
ALTER TABLE pill_dispenser
ADD COLUMN description TEXT NOT NULL DEFAULT '';

-- Insert descriptions for existing hardware types
UPDATE hardware_type
SET description = 'V1 is for REV1.1 and REV1.2 pill dispensers with 4 compartments and external buttons mount.'
WHERE id = 'HW_2X2_V1';

UPDATE hardware_type
SET description = 'V1 is for REV1.1 and REV1.2 pill dispensers with 28 compartments and external buttons mount.'
WHERE id = 'HW_4X7_V1';

INSERT INTO hardware_type (id, name, description)
VALUES ('HW_2X2_V2', '4 ячейки (REV2)',
        'V2 is for REV1.3 pill dispensers with 4 compartments and on board buttons.');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM hardware_type WHERE id = 'HW_2X2_V2';

ALTER TABLE hardware_type
DROP COLUMN description;

ALTER TABLE pill_dispenser
DROP COLUMN description;
-- +goose StatementEnd
