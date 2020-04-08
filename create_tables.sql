create table if not exists messages
(
	id serial not null constraint messages_pk primary key,
	scheduled_time timestamp not null,
	msg_type varchar,
	aggregation_id varchar,
	data bytea not null,
	attributes bytea
);

alter table messages owner to postgres;

create index if not exists messages_index_time
	on messages (scheduled_time);
