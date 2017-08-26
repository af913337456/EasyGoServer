CREATE database IF NOT EXISTS mattermost;
USE mattermost;

# 用户表
DROP TABLE IF EXISTS `LUser`;
CREATE TABLE `LUser`(
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `u_user_id` varchar(1024) COLLATE utf8mb4_bin DEFAULT '' COMMENT  '用户 id',
  `u_time` varchar(30) COLLATE utf8mb4_bin DEFAULT '' COMMENT '首次打开APP时间',
  `u_open_time` int(30) unsigned NOT NULL DEFAULT '0' COMMENT '一共打开的次数',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# comment 评论表
DROP TABLE IF EXISTS `LComment`;
CREATE TABLE `LComment`(
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `u_user_id` varchar(1024) COLLATE utf8mb4_bin DEFAULT '' COMMENT  '用户 id',
  `u_time` varchar(30) COLLATE utf8mb4_bin DEFAULT '' COMMENT '评论时间',
  `u_content` varchar(1024) COLLATE utf8mb4_bin DEFAULT '' COMMENT '内容',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;















