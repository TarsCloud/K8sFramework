
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

struct K8SPod {
    K8SPod(std::string podName, std::string podIp) : podName(std::move(podName)), podIp(std::move(podIp)) {}

    std::string podName;
    std::string podIp;
};

struct K8SEndpoint {
    std::vector<std::tuple<std::string, uint16_t, bool>> ports{};  //portName, portValue ,isTcp;
    std::vector<std::shared_ptr<K8SPod>> addresses{};
    std::vector<std::shared_ptr<K8SPod>> notReadyAddresses{};
};

class K8SEndpointInterface {
private:
    std::ostringstream strStream;
    std::mutex mutex;
    tars::TC_Mysql mysql;
    std::unordered_map<std::string, std::shared_ptr<K8SEndpoint>> serverEndpointMap;  //记录 ${ServerApp}-${ServerName} 与 K8SEndpoint 的对应关系

public:

    static K8SEndpointInterface &instance() {
        static K8SEndpointInterface endpointInterface;
        return endpointInterface;
    };

    void init(const tars::TC_Config &conf);

    void onPodAdd(const rapidjson::Value &pDocument);

    void onPodUpdate(const rapidjson::Value &pDocument);

    void onPodDeleted(const rapidjson::Value &pDocument);


    void onEndpointAdd(const rapidjson::Value &pDocument);

    void onEndpointUpdate(const rapidjson::Value &pDocument);

    void onEndpointDeleted(const rapidjson::Value &pDocument);


    void onEventAdd(const rapidjson::Value &pDocument);

    void onEventUpdate(const rapidjson::Value &pDocument);

    void onEventDeleted(const rapidjson::Value &pDocument);


    void onReceiveState(const std::string &sPodName, const std::string &settingState, const std::string &presentState);


    enum class FindEndpointRes {
        Success,
        NoServer //server不存在,便于 区分没查询到 endpoint 是因为没在 k8s集群部署该服务,还是因为服务的 endpoint 数本身就是 0个
    };

    FindEndpointRes findEndpoint(const string &sAppServer, const string &sPortName, vector<EndpointF> *pActiveEp,
                                 vector<tars::EndpointF> *pInactiveEp);

    FindEndpointRes findEndpoint(const string &id, vector<EndpointF> *pActiveEp, vector<tars::EndpointF> *pInactiveEp);

    std::string externalJceProxy() { return _externalJceProxy; }
 private:
    bool checkTafIOEndpoint(const rapidjson::Value &pDocument, std::string& sK8SEndpointName);
    std::shared_ptr<K8SEndpoint> parseEndpoint(const rapidjson::Value &pDocument);

 private:
    std::string _externalJceProxy;
};

inline void handleK8SPodsEvent(k8SEventTypeEnum eventType, const rapidjson::Value &pDocument) {

    assert(eventType == K8SEventTypeAdded || eventType == K8SEventTypeDeleted || eventType == K8SEventTypeModified);

    if (eventType == K8SEventTypeAdded) {
        return K8SEndpointInterface::instance().onPodAdd(pDocument);
    }

    if (eventType == K8SEventTypeModified) {
        return K8SEndpointInterface::instance().onPodUpdate(pDocument);
    }

    if (eventType == K8SEventTypeDeleted) {
        return K8SEndpointInterface::instance().onPodDeleted(pDocument);
    }

    assert(false);
    //should not reach here
}

inline void handleK8SEventsEvent(k8SEventTypeEnum eventType, const rapidjson::Value &pDocument) {

    assert(eventType == K8SEventTypeAdded || eventType == K8SEventTypeDeleted || eventType == K8SEventTypeModified);

    if (eventType == K8SEventTypeAdded) {
        return K8SEndpointInterface::instance().onEventAdd(pDocument);
    }

    if (eventType == K8SEventTypeModified) {
        return K8SEndpointInterface::instance().onEventUpdate(pDocument);
    }

    if (eventType == K8SEventTypeDeleted) {
        return K8SEndpointInterface::instance().onEventDeleted(pDocument);
    }

    assert(false);
    //should not reach here
}

inline void handleK8SEndpointsEvent(k8SEventTypeEnum eventType, const rapidjson::Value &pDocument) {

    assert(eventType == K8SEventTypeAdded || eventType == K8SEventTypeDeleted || eventType == K8SEventTypeModified);

    if (eventType == K8SEventTypeAdded) {
        return K8SEndpointInterface::instance().onEndpointAdd(pDocument);
    }

    if (eventType == K8SEventTypeModified) {
        return K8SEndpointInterface::instance().onEndpointUpdate(pDocument);
    }

    if (eventType == K8SEventTypeDeleted) {
        return K8SEndpointInterface::instance().onEndpointDeleted(pDocument);
    }

    assert(false);
}