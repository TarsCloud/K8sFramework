#include "RegistryImp.h"
#include "RegistryServer.h"
#include "util/tc_mysql.h"
#include "K8SEndpointInterface.h"

void RegistryImp::initialize() {
    LOG->debug() << "RegistryImp init begin." << endl;
    try {
        TC_DBConf tcDBConf;
        TC_Config &gConf = g_app.getConfig();
        tcDBConf.loadFromMap(gConf.getDomainMap("/tars/db"));
        _mysqlReg.init(tcDBConf);
    }
    catch (TC_Config_Exception &ex) {
        LOG->error() << __FUNCTION__ << " exception: " << ex.what() << endl;
        exit(-1);
    }
    catch (TC_Mysql_Exception &ex) {
        LOG->error() << __FUNCTION__ << " exception: " << ex.what() << endl;
        exit(-1);
    }
    LOG->debug() << "RegistryImp init ok." << endl;
}

int
RegistryImp::getServers(const std::string &app, const std::string &serverName, tars::ServerDescriptor &serverDescriptor,
                        tars::CurrentPtr current) {
    const std::string sPrepareSql =
            "SELECT f_server_app,f_server_name,f_async_thread,f_server_profile,f_start_script_path,f_stop_script_path,f_monitor_script_path,f_server_template,f_name,f_threads,"
            "b.f_is_tcp,b.f_is_taf,b.f_connections,b.f_capacity,b.f_port,b.f_timeout FROM "
            "t_server a LEFT JOIN t_server_adapter b USING (f_server_id) LEFT JOIN t_server_option c USING (f_server_id) WHERE 1=1 ";

    std::string sCondition;
    if (!app.empty()) sCondition += "and a.f_server_app='" + tars::TC_Mysql::escapeString(app) + "' ";
    if (!serverName.empty()) sCondition += "and a.f_server_name='" + tars::TC_Mysql::escapeString(serverName) + "' ";

    try {

        tars::TC_Mysql::MysqlData res = _mysqlReg.queryRecord(sPrepareSql + sCondition);

        if (res.size() == 0) {
            return -1;
        }

        for (size_t i = 0; i < res.size(); i++) {
            if (i == 0) {
                serverDescriptor.profile = res[i]["f_server_profile"];
                serverDescriptor.startScript = res[i]["f_start_script_path"];
                serverDescriptor.stopScript = res[i]["f_stop_script_path"];
                serverDescriptor.monitorScript = res[i]["f_monitor_script_path"];
                const std::string &sServerTemplate = res[i]["f_server_template"];
                if (!sServerTemplate.empty()) {
                    string sResult;
                    std::string sTemplateContent = getTemplateContent(sServerTemplate, sResult, nullptr);
                    if (sTemplateContent.empty()) {
                        return -1;
                    }
                    TC_Config tParent, tProfile;
                    tParent.parseString(sTemplateContent);
                    tProfile.parseString(serverDescriptor.profile);
                    tParent.joinConfig(tProfile, true);
                    int iDefaultAsyncThreadNum = 1;
                    int iConfigAsyncThreadNum = TC_Common::strto<int>(TC_Common::trim(res[i]["f_async_thread"]));
                    iDefaultAsyncThreadNum = iConfigAsyncThreadNum > iDefaultAsyncThreadNum ? iConfigAsyncThreadNum
                                                                                            : iDefaultAsyncThreadNum;
                    serverDescriptor.asyncThreadNum = TC_Common::strto<int>(
                            tProfile.get("/tars/application/client<asyncthread>",
                                         TC_Common::tostr(iDefaultAsyncThreadNum)));
                    serverDescriptor.profile = tParent.tostr();
                }
            }
            const string &f_servant_name = res[i]["f_name"];
            if (f_servant_name.empty()) {
                continue;
            }
            AdapterDescriptor adapter;
            adapter.servant.append(app).append(".").append(serverName).append(".").append(f_servant_name);
            adapter.adapterName.append(app).append(".").append(serverName).append(".").append(f_servant_name).append(
                    ".Adapter");
            bool isTars = TC_Common::strto<bool>(res[i]["f_is_taf"]);
            adapter.protocol = isTars ? "tars" : "not_taf";

            bool isTcp = TC_Common::strto<bool>(res[i]["f_is_tcp"]);
            string f_port = res[i]["f_port"];
            string f_timeout = res[i]["f_timeout"];

            adapter.endpoint.append(isTcp ? "tcp" : "udp").append(" -h localip -p ").append(f_port).append(
                    " -t ").append(f_timeout);

            adapter.threadNum = res[i]["f_threads"];
            adapter.maxConnections = TC_Common::strto<int>(res[i]["f_connections"]);

            adapter.queuecap = TC_Common::strto<int>(res[i]["f_capacity"]);
            adapter.queuetimeout = TC_Common::strto<int>(f_timeout);

            serverDescriptor.adapters[adapter.adapterName] = adapter;
        }
        TLOGDEBUG(__FUNCTION__ << " " << app << "." << serverName << " succ" << endl);
        return 0;
    }
    catch (TC_Mysql_Exception &ex) {
        TLOGERROR(__FUNCTION__ << " " << app << "." << serverName << " exception: " << ex.what() << "|"
                               << sPrepareSql + sCondition << endl);
        return -1;
    }
    catch (TC_Config_Exception &ex) {
        TLOGERROR(__FUNCTION__ << " " << app << "." << serverName << " TC_Config_Exception exception: " << ex.what()
                               << endl);
        return -1;
    }
}

int RegistryImp::getTemplateContent(std::vector<std::string> &contents, std::string &lastFindTemplateName) {
    const std::string TOP_LEVEL_TEMPLATE = "tars.default";
    const std::string sPrepareSql =
            "SELECT  a.f_template_name as selfName, a.f_template_content as selfContent, b.f_template_name as parentsName, b.f_template_content as parentsContent "
            "FROM t_template a "
            "JOIN t_template b ON a.f_template_parent =b.f_template_name "
            "WHERE a.f_template_name=";
    std::string i = sPrepareSql + "'" + lastFindTemplateName + "'";
    try {
        tars::TC_Mysql::MysqlData res = _mysqlReg.queryRecord(sPrepareSql + "'" + lastFindTemplateName + "'");
        if (res.size() == 0) {
            return -1;
        }
        if (contents.empty()) {
            contents.push_back(res[0]["selfContent"]);
            contents.push_back(res[0]["parentsContent"]);
        } else {
            contents.push_back(res[0]["parentsContent"]);
        }

        const std::string &parentsName = res[0]["parentsName"];
        const std::string &selfName = res[0]["selfName"];

        lastFindTemplateName = parentsName;
        if (parentsName == TOP_LEVEL_TEMPLATE || selfName == parentsName) {
            return 0;
        }
        return 1;
    }
    catch (TC_Mysql_Exception &ex) {
        TLOGERROR(__FUNCTION__ << " " << lastFindTemplateName << "." << " exception: " << ex.what() << "|"
                               << sPrepareSql + "'" + lastFindTemplateName + "'" << endl);
        return -1;
    }
}

string
RegistryImp::getTemplateContent(const std::string &templateName, std::string &result, tars::CurrentPtr current) {
    assert(!templateName.empty());
    std::vector<std::string> contents;
    std::string findTemplateName = templateName;
    while (true) {
        int res = getTemplateContent(contents, findTemplateName);
        if (res == 0)break;
        if (res == 1)continue;
        result = templateName + " not found or parentsTemplate not found";
        return "";
    }
    assert(!contents.empty());
    try {
        TC_Config config;
        config.parseString(contents[0]);
        for (size_t i = 1; i != contents.size(); ++i) {
            TC_Config tempConfig;
            tempConfig.parseString(contents[i]);
            config.joinConfig(tempConfig, false);
        }
        result = "succ";
        TLOGDEBUG(__FUNCTION__ << " " << templateName << " " << result << endl);
        return config.tostr();
    }
    catch (TC_Config_Exception &ex) {
        result = "(" + templateName + ":" + ex.what() + ")";
        TLOGERROR(__FUNCTION__ << " TC_Config_Exception exception: " << ex.what() << endl);
        return "";
    }
}

void
RegistryImp::updateState(const std::string &podName, const std::string &settingState, const std::string &presentState,
                         tars::CurrentPtr current) {
    return K8SEndpointInterface::instance().onReceiveState(podName, settingState, presentState);
}



