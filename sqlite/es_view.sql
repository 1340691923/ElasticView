DROP TABLE IF EXISTS `es_link`;
CREATE TABLE `es_link`  (
                            `id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
                            `ip` TEXT   NOT NULL,
                            `user` TEXT   NOT NULL,
                            `pwd` TEXT   NOT NULL,
                            `created` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            `updated` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            `remark` TEXT   DEFAULT '默认连接',
                            `version` tinyINTEGER NOT NULL DEFAULT 6
);
CREATE UNIQUE INDEX es_remark on es_link ( `remark`);
INSERT INTO `es_link` VALUES (10, 'http://127.0.0.1:9200', '', '', '2021-04-10 22:33:03', '2021-04-10 22:33:03', '测试', 6);
DROP TABLE IF EXISTS `gm_dsl_history`;
CREATE TABLE `gm_dsl_history`  (
                                   `id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
                                   `uid` INTEGER DEFAULT 0,
                                   `method` TEXT   DEFAULT '',
                                   `path` TEXT   DEFAULT '',
                                   `body` text,
                                       `created` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP
);
DROP TABLE IF EXISTS `gm_guid`;
CREATE TABLE `gm_guid`  (
                            `id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
                            `uid` INTEGER NOT NULL,
                            `guid_name` TEXT   NOT NULL,
                            `created` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX guid_name on gm_guid (`uid`, `guid_name`);
DROP TABLE IF EXISTS `gm_role`;
CREATE TABLE `gm_role`  (
                            `id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
                            `role_name` TEXT   DEFAULT NULL,
                            `description` TEXT   DEFAULT NULL,
                            `role_list` text
);
INSERT INTO `gm_role` VALUES (1, 'admin', '超级管理员', '[{"path":"/permission","component":"layout","redirect":"/permission/role","alwaysShow":true,"meta":{"title":"权限","icon":"el-icon-user-solid"},"children":[{"path":"role","component":"views/permission/role","name":"RolePermission","meta":{"title":"角色管理","icon":"el-icon-s-check"}},{"path":"user","component":"views/permission/user","name":"user","meta":{"title":"用户管理","icon":"el-icon-user"}}]},{"path":"/connect-tree","component":"layout","redirect":"/connect-tree/index","alwaysShow":false,"meta":{"title":"连接树管理","icon":"el-icon-link"},"children":[{"path":"/connect-tree/index","component":"views/connect-tree/index","name":"index","meta":{"title":"连接树管理","icon":"el-icon-link"}}]},{"path":"/cat","component":"layout","redirect":"/cat/index","alwaysShow":false,"meta":{"title":"ES状态","icon":"el-icon-pie-chart"},"children":[{"path":"/cat/index","component":"views/cat/index","name":"index","meta":{"title":"ES状态","icon":"el-icon-pie-chart"}}]},{"path":"/rest","component":"layout","redirect":"/rest/index","alwaysShow":false,"meta":{"title":"开发工具","icon":"el-icon-edit"},"children":[{"path":"/rest/index","component":"views/rest/index","name":"index","meta":{"title":"开发工具","icon":"el-icon-search"}}]},{"path":"/indices","component":"layout","redirect":"/indices/index","alwaysShow":true,"meta":{"title":"索引管理","icon":"el-icon-coin"},"children":[{"path":"index","component":"views/indices/index","name":"index","meta":{"title":"索引管理","icon":"el-icon-coin"}},{"path":"reindex","component":"views/indices/reindex","name":"reindex","meta":{"title":"重建索引","icon":"el-icon-document-copy"}}]},{"path":"/task","component":"layout","redirect":"/task/index","alwaysShow":false,"meta":{"title":"任务","icon":"el-icon-notebook-2"},"children":[{"path":"/task/index","component":"views/task/index","name":"index","meta":{"title":"任务","icon":"el-icon-notebook-2"}}]},{"path":"/back-up","component":"layout","redirect":"/back-up/index","alwaysShow":true,"meta":{"title":"备份","icon":"el-icon-copy-document"},"children":[{"path":"index","component":"views/back-up/index","name":"index","meta":{"title":"快照存储库","icon":"el-icon-first-aid-kit"}},{"path":"snapshot","component":"views/back-up/snapshot","name":"index","meta":{"title":"快照管理","icon":"el-icon-shopping-bag-2"}}]},{"path":"/navicat","component":"layout","redirect":"/navicat/index","alwaysShow":false,"meta":{"title":"Navicat","icon":"el-icon-copy-document"},"children":[{"path":"/navicat/index","component":"views/navicat/index","name":"index","meta":{"title":"Navicat","icon":"el-icon-first-aid-kit"}}]}]');
DROP TABLE IF EXISTS `gm_user`;
CREATE TABLE `gm_user`  (
                            `id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
                            `username` TEXT   DEFAULT NULL,
                            `password` TEXT   DEFAULT NULL,
                            `role_id` INTEGER DEFAULT NULL ,
                            `realname` TEXT   DEFAULT ''
);
CREATE UNIQUE INDEX gm_user_username on gm_user ( `username`);
INSERT INTO `gm_user` VALUES (1, 'admin', '21232f297a57a5a743894a0e4a801fc3', 1, '肖文龙');
DROP TABLE IF EXISTS `gm_timed_list`;
CREATE TABLE `gm_timed_list`  (
                            `id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
                            `action` TEXT   DEFAULT NULL,
                            `exec_time` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            `exec_time_format` TEXT NOT NULL DEFAULT '' ,
                            `status` INTEGER NOT NULL DEFAULT '0',
                            `task_id` TEXT NOT NULL DEFAULT '' ,
                            `msg` TEXT NOT NULL DEFAULT '' ,
                            `extra` TEXT NOT NULL DEFAULT '' ,
                            `updated` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            `created` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            `data` TEXT  NOT NULL   DEFAULT ''
);
DROP TABLE IF EXISTS `datax_link_info`;
CREATE TABLE `datax_link_info`  (
                                  `id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
                                  `ip` TEXT   DEFAULT NULL,
                                  `port` INTEGER NOT NULL DEFAULT '0',
                                  `db_name` TEXT NOT NULL DEFAULT '' ,
                                  `username` TEXT NOT NULL DEFAULT '' ,
                                  `pwd` TEXT NOT NULL DEFAULT '' ,
                                  `remark` TEXT NOT NULL DEFAULT '' ,
                                  `typ`  TEXT NOT NULL DEFAULT '' ,
                                  `updated` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                  `created` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX link_remark_uniq on datax_link_info ( `remark`,`typ`);
DROP TABLE IF EXISTS `datax_transfer_list`;
CREATE TABLE `datax_transfer_list`  (
                                  `id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
                                  `form_data` TEXT   DEFAULT NULL,
                                  `remark` TEXT NOT NULL DEFAULT '' ,
                                  `table_name` INTEGER NOT NULL DEFAULT '0',
                                  `index_name` TEXT NOT NULL DEFAULT '' ,
                                  `error_msg` TEXT NOT NULL DEFAULT '无报错' ,
                                    `crontab_spec` TEXT NOT NULL DEFAULT '' ,
                                 `dbcount` INTEGER NOT NULL DEFAULT '0',
                                  `escount` INTEGER NOT NULL DEFAULT '0',
                                  `status` TEXT NOT NULL DEFAULT '任务运行中...' ,
                                  `updated` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                  `created` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX datax_transfer_list_remark on datax_transfer_list ( `remark`);