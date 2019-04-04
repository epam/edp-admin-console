create table if not exists codebase
(
  id                integer default nextval('codebase_id_seq'::regclass) not null
    constraint codebase_id_pk
      primary key,
  type              codebase_type                                        not null,
  name              text                                                 not null,
  tenant_name       text                                                 not null
    constraint codebase_tenant_fk
      references edp_specification,

  language          language                                             not null,
  framework         framework                                            not null,
  build_tool        build_tool                                           not null,
  strategy          strategy                                             not null,

  repository_url    text,

  route_site        text,
  route_path        text,

  database_kind     text,
  database_version  text,
  database_capacity text,
  database_storage  text,

  status            status                                               not null,
  constraint codebase_pk
    unique (name, tenant_name, type)
);