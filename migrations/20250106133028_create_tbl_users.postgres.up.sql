create table users (
  id serial primary key,
  username varchar(255) not null,
  password varchar(255) not null,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  deleted_at timestamp
);