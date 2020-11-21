
#pragma once
#include <string>
#include "Registry.h"

using namespace taf;

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

    taf::Int32 getServerDescriptor(const std::string &serverApp, const std::string &serverName, taf::ServerDescriptor &serverDescriptor, taf::CurrentPtr current) override;

    void updateServerState(const std::string &podName, const std::string &settingState, const std::string &presentState, taf::CurrentPtr current) override;

};
