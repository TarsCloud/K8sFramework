
#pragma once

#include <mutex>
#include <unordered_map>
#include <algorithm>
#include <cassert>
#include <servant/EndpointF.h>
#include "K8SWatcher.h"
#include "util/tc_config.h"
#include "rapidjson/document.h"
#include "Registry.h"

struct PodStatus {
    std::string name;
    std::string podIP;
    std::string hostIP;
    std::string presentState;
};

struct Adapter {
    int port;
    int hostPort;
    std::string name;
    uint thread;
    uint connection;
    uint timeout;
    uint capacity;
    bool isTcp;
    bool isTaf;
};

struct UpChain {
    std::unordered_map<string, std::vector<taf::EndpointF>> customUpChain;
    std::vector<taf::EndpointF> defaultUpChain;
};

struct TafInfo {
    int asyncThread;
    std::string profileContent;
    std::string templateName;
    std::vector<std::shared_ptr<Adapter>> adapters;
};

enum class ServerSubType {
    Taf,
    Normal,
};

struct ServerInfo {
    ServerSubType subType{};
    std::string serverApp{};
    std::string serverName{};
    std::shared_ptr<TafInfo> tafInfo{};
    std::vector<std::shared_ptr<PodStatus>> pods{};
};

struct Template {
    std::string content_;
    std::string parent_;
};

class ServerInfoInterface {
private:
    std::mutex mutex_;
    std::shared_ptr<UpChain> upChainInfo_;
    std::unordered_map<std::string, std::shared_ptr<ServerInfo>> serverInfoMap_;  //记录 ${ServerApp}-${ServerName} 与 ServerInfo 的对应关系
    std::unordered_map<std::string, std::shared_ptr<Template>> templateMap_;

public:

    static ServerInfoInterface &instance() {
        static ServerInfoInterface endpointInterface;
        return endpointInterface;
    };

    void onEndpointAdd(const rapidjson::Value &pDocument);

    void onEndpointUpdate(const rapidjson::Value &pDocument);

    void onEndpointDeleted(const rapidjson::Value &pDocument);

    void onTemplateAdd(const rapidjson::Value &pDocument);

    void onTemplateUpdate(const rapidjson::Value &pDocument);

    void onTemplateDeleted(const rapidjson::Value &pDocument);

    void findEndpoint(const string &id, vector<taf::EndpointF> *pActiveEp, vector<taf::EndpointF> *pInactiveEp);

    int getServerDescriptor(const string &serverApp, const string &serverName, taf::ServerDescriptor &descriptor);

    void loadUpChainConf();

private:
    taf::TC_Config getTemplateContent(const std::string &sTemplateName, std::string &result);

    bool joinParentTemplate(const string &sTemplateName, taf::TC_Config &conf, std::string &result);

    int getTafServerDescriptor(const shared_ptr<ServerInfo> &serverInfo, taf::ServerDescriptor &descriptor);

    void findTafEndpoint(const std::shared_ptr<ServerInfo> &serverInfo, const string &sPortName, vector<taf::EndpointF> *pActiveEp, vector<taf::EndpointF> *pInactiveEp);

    void findUpChainEndpoint(const std::string &id, vector<taf::EndpointF> *pActiveEp, vector<taf::EndpointF> *pInactiveEp);
};

inline void handleEndpointsEvent(K8SWatchEvent eventType, const rapidjson::Value &pDocument) {

    assert(eventType == K8SWatchEventAdded || eventType == K8SWatchEventDeleted || eventType == K8SWatchEventUpdate);

    if (eventType == K8SWatchEventAdded) {
        return ServerInfoInterface::instance().onEndpointAdd(pDocument);
    }

    if (eventType == K8SWatchEventUpdate) {
        return ServerInfoInterface::instance().onEndpointUpdate(pDocument);
    }

    if (eventType == K8SWatchEventDeleted) {
        return ServerInfoInterface::instance().onEndpointDeleted(pDocument);
    }

    assert(false);
}

inline void handleTemplateEvent(K8SWatchEvent eventType, const rapidjson::Value &pDocument) {

    assert(eventType == K8SWatchEventAdded || eventType == K8SWatchEventDeleted || eventType == K8SWatchEventUpdate);

    if (eventType == K8SWatchEventAdded) {
        return ServerInfoInterface::instance().onTemplateAdd(pDocument);
    }

    if (eventType == K8SWatchEventUpdate) {
        return ServerInfoInterface::instance().onTemplateUpdate(pDocument);
    }

    if (eventType == K8SWatchEventDeleted) {
        return ServerInfoInterface::instance().onTemplateDeleted(pDocument);
    }

    assert(false);
}
