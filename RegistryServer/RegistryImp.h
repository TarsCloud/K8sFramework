
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

    int getServers(const std::string &app, const std::string &serverName, tars::ServerDescriptor &serverDesc,
                   tars::CurrentPtr current) override;

    string
    getTemplateContent(const std::string &templateName, std::string &result, tars::CurrentPtr current) override;

    void updateState(const std::string &podName, const std::string &settingState, const std::string &presentState,
                     tars::CurrentPtr current) override;

private:
    tars::TC_Mysql _mysqlReg;

    int getTemplateContent(vector<std::string> &contents, string &lastFindTemplateName);
};
