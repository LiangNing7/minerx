-- gateway_chain

CREATE TABLE `gateway_chain` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `namespace` varchar(253) NOT NULL DEFAULT '' COMMENT '命名空间',
  `name` varchar(253) NOT NULL DEFAULT '' COMMENT '区块链名',
  `displayName` varchar(253) NOT NULL DEFAULT '' COMMENT '区块链展示名',
  `minerType` varchar(16) NOT NULL DEFAULT '' COMMENT '区块链矿机机型',
  `image` varchar(253) NOT NULL DEFAULT '' COMMENT '区块链镜像 ID',
  `minMineIntervalSeconds` int(8) NOT NULL DEFAULT 0 COMMENT '矿机挖矿间隔',
  `createdAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq.gateway_chain.namespace_name` (`namespace`,`name`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COMMENT='区块链表';


-- gateway_minerset

CREATE TABLE `gateway_minerset` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `namespace` varchar(253) NOT NULL DEFAULT '' COMMENT '命名空间',
  `name` varchar(253) NOT NULL DEFAULT '' COMMENT '矿机池名',
  `replicas` int(8) NOT NULL DEFAULT 0 COMMENT '矿机副本数',
  `displayName` varchar(253) NOT NULL DEFAULT '' COMMENT '矿机池展示名',
  `deletePolicy` varchar(32) NOT NULL DEFAULT '' COMMENT '矿机池缩容策略',
  `minReadySeconds` int(8) NOT NULL DEFAULT 0 COMMENT '矿机 Ready 最小等待时间',
  `fullyLabeledReplicas` int(8) NOT NULL DEFAULT 0 COMMENT '所有标签匹配的副本数',
  `readyReplicas` int(8) NOT NULL DEFAULT 0 COMMENT 'Ready 副本数',
  `availableReplicas` int(8) NOT NULL DEFAULT 0 COMMENT '可用副本数',
  `failureReason` longtext DEFAULT NULL COMMENT '失败原因',
  `failureMessage` longtext DEFAULT NULL COMMENT '失败信息',
  `conditions` longtext DEFAULT NULL COMMENT '矿机池状态',
  `createdAt` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatedAt` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq.gateway_minerset.namespace_name` (`namespace`,`name`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COMMENT='矿机池表';

-- gateway_miner

CREATE TABLE `gateway_miner` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `namespace` varchar(253) NOT NULL DEFAULT '' COMMENT '命名空间',
  `name` varchar(253) NOT NULL DEFAULT '' COMMENT '矿机名',
  `displayName` varchar(253) NOT NULL DEFAULT '' COMMENT '矿机展示名',
  `phase` varchar(45) NOT NULL DEFAULT '' COMMENT '矿机状态',
  `minerType` varchar(16) NOT NULL DEFAULT '' COMMENT '矿机机型',
  `chainName` varchar(253) NOT NULL DEFAULT '' COMMENT '矿机所属的区块链名',
  `cpu` int(8) NOT NULL DEFAULT 0 COMMENT '矿机 CPU 规格',
  `memory` int(8) NOT NULL DEFAULT 0 COMMENT '矿机内存规格',
  `createdAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq.gateway_miner.namespace_name` (`namespace`,`name`),
  KEY `idx.gateway_miner.chainName` (`chainName`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COMMENT='矿机表';
