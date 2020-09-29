insert into jenkins_slave(name)
VALUES ('maven'),
       ('gradle'),
       ('dotnet')
on conflict (name) do nothing;