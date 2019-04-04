create table codebase_branch
(
  id          integer default nextval('codebase_branch_id_seq'::regclass) not null
    constraint codebase_branch_pk
      primary key,
  name        text                                                        not null,
  codebase_id integer                                                     not null
    constraint codebase_branch_app_fk
      references codebase,
  from_commit text                                                        not null,
  constraint codebase_branch_name_pk
    unique (name, codebase_id)
);