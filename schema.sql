create table users (name text) if not exists;
create table outflows (
  user references users,
  date text,
  id integer,
  product text,
  amount decimal,
  currency text,
  primary key (user, date, id)
) without rowid if not exists;
create table tags (
  user references users,
  name text,
  primary key (user, name)
) without rowid if not exists;
create table outflow_tags (
  user,
  name text,
  outflow_date text,
  outflow_id integer,
  primary key (user, name, outflow_date, outflow_id),
  foreign key (user, name) references tags (user, name),
  foreign key (user, outflow_date, outflow_id) references outflows (user, date, id)
) without rowid if not exists;
