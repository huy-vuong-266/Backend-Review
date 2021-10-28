Create table `tokens` (
	user_id varchar(36) not null,
	token varchar(255) not null,
	created_at int unsigned,
	primary key (user_id)
);


create Event deactive_token
on Schedule every 5 MINUTE 
Starts '2021-10-21 23:00:00'
Do delete from token where UNIX_TIMESTAMP(NOW()) - token.created_at > 14400
