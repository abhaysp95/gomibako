create table gomi (
	id integer not null primary key auto_increment,
	title varchar(100) not null,
	content text not null,
	created datetime not null,
	expires datetime not null
);

create index gomi_idx_created on gomi(created);

-- dummy data
INSERT INTO gomi (title, content, created, expires) VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);

INSERT INTO gomi (title, content, created, expires) VALUES (
	"Over the wintry forest",
	"Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n - I don't know the rest of it",
	UTC_TIMESTAMP(),
	DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);

INSERT INTO gomi (title, content, created, expires) VALUES (
	'First autumn morning',
	'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n',
	UTC_TIMESTAMP(),
	DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
);

create table users (
	id integer primary key not null auto_increment,
	name varchar(100) not null,
	email varchar(255) not null,
	passwd char(60) not null,
	created datetime not null
);

alter table users add constraint users_unq_email unique(email);

-- to check for a constraint in a table you can use following query
select constraint_name from information_schema.table_constraints where table_name='users' and constraint_type='unique';
