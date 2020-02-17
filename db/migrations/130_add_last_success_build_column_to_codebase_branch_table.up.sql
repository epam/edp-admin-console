alter table codebase_branch
    add column if not exists last_success_build text;