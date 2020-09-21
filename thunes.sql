create database thunes;
use thunes;

create table `thunes_user` (
    `id` INT(16) AUTO_INCREMENT COMMENT '用户ID',
    `name` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '用户姓名',
    `email` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '用户邮箱',
    `pwd` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '密码',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modified_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_name` (`name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户信息';

####
create table `thunes_account` (
    `id` INT(16) COMMENT '用户ID',
    `account_id` VARCHAR(16) NOT NULL DEFAULT '' COMMENT '账户ID',
    `balance` DECIMAL(64,12) NOT NULL DEFAULT 0 COMMENT '账户金额',
    `unit` int(4) NOT NULL DEFAULT 252 COMMENT '货币类型',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modified_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    UNIQUE KEY `idx_id` (`id`),
    UNIQUE KEY `idx_aid` (`account_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户账户信息';

#######
create table `thunes_transfer_order` (
	`id` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '订单ID',
	`ts` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '订单生成时间戳',
	`status` INT not NULL DEFAULT 0 COMMENT '订单状态',
	`status_info` VARCHAR(32) DEFAULT '' COMMENT '状态描述',
	
	`from` VARCHAR(16) NOT NULL DEFAULT '' COMMENT '转出账户account_id',
	`from_num` DECIMAL(64,12) NOT NULL DEFAULT 0 COMMENT '转出账户金额',
	`from_unit` INT NOT NULL DEFAULT 0 COMMENT '转出账户余额货币类别',
	
	`to` VARCHAR(16) NOT NULL DEFAULT '' COMMENT '转入账户account_id',
	`to_num` DECIMAL(64,12) NOT NULL DEFAULT 0 COMMENT '转入账户金额',
	`to_unit` INT NOT NULL DEFAULT 0 COMMENT '转入账户余额货币类别',
	
	`exchange_rate` DECIMAL(64, 12) NOT NULL DEFAULT 0 COMMENT '转出汇率，即1 from货币单位 = x to货币单位',
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modified_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
	
	UNIQUE `idx_id` (`id`)
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '转账交易信息';

CREATE TABLE `thunes_account_balance_log` (
	`lid` VARCHAR(64) NOT NULL DEFAULT '' COMMENT 'log的id',
	`account_id` VARCHAR(16) NOT NULL DEFAULT '' COMMENT '账户ID',
	`operation` int NOT NULL DEFAULT 0 COMMENT '操作类型, 0:default 1:set 2:update 3:incr 4:decr',
	`can_rollback` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '该log是否可以执行回滚操作',
	`raw` DECIMAL(64,12) not NULL DEFAULT 0 COMMENT 'balance操作前值',
	`delta` DECIMAL(64,12) not NULL DEFAULT 0 COMMENT 'balance操作的变量值',
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modified_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
	UNIQUE `idx_lid_account_id` (`lid`, `account_id`)
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '账户余额变更记录';