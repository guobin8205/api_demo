/*
 Navicat Premium Data Transfer

 Source Server         : 121.199.6.172
 Source Server Type    : MySQL
 Source Server Version : 50720
 Source Host           : 121.199.6.172:3306
 Source Schema         : live

 Target Server Type    : MySQL
 Target Server Version : 50720
 File Encoding         : 65001

 Date: 23/07/2018 14:26:01
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `uid` int(11) NOT NULL AUTO_INCREMENT,
  `account` varchar(100) NOT NULL,
  `openid` int(10) unsigned NOT NULL DEFAULT '0',
  `is_reg` tinyint(3) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  KEY `account` (`account`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1336789 DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
