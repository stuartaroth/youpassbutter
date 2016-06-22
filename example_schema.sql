drop database if exists youpassbutter;
create database youpassbutter;
use youpassbutter;

create table network (
	id int auto_increment primary key,
	name varchar (255) not null,
	unique (name)
);

create table series (
	id int auto_increment primary key,
	id_network int not null,
	name varchar (255) not null,
	unique (name),
	foreign key (id_network) references network (id)
);

insert into network (name) values ('Fox');
insert into network (name) values ('Cartoon Network');
insert into network (name) values ('Netflix');

insert into series (id_network, name) values (1, 'Futurama');
insert into series (id_network, name) values (1, 'Bob\'s Burgers');
insert into series (id_network, name) values (1, 'The Simpsons');

insert into series (id_network, name) values (2, 'Adventure Time');
insert into series (id_network, name) values (2, 'Rick and Morty');
insert into series (id_network, name) values (2, 'The Venture Bros.');

insert into series (id_network, name) values (3, 'BoJack Horseman');
