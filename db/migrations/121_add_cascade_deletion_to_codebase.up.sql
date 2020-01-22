alter table codebase_branch
    drop constraint codebase_branch_app_fk,
    add constraint codebase_branch_app_fk foreign key (codebase_id) references codebase (id) on delete cascade;
alter table codebase_action_log
    drop constraint codebase_fk,
    add constraint codebase_fk foreign key (codebase_id) references codebase (id) on delete cascade;
alter table applications_to_promote
    drop constraint codebase_fk,
    add constraint codebase_fk foreign key (codebase_id) references codebase (id) on delete cascade;
alter table codebase_docker_stream
    drop constraint codebase_branch_fk,
    add constraint codebase_branch_fk foreign key (codebase_branch_id) references codebase_branch (id) on delete cascade;
alter table cd_stage
    drop constraint codebase_branch_fk,
    add constraint codebase_branch_fk foreign key (codebase_branch_id) references codebase_branch (id) on delete cascade;