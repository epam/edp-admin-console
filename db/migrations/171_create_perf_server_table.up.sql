create table if not exists perf_server
(
	id serial not null
		constraint perf_server_pk
			primary key,
	name text,
	available bool
);