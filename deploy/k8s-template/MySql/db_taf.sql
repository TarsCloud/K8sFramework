-- MySQL dump 10.17  Distrib 10.3.22-MariaDB, for debian-linux-gnu (x86_64)
--
-- ------------------------------------------------------
-- Server version	5.7.24-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Current Database: `_DB_TAF_DATABASE_`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `_DB_TAF_DATABASE_` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

USE `_DB_TAF_DATABASE_`;

--
-- Table structure for table `t_affinity`
--

DROP TABLE IF EXISTS `t_affinity`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_affinity` (
  `f_affinity_id` int(11) NOT NULL AUTO_INCREMENT,
  `f_node_name` varchar(64) DEFAULT NULL,
  `f_app_server` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`f_affinity_id`),
  UNIQUE KEY `f_node_name_2` (`f_node_name`,`f_app_server`),
  KEY `f_node_name` (`f_node_name`),
  KEY `f_app_server` (`f_app_server`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_affinity`
--

LOCK TABLES `t_affinity` WRITE;
/*!40000 ALTER TABLE `t_affinity` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_affinity` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_app`
--

DROP TABLE IF EXISTS `t_app`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_app` (
  `f_app_name` varchar(64) NOT NULL DEFAULT '' COMMENT '应用名',
  `f_app_mark` text,
  `f_create_person` varchar(64) NOT NULL DEFAULT '' COMMENT '创建者',
  `f_create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `f_business_name` varchar(64) DEFAULT NULL COMMENT '业务名',
  PRIMARY KEY (`f_app_name`),
  KEY `f_business_name` (`f_business_name`),
  CONSTRAINT `t_app_ibfk_1` FOREIGN KEY (`f_business_name`) REFERENCES `t_business` (`f_business_name`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_app`
--

LOCK TABLES `t_app` WRITE;
/*!40000 ALTER TABLE `t_app` DISABLE KEYS */;
INSERT INTO `t_app` VALUES ('taf','','admin',current_timestamp,NULL);
/*!40000 ALTER TABLE `t_app` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_auth`
--

DROP TABLE IF EXISTS `t_auth`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_auth` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `flag` varchar(256) DEFAULT NULL,
  `role` varchar(256) DEFAULT NULL,
  `uid` varchar(256) DEFAULT NULL,
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_auth`
--

LOCK TABLES `t_auth` WRITE;
/*!40000 ALTER TABLE `t_auth` DISABLE KEYS */;
INSERT INTO `t_auth` VALUES (8,NULL,'admin','admin',current_timestamp,NULL),(9,'taflog','admin','krusli',current_timestamp,NULL);
/*!40000 ALTER TABLE `t_auth` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_business`
--

DROP TABLE IF EXISTS `t_business`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_business` (
  `f_business_name` varchar(64) NOT NULL DEFAULT '' COMMENT '业务名称',
  `f_business_show` varchar(64) NOT NULL DEFAULT '' COMMENT '业务显示名称',
  `f_business_mark` text COMMENT '业务描述',
  `f_business_order` int(11) DEFAULT '20' COMMENT '排序优先级',
  `f_create_person` varchar(64) NOT NULL DEFAULT '' COMMENT '创建者',
  `f_create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`f_business_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_business`
--

LOCK TABLES `t_business` WRITE;
/*!40000 ALTER TABLE `t_business` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_business` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_config`
--

DROP TABLE IF EXISTS `t_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_config` (
  `f_config_id` int(11) NOT NULL AUTO_INCREMENT,
  `f_config_name` varchar(128) DEFAULT NULL COMMENT '配置文件名',
  `f_config_version` int(11) NOT NULL COMMENT '版本号',
  `f_config_content` longtext COMMENT '配置内容',
  `f_create_person` varchar(16) DEFAULT NULL COMMENT '上传人',
  `f_create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  `f_config_mark` text COMMENT '上传描述',
  `f_app_server` varchar(64) NOT NULL COMMENT '应用名.服务名',
  `f_pod_seq` int(11) DEFAULT NULL,
  PRIMARY KEY (`f_config_id`),
  UNIQUE KEY `f_app_server` (`f_app_server`,`f_config_name`,`f_pod_seq`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_config`
--

LOCK TABLES `t_config` WRITE;
/*!40000 ALTER TABLE `t_config` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_config_history`
--

DROP TABLE IF EXISTS `t_config_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_config_history` (
  `f_history_id` int(11) NOT NULL AUTO_INCREMENT,
  `f_config_id` int(11) NOT NULL,
  `f_config_name` varchar(128) DEFAULT NULL COMMENT '配置文件名',
  `f_config_version` int(11) NOT NULL COMMENT '版本号',
  `f_pod_seq` int(11) DEFAULT NULL,
  `f_config_content` longtext COMMENT '配置内容',
  `f_create_person` varchar(16) DEFAULT NULL COMMENT '上传人',
  `f_create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  `f_config_mark` text COMMENT '上传描述',
  `f_app_server` varchar(64) NOT NULL COMMENT '应用名.服务名',
  PRIMARY KEY (`f_history_id`),
  UNIQUE KEY `f_app_server` (`f_app_server`,`f_config_name`,`f_config_version`,`f_pod_seq`),
  KEY `f_config_id` (`f_config_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_config_history`
--

LOCK TABLES `t_config_history` WRITE;
/*!40000 ALTER TABLE `t_config_history` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_config_history` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_default_value`
--

DROP TABLE IF EXISTS `t_default_value`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_default_value` (
  `f_id` int(11) NOT NULL AUTO_INCREMENT,
  `f_label` varchar(32) NOT NULL,
  `f_value` json NOT NULL,
  PRIMARY KEY (`f_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_default_value`
--

LOCK TABLES `t_default_value` WRITE;
/*!40000 ALTER TABLE `t_default_value` DISABLE KEYS */;
INSERT INTO `t_default_value` VALUES (1,'ServerServantElem','{\"Port\": 10000, \"IsTcp\": true, \"IsTaf\": true, \"Threads\": 3, \"Timeout\": 60000, \"Capacity\": 10000, \"HostPort\": 0, \"Connections\": 10000}'),(2,'ServerOption','{\"StopScript\": \"\", \"AsyncThread\": 3, \"StartScript\": \"\", \"MonitorScript\": \"\", \"ServerProfile\": \"\", \"ServerTemplate\": \"taf.default\", \"RemoteLogEnable\": true, \"ServerImportant\": 2, \"RemoteLogReserveTime\": 65, \"RemoteLogCompressTime\": 1}'),(3,'ServerK8S','{\"HostIpc\": false, \"HostPort\": false, \"Replicas\": 0, \"HostNetwork\": false, \"NodeSelector\": {\"Kind\": \"AbilityPool\", \"Value\": []}}'),(4,'K8SNodeSelectorKind','[\"AbilityPool\", \"PublicPool\", \"NodeBind\"]'),(5,'ServerTypeOptional','[\"taf_cpp\", \"taf_java_war\", \"taf_java_jar\", \"taf_node\", \"taf_node8\", \"taf_node10\", \"taf_node_pkg\", \"not_taf\"]');
/*!40000 ALTER TABLE `t_default_value` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_deploy_approval`
--

DROP TABLE IF EXISTS `t_deploy_approval`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_deploy_approval` (
  `f_approval_id` int(11) NOT NULL AUTO_INCREMENT,
  `f_approval_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '审批时间',
  `f_approval_person` varchar(16) NOT NULL DEFAULT '' COMMENT '审批人',
  `f_approval_result` tinyint(1) NOT NULL COMMENT '审批结果',
  `f_approval_mark` text COMMENT '审批描述',
  `f_request_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '申请时间',
  `f_request_person` varchar(16) NOT NULL COMMENT '申请人',
  `f_server_app` varchar(64) NOT NULL COMMENT '应用名',
  `f_server_name` varchar(64) NOT NULL COMMENT '服务名',
  `f_server_mark` text COMMENT '服务描述',
  `f_server_servant` json DEFAULT NULL COMMENT '服务servant参数',
  `f_server_option` json DEFAULT NULL COMMENT '服务conf参数',
  `f_server_k8s` json DEFAULT NULL COMMENT '服务k8s参数',
  PRIMARY KEY (`f_approval_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_deploy_approval`
--

LOCK TABLES `t_deploy_approval` WRITE;
/*!40000 ALTER TABLE `t_deploy_approval` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_deploy_approval` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_deploy_request`
--

DROP TABLE IF EXISTS `t_deploy_request`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_deploy_request` (
  `f_request_id` int(11) NOT NULL AUTO_INCREMENT,
  `f_request_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '申请时间',
  `f_request_person` varchar(16) NOT NULL COMMENT '申请人',
  `f_server_app` varchar(64) NOT NULL COMMENT '应用名',
  `f_server_name` varchar(64) NOT NULL COMMENT '服务名',
  `f_server_mark` text COMMENT '服务描述',
  `f_server_servant` json DEFAULT NULL COMMENT '服务servant参数',
  `f_server_option` json DEFAULT NULL COMMENT '服务option参数',
  `f_server_k8s` json DEFAULT NULL COMMENT '服务k8s参数',
  PRIMARY KEY (`f_request_id`),
  UNIQUE KEY `f_server_app` (`f_server_app`,`f_server_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_deploy_request`
--

LOCK TABLES `t_deploy_request` WRITE;
/*!40000 ALTER TABLE `t_deploy_request` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_deploy_request` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_login_temp_info`
--

DROP TABLE IF EXISTS `t_login_temp_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_login_temp_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ticket` varchar(256) DEFAULT NULL,
  `uid` varchar(256) DEFAULT NULL,
  `expire_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_login_temp_info`
--

LOCK TABLES `t_login_temp_info` WRITE;
/*!40000 ALTER TABLE `t_login_temp_info` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_login_temp_info` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_node`
--

DROP TABLE IF EXISTS `t_node`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_node` (
  `f_node_name` varchar(64) NOT NULL COMMENT 'NodeName',
  `f_ability` json DEFAULT NULL,
  `f_public` tinyint(4) DEFAULT NULL,
  `f_address` json DEFAULT NULL,
  `f_info` json DEFAULT NULL,
  `f_resource_version` bigint(20) NOT NULL DEFAULT '0' COMMENT 'k8s资源版本号',
  PRIMARY KEY (`f_node_name`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_node`
--

LOCK TABLES `t_node` WRITE;
/*!40000 ALTER TABLE `t_node` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_node` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_pod`
--

DROP TABLE IF EXISTS `t_pod`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_pod` (
  `f_pod_id` char(36) NOT NULL,
  `f_pod_name` varchar(64) NOT NULL,
  `f_pod_ip` char(15) DEFAULT '',
  `f_node_ip` char(15) DEFAULT '',
  `f_server_app` varchar(64) NOT NULL default '' COMMENT '应用名',
  `f_server_name` varchar(64) NOT NULL default '' COMMENT '服务名',
  `f_service_version` int(11) DEFAULT NULL COMMENT '服务版本',
  `f_create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `f_setting_state` enum('Inactive','Activating','Active','Deactivating') DEFAULT 'Active',
  `f_present_state` varchar(12) DEFAULT 'Unknown',
  `f_present_message` text,
  `f_update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `f_resource_version` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`f_pod_id`),
  KEY `f_server_app` (`f_server_app`,`f_server_name`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_pod`
--

LOCK TABLES `t_pod` WRITE;
/*!40000 ALTER TABLE `t_pod` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_pod` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_pod_history`
--

DROP TABLE IF EXISTS `t_pod_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_pod_history` (
  `f_pod_id` char(36) NOT NULL,
  `f_pod_name` varchar(64) NOT NULL,
  `f_pod_ip` char(15) DEFAULT NULL,
  `f_node_ip` char(15) DEFAULT NULL,
  `f_server_app` varchar(64) NOT NULL COMMENT '应用名',
  `f_server_name` varchar(64) NOT NULL COMMENT '服务名',
  `f_service_version` int(11) DEFAULT NULL COMMENT '服务版本',
  `f_create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '创建时间',
  `f_delete_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '死亡时间',
  PRIMARY KEY (`f_pod_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_pod_history`
--

LOCK TABLES `t_pod_history` WRITE;
/*!40000 ALTER TABLE `t_pod_history` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_pod_history` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_server`
--

DROP TABLE IF EXISTS `t_server`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_server` (
  `f_server_id` int(11) NOT NULL AUTO_INCREMENT,
  `f_server_app` varchar(64) NOT NULL COMMENT '应用名',
  `f_server_name` varchar(64) NOT NULL COMMENT '服务名',
  `f_server_mark` text COMMENT '服务描述',
  `f_server_type` varchar(32) NOT NULL DEFAULT '' COMMENT '服务类型',
  `f_deploy_person` varchar(16) NOT NULL DEFAULT '' COMMENT '创建人',
  `f_deploy_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `f_approval_person` varchar(32) NOT NULL DEFAULT '' COMMENT '审批人',
  `f_server_kind` enum('System','Normal','DCache') NOT NULL DEFAULT 'Normal',
  PRIMARY KEY (`f_server_id`),
  UNIQUE KEY `f_server_app` (`f_server_app`,`f_server_name`),
  CONSTRAINT `t_server_ibfk_1` FOREIGN KEY (`f_server_app`) REFERENCES `t_app` (`f_app_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_server`
--

LOCK TABLES `t_server` WRITE;
/*!40000 ALTER TABLE `t_server` DISABLE KEYS */;
INSERT INTO `t_server` VALUES (1,'taf','tafnotify','','taf_cpp','admin',current_timestamp,'admin','System'),(2,'taf','tafconfig','','taf_cpp','admin',current_timestamp,'admin','System'),(3,'taf','taflog','','taf_cpp','admin',current_timestamp,'admin','System'),(4,'taf','tafstat','','taf_cpp','admin',current_timestamp,'admin','System'),(5,'taf','tafquerystat','','taf_cpp','admin',current_timestamp,'admin','System'),(6,'taf','tafproperty','','taf_cpp','admin',current_timestamp,'admin','System'),(7,'taf','tafqueryproperty','','taf_cpp','admin',current_timestamp,'admin','System');
/*!40000 ALTER TABLE `t_server` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_server_adapter`
--

DROP TABLE IF EXISTS `t_server_adapter`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_server_adapter` (
  `f_adapter_id` int(11) NOT NULL AUTO_INCREMENT,
  `f_server_id` int(11) NOT NULL,
  `f_name` varchar(64) NOT NULL,
  `f_threads` int(11) NOT NULL DEFAULT '3',
  `f_connections` int(11) NOT NULL DEFAULT '1000',
  `f_port` int(11) NOT NULL,
  `f_capacity` int(11) NOT NULL DEFAULT '10000',
  `f_timeout` int(11) NOT NULL DEFAULT '6000',
  `f_is_taf` tinyint(1) NOT NULL,
  `f_is_tcp` tinyint(1) NOT NULL,
  PRIMARY KEY (`f_adapter_id`),
  UNIQUE KEY `f_server_id` (`f_server_id`,`f_name`),
  CONSTRAINT `t_server_adapter_ibfk_1` FOREIGN KEY (`f_server_id`) REFERENCES `t_server` (`f_server_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_server_adapter`
--

LOCK TABLES `t_server_adapter` WRITE;
/*!40000 ALTER TABLE `t_server_adapter` DISABLE KEYS */;
INSERT INTO `t_server_adapter` VALUES (1,1,'NotifyObj',3,10000,10000,10000,60000,1,1),(2,2,'ConfigObj',3,10000,10000,10000,60000,1,1),(3,3,'LogObj',3,10000,10000,10000,60000,1,1),(4,4,'StatObj',3,10000,10000,10000,60000,1,1),(5,5,'QueryObj',3,10000,10000,10000,60000,1,1),(6,6,'PropertyObj',3,10000,10000,10000,60000,1,1),(7,7,'QueryObj',3,10000,10000,10000,60000,1,1);
/*!40000 ALTER TABLE `t_server_adapter` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_server_k8s`
--

DROP TABLE IF EXISTS `t_server_k8s`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_server_k8s` (
  `f_server_id` int(11) NOT NULL,
  `f_server_app` varchar(64) NOT NULL COMMENT '应用名',
  `f_server_name` varchar(64) NOT NULL COMMENT '服务名',
  `f_replicas` int(11) NOT NULL COMMENT 'Pod 副本数',
  `f_host_network` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否使用 HostNetwork',
  `f_host_ipc` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否使用 HostIpc',
  `f_not_stacked` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否允许 Pod堆 叠部署',
  `f_node_selector` json NOT NULL COMMENT 'Pod的亲和性与反亲和性配置',
  `f_host_port` json DEFAULT NULL COMMENT 'Pod 的 HostPort 配置',
  PRIMARY KEY (`f_server_id`),
  CONSTRAINT `t_server_k8s_ibfk_1` FOREIGN KEY (`f_server_id`) REFERENCES `t_server` (`f_server_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_server_k8s`
--

LOCK TABLES `t_server_k8s` WRITE;
/*!40000 ALTER TABLE `t_server_k8s` DISABLE KEYS */;
INSERT INTO `t_server_k8s` VALUES (1,'taf','tafnotify',1,0,0,0,'{\"Kind\": \"AbilityPool\", \"Value\": []}',NULL),(2,'taf','tafconfig',1,0,0,0,'{\"Kind\": \"AbilityPool\", \"Value\": []}',NULL),(3,'taf','taflog',1,0,0,0,'{\"Kind\": \"AbilityPool\", \"Value\": []}',NULL),(4,'taf','tafstat',1,0,0,0,'{\"Kind\": \"AbilityPool\", \"Value\": []}',NULL),(5,'taf','tafquerystat',1,0,0,0,'{\"Kind\": \"AbilityPool\", \"Value\": []}',NULL),(6,'taf','tafproperty',1,0,0,0,'{\"Kind\": \"AbilityPool\", \"Value\": []}',NULL),(7,'taf','tafqueryproperty',1,0,0,0,'{\"Kind\": \"AbilityPool\", \"Value\": []}',NULL);
/*!40000 ALTER TABLE `t_server_k8s` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_server_notify`
--

DROP TABLE IF EXISTS `t_server_notify`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_server_notify` (
  `f_notify_id` int(11) NOT NULL AUTO_INCREMENT,
  `f_app_server` varchar(64) NOT NULL,
  `f_pod_name` varchar(64) NOT NULL DEFAULT '',
  `f_notify_level` varchar(12) DEFAULT '',
  `f_notify_message` text,
  `f_notify_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `f_notify_thread` varchar(16) DEFAULT '',
  `f_notify_source` varchar(12) DEFAULT '',
  PRIMARY KEY (`f_notify_id`),
  KEY `f_app_server` (`f_app_server`,`f_pod_name`,`f_notify_level`),
  KEY `f_notify_time` (`f_notify_time`)
) ENGINE=MyISAM AUTO_INCREMENT=341 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_server_notify`
--

LOCK TABLES `t_server_notify` WRITE;
/*!40000 ALTER TABLE `t_server_notify` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_server_notify` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_server_option`
--

DROP TABLE IF EXISTS `t_server_option`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_server_option` (
  `f_server_id` int(11) NOT NULL,
  `f_server_template` varchar(16) NOT NULL COMMENT '服务模板',
  `f_server_profile` text NOT NULL,
  `f_start_script_path` varchar(128) NOT NULL DEFAULT '',
  `f_stop_script_path` varchar(128) NOT NULL DEFAULT '',
  `f_monitor_script_path` varchar(128) DEFAULT NULL,
  `f_async_thread` int(11) NOT NULL DEFAULT '3',
  `f_important_type` tinyint(1) DEFAULT '0',
  `f_remote_log_reserve_time` int(11) NOT NULL DEFAULT '65',
  `f_remote_log_compress_time` int(11) NOT NULL DEFAULT '2',
  `f_remote_log_type` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`f_server_id`),
  CONSTRAINT `t_server_option_ibfk_1` FOREIGN KEY (`f_server_id`) REFERENCES `t_server` (`f_server_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_server_option`
--

LOCK TABLES `t_server_option` WRITE;
/*!40000 ALTER TABLE `t_server_option` DISABLE KEYS */;
INSERT INTO `t_server_option` VALUES (1,'taf.tafnotify','','','','',3,2,65,1,1),(2,'taf.tafconfig','','','','',3,2,65,1,1),(3,'taf.default','','','','',3,0,65,2,0),(4,'taf.default','<taf>\n	sql=CREATE TABLE `${TABLE}`( `stattime` timestamp NOT NULL default CURRENT_TIMESTAMP,`f_date` date NOT NULL default \'1970-01-01\', `f_tflag` varchar(8) NOT NULL default \'\',`source_id` varchar(15) NOT NULL default \'\',`master_name` varchar(128) NOT NULL default \'\',`slave_name` varchar(128) NOT NULL default \'\',`interface_name` varchar(128) NOT NULL default \'\',`taf_version` varchar(16) NOT NULL default \'\',`master_ip` varchar(15) NOT NULL default \'\',`slave_ip` varchar(21) NOT NULL default \'\',`slave_port` int(10) NOT NULL default 0,`return_value` int(11) NOT NULL default 0,`succ_count` int(10) unsigned default NULL,`timeout_count` int(10) unsigned default NULL,`exce_count` int(10) unsigned default NULL,`interv_count` varchar(128) default NULL,`total_time` bigint(20) unsigned default NULL,`ave_time` int(10) unsigned default NULL,`maxrsp_time` int(10) unsigned default NULL,`minrsp_time` int(10) unsigned default NULL,PRIMARY KEY (`source_id`,`f_date`,`f_tflag`,`master_name`,`slave_name`,`interface_name`,`master_ip`,`slave_ip`,`slave_port`,`return_value`,`taf_version`),KEY `IDX_TIME` (`stattime`),KEY `IDC_MASTER` (`master_name`),KEY `IDX_INTERFACENAME` (`interface_name`),KEY `IDX_FLAGSLAVE` (`f_tflag`,`slave_name`), KEY `IDX_SLAVEIP` (`slave_ip`),KEY `IDX_SLAVE` (`slave_name`),KEY `IDX_RETVALUE` (`return_value`),KEY `IDX_MASTER_IP` (`master_ip`),KEY `IDX_F_DATE` (`f_date`)) ENGINE\\=InnoDB DEFAULT CHARSET\\=utf8mb4\n        enWeighted=1\n	<masteripGroup>\n		taf.tafstat;1.1.1.1\n	</masteripGroup>\n	<hashmap>\n		masterfile=hashmap_master.txt\n		slavefile=hashmap_slave.txt\n		insertInterval=5\n		enableStatCount=0\n		size=256M\n		countsize=1M\n	</hashmap>\n	<reapSql>\n		interval=5\n		insertDbThreadNum=4\n	</reapSql>\n	<statdb>\n		<db1>\n           dbhost=_DB_TAF_STAT_HOST_\n           dbname=_DB_TAF_STAT_DATABASE_\n		   tbname=taf_stat_\n           dbuser=_DB_TAF_STAT_USER_\n           dbpass=_DB_TAF_STAT_PASSWORD_\n           dbport=_DB_TAF_STAT_PORT_            \n		   charset=utf8mb4\n		</db1>\n	</statdb>\n</taf>\n','','','',3,2,65,1,1),(5,'taf.default','<taf>\n  <statdb>\n    <db1>\n      dbhost=_DB_TAF_STAT_HOST_\n      dbname=_DB_TAF_STAT_DATABASE_\n      dbuser=_DB_TAF_STAT_USER_\n      dbpass=_DB_TAF_STAT_PASSWORD_\n      dbport=_DB_TAF_STAT_PORT_   \n      charset=utf8mb4\n    </db1>\n  </statdb>\n</taf>\n','','','',3,2,65,1,1),(6,'taf.default','<taf>\n	sql=CREATE TABLE `${TABLE}` (`stattime` timestamp NOT NULL default CURRENT_TIMESTAMP,`f_date` date NOT NULL default \'1970-01-01\', `f_tflag` varchar(8) NOT NULL default \'\',`master_name` varchar(128) NOT NULL default \'\',`master_ip` varchar(16) default NULL,`property_name` varchar(100) default NULL,`set_name` varchar(15) NOT NULL default \'\',`set_area` varchar(15) NOT NULL default \'\',`set_id` varchar(15) NOT NULL default \'\',`policy` varchar(20) default NULL,`value` varchar(255) default NULL, KEY (`f_date`,`f_tflag`,`master_name`,`master_ip`,`property_name`,`policy`),KEY `IDX_MASTER_NAME` (`master_name`),KEY `IDX_MASTER_IP` (`master_ip`),KEY `IDX_TIME` (`stattime`)) ENGINE\\=Innodb\n	<propertydb>\n		<db1>\n			tbname=taf_property_\n            dbhost=_DB_TAF_PROPERTY_HOST_\n            dbname=_DB_TAF_PROPERTY_DATABASE_\n            dbuser=_DB_TAF_PROPERTY_USER_\n            dbpass=_DB_TAF_PROPERTY_PASSWORD_\n            dbport=_DB_TAF_PROPERTY_PORT_\n			charset=utf8mb4\n		</db1>\n	</propertydb>\n	<hashmap>\n		factor=1.5\n		file=hashmap.txt\n		insertInterval=5\n		maxBlock=200\n		minBlock=100\n		size=128M\n	</hashmap>\n	<reapSql>\n		Interval=10\n		sql=insert ignore into t_master_property select  master_name, property_name, policy from ${TABLE}  group by  master_name, property_name, policy;\n	</reapSql>\n</taf>','','','',3,2,65,1,1),(7,'taf.default','<taf>\n  <propertydb>\n    <db1>\n      dbhost=_DB_TAF_PROPERTY_HOST_\n      dbname=_DB_TAF_PROPERTY_DATABASE_\n      dbuser=_DB_TAF_PROPERTY_USER_\n      dbpass=_DB_TAF_PROPERTY_PASSWORD_\n      dbport=_DB_TAF_PROPERTY_PORT_\n      charset=utf8mb4\n    </db1>\n  </propertydb>\n</taf>','','','',3,2,65,1,1);
/*!40000 ALTER TABLE `t_server_option` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_service_enabled`
--

DROP TABLE IF EXISTS `t_service_enabled`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_service_enabled` (
  `f_service_id` int(11) DEFAULT NULL,
  `f_server_id` int(11) NOT NULL,
  `f_enable_person` varchar(16) NOT NULL DEFAULT '' COMMENT '启用人',
  `f_enable_mark` text COMMENT '启用描述',
  `f_enable_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '启用时间',
  PRIMARY KEY (`f_server_id`),
  UNIQUE KEY `f_service_id` (`f_service_id`),
  CONSTRAINT `t_service_enabled_ibfk_1` FOREIGN KEY (`f_service_id`) REFERENCES `t_service_pool` (`f_service_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_service_enabled`
--

LOCK TABLES `t_service_enabled` WRITE;
/*!40000 ALTER TABLE `t_service_enabled` DISABLE KEYS */;
INSERT INTO `t_service_enabled` VALUES (1,1,'admin','',current_timestamp),(2,2,'admin','',current_timestamp),(3,3,'admin','',current_timestamp),(4,4,'admin','',current_timestamp),(5,5,'admin','',current_timestamp),(6,6,'admin','',current_timestamp),(7,7,'admin','',current_timestamp);
/*!40000 ALTER TABLE `t_service_enabled` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_service_pool`
--

DROP TABLE IF EXISTS `t_service_pool`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_service_pool` (
  `f_service_id` int(11) NOT NULL AUTO_INCREMENT,
  `f_service_version` varchar(12) NOT NULL COMMENT '版本号',
  `f_service_mark` text COMMENT '描述',
  `f_service_image` varchar(128) NOT NULL COMMENT '镜像地址',
  `f_image_detail` json DEFAULT NULL COMMENT '镜像元信息',
  `f_server_id` int(11) NOT NULL,
  `f_server_app` varchar(64) NOT NULL COMMENT '应用名',
  `f_server_name` varchar(64) NOT NULL COMMENT '服务名',
  `f_create_person` varchar(16) NOT NULL DEFAULT '' COMMENT '创建人',
  `f_create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`f_service_id`),
  UNIQUE KEY `f_server_id` (`f_server_id`,`f_service_version`),
  CONSTRAINT `t_service_pool_ibfk_1` FOREIGN KEY (`f_server_id`) REFERENCES `t_server` (`f_server_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_service_pool`
--

LOCK TABLES `t_service_pool` WRITE;
/*!40000 ALTER TABLE `t_service_pool` DISABLE KEYS */;
INSERT INTO `t_service_pool` VALUES (1,'10000',NULL,'_DOCKER_REGISTRY_URL_/taf.tafnotify:10000','{}',1,'taf','tafnotify','admin',current_timestamp),(2,'10000',NULL,'_DOCKER_REGISTRY_URL_/taf.tafconfig:10000','{}',2,'taf','tafconfig','admin',current_timestamp),(3,'10000',NULL,'_DOCKER_REGISTRY_URL_/taf.taflog:10000','{}',3,'taf','taflog','admin',current_timestamp),(4,'10000',NULL,'_DOCKER_REGISTRY_URL_/taf.tafstat:10000','{}',4,'taf','tafstat','admin',current_timestamp),(5,'10000',NULL,'_DOCKER_REGISTRY_URL_/taf.tafquerystat:10000','{}',5,'taf','tafquerystat','admin',current_timestamp),(6,'10000',NULL,'_DOCKER_REGISTRY_URL_/taf.tafproperty:10000','{}',6,'taf','tafproperty','admin',current_timestamp),(7,'10000',NULL,'_DOCKER_REGISTRY_URL_/taf.tafqueryproperty:10000','{}',7,'taf','tafqueryproperty','admin',current_timestamp);
/*!40000 ALTER TABLE `t_service_pool` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_template`
--

DROP TABLE IF EXISTS `t_template`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_template` (
  `f_template_id` int(11) NOT NULL AUTO_INCREMENT,
  `f_template_name` varchar(128) NOT NULL DEFAULT '',
  `f_template_parent` varchar(128) NOT NULL DEFAULT '',
  `f_template_content` text,
  `f_create_person` varchar(16) DEFAULT NULL COMMENT '创建人',
  `f_create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '创建时间',
  `f_create_mark` text COMMENT '创建描述',
  `f_update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `f_update_person` varchar(16) DEFAULT NULL,
  `f_update_mark` text COMMENT '更改描述',
  PRIMARY KEY (`f_template_id`),
  UNIQUE KEY `f_template_name` (`f_template_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_template`
--

LOCK TABLES `t_template` WRITE;
/*!40000 ALTER TABLE `t_template` DISABLE KEYS */;
INSERT INTO `t_template` VALUES (1,'taf.default','taf.default','<taf>\n    <application>\n    #是否启用SET分组\n    enableset=${enableset}\n    #SET分组的全名.(mtt.s.1)\n    setdivision=${setdivision}\n    <client>\n        #地址\n        locator =${locator}\n        #缺省3s(毫秒)\n        sync-invoke-timeout = 3000\n        #最大超时时间(毫秒)\n        async-invoke-timeout =5000\n        #重新获取服务列表时间间隔(毫秒)\n        refresh-endpoint-interval = 60000\n        #模块间调用[可选]\n        stat            = taf.tafstat.StatObj\n        #属性上报服务\n        property                    = taf.tafproperty.PropertyObj\n        #上报间隔时间(毫秒)\n        report-interval            = 60000\n        #stat采样比1:n 例如sample-rate为1000时 采样比为千分之一\n         sample-rate = 100000\n        #1分钟内stat最大采样条数\n         max-sample-count = 50\n\n        #网络发送线程个数\n        sendthread      = 1\n        #网络接收线程个数\n        recvthread      = 1\n        #网络异步回调线程个数\n        asyncthread      = ${asyncthread}\n        #模块名称\n        modulename      = ${modulename}\n    </client>\n        \n    #定义所有绑定的IP\n    <server>\n        #应用名称\n        app      = ${app}\n        #服务名称\n        server  = ${server}\n        #本地ip\n       localip  = ${localip}\n\n        #本地管理套接字[可选]\n        local  = ${local}\n        #服务的数据目录,可执行文件,配置文件等\n        basepath = ${basepath}\n        #\n        datapath = ${datapath}\n        #日志路径\n        logpath  = ${logpath}\n        #日志大小\n        logsize = 10M\n        #日志数量\n        #   lognum = 10\n        #配置中心的地址[可选]\n        config  = taf.tafconfig.ConfigObj\n        #信息中心的地址[可选]\n        notify  = taf.tafnotify.NotifyObj\n        #远程LogServer[可选]\n        log    = taf.taflog.LogObj\n        #关闭服务时等待时间\n         deactivating-timeout = 2000\n        #是否启用用户级线程切换（默认为0，不启用）\n\n         openthreadcontext = 0\n\n        #用户级线程上下文个数 (openthreadcontext为1时生效,默认10000)\n\n         threadcontextnum  = 10000\n\n        #用户级线程上下文栈大小 (openthreadcontext为1时生效,默认32k)\n\n        threadcontextstack = 32768    \n\n        #滚动日志等级默认值\n        logLevel=DEBUG\n    </server>          \n    </application>\n    <log>\n		logpath=/usr/local/app/taf/remote_app_log\n		logthread=2\n    </log>\n</taf>','admin',current_timestamp,'',current_timestamp,'admin',NULL),(2,'taf.cpp','taf.default','','admin',current_timestamp,'',current_timestamp,'admin',NULL),(3,'taf.java','taf.default','<taf>\n <application>\n   \n   <server>\n    \n      mainclass=com.taf.server.startup.Main\n     classpath=${basepath}/conf:${basepath}/WEB-INF/classes:${basepath}/WEB-INF/lib\n      #jvmparams=-Dcom.sun.management.jmxremote.ssl\\=false -Dcom.sun.management.jmxremote.authenticate\\=false -Xms2000m -Xmx2000m -Xmn1000m -Xss1000k -XX:PermSize\\=128M -XX:+UseConcMarkSweepGC -XX:CMSInitiatingOccupancyFraction\\=60 -XX:+PrintGCApplicationStoppedTime -XX:+PrintGCDateStamps -XX:+CMSParallelRemarkEnabled -XX:+CMSScavengeBeforeRemark -XX:+UseCMSCompactAtFullCollection -XX:CMSFullGCsBeforeCompaction\\=0 -verbosegc -XX:+PrintGCDetails -XX:ErrorFile\\=${logpath}/${app}/${server}/jvm_error.log\n\njvmparams=-XX:ErrorFile\\=${logpath}/${app}/${server}/jvm_error.log -Dtaf=true\n      sessiontimeout=120000\n     sessioncheckinterval=60000\n      tcpnodelay=true\n     udpbuffersize=8192\n      charsetname=UTF-8\n     backupfiles=conf\n      loglevel=DEBUG\n    </server>\n </application>\n</taf>','admin',current_timestamp,'',current_timestamp,'admin',NULL),(4,'taf.nodejs8','taf.default','<taf>\n    <application>\n    <server>\n        env=NODE_PATH=/usr/local/app/taf/tafnode/lib/node_modules\n\n    exefile=node8\n    </server>\n    <client>\n      \n        #网络异步回调线程个数\n        asyncthread      = 2\n\n    </client>\n        \n  </application>\n\n  \n</taf>','admin',current_timestamp,'',current_timestamp,'admin',NULL),(5,'taf.tafconfig','taf.default','<taf>\n    <application>\n    <server>\n        #log    = taf.taflog4other.LogObj\n        logLevel = DEBUG\n    </server>          \n    </application>\n    <db>\n        charset=utf8mb4\n        dbhost=_DB_TAF_HOST_\n        dbname=_DB_TAF_DATABASE_\n        dbpass=_DB_TAF_PASSWORD_\n        dbport=_DB_TAF_PORT_\n        dbuser=_DB_TAF_USER_\n    </db>\n</taf>','admin',current_timestamp,'',current_timestamp,'admin',NULL),(6,'taf.tafnotify','taf.default','<taf>\n    <application>\n    <server>\n        #log    = taf.taflog4other.LogObj\n        logLevel = DEBUG\n    </server>          \n    </application>\n    <db>\n        charset=utf8mb4\n        dbhost=_DB_TAF_HOST_\n        dbname=_DB_TAF_DATABASE_\n        dbpass=_DB_TAF_PASSWORD_\n        dbport=_DB_TAF_PORT_\n        dbuser=_DB_TAF_USER_\n    </db>\n    sql = create table t_server_notify(f_notify_id int auto_increment primary key,f_app_server varchar(64) not null,f_pod_name varchar(64) default \'\' not null,f_notify_level int null,f_notify_message text null,f_notify_time timestamp default CURRENT_TIMESTAMP not null,f_notify_thread  varchar(16) null,f_notify_source text null) engine = MyISAM DEFAULT CHARSET=utf8mb4\n</taf>','admin',current_timestamp,'',current_timestamp,'admin',NULL);
/*!40000 ALTER TABLE `t_template` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_token`
--

DROP TABLE IF EXISTS `t_token`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_token` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uid` varchar(128) DEFAULT NULL,
  `token` varchar(128) DEFAULT NULL,
  `valid` int(11) DEFAULT NULL,
  `expire_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  KEY `token` (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_token`
--

LOCK TABLES `t_token` WRITE;
/*!40000 ALTER TABLE `t_token` DISABLE KEYS */;
/*!40000 ALTER TABLE `t_token` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `t_user_info`
--

DROP TABLE IF EXISTS `t_user_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `t_user_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` varchar(128) DEFAULT NULL,
  `password` varchar(256) NOT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_name` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `t_user_info`
--

LOCK TABLES `t_user_info` WRITE;
/*!40000 ALTER TABLE `t_user_info` DISABLE KEYS */;
INSERT INTO `t_user_info` VALUES (1,'admin','d033e22ae348aeb5660fc2140aec35850c4da997',current_timestamp);
/*!40000 ALTER TABLE `t_user_info` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
