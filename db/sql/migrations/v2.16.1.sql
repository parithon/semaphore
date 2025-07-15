create table azure_config (
  id integer primary key autoincrement,
  created datetime not null,
  name varchar(100) not null,
  description text,
  tenant_id varchar(100) not null,
  client_id varchar(100) not null,
  client_secret varchar(500) not null
);