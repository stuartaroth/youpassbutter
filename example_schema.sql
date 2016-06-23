drop database if exists youpassbutter;
create database youpassbutter;
use youpassbutter;

create table author (
	id int auto_increment primary key,
	first_name varchar (255) not null,
	last_name varchar (255) not null,
	middle_name varchar (255),
	unique (first_name, middle_name, last_name)
);

create table genre (
	id int auto_increment primary key,
	name varchar (255) not null,
	unique (name)
);

create table book (
	id int auto_increment primary key,
	id_author int not null,
	id_genre int,
	title varchar (255) not null,
	year_published int not null,
	foreign key (id_author) references author (id),
	foreign key (id_genre) references genre (id),
	unique (id_author, title, year_published)
);
