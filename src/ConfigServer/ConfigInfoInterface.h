#pragma once

#include <mutex>
#include <vector>
#include <memory>
#include <unordered_map>
#include <rapidjson/document.h>

class ConfigInfoInterface {
private:

    ConfigInfoInterface() = default;

    std::mutex mutex_;
    std::unordered_map<std::string, int> ipPodSeqMap_;

public:
    static ConfigInfoInterface &instance() {
        static ConfigInfoInterface infoInterface;
        return infoInterface;
    }

    void onPodAdd(const rapidjson::Value &pDocument);

    void onPodUpdate(const rapidjson::Value &pDocument);

    void onPodDelete(const rapidjson::Value &pDocument);

    int listConfig(const std::string &sSeverApp, const std::string &sSeverName, std::vector <std::string> &vector);

    int loadConfig(const std::string &sServerApp, const std::string &sSeverName, const std::string &sConfigName, const std::string &sClientIP, std::string &sConfigContent);

private:
    int loadAppConfig(const std::string &sServerApp, const std::string &sConfigName, std::string &sConfigContent);

    int loadServerConfig(const std::string &sServerApp, const std::string &sSeverName, const std::string &sConfigName, const std::string &sClientIP, std::string &sConfigContent);

    int listAppConfig(const std::string &sServerApp, std::vector <std::string> &vector);

    int listServerConfig(const std::string &sServerApp, const std::string &sServerName, std::vector <std::string> &vector);
};
