alter table codebase
  add constraint codebase_pk
    unique (name, type);