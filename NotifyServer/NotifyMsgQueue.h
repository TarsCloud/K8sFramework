#ifndef _NOTIFY_MSG_QUEUE_H_
#define _NOTIFY_MSG_QUEUE_H_

#include "servant/NotifyF.h"
#include "util/tc_common.h"
#include "util/tc_config.h"
#include "util/tc_monitor.h"
#include "util/tc_mysql.h"
#include "util/tc_singleton.h"
#include "util/tc_thread.h"
#include "util/tc_thread_queue.h"

using namespace tars;

struct MsgData {
    string sServer;
    std::string sPodName;
    string sTitle;
    string sMessage;
    NOTIFYLEVEL level;
};

class FreqLimit {
public:
    struct LimitData {
        unsigned int t;
        int count;
    };

    void initLimit(TC_Config *conf);

    // return true 表示检测通过，没有被频率限制，  fasle: 被频率限制了
    bool checkLimit(const string &sServer);

protected:
    unordered_map<string, LimitData> _limit;

    unsigned int _interval;
    int _count;
};

class NotifyMsgQueue : public TC_Singleton<NotifyMsgQueue>,
                       public TC_ThreadLock,
                       public TC_Thread,
                       public FreqLimit {
public:
    NotifyMsgQueue() = default;

    ~NotifyMsgQueue() = default;

    void init();

    void
    add(const string &sServer, NOTIFYLEVEL level, const string &sTitle, const string &sMessage, const string &sPodName);

    /**
     * stop
     */
    void terminate();

protected:
    virtual void run();

    void writeDB(const MsgData &data);

protected:
    bool _terminate;
    TC_ThreadQueue<MsgData> _qMsg;

    TC_Mysql _mysql;
};

#endif // _NOTIFY_MSG_QUEUE_H_
