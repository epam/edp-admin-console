alter table cd_stage
    add codebase_branch_id integer
        constraint codebase_branch_fk
            references codebase_branch;