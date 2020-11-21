#include "ConfigImp.h"
#include "servant/taf_logger.h"
#include "ConfigServer.h"
#include "ConfigInfoInterface.h"

using namespace taf;

void ConfigImp::initialize() {
}

int ConfigImp::ListConfig(const string &app, const string &server, vector <string> &vf, taf::JceCurrentPtr current) {
    LOG->debug() << "ListConfig|" << app << "." << server << "|" << endl;
    return ConfigInfoInterface::instance().listConfig(app, server, vf);
}

int ConfigImp::loadConfigByHost(const std::string &appServerName, const std::string &fileName, const string &host, string &config, taf::JceCurrentPtr current) {
    auto v = taf::TC_Common::sepstr<string>(appServerName, ".");

    if (v.size() == 1) {
        return ConfigInfoInterface::instance().loadConfig(v[0], "", fileName, host, config);
    }
    if (v.size() == 2) {
        return ConfigInfoInterface::instance().loadConfig(v[0], v[1], fileName, host, config);
    }

    return 0;
}

int ConfigImp::loadConfig(const std::string &app, const std::string &server, const std::string &fileName, string &config, taf::JceCurrentPtr current) {
    std::string sClientIP = current->getIp();
    return ConfigInfoInterface::instance().loadConfig(app, server, fileName, sClientIP, config);
}

int ConfigImp::checkConfig(const std::string &appServerName, const std::string &fileName, const string &host, string &result, taf::JceCurrentPtr current) {

    int ret = -1;

    auto v = taf::TC_Common::sepstr<string>(appServerName, ".");
    if (v.size() == 1) {
        ret = ConfigInfoInterface::instance().loadConfig(v[0], "", fileName, host, result);
    }
    if (v.size() == 2) {
        ret = ConfigInfoInterface::instance().loadConfig(v[0], v[1], fileName, host, result);
    }
    if (ret != 0) {
        result = "get config error";
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

int ConfigImp::ListConfigByInfo(const ConfigInfo &configInfo, vector <string> &vf, taf::JceCurrentPtr current) {
    LOG->debug() << "ListAllConfigByInfo|" << configInfo.appname << "|" << configInfo.servername << endl;
    if (configInfo.bAppOnly) {
        return ConfigInfoInterface::instance().listConfig(configInfo.appname, "", vf);
    }
    return ConfigInfoInterface::instance().listConfig(configInfo.appname, configInfo.servername, vf);
}

int ConfigImp::loadConfigByInfo(const ConfigInfo &configInfo, string &config, taf::JceCurrentPtr current) {
    LOG->debug() << "loadConfigByInfo|" << configInfo.appname << "|" << configInfo.servername << "|" << configInfo.filename << endl;
    if (configInfo.bAppOnly) {
        return ConfigInfoInterface::instance().loadConfig(configInfo.appname, "", configInfo.filename, configInfo.host, config);
    }
    return ConfigInfoInterface::instance().loadConfig(configInfo.appname, configInfo.servername, configInfo.filename, configInfo.host, config);
}

taf::Int32 ConfigImp::ListAllConfigByInfo(const taf::GetConfigListInfo &configInfo, vector <std::string> &vf, taf::JceCurrentPtr current) {
    LOG->debug() << "ListAllConfigByInfo|" << configInfo.appname << "|" << configInfo.servername << endl;
    if (configInfo.bAppOnly) {
        return ListConfig(configInfo.appname, "", vf, current);
    }
    return ListConfig(configInfo.appname, configInfo.servername, vf, current);
}

int ConfigImp::checkConfigByInfo(const ConfigInfo &configInfo, string &result, taf::JceCurrentPtr current) {
    int ret = ConfigInfoInterface::instance().loadConfig(configInfo.appname, configInfo.servername, configInfo.filename, configInfo.host, result);
    if (ret != 0) {
        result = "get config error";
        return -1;
    }
    try {
        TC_Config conf;
        conf.parseString(result);
    }
    catch (TC_Config_Exception &ex) {
        result = ex.what();
        return -1;
    }
    return 0;
}