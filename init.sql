USE 'go-note'

CREATE TABLE `account` (
	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`username` varchar(32) NOT NULL UNIQUE,
	`password` varchar(64) NOT NULL,
	`create_time` bigint(20) NULL DEFAULT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `note` (
	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`user_id` int(11) unsigned NOT NULL,
	`topic` varchar(256) NOT NULL,
	`content` TEXT NULL DEFAULT NULL,
	PRIMARY KEY (`id`),
	UNIQUE KEY `user_topic` (`user_id`, `topic`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;