﻿
#pragma once

#include "Node.h"
#include <servant/NodeF.h>
#include "util/tc_common.h"

using namespace tars;
using namespace std;

class ServerImp : public ServerF {
public:
    /**
     *
     */
    ServerImp() = default;;

    /**
     * 销毁服务
     * @param k
     * @param v
     *
     * @return int
     */
    ~ServerImp() override = default;

    /**
    * 初始化
    */
    void initialize() override {
    };

    /**
    * 退出
    */
    void destroy() override {
    };

    /**
    * 退出
    */

    int keepActiving(const tars::ServerInfo &serverInfo, tars::CurrentPtr current) override;

    int keepAlive(const tars::ServerInfo &serverInfo, tars::CurrentPtr current) override;

    int reportVersion(const string &app, const string &serverName, const string &version, tars::CurrentPtr current) override;

   /**
    * 获取最近keepalive的时间戳
    * @return 最后一次keepalive的时间戳
    */
    unsigned int getLatestKeepAliveTime(tars::CurrentPtr current);

};


