#include "ConfigImp.h"
#include "servant/taf_logger.h"
#include "ConfigServer.h"

using namespace taf;

//多配置文件的分割符
constexpr char MultiConfigSeparator[] = "\r\n\r\n";

static int extractPodSeqFromHost(const std::string &sApp, const std::string &sServer, const std::string &host) {
    std::vector<std::string> v = taf::TC_Common::sepstr<string>(host, "-");
    if (v.size() != 3) {
        return -1;
    }
    if (v[0] != taf::TC_Common::lower(sApp) || (v[1]) != taf::TC_Common::lower(sServer)) {
        return -1;
    }
    size_t pos = {};
    int podSeq = std::stoi(v[2], &pos, 10);
    return v[2].size() != pos ? -1 : podSeq;
}

static int extractPodSeqFromHost(const std::string &sAppServer, const std::string &host) {
    std::vector<std::string> v = taf::TC_Common::sepstr<string>(host, "-");
    if (v.size() != 3) {
        return -1;
    }
    std::string sAppServerFromHost = v[0] + "." + v[1];
    if (sAppServerFromHost != taf::TC_Common::lower(sAppServer)) {
        return -1;
    }
    size_t pos = {};
    int podSeq = std::stoi(v[2], &pos, 10);
    if (pos == v[2].size()) {
        return podSeq;
    }

    return -1;
}

void ConfigImp::initialize() {
    const TC_Config &serverConf = ConfigServer::getConfig();
    TC_DBConf tcDBConf;
    tcDBConf.loadFromMap(serverConf.getDomainMap("/taf/db"));
    LOG->debug() << "db conf:" << TC_Common::tostr(serverConf.getDomainMap("/taf/db")) << endl;
    _mysqlConfig.init(tcDBConf);
}

int ConfigImp::ListConfig(const string &app, const string &server, vector<string> &vf, taf::JceCurrentPtr current) {
    LOG->debug() << "ListConfig|" << app << "." << server << "|" << endl;

    std::string sAppServer = app + (server.empty() ? "" : "." + server);

    string sSql = "select distinct f_config_name from t_config where f_app_server='" + taf::TC_Mysql::escapeString(sAppServer) + "'";

    try {
        std::lock_guard<std::mutex> lockGuard(_mutex);
        TC_Mysql::MysqlData res = _mysqlConfig.queryRecord(sSql);
        LOG->debug() << "sql|" << sSql << "|" << res.size() << endl;
        for (size_t i = 0; i < res.size(); i++) {
            vf.push_back(res[i]["f_config_name"]);
        }
    }
    catch (TC_Mysql_Exception &ex) {
        LOG->error() << "exception: " << ex.what() << endl;
        return -1;
    }
    return 0;
}

int ConfigImp::loadConfigByHost(const std::string &appServerName, const std::string &fileName, const string &host, string &config, taf::JceCurrentPtr current) {
    int podSeq = extractPodSeqFromHost(appServerName, host);
    return loadConfigByPodSeq(appServerName, fileName, podSeq, config);
}

int ConfigImp::loadConfig(const std::string &app, const std::string &server, const std::string &fileName, string &config, taf::JceCurrentPtr current) {
    std::string sAppServer = app + (server.empty() ? "" : "." + server);
    std::string hostName = current->getContext()["SERVER_HOST_NAME"];
    int podSeq = extractPodSeqFromHost(app, server, hostName);
    return loadConfigByPodSeq(sAppServer, fileName, podSeq, config);
}

int ConfigImp::checkConfig(const std::string &appServerName, const std::string &fileName, const string &host, string &result, taf::JceCurrentPtr current) {

    int podSeq = extractPodSeqFromHost(appServerName, host);
    int ret = loadConfigByPodSeq(appServerName, fileName, podSeq, result);

    if (ret != 0) {
        return -1;
    }

    try {
        TC_Config conf;
        conf.parseString(result);
    }
    catch (exception &ex) {
        result = ex.what();
        return -1;
    }

    return 0;
}

int ConfigImp::ListConfigByInfo(const ConfigInfo &configInfo, vector<string> &vf, taf::JceCurrentPtr current) {
    LOG->debug() << "ListAllConfigByInfo|" << configInfo.appname << "|" << configInfo.servername << endl;
    if (configInfo.bAppOnly) {
        return ListConfig(configInfo.appname, "", vf, current);
    }
    return ListConfig(configInfo.appname, configInfo.servername, vf, current);
}

int ConfigImp::loadConfigByInfo(const ConfigInfo &configInfo, string &config, taf::JceCurrentPtr current) {
    LOG->debug() << "loadConfigByInfo|" << configInfo.appname << "|" << configInfo.servername << "|" << configInfo.filename << endl;
    int podSeq = -1;
    if (!configInfo.host.empty()) {
        podSeq = extractPodSeqFromHost(configInfo.appname, configInfo.servername, configInfo.host);
    }

    if (podSeq == -1) {
        std::string hostName = current->getHostName();
        podSeq = extractPodSeqFromHost(configInfo.appname, configInfo.servername, hostName);
    }

    std::string sAppServer;
    if (configInfo.bAppOnly || configInfo.servername.empty()) {
        sAppServer = configInfo.appname;
    } else {
        sAppServer = configInfo.appname + "." + configInfo.servername;
    }

    return loadConfigByPodSeq(sAppServer, configInfo.filename, podSeq, config);
}


taf::Int32 ConfigImp::ListAllConfigByInfo(const taf::GetConfigListInfo &configInfo, vector<std::string> &vf, taf::JceCurrentPtr current) {
    LOG->debug() << "ListAllConfigByInfo|" << configInfo.appname << "|" << configInfo.servername << endl;
    if (configInfo.bAppOnly) {
        return ListConfig(configInfo.appname, "", vf, current);
    }
    return ListConfig(configInfo.appname, configInfo.servername, vf, current);
}

int ConfigImp::checkConfigByInfo(const ConfigInfo &configInfo, string &result, taf::JceCurrentPtr current) {
    int podSeq = extractPodSeqFromHost(configInfo.appname, configInfo.servername, configInfo.host);
    std::string sAppServer = configInfo.appname + (configInfo.servername.empty() ? "" : "." + configInfo.servername);
    int ret = loadConfigByPodSeq(sAppServer, configInfo.filename, podSeq, result);

    if (ret != 0) {
        return -1;
    }

    try {
        TC_Config conf;
        conf.parseString(result);
    }
    catch (exception &ex) {
        result = ex.what();
        return -1;
    }

    return 0;
}

int ConfigImp::loadConfigByPodSeq(const string &appServerName, const string &filename, int podSeq, string &result) {
    std::ostringstream stream;
    stream << "select f_config_content from t_config where 1=1";
    stream << " and f_app_server=" << "'" << taf::TC_Mysql::escapeString(appServerName) << "'";
    stream << " and f_config_name=" << "'" << taf::TC_Mysql::escapeString(filename) << "'";
    stream << " and (f_pod_seq=-1 or f_pod_seq=" << podSeq << ")";
    stream << " order by f_pod_seq";

    const std::string sSql = stream.str();
    stream.str("");
    try {
        std::lock_guard<std::mutex> lockGuard(_mutex);
        TC_Mysql::MysqlData res = _mysqlConfig.queryRecord(sSql);
        LOG->debug() << "sql|" << sSql << "|" << res.size() << endl;
        for (size_t i = 0; i < res.size(); ++i) {
            if (i > 0) {
                stream << MultiConfigSeparator;
            }
            stream << res[i]["f_config_content"];
        }
    }
    catch (TC_Mysql_Exception &ex) {
        result = ex.what();
        LOG->error() << "exception: " << ex.what() << endl;
        return -1;
    }
    result = stream.str();
    return 0;
}
