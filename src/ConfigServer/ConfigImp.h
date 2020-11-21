#ifndef _CONFIG_IMP_H_
#define _CONFIG_IMP_H_

#include "util/tc_common.h"
#include "util/tc_config.h"
#include "util/tc_mysql.h"
#include "servant/taf_logger.h"
#include "servant/ConfigF.h"

using namespace taf;

class ConfigImp final : public Config {
public:
    /**
     *
     */
    ConfigImp() = default;;

    /**
     *
     */
    ~ConfigImp() final = default;;

    /**
     * 初始化
     *
     * @return int
     */
    void initialize() final;

    /**
     * 退出
     */
    void destroy() final {};

    /**
    * 获取配置文件列表
    * param app :应用
    * param server:  server名
    * param vf: 配置文件名
    *
    * return  : 配置文件内容
    */
    int ListConfig(const string &app, const string &server, vector <string> &vf, taf::JceCurrentPtr current) final;

    /**
     * 加载配置文件
     * param app :应用
     * param server:  server名
     * param filename:  配置文件名
     *
     * return  : 配置文件内容
     */
    int loadConfig(const std::string &app, const std::string &server, const std::string &filename, string &config,
                   taf::JceCurrentPtr current) final;

    /**
     * 根据ip获取配置
     * @param appServerName
     * @param filename
     * @param host
     * @param config
     *
     * @return int
     */
    int loadConfigByHost(const string &appServerName, const string &filename, const string &host, string &config,
                         taf::JceCurrentPtr current) final;

    /**
     *
     * @param appServerName
     * @param filename
     * @param host
     * @param current
     *
     * @return int
     */
    int checkConfig(const string &appServerName, const string &filename, const string &host, string &result,
                    taf::JceCurrentPtr current) final;

    /**
    * 获取配置文件列表
    * param configInfo ConfigInfo
    * param vf: 配置文件名
    *
    * return  : 配置文件内容
    */
    int ListConfigByInfo(const ConfigInfo &configInfo, vector <string> &vf, taf::JceCurrentPtr current) final;

    /**
     * 加载配置文件
     * param configInfo ConfigInfo
     * param config:  配置文件内容
     *
     * return  :
     */

    int loadConfigByInfo(const ConfigInfo &configInfo, string &config, taf::JceCurrentPtr current) final;

    /**
     *
     * @param configInfo ConfigInfo
     *
     * @return int
     */

    int checkConfigByInfo(const ConfigInfo &configInfo, string &result, taf::JceCurrentPtr current) final;

    /**
    * 获取服务的所有配置文件列表，
    * @param configInfo 支持拉取应用配置列表，服务配置列表，机器配置列表和容器配置列表
    * @param[out] vf  获取到的文件名称列表
    * @return int 0: 成功, -1:失败
    **/
    taf::Int32 ListAllConfigByInfo(const taf::GetConfigListInfo &configInfo, vector <std::string> &vf,
                                   taf::JceCurrentPtr current) final;
};

#endif

