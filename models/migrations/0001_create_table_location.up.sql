CREATE TABLE location (
	id serial PRIMARY KEY,
	name varchar(100) NOT NULL,
	description varchar(300) NOT NULL,
	relevant boolean default FALSE,
	coords point NOT NULL,
	url varchar(100) NOT NULL,
	country varchar(5) NOT NULL,
	created timestamp NOT NULL default current_timestamp
);