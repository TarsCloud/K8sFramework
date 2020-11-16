#include "NotifyMsgQueue.h"
#include "servant/Application.h"
#include "NotifyServer.h"
void NotifyMsgQueue::init() {
    TC_Config &g_pconf = g_app.getConfig();

    TC_DBConf tcDBConf;
    tcDBConf.loadFromMap(g_pconf.getDomainMap("/tars/db"));
    _mysql.init(tcDBConf);
    initLimit(&g_pconf);
    start();
}

void NotifyMsgQueue::terminate() {
    _terminate = true;

    TC_ThreadLock::Lock lock(*this);
    notifyAll();
}

void NotifyMsgQueue::add(const string &sServer, NOTIFYLEVEL level, const string &sTitle, const string &sMessage,
                         const string &sPodName) {
    MsgData data;
    data.sServer = sServer;
    data.level = level;
    data.sTitle = sTitle;
    data.sMessage = sMessage;
    data.sPodName = sPodName;
    _qMsg.push_back(data);
}

void NotifyMsgQueue::run() {
    while (!_terminate) {
        try {
            vector<MsgData> vQData;
            do {
                MsgData qData;
                _qMsg.pop_front(qData, -1);
                vQData.push_back(qData);

            } while ((!_qMsg.empty()) && (vQData.size() < 50));

            for (size_t i = 0; i < vQData.size(); i++) {
                writeDB(vQData[i]);
            }
        }
        catch (exception &ex) {
            TLOGERROR("exception:" << ex.what() << endl);
        }
        catch (...) {
            TLOGERROR("exception:e unknown error." << endl);
        }
    }
}

void NotifyMsgQueue::writeDB(const MsgData &data) {
    try {
        if (!checkLimit(data.sServer)) {
            TLOGERROR("limit fail|" << data.sServer << "|" << data.sPodName << "|" << data.level << "|" << data.sTitle
                                    << "|" << data.sMessage << endl);
            return;
        }
        string cmd = (data.level == NOTIFYERROR ? "[error]" : (data.level == NOTIFYWARN ? "[warn]" : "[normal]"));
        string sql;

        TC_Mysql::RECORD_DATA rd;
        rd["f_app_server"] = make_pair(TC_Mysql::DB_STR, data.sServer);
        rd["f_pod_name"] = make_pair(TC_Mysql::DB_STR, data.sPodName);
        rd["f_notify_thread"] = make_pair(TC_Mysql::DB_INT, "-1");
        rd["f_notify_level"] = make_pair(TC_Mysql::DB_STR, cmd);
        rd["f_notify_message"] = make_pair(TC_Mysql::DB_STR, data.sTitle + "|" + data.sMessage);
        rd["f_notify_time"] = make_pair(TC_Mysql::DB_INT, "now()");
        rd["f_notify_source"] = make_pair(TC_Mysql::DB_STR, "program");
        string sTable = "t_server_notify";

        _mysql.insertRecord(sTable, rd);
    }
    catch (TC_Mysql_Exception &ex) {
        TLOGERROR("insert2Db exception: " << ex.what() << endl);
    }
    catch (exception &ex) {
        TLOGERROR("insert2Db exception: " << ex.what() << endl);
    }
}

void FreqLimit::initLimit(TC_Config *conf) {
    string limitConf = conf->get("/tars/server<notify_limit>", "300:5");
    vector<int> vi = TC_Common::sepstr<int>(limitConf, ":,|");
    if (vi.size() != 2) {
        _interval = 300;
        _count = 5;
    } else {
        _interval = (unsigned int) vi[0];
        _count = vi[1];
        if (_count <= 1) {
            _count = 1;
        }
    }
}

bool FreqLimit::checkLimit(const string &sServer) {
    auto it = _limit.find(sServer);
    time_t t = TNOW;
    if (it != _limit.end()) {
        if (t > _limit[sServer].t + _interval) {
            _limit[sServer].t = t;
            _limit[sServer].count = 1;
            return true;
        } else if (_limit[sServer].count >= _count) {
            return false;
        } else {
            _limit[sServer].count++;
            return true;
        }
    } else {
        LimitData ld;
        ld.t = t;
        ld.count = 1;
        _limit[sServer] = ld;
        return true;
    }
}