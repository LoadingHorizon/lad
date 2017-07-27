
CREATE TABLE IF NOT EXISTS `advertiser` (
    `id` int(10) NOT NULL AUTO_INCREMENT,
    `user` varchar(20) NOT NULL,
    `mail` varchar(100) NOT NULL,
    `phone` varchar(20),
    `country` varchar(20),
    PRIMARY KEY (`id`),
    UNIQUE `user_index` (`user`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

CREATE TABLE IF NOT EXISTS `advapp` (
    `id` int(10) NOT NULL AUTO_INCREMENT,
    `advertiser_id` int(10),
    `platform` tinyint(2) NOT NULL,
    `package` varchar(128) NOT NULL,
    `third` varchar(1024),
    `icon` varchar(1024),
    PRIMARY KEY (`id`),
    FOREIGN KEY (`advertiser_id`) REFERENCES `advertiser` (`id`),
    UNIQUE `advertiser_package_index` (`advertiser_id`, `package`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

CREATE TABLE IF NOT EXISTS `creative` (
    `id` int(10) NOT NULL AUTO_INCREMENT,
    `advertiser_id` int(10),
    `title` varchar(1024),
    FOREIGN KEY (`advertiser_id`) REFERENCES `advertiser` (`id`),
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

CREATE TABLE IF NOT EXISTS `campaign` (
    `id` int(20) NOT NULL AUTO_INCREMENT,
    `advapp_id` int(10),
    `creative_id` int(10),

    `price` float NOT NULL DEFAULT 0,
    `price_out` float NOT NULL DEFAULT 0,

    `daily_cap` bigint(20) NOT NULL DEFAULT -1,
    `budget` bigint(20) NOT NULL DEFAULT 0,
    `regions` varchar(2048) NOT NULL DEFAULT "ALL",
    `stime` bigint(20) NOT NULL DEFAULT 0,
    `etime` bigint(20) NOT NULL DEFAULT -1,
    `status` tinyint(2) NOT NULL DEFAULT 1,
    `osmin` int(10) NOT NULL DEFAULT 0,
    `osmax` int(10) NOT NULL DEFAULT -1,
    `black` varchar(4096),
    `white` varchar(4096),

    PRIMARY KEY (`id`),
    FOREIGN KEY (`advapp_id`) REFERENCES `advapp` (`id`),
    FOREIGN KEY (`creative_id`) REFERENCES `creative` (`id`),
    UNIQUE `advapp_creative_index` (`advapp_id`, `creative_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

CREATE TABLE IF NOT EXISTS `publisher` (
    `id` int(10) NOT NULL AUTO_INCREMENT,
    `user` varchar(20) NOT NULL,
    `mail` varchar(100) NOT NULL,
    `phone` varchar(20),
    `country` varchar(20),
    PRIMARY KEY (id),
    UNIQUE `user_index` (`user`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

CREATE TABLE IF NOT EXISTS `pubapp` (
    `id` int(10) NOT NULL AUTO_INCREMENT,
    `publisher_id` int(10),
    `platform` tinyint(2) NOT NULL,
    `package` varchar(128) NOT NULL,
    `icon` varchar(1024),
    PRIMARY KEY (`id`),
    FOREIGN KEY (`publisher_id`) REFERENCES `publisher` (`id`),
    UNIQUE `publisher_package_index` (`publisher_id`, `package`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

CREATE TABLE IF NOT EXISTS `position` (
    `id` int(10) NOT NULL AUTO_INCREMENT,
    `pubapp_id` int(10),
    `type` varchar(20),
    PRIMARY KEY (`id`),
    FOREIGN KEY (`pubapp_id`) REFERENCES `pubapp` (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

-- DROP TABLE campaign;
-- DROP TABLE creative;
-- DROP TABLE position;
-- DROP TABLE advapp;
-- DROP TABLE pubapp;
-- DROP TABLE advertiser;
-- DROP TABLE publisher;


