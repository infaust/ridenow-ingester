CREATE TABLE forecast (
	id varchar(100) PRIMARY KEY,
	location serial references location(id),
	wave_height_m real NOT NULL,
	swell_period_sec real NOT NULL,
	created timestamp NOT NULL default current_timestamp,
	modified timestamp NOT NULL,
	time timestamp NOT NULL
);