CREATE DATABASE `uai_common` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE `counter` (
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
    `msg_id`      bigint(20) unsigned NOT NULL COMMENT '业务计数场景的消息(业务)id',
    `counter_key` varchar(512) NOT NULL COMMENT '业务计数场景的唯一key',
    `counter_num` int(11) NOT NULL DEFAULT '0' COMMENT '计数器数量',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at`  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_counter` (`counter_key`,`msg_id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COMMENT='计数服务层 计数器表';


CREATE TABLE `dconfig` (
    `id`             bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
    `namespace` varchar(64) NOT NULL COMMENT '业务命名空间',
    `key` varchar(64) NOT NULL COMMENT '配置的key，格式模板为xxx.xxx.xxx...',
    `content`        varchar(128) NOT NULL COMMENT '配置内容',
    `version`        int(11) NOT NULL DEFAULT '0' COMMENT '配置版本',
    `deleted_at` timestamp NULL DEFAULT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_at`  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_dconfig`(`namespace`, `key`, `version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='动态配置层 配置表';