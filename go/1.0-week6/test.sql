create database if not exists testdb;

use testdb;

create TABLE
	if not exists test_table (
		id int auto_increment primary key,
		name char(50) not NULL,
		age int not NULL,
		birthdate DATE,
		is_active boolean default true
		-- 这是注释喵，mysql中boolean和bool一样的，都是tinyint(1)的别名
	);