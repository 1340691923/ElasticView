/*
 Navicat Premium Data Transfer

 Source Server         : local
 Source Server Type    : MySQL
 Source Server Version : 50717
 Source Host           : localhost:3306
 Source Schema         : es_view

 Target Server Type    : MySQL
 Target Server Version : 50717
 File Encoding         : 65001

 Date: 22/04/2021 01:31:30
*/

SET NAMES utf8mb4;
SET
FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for es_link
-- ----------------------------
DROP TABLE IF EXISTS `es_link`;
CREATE TABLE `es_link`
(
    `id`      int(11) NOT NULL AUTO_INCREMENT,
    `ip`      varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `user`    varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `pwd`     varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `created` timestamp(0)                                                  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated` timestamp(0)                                                  NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP (0),
    `remark`  varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci          DEFAULT '默认连接',
    `version` tinyint(10) NOT NULL DEFAULT 6,
    `rootpem` text,
    `certpem` text,
    `keypem`  text,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE INDEX `es_remark`(`remark`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of es_link
-- ----------------------------
INSERT INTO `es_link`
VALUES (10, 'http://127.0.0.1:9200', '', '', '2021-04-10 22:33:03', '2021-04-10 22:33:03', '测试', 6, '', '', '');

-- ----------------------------
-- Table structure for gm_dsl_history
-- ----------------------------
DROP TABLE IF EXISTS `gm_dsl_history`;
CREATE TABLE `gm_dsl_history`
(
    `id`      int(11) NOT NULL AUTO_INCREMENT,
    `uid`     int(11) DEFAULT 0,
    `method`  varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  DEFAULT '',
    `path`    varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '',
    `body`    text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
    `created` timestamp(0)                                                  DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP (0),
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 42 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for gm_guid
-- ----------------------------
DROP TABLE IF EXISTS `gm_guid`;
CREATE TABLE `gm_guid`
(
    `id`        int(11) NOT NULL AUTO_INCREMENT,
    `uid`       int(11) NOT NULL,
    `guid_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `created`   timestamp(0)                                                  NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP (0),
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE INDEX `guid_name`(`uid`, `guid_name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 15 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for gm_role
-- ----------------------------
DROP TABLE IF EXISTS `gm_role`;
CREATE TABLE `gm_role`
(
    `id`          int(11) NOT NULL AUTO_INCREMENT,
    `role_name`   varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `role_list`   text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of gm_role
-- ----------------------------
INSERT INTO gm_role
(id, role_name, description, role_list)
VALUES(1, 'admin', '超级管理员', '[{"path":"/permission","component":"layout","redirect":"/permission/role","alwaysShow":true,"meta":{"title":"权限","icon":"el-icon-user-solid"},"children":[{"path":"role","component":"views/permission/role","name":"role","meta":{"title":"角色管理","icon":"el-icon-s-check"}},{"path":"user","component":"views/permission/user","name":"user","meta":{"title":"用户管理","icon":"el-icon-user"}},{"path":"operater_log","component":"views/permission/operater_log","name":"operater_log","meta":{"title":"操作日志列表","icon":"el-icon-s-order"}}]},{"path":"/connect-tree","component":"layout","redirect":"/connect-tree/index","alwaysShow":false,"meta":{"title":"连接树管理","icon":"el-icon-link"},"children":[{"path":"index","component":"views/connect-tree/index","name":"index","meta":{"title":"连接树管理","icon":"el-icon-link"}}]},{"path":"/cat","component":"layout","redirect":"/cat/index","alwaysShow":false,"meta":{"title":"ES状态","icon":"el-icon-pie-chart"},"children":[{"path":"index","component":"views/cat/index","name":"index","meta":{"title":"ES状态","icon":"el-icon-pie-chart"}}]},{"path":"/rest","component":"layout","redirect":"/rest/index","alwaysShow":false,"meta":{"title":"开发工具","icon":"el-icon-edit"},"children":[{"path":"index","component":"views/rest/index","name":"index","meta":{"title":"开发工具","icon":"el-icon-search"}}]},{"path":"/indices","component":"layout","redirect":"/indices/index","alwaysShow":true,"meta":{"title":"索引管理","icon":"el-icon-coin"},"children":[{"path":"index","component":"views/indices/index","name":"index","meta":{"title":"索引管理","icon":"el-icon-coin"}},{"path":"reindex","component":"views/indices/reindex","name":"reindex","meta":{"title":"重建索引","icon":"el-icon-document-copy"}}]},{"path":"/task","component":"layout","redirect":"/task/index","alwaysShow":false,"meta":{"title":"任务","icon":"el-icon-notebook-2"},"children":[{"path":"index","component":"views/task/index","name":"index","meta":{"title":"任务","icon":"el-icon-notebook-2"}}]},{"path":"/back-up","component":"layout","redirect":"/back-up/index","alwaysShow":true,"meta":{"title":"备份","icon":"el-icon-copy-document"},"children":[{"path":"index","component":"views/back-up/index","name":"index","meta":{"title":"快照存储库","icon":"el-icon-first-aid-kit"}},{"path":"snapshot","component":"views/back-up/snapshot","name":"index","meta":{"title":"快照管理","icon":"el-icon-shopping-bag-2"}}]},{"path":"/navicat","component":"layout","redirect":"/navicat/index","alwaysShow":false,"meta":{"title":"Navicat","icon":"el-icon-copy-document"},"children":[{"path":"index","component":"views/navicat/index","name":"index","meta":{"title":"Navicat","icon":"el-icon-first-aid-kit"}}]}]');

-- ----------------------------
-- Table structure for gm_user
-- ----------------------------
DROP TABLE IF EXISTS `gm_user`;
CREATE TABLE `gm_user`
(
    `id`       int(11) NOT NULL AUTO_INCREMENT,
    `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
    `role_id`  int(11) DEFAULT NULL COMMENT '角色id',
    `realname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '真实姓名',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE INDEX `gm_user_username`(`username`) USING BTREE COMMENT '角色名唯一索引'
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of gm_user
-- ----------------------------
INSERT INTO `gm_user`
VALUES (1, 'admin', '21232f297a57a5a743894a0e4a801fc3', 1, '肖文龙');
DROP TABLE IF EXISTS `datax_link_info`;
CREATE TABLE `datax_link_info`
(
    `id`       int(11) NOT NULL AUTO_INCREMENT,
    `ip`       varchar(255) DEFAULT '',
    `port`     int(11) DEFAULT '0',
    `db_name`  varchar(255) DEFAULT '',
    `username` varchar(255) DEFAULT '',
    `pwd`      varchar(255) DEFAULT '',
    `remark`   varchar(50)  DEFAULT '',
    `typ`      varchar(50)  DEFAULT '',
    `updated`  timestamp NULL DEFAULT NULL,
    `created`  timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `link_remark_uniq` (`remark`,`typ`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
DROP TABLE IF EXISTS `datax_transfer_list`;
CREATE TABLE datax_transfer_list
(
    `id`           int(11) NOT NULL AUTO_INCREMENT,
    `form_data`    text         NOT NULL,
    `remark`       varchar(50)  NOT NULL DEFAULT '',
    `table_name`   varchar(255) NOT NULL DEFAULT '',
    `index_name`   varchar(255) NOT NULL DEFAULT '',
    `error_msg`    varchar(255) NOT NULL DEFAULT '',
    `crontab_spec` varchar(255) NOT NULL DEFAULT '',
    `dbcount`      int(11) NOT NULL DEFAULT 0,
    `es_connect`   int(11) NOT NULL DEFAULT 0,
    `escount`      int(11) NOT NULL DEFAULT 0,
    `status`       varchar(255) NOT NULL DEFAULT '',
    `updated`      timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP (0),
    `created`      timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),
    PRIMARY KEY (`id`),
    UNIQUE INDEX `datax_transfer_list_remark`(`remark`, `es_connect`) USING BTREE

);
CREATE TABLE `gm_operater_log` (
                                   `id` int(11) NOT NULL AUTO_INCREMENT,
                                   `operater_name` varchar(50) DEFAULT '' COMMENT '操作者名字',
                                   `operater_id` int(11) DEFAULT '0' COMMENT '操作者id',
                                   `operater_action` varchar(50) DEFAULT '' COMMENT '请求路由',
                                   `created` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                   `method` varchar(500) DEFAULT NULL COMMENT '请求方法',
                                   `body` blob NOT NULL COMMENT '请求body',
                                   `operater_role_id` int(11) NOT NULL,
                                   PRIMARY KEY (`id`) USING BTREE,
                                   KEY `operater_action` (`operater_action`) USING BTREE,
                                   KEY `operater_id` (`operater_id`) USING BTREE,
                                   KEY `operater_role_id` (`operater_role_id`) USING BTREE,
                                   KEY `operater_id_act_role` (`operater_action`,`operater_id`,`operater_role_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;


SET
FOREIGN_KEY_CHECKS = 1;
