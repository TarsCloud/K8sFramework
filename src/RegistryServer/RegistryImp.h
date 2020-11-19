
#pragma once

#include <util/tc_mysql.h>
#include "Registry.h"

using namespace tars;

/*
 * 提供给node调用的接口类
 */
class RegistryImp : public Registry {
public:
    /**
     * 构造函数
     */

    RegistryImp() = default;;

    /**
     * 初始化
     */

    void initialize() override;

    /**
     ** 退出
     */
    void destroy() override {};

    tars::Int32 getServerDescriptor(const std::string &serverApp, const std::string &serverName, tars::ServerDescriptor &serverDescriptor, tars::CurrentPtr current) override;

    void updateServerState(const std::string &podName, const std::string &settingState, const std::string &presentState, tars::CurrentPtr current) override;

};
