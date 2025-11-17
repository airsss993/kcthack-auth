CREATE TABLE users
(
    id          uuid                    not null primary key,
    first_name  varchar(255)            not null,
    last_name   varchar(255)            not null,
    role        varchar(50)             not null,
    email       varchar(100) unique     not null,
    tg_name     varchar(100) unique     null,
    birth_date  date                    null,
    bio         text                    null,
    pass_hash   varchar(255)            not null,
    is_verified bool      default false not null,
    updated_at  timestamp DEFAULT NOW() not null,
    created_at  timestamp DEFAULT NOW() not null
);
CREATE INDEX userEmail_index ON users (email);