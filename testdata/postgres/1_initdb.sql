CREATE DATABASE scheman;
\c scheman;

CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TYPE workday AS ENUM('monday', 'tuesday', 'wednesday', 'thursday', 'friday');

CREATE TABLE event_zero (
  id  uuid DEFAULT gen_random_uuid() NOT NULL,
  id2 serial NOT NULL,
  day workday NOT NULL,

  PRIMARY KEY (id)
);

CREATE TABLE multi_pks (
  id1 integer NOT NULL,
  id2 integer NOT NULL,

  PRIMARY KEY(id1, id2)
);

-- CREATE VIEW magic_views AS
-- SELECT e.id event_id, m.id1 multi_id
-- FROM event_zero e
-- INNER JOIN multi_pks m ON e.id2 = m.id2;

CREATE TABLE magic (
  id       serial PRIMARY KEY NOT NULL,
  id_two   serial NOT NULL,
  id_three serial,

  bool_zero   bool,
  bool_one    bool NULL,
  bool_two    bool NOT NULL,
  bool_three  bool NULL DEFAULT FALSE,
  bool_four   bool NULL DEFAULT TRUE,
  bool_five   bool NOT NULL DEFAULT FALSE,
  bool_six    bool NOT NULL DEFAULT TRUE,

  string_zero   VARCHAR(1),
  string_one    VARCHAR(1) NULL,
  string_two    VARCHAR(1) NOT NULL,
  string_three  VARCHAR(1) NULL DEFAULT 'a',
  string_four   VARCHAR(1) NOT NULL DEFAULT 'b',
  string_five   VARCHAR(1000),
  string_six    VARCHAR(1000) NULL,
  string_seven  VARCHAR(1000) NOT NULL,
  string_eight  VARCHAR(1000) NULL DEFAULT 'abcdefgh',
  string_nine   VARCHAR(1000) NOT NULL DEFAULT 'abcdefgh',
  string_ten    VARCHAR(1000) NULL DEFAULT '',
  string_eleven VARCHAR(1000) NOT NULL DEFAULT '',

  nonbyte_zero   CHAR(1),
  nonbyte_one    CHAR(1) NULL,
  nonbyte_two    CHAR(1) NOT NULL,
  nonbyte_three  CHAR(1) NULL DEFAULT 'a',
  nonbyte_four   CHAR(1) NOT NULL DEFAULT 'b',
  nonbyte_five   CHAR(1000),
  nonbyte_six    CHAR(1000) NULL,
  nonbyte_seven  CHAR(1000) NOT NULL,
  nonbyte_eight  CHAR(1000) NULL DEFAULT 'a',
  nonbyte_nine   CHAR(1000) NOT NULL DEFAULT 'b',

  byte_zero   "char",
  byte_one    "char" NULL,
  byte_two    "char" NULL DEFAULT 'a',
  byte_three  "char" NOT NULL,
  byte_four   "char" NOT NULL DEFAULT 'b',

  big_int_zero  bigint,
  big_int_one   bigint NULL,
  big_int_two   bigint NOT NULL,
  big_int_three bigint NULL DEFAULT 111111,
  big_int_four  bigint NOT NULL DEFAULT 222222,
  big_int_five  bigint NULL DEFAULT 0,
  big_int_six   bigint NOT NULL DEFAULT 0,

  int_zero  int,
  int_one   int NULL,
  int_two   int NOT NULL,
  int_three int NULL DEFAULT 333333,
  int_four  int NOT NULL DEFAULT 444444,
  int_five  int NULL DEFAULT 0,
  int_six   int NOT NULL DEFAULT 0,

  float_zero  decimal,
  float_one   numeric,
  float_two   numeric(2,1),
  float_three numeric(2,1),
  float_four  numeric(2,1) NULL,
  float_five  numeric(2,1) NOT NULL,
  float_six   numeric(2,1) NULL DEFAULT 1.1,
  float_seven numeric(2,1) NOT NULL DEFAULT 1.1,
  float_eight numeric(2,1) NULL DEFAULT 0.0,
  float_nine  numeric(2,1) NULL DEFAULT 0.0,

  bytea_zero  bytea,
  bytea_one   bytea NULL,
  bytea_two   bytea NOT NULL,
  bytea_three bytea NOT NULL DEFAULT 'a',
  bytea_four  bytea NULL DEFAULT 'b',
  bytea_five  bytea NOT NULL DEFAULT 'abcdefghabcdefghabcdefgh',
  bytea_six   bytea NULL DEFAULT 'hgfedcbahgfedcbahgfedcba',
  bytea_seven bytea NOT NULL DEFAULT '',
  bytea_eight bytea NOT NULL DEFAULT '',

  time_zero      timestamp,
  time_one       date,
  time_two       timestamp NULL DEFAULT NULL,
  time_three     timestamp NULL,
  time_four      timestamp NOT NULL,
  time_five      timestamp NULL DEFAULT '1999-01-08 04:05:06.789',
  time_six       timestamp NULL DEFAULT '1999-01-08 04:05:06.789 -8:00',
  time_seven     timestamp NULL DEFAULT 'January 8 04:05:06 1999 PST',
  time_eight     timestamp NOT NULL DEFAULT '1999-01-08 04:05:06.789',
  time_nine      timestamp NOT NULL DEFAULT '1999-01-08 04:05:06.789 -8:00',
  time_ten       timestamp NOT NULL DEFAULT 'January 8 04:05:06 1999 PST',
  time_eleven    date NULL,
  time_twelve    date NOT NULL,
  time_thirteen  date NULL DEFAULT '1999-01-08',
  time_fourteen  date NULL DEFAULT 'January 8, 1999',
  time_fifteen   date NULL DEFAULT '19990108',
  time_sixteen   date NOT NULL DEFAULT '1999-01-08',
  time_seventeen date NOT NULL DEFAULT 'January 8, 1999',
  time_eighteen  date NOT NULL DEFAULT '19990108',

  uuid_zero  uuid,
  uuid_one   uuid NULL,
  uuid_two   uuid NULL DEFAULT NULL,
  uuid_three uuid NOT NULL,
  uuid_four  uuid NULL DEFAULT '6ba7b810-9dad-11d1-80b4-00c04fd430c8',
  uuid_five  uuid NOT NULL DEFAULT '6ba7b810-9dad-11d1-80b4-00c04fd430c8',

  strange_one   integer DEFAULT '5'::integer,
  strange_two   varchar(1000) DEFAULT 5::varchar,
  strange_three timestamp without time zone default (now() at time zone 'utc'),
  strange_four  timestamp with time zone default (now() at time zone 'utc'),
  strange_five  interval NOT NULL DEFAULT '21 days',
  strange_six   interval NULL DEFAULT '23 hours',

  aa  json NULL,
  bb  json NOT NULL,
  cc  jsonb NULL,
  dd  jsonb NOT NULL,
  ee  box NULL,
  ff  box NOT NULL,
  gg  cidr NULL,
  hh  cidr NOT NULL,
  ii  circle NULL,
  jj  circle NOT NULL,
  kk  double precision NULL,
  ll  double precision NOT NULL,
  mm  inet NULL,
  nn  inet NOT NULL,
  oo  line NULL,
  pp  line NOT NULL,
  qq  lseg NULL,
  rr  lseg NOT NULL,
  ss  macaddr NULL,
  tt  macaddr NOT NULL,
  uu  money NULL,
  vv  money NOT NULL,
  ww  path NULL,
  xx  path NOT NULL,
  yy  pg_lsn NULL,
  zz  pg_lsn NOT NULL,
  aaa point NULL,
  bbb point NOT NULL,
  ccc polygon NULL,
  ddd polygon NOT NULL,
  eee tsquery NULL,
  fff tsquery NOT NULL,
  ggg tsvector NULL,
  hhh tsvector NOT NULL,
  iii txid_snapshot NULL,
  jjj txid_snapshot NOT NULL,
  kkk xml NULL,
  lll xml NOT NULL,
  mmm citext NULL,
  nnn citext NOT NULL
);

create table fun_arrays (
  id serial,
  fun_one integer[] null,
  fun_two integer[] not null,
  fun_three boolean[] null,
  fun_four boolean[] not null,
  fun_five varchar[] null,
  fun_six varchar[] not null,
  fun_seven decimal[] null,
  fun_eight decimal[] not null,
  fun_nine bytea[] null,
  fun_ten bytea[] not null,
  fun_eleven jsonb[] null,
  fun_twelve jsonb[] not null,
  fun_thirteen json[] null,
  fun_fourteen json[] not null,
  primary key (id)
);
