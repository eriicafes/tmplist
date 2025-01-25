-- +goose Up
-- +goose StatementBegin
create table users (
    id serial primary key,
    email text unique not null,
    password_hash text not null,
    created_at timestamp default current_timestamp
);

create table sessions (
    id text primary key,
    user_id integer not null,
    expires_at timestamp,
    foreign key (user_id) references users(id) on delete cascade
);

create table topics (
    id serial primary key,
    user_id integer not null,
    title text not null,
    created_at timestamp default current_timestamp,
    foreign key (user_id) references users(id) on delete cascade
);

create table todos (
    id serial primary key,
    topic_id integer not null,
    body text not null,
    done boolean not null,
    created_at timestamp default current_timestamp,
    foreign key (topic_id) references topics(id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users, sessions, topics, items;
-- +goose StatementEnd
