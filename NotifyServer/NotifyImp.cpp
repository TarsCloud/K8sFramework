#include "NotifyImp.h"
#include "NotifyMsgQueue.h"
#include "NotifyServer.h"

static std::string getNotifyLevel(const std::string &sNotifyMessage) {
    std::string Alarm = "[alarm]";
    std::string Error = "[error]";
    std::string Warn  = "[warn]";
    std::string Fail  = "[fail]";

    if (sNotifyMessage.find(Alarm) != std::string::npos) return Alarm;
    if (sNotifyMessage.find(Error) != std::string::npos) return Error;
    if (sNotifyMessage.find(Warn)  != std::string::npos) return Warn;
    if (sNotifyMessage.find(Fail)  != std::string::npos) return Error;

    return "[normal]";
}

void NotifyImp::loadConf() {
    TC_Config &g_pconf = g_app.getConfig();
    TC_DBConf tcDBConf;
    tcDBConf.loadFromMap(g_pconf.getDomainMap("/tars/db"));
    _mysqlConfig.init(tcDBConf);
}

void NotifyImp::initialize() {
    loadConf();
}

void NotifyImp::reportServer(const string &sServerName,
                             const string &sThreadId,
                             const string &sResult,
                             tars::CurrentPtr current) {
    std::string sPodId = current->getContext()["SERVER_HOST_NAME"];
    LOG->debug() << "reportServer|" << sServerName << "|" << sPodId << "|" << sThreadId << "|" << sResult << endl;
    DLOG << "reportServer|" << sServerName << "|" << sPodId << "|" << sThreadId << "|" << sResult << endl;

    string sql;
    TC_Mysql::RECORD_DATA rd;
    rd["f_app_server"] = make_pair(TC_Mysql::DB_STR, sServerName);
    rd["f_pod_name"] = make_pair(TC_Mysql::DB_STR, sPodId);
    rd["f_notify_message"] = make_pair(TC_Mysql::DB_STR, sResult);
    rd["f_notify_level"] = make_pair(TC_Mysql::DB_STR, getNotifyLevel(sResult));
    rd["f_notify_thread"] = make_pair(TC_Mysql::DB_STR, sThreadId);
    rd["f_notify_source"] = make_pair(TC_Mysql::DB_STR, "program");

    string sTable = "t_server_notify";

    try {
        _mysqlConfig.insertRecord(sTable, rd);
    } catch (TC_Mysql_Exception &ex) {
        TLOGERROR("insert2Db exception: " << ex.what() << endl);
    }
    catch (exception &ex) {
        TLOGERROR("insert2Db exception: " << ex.what() << endl);
    }
}

void NotifyImp::notifyServerEx(const string &sServerName, tars::NOTIFYLEVEL level, const string &sTitle,
                               const string &sMessage, tars::CurrentPtr current) {
    std::string sPodId = current->getContext()["SERVER_HOST_NAME"];
    NotifyMsgQueue::getInstance()->add(sServerName, level, sTitle, sMessage, sPodId);
}
