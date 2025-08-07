-- +goose Up
-- +goose StatementBegin
create table battery_status
(
    id          serial primary key,
    serial_nu   varchar(255) not null,
    voltage     integer not null, -- in millivolts
    created_at  timestamp default current_timestamp not null,
    foreign key (serial_nu) references pill_dispenser (serial_number)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table battery_status;
-- +goose StatementEnd
