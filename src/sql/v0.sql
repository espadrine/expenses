PRAGMA journal_mode=WAL;
create table if not exists versions (
  id integer primary key,
  applied_at text
) without rowid;
insert into versions (id, applied_at) values (1, datetime());
create table if not exists users (
  id text primary key,  -- Random base32 ID.
  name text
) without rowid;
create table if not exists entry (
  user text references users,
  date text,
  id integer,
  amount decimal,
  currency text,
  label text,
  primary key (user, date, id)
) without rowid;
create table if not exists tags (
  user text references users,
  id integer,
  name text,
  primary key (user, id)
) without rowid;
create table if not exists entry_tags (
  user text,
  tag integer,
  entry_date text,
  entry integer,
  primary key (user, tag, entry_date, entry),
  foreign key (user, tag) references tags (user, id),
  foreign key (user, entry_date, entry) references entry (user, date, id)
) without rowid;
