insert into jenkins_slave(name)
VALUES ('maven'),
       ('gradle'),
       ('npm'),
       ('dotnet')
on conflict (name) do nothing;