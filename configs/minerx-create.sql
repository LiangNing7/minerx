-- usercenter_user

CREATE TABLE `usercenter_user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `userID` varchar(253) NOT NULL DEFAULT '' COMMENT '用户 ID',
  `username` varchar(253) NOT NULL DEFAULT '' COMMENT '用户名称',
  `status` varchar(64) NOT NULL DEFAULT '' COMMENT '用户状态：registered,active,disabled,blacklisted,locked,deleted',
  `nickname` varchar(253) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `password` varchar(64) NOT NULL DEFAULT '' COMMENT '用户加密后的密码',
  `email` varchar(253) NOT NULL DEFAULT '' COMMENT '用户电子邮箱',
  `phone` varchar(16) NOT NULL DEFAULT '' COMMENT '用户手机号',
  `createdAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq.usercenter_user.username` (`username`),
  UNIQUE KEY `uniq.usercenter_user.userID` (`userID`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- usercenter_secret

CREATE TABLE `usercenter_secret` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `userID` varchar(253) NOT NULL DEFAULT '' COMMENT '用户 ID',
  `name` varchar(253) NOT NULL DEFAULT '' COMMENT '密钥名称',
  `secretID` varchar(36) NOT NULL DEFAULT '' COMMENT '密钥 ID',
  `secretKey` varchar(36) NOT NULL DEFAULT '' COMMENT '密钥 Key',
  `status` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '密钥状态，0-禁用；1-启用',
  `expires` bigint(64) NOT NULL DEFAULT 0 COMMENT '0 永不过期',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '密钥描述',
  `createdAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq.usercenter_secret.secretID` (`secretID`),
  KEY `idx.usercenter_secret.userID` (`userID`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COMMENT='密钥表';
