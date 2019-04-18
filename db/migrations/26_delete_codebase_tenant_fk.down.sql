alter table codebase
  add constraint codebase_tenant_fk
    foreign key (tenant_name) references edp_specification;
