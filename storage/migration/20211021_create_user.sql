Create table `users` (
	user_id varchar(36) not null,
	created_at int unsigned not null,
	updated_at int unsigned not null,
	fullname varchar(70) not null,
	phone varchar(15) not null,
	email varchar(255) not null,
	encrypted_pw varchar(255) not null,
	salt varchar(255) not null,
	budget BigInt(20) default 0,
	status TinyInt,
	Primary Key (user_id)
);