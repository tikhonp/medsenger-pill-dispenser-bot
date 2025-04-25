-- +goose Up
-- +goose StatementBegin

-- drop not null constraint on pill_dispenser_sn
alter table schedule
    rename to schedule_old;
create table schedule
(
    id                               integer primary key autoincrement,
    is_offline_notifications_allowed boolean                            not null,
    refresh_rate_interval            integer, -- durations in seconds
    contract_id                      integer                            not null,
    pill_dispenser_sn                varchar(255),
    created_at                       datetime default current_timestamp not null,
    foreign key (contract_id) references contract (id),
    foreign key (pill_dispenser_sn) references pill_dispenser (serial_number)
);
insert into schedule(id, is_offline_notifications_allowed, refresh_rate_interval, contract_id, pill_dispenser_sn,
                     created_at)
select id, is_offline_notifications_allowed, refresh_rate_interval, contract_id, pill_dispenser_sn, created_at
from schedule_old;
drop table schedule_old;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table schedule
    rename to schedule_old;
create table schedule
(
    id                               integer primary key autoincrement,
    is_offline_notifications_allowed boolean                            not null,
    refresh_rate_interval            integer, -- durations in seconds
    contract_id                      integer                            not null,
    pill_dispenser_sn                varchar(255)                       not null,
    created_at                       datetime default current_timestamp not null,
    foreign key (contract_id) references contract (id),
    foreign key (pill_dispenser_sn) references pill_dispenser (serial_number)
);
insert into schedule(id, is_offline_notifications_allowed, refresh_rate_interval, contract_id, pill_dispenser_sn,
                     created_at)
select id, is_offline_notifications_allowed, refresh_rate_interval, contract_id, pill_dispenser_sn, created_at
from schedule_old;
drop table schedule_old;
-- +goose StatementEnd
