create table jira_server
(
	id serial not null
		constraint jira_server_pk
			primary key,
	name text,
	available bool
);