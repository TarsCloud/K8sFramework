const PageController = require('../controller/page/PageController');
const LoginController = require('../controller/login/LoginController');
const LocaleController = require('../controller/locale/LocaleController');
const DemoLocaleController = require('../../demo/app/controller/locale/LocaleController');
const logViewController = require('../controller/logview/logviewController')
const RpcController = require('../controller/rpc/RpcController');
const MonitorController = require('../controller/monitor/MonitorController');
const GatewayController = require('../controller/gateway/GatewayController');
const ServerController = require('../controller/server/ServerController');
const ResourceController = require('../controller/resource/ResourceController');
const AuthController = require('../controller/auth/AuthController');

const pageConf = [
	//首页
	['get', '/', PageController.index]
];

const apiConf = [
	// 目录树接口
	['get', '/tree', RpcController.ServerTree],

	// 默认参数
	['get', '/default', RpcController.DefaultValue],

    // 服务状态
    ['get', '/server_notify_list', RpcController.NotifySelect, { ServerId: 'notEmpty' }],

    // 应用管理 ( 创建、列表、更新、删除 )
    ['post', '/application_create', RpcController.ApplicationCreate],
    ['get', '/application_select', RpcController.ApplicationSelect],
    ['post', '/application_update', RpcController.ApplicationUpdate, { AppName: 'notEmpty' }],
    ['get', '/application_delete', RpcController.ApplicationDelete, { AppName: 'notEmpty' }],

    // 业务管理 ( 创建、列表、更新、删除 )
    ['post', '/business_create', RpcController.BusinessCreate],
    ['get', '/business_select', RpcController.BusinessSelect],
    ['post', '/business_update', RpcController.BusinessUpdate, { BusinessName: 'notEmpty' }],
    ['get', '/business_delete', RpcController.BusinessDelete, { BusinessName: 'notEmpty' }],
    ['post', '/business_add_app', RpcController.BusinessAddApp, { BusinessName: 'notEmpty' }],
    ['post', '/business_list_app', RpcController.BusinessListByApp, { BusinessName: 'notEmpty' }],

    // 模版管理 ( 创建、列表、更新、删除 )
    ['post', '/template_create', RpcController.TemplateCreate],
    ['get', '/template_select', RpcController.TemplateSelect],
    ['post', '/template_update', RpcController.TemplateUpdate, { TemplateId: 'notEmpty' }],
    ['get', '/template_delete', RpcController.TemplateDelete, { TemplateId: 'notEmpty' }],

    // 服务部署 ( 创建、创建列表、删除、审批、审批列表 )
	['post', '/deploy_create', RpcController.ServerDeployCreate, { ServerApp: 'notEmpty', ServerName: 'notEmpty' }],
	['get', '/deploy_select', RpcController.ServerDeploySelect],
	['post', '/deploy_update', RpcController.ServerDeployUpdate, { DeployId: 'notEmpty' }],
	['get', '/deploy_delete', RpcController.ServerDeployDelete, { DeployId: 'notEmpty' }],
	['post', '/approval_create', RpcController.ServerApprovalCreate, { ServerApp: 'notEmpty', ServerName: 'notEmpty' }],
	['get', '/approval_select', RpcController.ServerApprovalSelect],

    // 服务配置文件 ( 创建、列表、更新、删除 )
    ['post', '/server_config_create', RpcController.ServerConfigCreate],
    ['get', '/server_config_select', RpcController.ServerConfigSelect, { ServerId: 'notEmpty' }],
    ['post', '/server_config_update', RpcController.ServerConfigUpdate, { ConfigId: 'notEmpty' }],
    ['get', '/server_config_delete', RpcController.ServerConfigDelete, { ConfigId: 'notEmpty' }],
    ['get', '/merged_node_config', RpcController.ServerConfigContent, { ServerId: 'notEmpty', ConfigName: 'notEmpty' }],

    // 服务配置文件历史记录 ( 列表、删除 )
    ['get', '/server_config_history_select', RpcController.ServerConfigHistroySelect, { ConfigId: 'notEmpty' }],
    ['get', '/server_config_history_delete', RpcController.ServerConfigHistroyDelete, { HistoryId: 'notEmpty' }],
    ['post', '/server_config_history_back', RpcController.ServerConfigHistoryBack, { HistoryId: 'notEmpty' }],

    // 服务管理 ( pod列表、pod历史列表、服务列表、服务更新、状态(重启、停止)、编辑、更新 )
    ['get', '/pod_list', RpcController.PodAliveSelect, { ServerId: 'notEmpty' }],
    ['get', '/pod_history_list', RpcController.PodPerishedSelect, { ServerId: 'notEmpty' }],

    ['get', '/server_list', RpcController.ServerSelect],
    ['post', '/server_update', RpcController.ServerUpdate, { ServerId: 'notEmpty' }],
    ['post', '/server_undeploy', RpcController.ServerUndeploy, { ServerId: 'notEmpty' }],

    ['get', '/server_option_select', RpcController.ServerOptionSelect, { ServerId: 'notEmpty' }],
    ['get', '/server_option_template', RpcController.ServerOptionTemplate, { ServerId: 'notEmpty' }],
    ['post', '/server_option_update', RpcController.ServerOptionUpdate, { ServerId: 'notEmpty' }],

    // 服务ServerAdapter ( 创建、列表、更新、删除 )
    ['post', '/server_adapter_create', RpcController.ServerAdapterCreate],
    ['get', '/server_adapter_select', RpcController.ServerAdapterSelect, { ServerId: 'notEmpty' }],
    ['post', '/server_adapter_update', RpcController.ServerAdapterUpdate, { AdapterId: 'notEmpty' }],
    ['get', '/server_adapter_delete', RpcController.ServerAdapterDelete, { AdapterId: 'notEmpty' }],

    // 节点管理 ( 创建、列表、更新、删除 )
    //['post', '/node_create', RpcController.NodeCreate],
    //['post', '/node_update', RpcController.NodeUpdate, { NodeName: 'notEmpty' }],
    //['get', '/node_delete', RpcController.NodeDelete, { NodeName: 'notEmpty' }],
    //['post', '/node_add', RpcController.NodeAdd],
    ['get', '/node_select', RpcController.NodeSelect],
    ['get', '/node_list', RpcController.NodeList],
    ['post', '/node_start', RpcController.NodeStartPublic, { NodeName: 'notEmpty' }],
    ['post', '/node_stop', RpcController.NodeStopPublic, { NodeName: 'notEmpty' }],

    // 亲和性管理 ( 创建、列表、更新、删除 )
    //['post', '/affinity_create', RpcController.AffinityCreate],
    //['get', '/affinity_select', RpcController.AffinitySelect],
    //['get', '/affinity_delete', RpcController.AffinityDelete, { NodeName: 'notEmpty', ServerApp: 'notEmpty' }],
    ['get', '/affinity_list_node', RpcController.AffinityListByNode],
    ['get', '/affinity_list_server', RpcController.AffinityListByServer],
    ['post', '/affinity_add_server', RpcController.AffinityAddServer, { NodeName: 'notEmpty', ServerApp: 'notEmpty' }],
    ['post', '/affinity_add_node', RpcController.AffinityAddNode, { NodeName: 'notEmpty', ServerApp: 'notEmpty' }],
    ['post', '/affinity_del_server', RpcController.AffinityDeleteServer, { NodeName: 'notEmpty', ServerApp: 'notEmpty' }],
    ['post', '/affinity_del_node', RpcController.AffinityDeleteNode, { NodeName: 'notEmpty', ServerApp: 'notEmpty' }],

    // 服务K8S ( 列表、更新 )
    ['get', '/server_k8s_select', RpcController.ServerK8SSelect, { ServerId: 'notEmpty' }],
    ['post', '/server_k8s_update', RpcController.ServerK8SUpdate, { ServerId: 'notEmpty' }],
    ['get', '/generate_host_port', RpcController.ServerK8SGenerateHostPort, { NodeList: 'notEmpty', NodePort: 'notEmpty' }],

    // 发布包 ( 上传及编译、上传状态、版本列表、版本发布 )
    ['post', '/patch_upload', RpcController.uploadPatchPackage, { ServerId: 'notEmpty' }],
    ['get', '/patch_upload_status', RpcController.uploadPatchStatus, { BuildId: 'notEmpty' }],
    ['get', '/patch_list', RpcController.ServicePoolSelect, { ServerId: 'notEmpty' }],
    ['get', '/patch_enabled', RpcController.ServiceEnabledSelect, { ServerId: 'notEmpty' }],
    ['post', '/patch_publish', RpcController.ServicePoolUpdate, { ServiceId: 'notEmpty' }],

	//shell
	['get', '/shell_domain', logViewController.getShellDomain],

	//权限管理
    ['get', '/is_enable_auth', AuthController.isEnableAuth],
    ['get', '/get_roles', AuthController.getRoles],
    ['get', '/is_admin', AuthController.isAdmin],
    ['get', '/get_auth_list', AuthController.getAuthList],
    ['post', '/update_auth', AuthController.updateAuth, {
        application: 'notEmpty',
        server_name: 'notEmpty',
        operator: 'notEmpty',
        developer: 'notEmpty'
    }],
    ['get', '/has_auth', AuthController.hasAuth, { application: 'notEmpty', role: 'notEmpty' }],
    ['get', '/userCenter', AuthController.userCenter],

    //登录管理
    ['get', '/get_login_uid', LoginController.getLoginUid],
    ['get', '/is_enable_login', LoginController.isEnableLogin],

	//语言包接口
	['get', '/get_locale', LocaleController.getLocale],
	['get', '/get_demo_locale', DemoLocaleController.getLocale],

	// 命令
	['get', '/send_command', ServerController.sendCommand, { serverApp: 'notEmpty', serverName: 'notEmpty', podIp: 'notEmpty', command: 'notEmpty' }],

	// 监控
    ['get', '/tarsstat_monitor_data', MonitorController.tarsstat],
	['get', '/tarsproperty_monitor_data', MonitorController.tarsproperty],
	
	//网关配置
    ['get', '/station_list', GatewayController.getStationList], 
    ['post', '/add_station', GatewayController.addStation,{
        f_station_id: 'notEmpty',
        f_name_cn: 'notEmpty'
    }],
    ['post', '/update_station', GatewayController.updateStation,{
        f_id: 'notEmpty',
        f_station_id: 'notEmpty',
        f_name_cn: 'notEmpty'
    }],
    ['post', '/delete_station', GatewayController.deleteStation, { f_id: 'notEmpty' }],

    ['get', '/upstream_list', GatewayController.getUpstreamList], 
    ['post', '/add_upstream', GatewayController.addUpstream,{
        f_upstream: 'notEmpty',
        f_addr: 'notEmpty',
        f_weight: 'notEmpty',
        f_fusing_onoff: 'notEmpty'
    }],
    ['post', '/update_upstream', GatewayController.updateUpstream,{
        f_id: 'notEmpty',
        f_upstream: 'notEmpty',
        f_addr: 'notEmpty',
        f_weight: 'notEmpty',
        f_fusing_onoff: 'notEmpty'
    }],
    ['post', '/delete_upstream', GatewayController.deleteUpstream, { f_id: 'notEmpty' }],

    ['get', '/httprouter_list', GatewayController.getHttpRouterList, { f_id: 'f_station_id' }], 
    ['post', '/add_httprouter', GatewayController.addHttpRouter,{
        f_station_id: 'notEmpty',
        f_server_name: 'notEmpty',
        f_path_rule: 'notEmpty',
        f_proxy_pass: 'notEmpty'
    }],
    ['post', '/update_httprouter', GatewayController.updateHttpRouter,{
        f_id: 'notEmpty',
        f_station_id: 'notEmpty',
        f_server_name: 'notEmpty',
        f_path_rule: 'notEmpty',
        f_proxy_pass: 'notEmpty'
    }],
    ['post', '/delete_httprouter', GatewayController.deleteHttpRouter, { f_id: 'notEmpty' }],
    
    ['get', '/bwlist', GatewayController.getBWList, { 
        type: 'notEmpty' 
    }], 
    ['post', '/add_bwlist', GatewayController.addBWList,{
        f_ip: 'notEmpty',
        type: 'notEmpty'
    }],
    ['post', '/delete_bwlist', GatewayController.deleteBWList, { f_id: 'notEmpty',type: 'notEmpty' }],
    ['get', '/get_flowcontrol', GatewayController.getFlowControl, { f_station_id: 'notEmpty'}],
    ['post', '/upsert_flowcontrol', GatewayController.upsertFlowControl, {
        f_station_id: 'notEmpty',
        f_duration: 'notEmpty',
        f_max_flow: 'notEmpty'
    }],
];

const clientConf = [
    ['get', '/get_tarsnode', ResourceController.getTarsNode],
];

module.exports = {pageConf, apiConf, clientConf};
