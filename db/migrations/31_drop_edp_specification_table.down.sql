create table if not exists edp_specification
(
  name    text not null
    constraint edp_specification_pk
      primary key,
  version text not null
);