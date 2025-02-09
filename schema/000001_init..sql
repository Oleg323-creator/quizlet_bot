create table if not exists users(
    id serial primary key unique,
    tg_id integer unique,
    username varchar(50) not null,
    first_name varchar(50) not null,
    last_name varchar(50) not null
);

create table if not exists topics(
    id serial primary key unique,
    topic varchar(20) not null,
    user_id  integer references users(id) on delete cascade
);

create table if not exists words(
    word varchar(50) not null,
    translate varchar(50) not null,
    topic_id integer references topics(id) on delete cascade
);

create table if not exists stats(
    id serial primary key unique,
    user_id integer references users(id) on delete cascade,
    topic_id integer references topics(id) on delete cascade,
    stats integer
);

---- create above / drop below ----

drop table if exists users;

drop table  if exists topics;

drop table if exists words;

drop table if exists stats;