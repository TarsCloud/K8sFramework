
#pragma once

#include <ctime>
#include <utility>
#include <vector>
#include <map>
#include <unordered_map>
#include <mutex>
#include <algorithm>
#include <cassert>
#include <iostream>
#include <util/tc_mysql.h>
#include <servant/EndpointF.h>
#include "NodeDescriptor.h"
#include "Registry.h"


#include "rapidjson/document.h"
#include "K8SListWatchSession.h"
#include "util/tc_config.h"

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
    bool isTars;
};

struct ExternalUpstream {
    string name;
    bool isTcp;
    std::vector<std::pair<std::string, int>> addresses; //ip,port
};

struct NormalPort {
    string name;
    int port;
    bool isTcp;
};

struct UpChain {
    std::unordered_map<string, std::vector<tars::EndpointF>> customUpChain;
    std::vector<tars::EndpointF> defaultUpChain;
};


struct TarsInfo {
    int asyncThread;
    std::string profileContent;
    std::string templateName;
    std::vector<std::shared_ptr<Adapter>> adapters;
};

struct DCacheInfo {
};

struct DCacheDBAccess {
};

struct DCacheProxyInfo {
};

struct DCacheRouteInfo {
};

struct NormalInfo {
    std::vector<std::shared_ptr<NormalPort>> ports;
};

struct ExternalInfo {
    std::vector<std::shared_ptr<ExternalUpstream>> upstreams;
};

enum class ServerSubType {
    Tars,
    DCache,
    DCacheDBAccess,
    DCacheProxy,
    DCacheRoute,
    External,
    Normal,
};

struct ServerInfo {
    ServerSubType subType{};
    std::string serverApp{};
    std::string serverName{};
    std::shared_ptr<TarsInfo> tarsInfo{};
    std::shared_ptr<DCacheInfo> dCacheInfo{};
    std::shared_ptr<DCacheProxyInfo> dCacheProxyInfo{};
    std::shared_ptr<DCacheRouteInfo> dCacheRouteInfo{};
    std::shared_ptr<ExternalInfo> externalInfo{};
    std::shared_ptr<NormalInfo> normalInfo{};
    std::vector<std::shared_ptr<PodStatus>> pods{};
};

struct Template {
    std::string content;
    std::string parent;
};

class ServerInfoInterface {
private:
    std::mutex _mutex;
    std::shared_ptr<UpChain> _upChainInfo;
    std::unordered_map<std::string, std::shared_ptr<ServerInfo>> _serverInfoMap;  //记录 ${ServerApp}-${ServerName} 与 ServerInfo 的对应关系
    std::unordered_map<std::string, std::shared_ptr<Template>> _templateMap;

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

    void findEndpoint(const string &id, vector<tars::EndpointF> *pActiveEp, vector<tars::EndpointF> *pInactiveEp);

    int getServerDescriptor(const string &serverApp, const string &serverName, tars::ServerDescriptor &descriptor);

    void loadUpChainConf();

private:
    TC_Config getTemplateContent(const std::string &sTemplateName, std::string &result);

    bool joinParentTemplate(const string &sTemplateName, TC_Config &conf, std::string &result);

    int getTarsServerDescriptor(const shared_ptr<ServerInfo> &serverInfo, ServerDescriptor &descriptor);

    void findTarsEndpoint(const std::shared_ptr<ServerInfo> &serverInfo, const string &sPortName, vector<tars::EndpointF> *pActiveEp, vector<tars::EndpointF> *pInactiveEp);

    void findUpChainEndpoint(const std::string &id, vector<tars::EndpointF> *pActiveEp, vector<tars::EndpointF> *pInactiveEp);
};

inline void handleEndpointsEvent(k8SEventTypeEnum eventType, const rapidjson::Value &pDocument) {

    assert(eventType == K8SEventTypeAdded || eventType == K8SEventTypeDeleted || eventType == K8SEventTypeUpdate);

    if (eventType == K8SEventTypeAdded) {
        return ServerInfoInterface::instance().onEndpointAdd(pDocument);
    }

    if (eventType == K8SEventTypeUpdate) {
        return ServerInfoInterface::instance().onEndpointUpdate(pDocument);
    }

    if (eventType == K8SEventTypeDeleted) {
        return ServerInfoInterface::instance().onEndpointDeleted(pDocument);
    }

    assert(false);
}

inline void handleTemplateEvent(k8SEventTypeEnum eventType, const rapidjson::Value &pDocument) {

    assert(eventType == K8SEventTypeAdded || eventType == K8SEventTypeDeleted || eventType == K8SEventTypeUpdate);

    if (eventType == K8SEventTypeAdded) {
        return ServerInfoInterface::instance().onTemplateAdd(pDocument);
    }

    if (eventType == K8SEventTypeUpdate) {
        return ServerInfoInterface::instance().onTemplateUpdate(pDocument);
    }

    if (eventType == K8SEventTypeDeleted) {
        return ServerInfoInterface::instance().onTemplateDeleted(pDocument);
    }

    assert(false);
}
