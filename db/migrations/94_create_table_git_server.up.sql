create table git_server
(
	id serial not null
		constraint git_server_pk
			primary key,
	name text,
	available bool
);