create table if not exists rates
(
    ask         numeric(20, 10) not null,
    bid         numeric(20, 10) not null,
    received_at timestamptz     not null default now(),

    primary key (received_at)
)