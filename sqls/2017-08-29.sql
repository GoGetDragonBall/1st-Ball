-- 2017-08-29

drop table if exists `users`;
create table users (
       id int not null auto_increment,
       name varchar(50) character set utf8 not null,
       nickname varchar(50) character set utf8 not null,
       email varchar(100) not null,
       created_at datetime not null default current_timestamp,
       updated_at datetime not null default current_timestamp,
       deleted_at datetime,
       PRIMARY KEY(id),
       KEY `email` (`email`)
)  ENGINE=InnoDB DEFAULT CHARSET=utf8;

drop table if exists `user_authorization_key_types`;
create table user_authorization_key_types (
       id int not null auto_increment,
       name varchar(100) not null,
       created_at datetime not null default current_timestamp,
       updated_at datetime not null default current_timestamp,
       deleted_at datetime,
       PRIMARY KEY(id),
       KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

drop table if exists `user_authorization_keys`;
create table user_authorization_keys (
       id int not null auto_increment,
       user_id int not null,
       type_id int not null,
       authorization_key varchar(300) not null,
       created_at datetime not null default current_timestamp,
       updated_at datetime not null default current_timestamp,
       deleted_at datetime,
       PRIMARY KEY(id),
       CONSTRAINT `user_authorization_keys_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE NO ACTION ON UPDATE CASCADE,
       CONSTRAINT `user_authorization_keys_ibfk_2` FOREIGN KEY (`type_id`) REFERENCES `user_authorization_key_types` (`id`) ON DELETE NO ACTION ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
