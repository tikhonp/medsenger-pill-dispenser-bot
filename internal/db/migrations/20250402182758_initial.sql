-- +goose Up
-- +goose StatementBegin

-- List of hardware types
-- For this time there are two: 
--   2x2 and 4x7 cells
create table hardware_type
(
    id   varchar(255) primary key not null,
    name varchar(255) not null
);
insert into hardware_type(id, name)
values ('HW_2x2_V1', '4 ячейки'),
       ('HW_4x7_V1', '28 ячеек');

create table pill_dispenser
(
    serial_number   varchar(255) primary key not null,
    hw_type_id  varchar(255) not null,
    last_fetch_time datetime default null,
    contract_id int default null,
    foreign key (hw_type_id) references hardware_type (id),
    foreign key (contract_id) references contract (id)
);

create table schedule
(
    id                               integer primary key autoincrement,
    is_offline_notifications_allowed boolean                              not null,
    refresh_rate_interval            integer, -- durations in seconds
    contract_id                      integer                            not null,
    pill_dispenser_sn                varchar(255)                       not null,
    created_at                       datetime default current_timestamp not null,
    foreign key (contract_id) references contract (id),
    foreign key (pill_dispenser_sn) references pill_dispenser (serial_number)
);

create table schedule_cell
(
    idx                  integer not null,
    schedule_id          integer not null,
    start_time           datetime,
    end_time             datetime,
    contents_description varchar(255),
    foreign key (schedule_id) references schedule (id),
    primary key (schedule_id, idx)
);

create table contract
(
    id                  integer primary key not null,
    is_active           boolean             not null,
    clinic_id           integer             not null,
    agent_token         varchar(255)        not null,
    patient_agent_token varchar(255)        not null,
    doctor_agent_token  varchar(255)        not null,
    locale              varchar(5)          not null,
    patient_name        varchar(255),
    patient_email       varchar(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table contract;
drop table schedule_cell;
drop table schedule;
drop table pill_dispenser;
drop table hardware_type;
-- +goose StatementEnd
