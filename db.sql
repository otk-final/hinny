CREATE TABLE `` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `kid` bigint(20) NOT NULL COMMENT '业务唯一主键',
  `application` varchar(64) NOT NULL COMMENT '应用名称',
  `ws_name` varchar(64) NOT NULL COMMENT '空间名称',
  `ws_key` varchar(64) NOT NULL COMMENT '工作空间标识',
  `api_url` varchar(256) NOT NULL COMMENT '请求地址',
  `doc_url` varchar(256) NOT NULL COMMENT '文档接口',
  `script` text,
  `script_type` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工作空间';


CREATE TABLE `hinny_case_template` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `kid` bigint(20) NOT NULL COMMENT '业务唯一主键',
  `application` varchar(128) NOT NULL COMMENT '项目名',
  `module` varchar(64) DEFAULT NULL,
  `case_type` varchar(64) DEFAULT NULL COMMENT '案例类型',
  `case_name` varchar(255) DEFAULT NULL,
  `description` text COMMENT '备注',
  `path` varchar(256) NOT NULL COMMENT '请求路径标记',
  `request` text COMMENT '响应相关参数',
  `script_type` varchar(32) DEFAULT NULL COMMENT '验证脚本类型',
  `script` text COMMENT '验证脚本',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `hinny_workspace` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `kid` bigint(20) NOT NULL COMMENT '业务唯一主键',
  `application` varchar(64) NOT NULL COMMENT '应用名称',
  `ws_name` varchar(64) NOT NULL COMMENT '空间名称',
  `ws_key` varchar(64) NOT NULL COMMENT '工作空间标识',
  `api_url` varchar(256) NOT NULL COMMENT '请求地址',
  `doc_url` varchar(256) NOT NULL COMMENT '文档接口',
  `script` text,
  `script_type` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='工作空间';

