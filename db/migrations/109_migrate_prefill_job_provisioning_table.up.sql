insert into job_provisioning(name)
VALUES ('default')
on conflict (name) do nothing;