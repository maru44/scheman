CREATE DATABASE scheman;
GO

USE scheman;
GO

CREATE TABLE magic
(
  id int NOT NULL IDENTITY (1,1) PRIMARY KEY,
  id_two int NOT NULL,
  id_three int,
  bit_zero bit,
  bit_one bit NULL,
  bit_two bit NOT NULL,
  bit_three bit NULL DEFAULT 0,
  bit_four bit NULL DEFAULT 1,
  bit_five bit NOT NULL DEFAULT 0,
  bit_six bit NOT NULL DEFAULT 1,
  string_zero VARCHAR(1),
  string_one VARCHAR(1) NULL,
  string_two VARCHAR(1) NOT NULL,
  string_three VARCHAR(1) NULL DEFAULT 'a',
  string_four VARCHAR(1) NOT NULL DEFAULT 'b',
  string_five VARCHAR(1000),
  string_six VARCHAR(1000) NULL,
  string_seven VARCHAR(1000) NOT NULL,
  string_eight VARCHAR(1000) NULL DEFAULT 'abcdefgh',
  string_nine VARCHAR(1000) NOT NULL DEFAULT 'abcdefgh',
  string_ten VARCHAR(1000) NULL DEFAULT '',
  string_eleven VARCHAR(1000) NOT NULL DEFAULT '',
  big_int_zero bigint,
  big_int_one bigint NULL,
  big_int_two bigint NOT NULL,
  big_int_three bigint NULL DEFAULT 111111,
  big_int_four bigint NOT NULL DEFAULT 222222,
  big_int_five bigint NULL DEFAULT 0,
  big_int_six bigint NOT NULL DEFAULT 0,
  int_zero int,
  int_one int NULL,
  int_two int NOT NULL,
  int_three int NULL DEFAULT 333333,
  int_four int NOT NULL DEFAULT 444444,
  int_five int NULL DEFAULT 0,
  int_six int NOT NULL DEFAULT 0,
  float_zero float,
  float_one float,
  float_two float(24),
  float_three float(24),
  float_four float(24) NULL,
  float_five float(24) NOT NULL,
  float_six float(24) NULL DEFAULT 1.1,
  float_seven float(24) NOT NULL DEFAULT 1.1,
  float_eight float(24) NULL DEFAULT 0.0,
  float_nine float(24) NULL DEFAULT 0.0,
  bytea_zero binary NOT NULL,
  bytea_one binary NOT NULL,
  bytea_two binary NOT NULL,
  bytea_three binary NOT NULL DEFAULT CONVERT(VARBINARY(MAX),'a'),
  bytea_four binary NOT NULL DEFAULT CONVERT(VARBINARY(MAX),'b'),
  bytea_five binary(100) NOT NULL DEFAULT CONVERT(VARBINARY(MAX),'abcdefghabcdefghabcdefgh'),
  bytea_six binary(100) NOT NULL DEFAULT  CONVERT(VARBINARY(MAX),'hgfedcbahgfedcbahgfedcba'),
  bytea_seven binary NOT NULL DEFAULT CONVERT(VARBINARY(MAX),''),
  bytea_eight binary NOT NULL DEFAULT CONVERT(VARBINARY(MAX),''),
  time_zero timestamp NOT NULL,
  time_one date,
  time_eleven date NULL,
  time_twelve date NOT NULL,
  time_fifteen date NULL DEFAULT '19990108',
  time_sixteen date NOT NULL DEFAULT '1999-01-08'
);
GO

CREATE TABLE magicest
(
  id int NOT NULL IDENTITY (1,1) PRIMARY KEY,
  kk float NULL,
  ll float NOT NULL,
  mm tinyint NULL,
  nn tinyint NOT NULL,
  oo bit NULL,
  pp bit NOT NULL,
  qq smallint NULL,
  rr smallint NOT NULL,
  ss int NULL,
  tt int NOT NULL,
  uu bigint NULL,
  vv bigint NOT NULL,
  ww float NULL,
  xx float NOT NULL,
  yy float NULL,
  zz float NOT NULL,
  aaa double precision NULL,
  bbb double precision NOT NULL,
  ccc real NULL,
  ddd real NOT NULL,
  ggg date NULL,
  hhh date NOT NULL,
  iii datetime NULL,
  jjj datetime NOT NULL,
  kkk timestamp NOT NULL,
  mmm binary NOT NULL,
  nnn binary NOT NULL,
  ooo varbinary(100) NOT NULL,
  ppp varbinary(100) NOT NULL,
  qqq varbinary NOT NULL,
  rrr varbinary NOT NULL,
  www varbinary(max) NOT NULL,
  xxx varbinary(max) NOT NULL,
  yyy varchar(100) NULL,
  zzz varchar(100) NOT NULL,
  aaaa char NULL,
  bbbb char NOT NULL,
  cccc VARCHAR(MAX) NULL,
  dddd VARCHAR(MAX) NOT NULL,
  eeee tinyint NULL,
  ffff tinyint NOT NULL
);
GO

create table owner
(
  id int NOT NULL IDENTITY (1,1) PRIMARY KEY,
  name varchar(255) not null
);
GO

create table cats
(
  id int NOT NULL IDENTITY (1,1) PRIMARY KEY,
  name varchar(255) not null,
  owner_id int
);
GO

ALTER TABLE cats ADD CONSTRAINT cats_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES owner(id);
GO

create table toys
(
  id int NOT NULL IDENTITY (1,1) PRIMARY KEY,
  name varchar(255) not null
);
GO

create table cat_toys
(
  cat_id int not null references cats (id),
  toy_id int not null references toys (id),
  primary key (cat_id, toy_id)
);
GO
