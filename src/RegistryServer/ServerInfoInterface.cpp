
#include <thread>
#include "ServerInfoInterface.h"
#include "RegistryServer.h"


static inline std::string SFromP(const rapidjson::Value *p) {
    assert(p != nullptr);
    return {p->GetString(), p->GetStringLength()};
}

void ServerInfoInterface::findEndpoint(const string &id, vector<EndpointF> *pActiveEp, vector<taf::EndpointF> *pInactiveEp) {
    std::vector<std::string> v = taf::TC_Common::sepstr<string>(id, ".");
    if (v.size() != 3) {
        return;
    }

    const auto sAppServer = TC_Common::lower(v[0]) + "-" + TC_Common::lower(v[1]);
    const auto &sPortName = v[2];

    assert(pActiveEp != nullptr);

    pActiveEp->clear();
    if (pInactiveEp != nullptr) {
        pInactiveEp->clear();
    }

    std::lock_guard<std::mutex> lockGuard(_mutex);

    auto iterator = _serverInfoMap.find(sAppServer);

    if (iterator == _serverInfoMap.end()) {
        return findUpChainEndpoint(id, pActiveEp, pInactiveEp);
    }

    const auto &serverInfo = iterator->second;

    if (serverInfo == nullptr) {
        LOG->debug() << iterator->first << "->serverInfo is nullptr" << endl;
        return;
    }

    switch (serverInfo->subType) {
        case ServerSubType::Taf:
            return findTafEndpoint(serverInfo, sPortName, pActiveEp, pInactiveEp);
        case ServerSubType::DCache:
            break;
        case ServerSubType::DCacheProxy:
            break;
        case ServerSubType::DCacheRoute:
            break;
        case ServerSubType::DCacheDBAccess:
            break;
        case ServerSubType::External:
            return;
        case ServerSubType::Normal:
            break;
    }
}

static std::shared_ptr<TafInfo> buildTafInfoFromDocument(const rapidjson::Value &pDocument) {

    auto pTafInfo = std::make_shared<TafInfo>();

    auto pAsyncThread = rapidjson::GetValueByPointer(pDocument, "/spec/taf/asyncThread");
    if (pAsyncThread == nullptr) {
        //fixme  should log
        return nullptr;
    }
    assert(pAsyncThread->IsInt());
    pTafInfo->asyncThread = pAsyncThread->GetInt();

    auto pProfile = rapidjson::GetValueByPointer(pDocument, "/spec/taf/profile");
    if (pProfile == nullptr) {
        //fixme  should log
        return nullptr;
    }
    assert(pProfile->IsString());
    pTafInfo->profileContent = SFromP(pProfile);

    auto pTemplate = rapidjson::GetValueByPointer(pDocument, "/spec/taf/template");
    if (pTemplate == nullptr) {
        //fixme  should log
        return nullptr;
    }
    assert(pTemplate->IsString());
    pTafInfo->templateName = SFromP(pTemplate);

    auto pServants = rapidjson::GetValueByPointer(pDocument, "/spec/taf/servants");
    if (pServants == nullptr) {
        //fixme  should log
        return nullptr;
    }

    assert(pServants->IsArray());
    for (const auto &v :pServants->GetArray()) {
        auto pAdapter = std::make_shared<Adapter>();
        auto pName = rapidjson::GetValueByPointer(v, "/name");
        assert(pName != nullptr && pName->IsString());
        pAdapter->name = SFromP(pName);

        auto pPort = rapidjson::GetValueByPointer(v, "/port");
        assert(pPort != nullptr && pPort->IsInt());
        pAdapter->port = pPort->GetInt();

        auto pThread = rapidjson::GetValueByPointer(v, "/thread");
        assert(pThread != nullptr && pThread->IsInt());
        pAdapter->thread = pThread->GetInt();

        auto pConnection = rapidjson::GetValueByPointer(v, "/connection");
        assert(pConnection != nullptr && pConnection->IsInt());
        pAdapter->connection = pConnection->GetInt();

        auto pTimeout = rapidjson::GetValueByPointer(v, "/timeout");
        assert(pTimeout != nullptr && pTimeout->IsInt());
        pAdapter->timeout = pTimeout->GetInt();

        auto pCapacity = rapidjson::GetValueByPointer(v, "/capacity");
        assert(pCapacity != nullptr && pCapacity->IsInt());
        pAdapter->capacity = pCapacity->GetInt();

        auto pIsTaf = rapidjson::GetValueByPointer(v, "/isTaf");
        assert(pIsTaf != nullptr && pIsTaf->IsBool());
        pAdapter->isTaf = pIsTaf->GetBool();

        auto pIsTCP = rapidjson::GetValueByPointer(v, "/isTcp");
        assert(pIsTCP != nullptr && pIsTCP->IsBool());
        pAdapter->isTcp = pIsTCP->GetBool();

        pTafInfo->adapters.push_back(pAdapter);
    }

    auto pHostPorts = rapidjson::GetValueByPointer(pDocument, "/spec/hostPorts");

    if (pHostPorts != nullptr) {

        assert(pHostPorts->IsArray());

        for (const auto &hostPort :pHostPorts->GetArray()) {

            auto pNameRef = rapidjson::GetValueByPointer(hostPort, "/nameRef");
            assert(pNameRef != nullptr && pNameRef->IsString());
            auto nameRef = SFromP(pNameRef);

            for (auto &adapter:pTafInfo->adapters) {

                if (adapter->name == nameRef) {
                    auto pPort = rapidjson::GetValueByPointer(hostPort, "/port");
                    assert(pPort != nullptr && pPort->IsInt());
                    adapter->hostPort = pPort->GetInt();
                    break;
                }
            }
        }
    }

    return pTafInfo;
}

static std::shared_ptr<ExternalInfo> buildExternalInfoFromDocument(const rapidjson::Value &pDocument) {

    auto pExternalInfo = std::make_shared<ExternalInfo>();

    auto pUpstreams = rapidjson::GetValueByPointer(pDocument, "/spec/external/upstreams");
    if (pUpstreams == nullptr) {
        return nullptr;
    }

    assert(pUpstreams->IsArray());

    for (const auto &upstream: pUpstreams->GetArray()) {

        auto upstreamInfo = std::make_shared<ExternalUpstream>();

        auto pName = rapidjson::GetValueByPointer(upstream, "/name");
        assert(pName != nullptr && pName->IsString());
        upstreamInfo->name = SFromP(pName);

        auto pIsTcp = rapidjson::GetValueByPointer(upstream, "/isTcp");
        assert(pIsTcp != nullptr && pIsTcp->IsBool());
        upstreamInfo->isTcp = pIsTcp->GetBool();

        auto pAddresses = rapidjson::GetValueByPointer(upstream, "/addresses");
        if (pAddresses != nullptr) {
            assert(pAddresses->IsArray());

            for (const auto &address: pAddresses->GetArray()) {

                auto pIP = rapidjson::GetValueByPointer(address, "ip");
                assert(pIP != nullptr && pIP->IsString());

                auto pPort = rapidjson::GetValueByPointer(address, "port");
                assert(pPort != nullptr && pPort->IsInt());

                upstreamInfo->addresses.emplace_back(std::make_pair(SFromP(pIP), pPort->GetInt()));
            }
        }

        pExternalInfo->upstreams.emplace_back(upstreamInfo);
    }

    return pExternalInfo;
}

static shared_ptr<NormalInfo> buildNormalInfoFromDocument(const rapidjson::Value &pDocument) {
    auto pNormalInfo = std::make_shared<NormalInfo>();

    auto pPorts = rapidjson::GetValueByPointer(pDocument, "/spec/normal/ports");
    if (pPorts == nullptr) {
        return nullptr;
    }
    assert(pPorts->IsArray());

    for (const auto &port: pPorts->GetArray()) {
        auto portInfo = std::make_shared<NormalPort>();

        auto pPortName = rapidjson::GetValueByPointer(port, "/name");
        assert(pPortName != nullptr && pPortName->IsString());
        portInfo->name = SFromP(pPortName);

        auto pPortValue = rapidjson::GetValueByPointer(port, "/port");
        assert(pPortValue != nullptr && pPortValue->IsInt());
        portInfo->port = pPortValue->GetInt();

        auto pIsTcp = rapidjson::GetValueByPointer(port, "/isTcp");
        assert(pIsTcp != nullptr && pIsTcp->IsBool());
        portInfo->isTcp = pIsTcp->GetBool();

        pNormalInfo->ports.emplace_back(portInfo);
    }

    return pNormalInfo;
}

static std::shared_ptr<ServerInfo> buildServerInfoFromDocument(const rapidjson::Value &pDocument) {

    auto pServerInfo = std::make_shared<ServerInfo>();

    auto pServerApp = rapidjson::GetValueByPointer(pDocument, "/spec/app");
    assert(pServerApp != nullptr && pServerApp->IsString());
    pServerInfo->serverApp = SFromP(pServerApp);

    auto pServerName = rapidjson::GetValueByPointer(pDocument, "/spec/server");
    assert(pServerName != nullptr && pServerName->IsString());
    pServerInfo->serverName = SFromP(pServerName);

    auto pSubType = rapidjson::GetValueByPointer(pDocument, "/spec/subType");
    assert(pSubType != nullptr && pSubType->IsString());

    std::string subTypeStr = SFromP(pSubType);

    constexpr char TafType[] = "taf";
    constexpr char DCacheType[] = "dCache";
    constexpr char DCacheProxyType[] = "dCacheProxy";
    constexpr char DCacheRouteType[] = "dCacheRouter";
    constexpr char DCacheAccessType[] = "dCacheDBAccess";
    constexpr char NormalType[] = "normal";
    constexpr char ExternalType[] = "external";

    if (subTypeStr == TafType) {
        pServerInfo->subType = ServerSubType::Taf;
    } else if (subTypeStr == DCacheType) {
        pServerInfo->subType = ServerSubType::DCache;
    } else if (subTypeStr == DCacheProxyType) {
        pServerInfo->subType = ServerSubType::DCacheProxy;
    } else if (subTypeStr == DCacheRouteType) {
        pServerInfo->subType = ServerSubType::DCacheRoute;
    } else if (subTypeStr == DCacheAccessType) {
        pServerInfo->subType = ServerSubType::DCacheDBAccess;
    } else if (subTypeStr == NormalType) {
        pServerInfo->subType = ServerSubType::Normal;
    } else if (subTypeStr == ExternalType) {
        pServerInfo->subType = ServerSubType::External;
    } else {
        assert(false);
        return nullptr;
    }

    switch (pServerInfo->subType) {
        case ServerSubType::Taf:
            pServerInfo->tafInfo = buildTafInfoFromDocument(pDocument);
            break;
        case ServerSubType::DCache:
            break;
        case ServerSubType::DCacheProxy:
            break;
        case ServerSubType::DCacheRoute:
            break;
        case ServerSubType::External:
            pServerInfo->externalInfo = buildExternalInfoFromDocument(pDocument);
            break;
        case ServerSubType::Normal:
            pServerInfo->normalInfo = buildNormalInfoFromDocument(pDocument);
            break;
        case ServerSubType::DCacheDBAccess:
            break;
    }

    auto pPods = rapidjson::GetValueByPointer(pDocument, "/status/pods");
    if (pPods != nullptr) {
        assert(pPods->IsArray());
        for (const auto &pod :pPods->GetArray()) {
            auto pPod = std::make_shared<PodStatus>();

            auto pName = rapidjson::GetValueByPointer(pod, "/name");
            assert(pName != nullptr && pName->IsString());
            pPod->name = SFromP(pName);

            auto pPodIP = rapidjson::GetValueByPointer(pod, "/podIP");
            assert(pPodIP != nullptr && pPodIP->IsString());
            pPod->podIP = SFromP(pPodIP);

            auto pHostIP = rapidjson::GetValueByPointer(pod, "/hostIP");
            assert(pHostIP != nullptr && pHostIP->IsString());
            pPod->hostIP = SFromP(pHostIP);

            auto pPresentState = rapidjson::GetValueByPointer(pod, "/presentState");
            assert(pPresentState != nullptr && pPresentState->IsString());
            pPod->presentState = SFromP(pPresentState);

            pServerInfo->pods.push_back(pPod);
        }
    }
    return pServerInfo;
}

int ServerInfoInterface::getTafServerDescriptor(const std::shared_ptr<ServerInfo> &serverInfo, taf::ServerDescriptor &descriptor) {

    const auto &tafInfo = serverInfo->tafInfo;
    if (tafInfo == nullptr) {
        return -1;
    }

    const auto &adapters = tafInfo->adapters;
    descriptor.asyncThreadNum = tafInfo->asyncThread;

    const auto &sTemplateName = tafInfo->templateName;
    assert(!sTemplateName.empty());

    string sResult;

    TC_Config templateConf = getTemplateContent(sTemplateName, sResult);

    const auto &profileContent = tafInfo->profileContent;

    if (profileContent.empty()) {
        descriptor.profile = templateConf.tostr();
    } else {
        TC_Config profileConf{};
        profileConf.parseString(tafInfo->profileContent);
        profileConf.joinConfig(templateConf, false);
        descriptor.profile = profileConf.tostr();
    }

    for (const auto &adapter:adapters) {
        AdapterDescriptor adapterDescriptor;
        adapterDescriptor.adapterName.append(serverInfo->serverApp).append(".").append(serverInfo->serverName).append(".").append(adapter->name).append(".Adapter");
        adapterDescriptor.servant.append(serverInfo->serverApp).append(".").append(serverInfo->serverName).append(".").append(adapter->name);
        adapterDescriptor.protocol = adapter->isTaf ? "taf" : "not_taf";
        adapterDescriptor.endpoint.append(adapter->isTcp ? "tcp" : "udp").append(" -h ${localip} -p ").append(to_string(adapter->port)).append(" -t ").append(
                to_string(adapter->timeout));
        adapterDescriptor.threadNum = adapter->thread;
        adapterDescriptor.maxConnections = adapter->connection;
        adapterDescriptor.queuecap = adapter->capacity;
        adapterDescriptor.queuetimeout = adapter->timeout;
        descriptor.adapters[adapterDescriptor.adapterName] = adapterDescriptor;
    }
    return 0;
}

int ServerInfoInterface::getServerDescriptor(const string &serverApp, const string &serverName, taf::ServerDescriptor &descriptor) {
    std::lock_guard<std::mutex> lockGuard(_mutex);
    const std::string sAppServer = TC_Common::lower(serverApp) + "-" + TC_Common::lower(serverName);
    auto iterator = _serverInfoMap.find(sAppServer);
    if (iterator == _serverInfoMap.end()) {
        LOG->error() << "not found" << serverApp << "-" << serverName << endl;
        return -1;
    }

    if (iterator->second == nullptr) {
        LOG->error() << "null point" << serverApp << "-" << serverName << endl;
        return -1;
    }

    const auto &serverInfo = iterator->second;

    switch (serverInfo->subType) {
        case ServerSubType::Taf:
        case ServerSubType::DCache:
        case ServerSubType::DCacheProxy:
        case ServerSubType::DCacheRoute:
        case ServerSubType::DCacheDBAccess:
            return getTafServerDescriptor(serverInfo, descriptor);
        case ServerSubType::External:
        case ServerSubType::Normal:
            return -1;
    }

    assert(false); //should not read here
    return 0;
}

TC_Config ServerInfoInterface::getTemplateContent(const string &sTemplateName, std::string &result) {
    assert(!sTemplateName.empty());
    TC_Config conf{};
    auto iterator = _templateMap.find(sTemplateName);
    if (iterator != _templateMap.end()) {
        const auto &content = iterator->second->content;
        conf.parseString(content);
        const auto &parent = iterator->second->parent;

        if (sTemplateName != parent) {
            joinParentTemplate(parent, conf, result);
        }
    }
    return conf;
}

bool ServerInfoInterface::joinParentTemplate(const string &sTemplateName, TC_Config &conf, std::string &result) {
    assert(!sTemplateName.empty());
    auto currentTemplateName = sTemplateName;

    while (true) {
        auto iterator = _templateMap.find(sTemplateName);

        if (iterator == _templateMap.end()) {
            //todo set result
            return false;
        }

        const auto &content = iterator->second->content;

        TC_Config currentConf;
        try {
            currentConf.parseString(content);
            conf.joinConfig(currentConf, false);
        } catch (...) {
            //todo set result;
            return false;
        }

        auto parentTemplateName = iterator->second->parent;
        if (currentTemplateName == parentTemplateName) {
            return true;
        }

        currentTemplateName = std::move(parentTemplateName);
    }
}

void ServerInfoInterface::onTemplateAdd(const rapidjson::Value &pDocument) {
    auto pName = rapidjson::GetValueByPointer(pDocument, "/metadata/name");
    assert(pName != nullptr && pName->IsString());

    auto pContent = rapidjson::GetValueByPointer(pDocument, "/spec/content");
    assert(pContent != nullptr && pContent->IsString());

    auto pParent = rapidjson::GetValueByPointer(pDocument, "/spec/parent");
    assert(pContent != nullptr && pContent->IsString());

    auto pTemplate = std::make_shared<Template>();

    pTemplate->content = SFromP(pContent);
    pTemplate->parent = SFromP(pParent);

    auto name = SFromP(pName);

    std::lock_guard<std::mutex> lockGuard(_mutex);
    auto iterator = _templateMap.find(name);
    if (iterator == _templateMap.end()) {
        _templateMap[name] = pTemplate;
    } else {
        iterator->second.swap(pTemplate);
    }
}

void ServerInfoInterface::onTemplateUpdate(const rapidjson::Value &pDocument) {
    onTemplateAdd(pDocument);
}

void ServerInfoInterface::onTemplateDeleted(const rapidjson::Value &pDocument) {
    auto pName = rapidjson::GetValueByPointer(pDocument, "/metadata/name");
    assert(pName != nullptr && pName->IsString());
    auto name = SFromP(pName);
    std::lock_guard<std::mutex> lockGuard(_mutex);
    _templateMap.erase(name);
}

void ServerInfoInterface::onEndpointAdd(const rapidjson::Value &pDocument) {
    auto endpointName = rapidjson::GetValueByPointer(pDocument, "/metadata/name");
    assert(endpointName != nullptr && endpointName->IsString());
    std::string sServerName = SFromP(endpointName);

    auto pServerInfo = buildServerInfoFromDocument(pDocument);

    std::lock_guard<std::mutex> lockGuard(_mutex);
    auto iterator = _serverInfoMap.find(sServerName);
    if (iterator == _serverInfoMap.end()) {
        _serverInfoMap[sServerName] = pServerInfo;
    } else {
        iterator->second.swap(pServerInfo);
    }
}

void ServerInfoInterface::onEndpointUpdate(const rapidjson::Value &pDocument) {
    onEndpointAdd(pDocument);
}

void ServerInfoInterface::onEndpointDeleted(const rapidjson::Value &pDocument) {
    auto endpointName = rapidjson::GetValueByPointer(pDocument, "/metadata/name");
    assert(endpointName != nullptr && endpointName->IsString());
    std::string sEndpointName = SFromP(endpointName);
    std::lock_guard<std::mutex> lockGuard(_mutex);
    _serverInfoMap.erase(sEndpointName);
}

void
ServerInfoInterface::findTafEndpoint(const std::shared_ptr<ServerInfo> &serverInfo, const string &sPortName, vector<taf::EndpointF> *pActiveEp,
                                     vector<taf::EndpointF> *pInactiveEp) {

    const auto &tafInfo = serverInfo->tafInfo;

    if (tafInfo == nullptr) {
        LOG->debug() << serverInfo->serverApp << "." << serverInfo->serverName << "->tafInfo is nullptr" << endl;
        return;
    }

    const auto &adapters = tafInfo->adapters;

    const auto &pods = serverInfo->pods;

    for (const auto &port:adapters) {
        if (port->name == sPortName) {
            for (const auto &pod : pods) {
                if (pod->presentState == "Active") {
                    taf::EndpointF endpointF;
                    endpointF.port = port->port;
                    endpointF.istcp = port->isTcp;
                    endpointF.timeout = port->timeout;
                    endpointF.host.append(pod->name).append(".").append(TC_Common::lower(serverInfo->serverApp)).append("-").append(TC_Common::lower(serverInfo->serverName));
                    pActiveEp->push_back(endpointF);
                } else if (pInactiveEp != nullptr) {
                    taf::EndpointF endpointF;
                    endpointF.port = port->port;
                    endpointF.istcp = port->isTcp;
                    endpointF.timeout = port->timeout;
                    endpointF.host.append(pod->name).append(".").append(TC_Common::lower(serverInfo->serverApp)).append("-").append(TC_Common::lower(serverInfo->serverName));
                    pInactiveEp->push_back(endpointF);
                }
            }
        }
    }
}

void ServerInfoInterface::loadUpChainConf() {

    constexpr char UpChainConfFile[] = "/etc/upchain/upchain.conf";

    int fok = ::access(UpChainConfFile, F_OK);
    if (fok != 0) {
        std::lock_guard<std::mutex> lockGuard(_mutex);
        if (_upChainInfo != nullptr) {
            _upChainInfo = nullptr;
            LOG->debug() << "clear upchainInfo because file \"" << UpChainConfFile << "\"not exist";
        }
        return;
    }

    int rok = ::access(UpChainConfFile, R_OK);
    if (rok != 0) {
        LOG->error() << "permission denied to read " << UpChainConfFile << endl;
        return;
    }

    auto upChainConfContent = TC_File::load2str(UpChainConfFile);
    if (upChainConfContent.empty()) {
        std::lock_guard<std::mutex> lockGuard(_mutex);
        if (_upChainInfo != nullptr) {
            _upChainInfo = nullptr;
            LOG->debug() << "clear upchainInfo because file \"" << UpChainConfFile << "\"is empty";
        }
        return;
    }

    TC_Config tcConfig;
    try {
        tcConfig.parseString(upChainConfContent);
    } catch (TC_Config_Exception &e) {
        LOG->error() << "parser file \"" << UpChainConfFile << "\" content catch exception : " << e.what() << endl;
        return;
    }

    auto upChainInfo = std::make_shared<UpChain>();

    std::vector<std::string> domains = tcConfig.getDomainVector("/upchain");
    for (const auto &domain:domains) {
        auto absDomain = string("/upchain/" + domain);
        auto lines = tcConfig.getDomainLine(absDomain);
        std::vector<taf::EndpointF> ev;
        ev.reserve(lines.size());
        for (auto &&line: lines) {
            taf::TC_Endpoint endpoint(line);
            taf::EndpointF f;
            f.host = endpoint.getHost();
            f.port = endpoint.getPort();
            f.timeout = endpoint.getTimeout();
            f.istcp = endpoint.isTcp();
            ev.emplace_back(f);
        }
        if (domain == "default") {
            upChainInfo->defaultUpChain.swap(ev);
        }
        upChainInfo->customUpChain[domain] = std::move(ev);
    }

    std::lock_guard<std::mutex> lockGuard(_mutex);
    if (_upChainInfo != nullptr) {
        _upChainInfo.swap(upChainInfo);
        return;
    }
    _upChainInfo = upChainInfo;
    LOG->debug() << "update upchainInfo success" << endl;
}

void ServerInfoInterface::findUpChainEndpoint(const string &id, vector<EndpointF> *pActiveEp, vector<taf::EndpointF> *pInactiveEp) {
    assert(pActiveEp != nullptr);
    pActiveEp->clear();

    if (_upChainInfo == nullptr) {
        return;
    }

    auto customIterator = _upChainInfo->customUpChain.find(id);
    if (customIterator != _upChainInfo->customUpChain.end()) {
        *pActiveEp = customIterator->second;
        return;
    }

    if (!_upChainInfo->defaultUpChain.empty()) {
        *pActiveEp = _upChainInfo->defaultUpChain;
    }
}


