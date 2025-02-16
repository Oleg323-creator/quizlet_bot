create table if not exists users(
    id serial primary key,
    tg_id integer unique,
    username varchar(50) not null,
    first_name varchar(50) not null,
    last_name varchar(50) not null
);

create table if not exists topics(
    id serial primary key,
    topic varchar(20) not null,
    user_id  integer references users(id) on delete cascade
);

create table if not exists words(
    id serial primary key,
    word varchar(50) not null,
    translation varchar(50) not null,
    user_id  integer references users(id) on delete cascade,
    topic_id integer references topics(id) on delete cascade
);

create table if not exists stats(
    id serial primary key,
    user_id integer unique references users(id) on delete cascade,
    stat integer default 0
);

create unique index  stats_user_id_idx on stats(user_id);


---- create above / drop below ----

drop table if exists users cascade;

drop table  if exists topics cascade;

drop table if exists words cascade;

drop table if exists stats cascade;