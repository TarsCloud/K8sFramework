
#include "ConfigInfoInterface.h"
#include "K8SClient.h"
#include "K8SParams.h"
#include <rapidjson/pointer.h>

//多配置文件的分割符
constexpr char MultiConfigSeparator[] = "\r\n\r\n";

static int extractPodSeq(const std::string &sPodName, const std::string &sGenerateName) {
    try {
        auto sPodSeq = sPodName.substr(sGenerateName.size());
        return std::stoi(sPodSeq, nullptr, 10);
    } catch (std::exception &exception) {
        return -1;
    }
}

static inline std::string SFromP(const rapidjson::Value *p) {
    assert(p != nullptr);
    return {p->GetString(), p->GetStringLength()};
}

void ConfigInfoInterface::onPodAdd(const rapidjson::Value &pDocument) {
    auto pGenerateName = rapidjson::GetValueByPointer(pDocument, "/metadata/generateName");
    if (pGenerateName == nullptr) { return; }
    assert(pGenerateName->IsString());
    std::string sGenerateName = SFromP(pGenerateName);

    auto pPodIP = rapidjson::GetValueByPointer(pDocument, "/status/podIP");
    if (pPodIP == nullptr) { return; }
    assert(pPodIP->IsString());
    std::string podIP = SFromP(pPodIP);

    auto pPodName = rapidjson::GetValueByPointer(pDocument, "/metadata/name");
    assert(pPodName != nullptr && pPodName->IsString());
    std::string sPodName = SFromP(pPodName);

    int iPodSeq = extractPodSeq(sPodName, sGenerateName);

    std::lock_guard <std::mutex> lockGuard(mutex_);
    ipPodSeqMap_[podIP] = iPodSeq;
}

void ConfigInfoInterface::onPodUpdate(const rapidjson::Value &pDocument) {
    return onPodAdd(pDocument);
}

void ConfigInfoInterface::onPodDelete(const rapidjson::Value &pDocument) {
    auto pPodIP = rapidjson::GetValueByPointer(pDocument, "/status/podIP");
    if (pPodIP == nullptr) { return; }
    assert(pPodIP->IsString());
    std::string sPodIP = SFromP(pPodIP);

    std::lock_guard <std::mutex> lockGuard(mutex_);
    ipPodSeqMap_.erase(sPodIP);
}

int ConfigInfoInterface::listConfig(const std::string &sServerApp, const std::string &sSeverName, std::vector <std::string> &vector) {
    if (sSeverName.empty()) {
        return listAppConfig(sServerApp, vector);
    }
    return listServerConfig(sServerApp, sSeverName, vector);
}

int ConfigInfoInterface::loadConfig(const std::string &sSeverApp, const std::string &sSeverName, const std::string &sConfigName, const std::string &sClientIP,
                                    std::string &sConfigContent) {
    if (sSeverName.empty()) {
        return loadAppConfig(sSeverApp, sConfigName, sConfigContent);
    }
    return loadServerConfig(sSeverApp, sSeverName, sConfigName, sClientIP, sConfigContent);
}

int ConfigInfoInterface::loadAppConfig(const std::string &sServerApp, const std::string &sConfigName, std::string &sConfigContent) {
    sConfigContent.clear();
    std::ostringstream stream;
    stream << "/apis/k8s.taf.io/v1alpha1/namespaces/" << K8SParams::instance().bindNamespace() << "/tconfigs?labelSelector=taf.io/ServerApp=" << sServerApp
           << ",taf.io/ServerName=" << "";
    auto k8sClientRequest = K8SClient::instance().postRequest(K8SClientRequestMethod::Get, stream.str(), "");
    bool bTaskFinish = k8sClientRequest->waitFinish(std::chrono::seconds(1));
    if (!bTaskFinish) {
        return -1;
    }

    if (k8sClientRequest->state() != Done) {
        //error
        return -1;
    }

    const auto &responseJson = k8sClientRequest->responseJson();
    auto pItem = rapidjson::GetValueByPointer(responseJson, "/items");
    if (pItem == nullptr) {
        // no value;
        return 0;
    }

    assert(pItem->IsArray());
    for (auto &&config:pItem->GetArray()) {
        auto pConfigContent = rapidjson::GetValueByPointer(config, "/appConfig/configContent");
        assert(pConfigContent != nullptr && pConfigContent->IsString());
        sConfigContent = SFromP(pConfigContent);
    }
    return 0;
}

int ConfigInfoInterface::loadServerConfig(const std::string &sServerApp, const std::string &sSeverName, const std::string &sConfigName, const std::string &sClientIP,
                                          std::string &sConfigContent) {
    int podSeq = -1;
    {
        std::lock_guard <std::mutex> lockGuard(mutex_);
        if (!sClientIP.empty()) {
            auto podSeqIterator = ipPodSeqMap_.find(sClientIP);
            if (podSeqIterator != ipPodSeqMap_.end()) {
                podSeq = podSeqIterator->second;
            }
        }
    }
    std::ostringstream stream;
    stream << "/apis/k8s.taf.io/v1alpha1/namespaces/" << K8SParams::instance().bindNamespace() << "/tconfigs?labelSelector=taf.io/ServerApp=" << sServerApp
           << ",taf.io/ServerName=" << sSeverName << ",taf.io/ConfigName=" << sConfigName << ",taf.io/PodSeq+in+(m";
    if (podSeq != -1) {
        stream << "," << podSeq;
    }
    stream << ")";
    auto k8sClientRequest = K8SClient::instance().postRequest(K8SClientRequestMethod::Get, stream.str(), "");
    bool bTaskFinish = k8sClientRequest->waitFinish(std::chrono::seconds(1));
    if (!bTaskFinish) {
        return -1;
    }

    if (k8sClientRequest->state() != Done) {
        //todo log;
        return -1;
    }

    const auto &responseJson = k8sClientRequest->responseJson();

    auto pItem = rapidjson::GetValueByPointer(responseJson, "/items");
    if (pItem == nullptr) {
        sConfigContent = "";
        return 0;
    }

    assert(pItem->IsArray());

    const char *masterConfigContent{};
    size_t masterConfigContentLength{};

    const char *podConfigContent{};
    size_t podConfigContentLength{};

    auto &&jsonArray = pItem->GetArray();
    for (auto &&item : jsonArray) {
        auto pConfigContent = rapidjson::GetValueByPointer(item, "/serverConfig/configContent");
        assert(pConfigContent != nullptr && pConfigContent->IsString());

        auto pPodSeq = rapidjson::GetValueByPointer(item, "/serverConfig/podSeq");
        if (pPodSeq == nullptr) {
            masterConfigContent = pConfigContent->GetString();
            masterConfigContentLength = pConfigContent->GetStringLength();
            continue;
        }

        podConfigContent = pConfigContent->GetString();
        podConfigContentLength = pConfigContent->GetStringLength();
    }

    assert(masterConfigContent != nullptr);
    stream.str("");

    stream << std::string(masterConfigContent, masterConfigContentLength);
    if (podConfigContent != nullptr && podConfigContentLength != 0) {
        stream << MultiConfigSeparator << std::string(podConfigContent, podConfigContentLength);
    }
    sConfigContent = stream.str();
    return 0;
}

int ConfigInfoInterface::listAppConfig(const std::string &sServerApp, std::vector <std::string> &vector) {

    assert(!sServerApp.empty());

    vector.clear();
    std::ostringstream stream;
    stream << "/apis/k8s.taf.io/v1alpha1/namespaces/" << K8SParams::instance().bindNamespace() << "/tconfigs?labelSelector=taf.io/ServerApp=" << sServerApp
           << ",taf.io/ServerName=" << "";
    auto k8sClientRequest = K8SClient::instance().postRequest(K8SClientRequestMethod::Get, stream.str(), "");
    bool bTaskFinish = k8sClientRequest->waitFinish(std::chrono::seconds(1));
    if (!bTaskFinish) {
        return -1;
    }

    if (k8sClientRequest->state() != Done) {
        // todo log;
        return -1;
    }

    const auto &responseJson = k8sClientRequest->responseJson();
    auto pItem = rapidjson::GetValueByPointer(responseJson, "/items");
    if (pItem == nullptr) {// no value;
        return 0;
    }

    assert(pItem->IsArray());
    for (auto &&config:pItem->GetArray()) {
        auto pConfigName = rapidjson::GetValueByPointer(config, "/appConfig/configName");
        assert(pConfigName != nullptr && pConfigName->IsString());
        vector.emplace_back(SFromP(pConfigName));
    }
    return 0;
}

int ConfigInfoInterface::listServerConfig(const std::string &sServerApp, const std::string &sServerName, std::vector <std::string> &vector) {

    assert(!sServerApp.empty());
    assert(!sServerName.empty());

    vector.clear();
    std::ostringstream stream;
    stream << "/apis/k8s.taf.io/v1alpha1/namespaces/" << K8SParams::instance().bindNamespace() << "/tconfigs?labelSelector=taf.io/ServerApp=" << sServerApp
           << ",taf.io/ServerName=" << sServerName;
    auto k8sClientRequest = K8SClient::instance().postRequest(K8SClientRequestMethod::Get, stream.str(), "");
    bool bTaskFinish = k8sClientRequest->waitFinish(std::chrono::seconds(1));
    if (!bTaskFinish) {
        return -1;
    }

    if (k8sClientRequest->state() != Done) {
        //todo log;
        return -1;
    }

    const auto &responseJson = k8sClientRequest->responseJson();
    auto pItem = rapidjson::GetValueByPointer(responseJson, "/items");
    if (pItem == nullptr) { // no value;
        return 0;
    }

    assert(pItem->IsArray());
    for (auto &&config:pItem->GetArray()) {
        auto pConfigName = rapidjson::GetValueByPointer(config, "/serverConfig/configName");
        assert(pConfigName != nullptr && pConfigName->IsString());
        vector.emplace_back(SFromP(pConfigName));
    }
    return 0;
}