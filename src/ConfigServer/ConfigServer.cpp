#include "ConfigServer.h"
#include "ConfigImp.h"

void ConfigServer::initialize() {
    //滚动日志也打印毫秒
    LocalRollLogger::getInstance()->logger()->modFlag(TC_DayLogger::HAS_MTIME);
    //增加对象
    addServant<ConfigImp>(ServerConfig::Application + "." + ServerConfig::ServerName + ".ConfigObj");

    TLOGDEBUG("ConfigServer::initialize OK!" << endl);
}

void ConfigServer::destroyApp() {
}

