CREATE TABLE symbols
(
    id      serial not null unique,
    symbol  varchar(255) not null,
    price   decimal not null,
    exchange varchar(255) not null,
    time    timestamp not null
);