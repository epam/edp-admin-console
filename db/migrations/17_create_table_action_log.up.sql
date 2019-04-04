create table action_log
(
  id               integer default nextval('action_log_id_seq'::regclass) not null
    constraint action_log_pk
      primary key,
  event            event                                                  not null,
  detailed_message text,
  username         text,
  updated_at       timestamp                                              not null
);