create table qrscans (
  id serial primary key,
  uuid varchar(255) not null,
  user_id integer default null,
  is_valid boolean not null default false,
  valid_until timestamp not null,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  deleted_at timestamp
)