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
	`sid` varchar(36) NOT NULL UNIQUE,
	`topic` varchar(256) NOT NULL,
	`content` TEXT NULL DEFAULT NULL,
	`create_time` bigint(20) NULL DEFAULT NULL,
	`update_time` bigint(20) NULL DEFAULT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `account_note` (
	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`account_id` int(11) unsigned NOT NULL,
	`note_id` int(11) unsigned NOT NULL,
	`is_owner` tinyint(1) NOT NULL DEFAULT 0,
	`create_time` bigint(20) NULL DEFAULT NULL,
	`update_time` bigint(20) NULL DEFAULT NULL,
	PRIMARY KEY (`id`),
	UNIQUE KEY `account_note_id` (`account_id`, `note_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;