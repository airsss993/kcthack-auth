CREATE TABLE users_sessions
(
    id         uuid                    not null primary key,
    user_id    uuid                    not null references users (id),
    token_hash varchar(255) unique      not null,
    expires_at timestamp               not null,
    created_at timestamp DEFAULT NOW() not null
);
CREATE INDEX userSessionsUserID_index ON users_sessions (user_id);
CREATE INDEX userSessionsTokenHash_index ON users_sessions (token_hash);
CREATE INDEX userSessionsExpAt_index ON users_sessions (expires_at);