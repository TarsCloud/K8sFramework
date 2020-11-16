
#include <thread>
#include "K8SEndpointInterface.h"
#include "RegistryServer.h"
#include "K8SRestfulClient.h"
#include "SqlBuildWrapper.h"

static inline std::string SFromP(const rapidjson::Value *p) {
    assert(p != nullptr);
    return {p->GetString(), p->GetStringLength()};
}

K8SEndpointInterface::FindEndpointRes
K8SEndpointInterface::findEndpoint(const string &sAppServer, const string &sPortName, vector<EndpointF> *pActiveEp,
                                   vector<tars::EndpointF> *pInactiveEp) {
    assert(pActiveEp != nullptr);

    pActiveEp->clear();
    if (pInactiveEp != nullptr) {
        pInactiveEp->clear();
    }

    std::lock_guard<std::mutex> lockGuard(mutex);

    auto serverIterator = serverEndpointMap.find(sAppServer);

    if (serverIterator == serverEndpointMap.end()) {
        //  serverIterator == serverEndpointMap.end() 表示 在 k8s集群中不存在此服务
        return FindEndpointRes::NoServer;
    }

    if (serverIterator->second == nullptr) {
        //  serverIterator->second == nullptr 表示 在 k8s集群中存在此服务,只是endpoints数量为0个
        return FindEndpointRes::Success;
    }

    const auto &ports = serverIterator->second->ports;

    constexpr size_t FIXED_ENDPOINT_TIMEOUT = 100000;

    //fixme k8s endpoint 没有携带 tars obj 的 timeout 信息, 所以此处只能填一个 timeout 经验值,这可能是潜在的隐患

    for (const auto &tuple:ports) {
        const auto &portName = std::get<0>(tuple);

        if (sPortName == portName) {
            const auto &portValue = std::get<1>(tuple);
            const auto &isTcp = std::get<2>(tuple);
            const auto &addresses = serverIterator->second->addresses;
            for (const auto &item:addresses) {
                tars::EndpointF endpointF;
                endpointF.istcp = isTcp;
                endpointF.port = portValue;
                if (item->podName.empty()) {
                    endpointF.host = item->podIp;
                } else {
                    endpointF.host = item->podName + "." + tars::TC_Common::lower(sAppServer);
                }
                endpointF.timeout = FIXED_ENDPOINT_TIMEOUT;
                pActiveEp->emplace_back(std::move(endpointF));
            }

            if (pInactiveEp == nullptr) {
                break;
            }

            const auto &notReadyAddresses = serverIterator->second->notReadyAddresses;
            for (const auto &item:notReadyAddresses) {
                tars::EndpointF endpointF;
                endpointF.istcp = isTcp;
                endpointF.port = portValue;
                if (item->podName.empty()) {
                    endpointF.host = item->podIp;
                } else {
                    endpointF.host = item->podName + "." + tars::TC_Common::lower(sAppServer);
                }
                endpointF.timeout = FIXED_ENDPOINT_TIMEOUT;
                pInactiveEp->emplace_back(std::move(endpointF));
            }
            break;
        }
    }
    return FindEndpointRes::Success;
}

K8SEndpointInterface::FindEndpointRes
K8SEndpointInterface::findEndpoint(const string &id, vector<EndpointF> *pActiveEp, vector<tars::EndpointF> *pInactiveEp) {
	assert(pActiveEp != nullptr);

	pActiveEp->clear();
	if (pInactiveEp != nullptr) {
		pInactiveEp->clear();
	}

	if (_externalJceProxy.empty()) {
		// 没有注入外联服务Name
		return FindEndpointRes::NoServer;
	}

	std::lock_guard<std::mutex> lockGuard(mutex);

	auto serverIterator = serverEndpointMap.find(_externalJceProxy);

	if (serverIterator == serverEndpointMap.end()) {
		//  serverIterator == serverEndpointMap.end() 表示 在 k8s集群中不存在此服务
		return FindEndpointRes::NoServer;
	}

	if (serverIterator->second == nullptr) {
		//  serverIterator->second == nullptr 表示 在 k8s集群中存在此服务,只是endpoints数量为0个
		return FindEndpointRes::Success;
	}

	const auto &ports = serverIterator->second->ports;

	constexpr size_t FIXED_ENDPOINT_TIMEOUT = 100000;

	//fixme k8s endpoint 没有携带 tars obj 的 timeout 信息, 所以此处只能填一个 timeout 经验值,这可能是潜在的隐患

	for (const auto &tuple:ports) {
		const auto &portValue = std::get<1>(tuple);
		const auto &isTcp = std::get<2>(tuple);
		const auto &addresses = serverIterator->second->addresses;
		for (const auto &item:addresses) {
			tars::EndpointF endpointF;
			endpointF.istcp = isTcp;
			endpointF.port = portValue;
			endpointF.host = item->podIp;
			endpointF.timeout = FIXED_ENDPOINT_TIMEOUT;
			pActiveEp->emplace_back(std::move(endpointF));
		}

		if (pInactiveEp == nullptr) {
			break;
		}

		const auto &notReadyAddresses = serverIterator->second->notReadyAddresses;
		for (const auto &item:notReadyAddresses) {
			tars::EndpointF endpointF;
			endpointF.istcp = isTcp;
			endpointF.port = portValue;
			endpointF.host = item->podIp;
			endpointF.timeout = FIXED_ENDPOINT_TIMEOUT;
			pInactiveEp->emplace_back(std::move(endpointF));
		}
		break;
	}
	return FindEndpointRes::Success;
}

void K8SEndpointInterface::init(const tars::TC_Config &conf) {
    try {
        _externalJceProxy = conf.get("/tars<externalJceProxy>", "");

        tars::TC_DBConf tcDBConf;
        tcDBConf.loadFromMap(conf.getDomainMap("/tars/db"));
        mysql.init(tcDBConf);
    }
    catch (TC_Config_Exception &ex) {
        LOG->error() << __FUNCTION__ << " exception: " << ex.what() << endl;
        LOG->flush();
        exit(-1);
    }
    catch (TC_Mysql_Exception &ex) {
        LOG->error() << __FUNCTION__ << " exception: " << ex.what() << endl;
        LOG->flush();
        exit(-1);
    }
}

void K8SEndpointInterface::onReceiveState(const std::string &sPodName, const std::string &settingState,
                                          const std::string &presentState) {

    std::lock_guard<std::mutex> lockGuard(mutex);

    strStream.str("");
    strStream << "/api/v1/namespaces/" << K8SRuntimeParams::interface().bindNamespace() << "/pods/" << sPodName
              << "/status";
    const std::string url = strStream.str();

    strStream.str("");
    //strStream << R"({"status":{"conditions":[{"type":"tars.io/active")" << ","
    strStream << R"({"status":{"conditions":[{"type":"tars.io/service")" << ","
              << R"("status":")" << (presentState == "Active" ? "True" : "False") << R"(",)"
              << R"("reason":")" << settingState << "/" << presentState << R"("}]}})";
    K8SRestfulClient::interface().postTask(StrategicMergePatch, url, strStream.str());

    strStream.str("");
    strStream << "update t_pod set ";
    strStream << "f_setting_state=" << SqlTr(settingState);
    strStream << ",f_present_state=" << SqlTr(presentState);
    strStream << ",f_present_message=" << SqlTr("");
    strStream << ",f_update_time=CURRENT_TIMESTAMP";
    strStream << " where f_pod_name=" << SqlTr(sPodName);

    try {
        LOG->debug() << __FUNCTION__ << strStream.str() << std::endl;
        mysql.execute(strStream.str());
    }
    catch (tars::TC_Mysql_Exception &ex) {
        LOG->error() << __FUNCTION__ << ex.what() << std::endl;
    }
}

void K8SEndpointInterface::onPodAdd(const rapidjson::Value &pDocument) {
    /*
         在 onEventAdd 的逻辑中, 会插入所有的 "Scheduled" 和 "FailedScheduling" 事件对应的Pod,这会导致非taf服务的Pod也被插入到数据库
         k8s 的 events 比较特殊,不保证一定会推送到 watcher,完全依靠 onEventAdd 插入Pod记录可能导致一些Pod 不能插入数据库
         所以,此处需要的逻辑操作有两条:
         1. 如果是 taf服务的Pod , 确保该Pod记录插入到数据库.
         2. 如果是 非taf服务的Pod, 从数据库总删除
         判断是否为 taf服务的标准为是否同时存在 tars.io/ServerApp ,tars.io/ServerName ,tars.io/ServerVersion label
     */

    do {
        auto pServerApp = rapidjson::GetValueByPointer(pDocument, "/metadata/labels/tars.io~1ServerApp");
        if (pServerApp != nullptr) {
            break;
        }

        auto pServerName = rapidjson::GetValueByPointer(pDocument, "/metadata/labels/tars.io~1ServerName");
        if (pServerName == nullptr) {
            break;
        }

        auto pServerVersion = rapidjson::GetValueByPointer(pDocument, "/metadata/labels/tars.io~1ServerVersion");
        if (pServerVersion == nullptr) {
            break;
        }

        assert(pServerApp->IsString());
        auto sServerApp = SFromP(pServerApp);

        assert(pServerName->IsString());
        auto sServerName = SFromP(pServerName);

        assert(pServerVersion->IsString());
        auto sServerVersion = SFromP(pServerVersion);

        auto pPodUid = rapidjson::GetValueByPointer(pDocument, "/metadata/uid");
        assert(pPodUid != nullptr && pPodUid->IsString());
        auto sPodUid = SFromP(pPodUid);

        auto pPodName = rapidjson::GetValueByPointer(pDocument, "/metadata/name");
        assert(pPodName != nullptr && pPodName->IsString());
        auto sPodName = SFromP(pPodName);

        auto pPodHostIp = rapidjson::GetValueByPointer(pDocument, "/status/hostIP");
        auto sNodeIp = (pPodHostIp == nullptr ? "" : SFromP(pPodHostIp));

        auto pPodPodIp = rapidjson::GetValueByPointer(pDocument, "/status/podIP");
        auto sPodIp = (pPodPodIp == nullptr ? "" : SFromP(pPodPodIp));

        //todo 构造 sql 插入 pod;
//        std::lock_guard<std::mutex> lockGuard(mutex);
//        strStream.str("");
//
//        try {
//            mysql.execute(strStream.str());
//        }
//        catch (tars::TC_Mysql_Exception &ex) {
//            LOG->error() << ex.what() << std::endl;
//        }
        return;
    } while (false);
}

void K8SEndpointInterface::onPodUpdate(const rapidjson::Value &pDocument) {
    auto pPhase = rapidjson::GetValueByPointer(pDocument, "/status/phase");
    assert(pPhase != nullptr && pPhase->IsString());
    if (SFromP(pPhase) != "Running") {
        return;
    }

    auto pServerApp = rapidjson::GetValueByPointer(pDocument, "/metadata/labels/tars.io~1ServerApp");

    if (pServerApp == nullptr) {
        return;
    }
    assert(pServerApp->IsString());
    auto sServerApp = SFromP(pServerApp);


    auto pServerName = rapidjson::GetValueByPointer(pDocument, "/metadata/labels/tars.io~1ServerName");
    if (pServerName == nullptr) {
        return;
    }
    assert(pServerName->IsString());
    auto sServerName = SFromP(pServerName);

    auto pServerVersion = rapidjson::GetValueByPointer(pDocument, "/metadata/labels/tars.io~1ServerVersion");
    if (pServerVersion == nullptr) {
        return;
    }
    assert(pServerVersion->IsString());
    auto sServerVersion = SFromP(pServerVersion);

    auto pPodUid = rapidjson::GetValueByPointer(pDocument, "/metadata/uid");
    assert(pPodUid != nullptr && pPodUid->IsString());
    auto sPodUid = SFromP(pPodUid);

    auto pPodName = rapidjson::GetValueByPointer(pDocument, "/metadata/name");
    assert(pPodName != nullptr && pPodName->IsString());
    auto sPodName = SFromP(pPodName);

    auto pPodHostIp = rapidjson::GetValueByPointer(pDocument, "/status/hostIP");
    auto sNodeIp = (pPodHostIp == nullptr ? "" : SFromP(pPodHostIp));

    auto pPodPodIp = rapidjson::GetValueByPointer(pDocument, "/status/podIP");
    auto sPodIp = (pPodPodIp == nullptr ? "" : SFromP(pPodPodIp));


    std::lock_guard<std::mutex> lockGuard(mutex);
    strStream.str("");
    strStream << "update t_pod set ";
    strStream << "f_node_ip=" << SqlTr(sNodeIp);
    strStream << ",f_pod_ip=" << SqlTr(sPodIp);
    strStream << ",f_service_version=" << SqlTr(sServerVersion);
    strStream << ",f_update_time=" << "CURRENT_TIMESTAMP";
    strStream << " where f_pod_name=" << SqlTr(sPodName);

    try {
        mysql.execute(strStream.str());
    }
    catch (tars::TC_Mysql_Exception &ex) {
        LOG->error() << ex.what() << std::endl;
    }
}

void K8SEndpointInterface::onPodDeleted(const rapidjson::Value &pDocument) {

    auto pPodUid = rapidjson::GetValueByPointer(pDocument, "/metadata/uid");
    auto sPodUid = SFromP(pPodUid);

    std::lock_guard<std::mutex> lockGuard(mutex);

    strStream.str("");
    strStream
            << "insert into t_pod_history (f_pod_id, f_pod_name, f_pod_ip, f_node_ip, f_server_app, f_server_name, f_service_version,f_create_time) "
               "select f_pod_id,f_pod_name,f_pod_ip,f_node_ip,f_server_app,f_server_name,f_service_version,f_create_time from t_pod where f_pod_id="
            << SqlTr(sPodUid) << " on duplicate key update t_pod_history.f_pod_id=t_pod_history.f_pod_id";
    try {
        mysql.execute(strStream.str());
    }
    catch (tars::TC_Mysql_Exception &ex) {
        LOG->error() << ex.what() << std::endl;
    }

    strStream.str("");
    strStream << "delete from t_pod where f_pod_id=" << SqlTr(sPodUid);
    try {
        mysql.execute(strStream.str());
    }
    catch (tars::TC_Mysql_Exception &ex) {
        LOG->error() << "Exec Sql Get Exception : " << ex.what() << std::endl;
    }
}

void K8SEndpointInterface::onEndpointAdd(const rapidjson::Value &pDocument) {
    std::string sK8SEndpointName;
    if (!checkTafIOEndpoint(pDocument, sK8SEndpointName)) {
        return;
    }

    std::shared_ptr<K8SEndpoint> pK8SEndpoint = parseEndpoint(pDocument);

    std::lock_guard<std::mutex> lockGuard(mutex);
    serverEndpointMap[sK8SEndpointName] = pK8SEndpoint;
}

void K8SEndpointInterface::onEndpointUpdate(const rapidjson::Value &pDocument) {
    std::string sK8SEndpointName;
    if (!checkTafIOEndpoint(pDocument, sK8SEndpointName)) {
        return;
    }

    std::shared_ptr<K8SEndpoint> pK8SEndpoint = parseEndpoint(pDocument);

    std::lock_guard<std::mutex> lockGuard(mutex);

    auto iterator = serverEndpointMap.find(sK8SEndpointName);
    if (iterator == serverEndpointMap.end()) {
        serverEndpointMap[sK8SEndpointName] = pK8SEndpoint;
    } else {
        std::swap(iterator->second, pK8SEndpoint);
    }
}

void K8SEndpointInterface::onEndpointDeleted(const rapidjson::Value &pDocument) {
    std::string sK8SEndpointName;
    if (!checkTafIOEndpoint(pDocument, sK8SEndpointName)) {
        return;
    }

    std::lock_guard<std::mutex> lockGuard(mutex);
    serverEndpointMap.erase(sK8SEndpointName);
}

void K8SEndpointInterface::onEventAdd(const rapidjson::Value &pDocument) {

    auto pInvolveObjectKindLabel = rapidjson::GetValueByPointer(pDocument, "/involvedObject/kind");
    assert(pInvolveObjectKindLabel != nullptr && pInvolveObjectKindLabel->IsString());

    constexpr char ExpectObject[] = "Pod";
    if (strcmp(ExpectObject, pInvolveObjectKindLabel->GetString()) != 0) {
        return;
    }

    auto pResourceVersionLabel = rapidjson::GetValueByPointer(pDocument, "/metadata/resourceVersion");
    assert(pResourceVersionLabel != nullptr);
    assert(pResourceVersionLabel->IsString());

    auto iResourceVersion = std::stoull(pResourceVersionLabel->GetString(), nullptr, 10);

    auto pEventReason = rapidjson::GetValueByPointer(pDocument, "/reason");
    assert(pEventReason != nullptr && pEventReason->IsString());
    auto sReason = SFromP(pEventReason);

    auto pEventMessage = rapidjson::GetValueByPointer(pDocument, "/message");
    assert(pEventMessage != nullptr && pEventMessage->IsString());
    auto sMessage = SFromP(pEventMessage);

    // message 有时候包含了双引号，单引号,影响 Sql的构造
    sMessage = tars::TC_Common::replace(sMessage, {
        {R"(")", R"(\")"},
        {R"(')", R"(\')"}
    });

    auto pPodUidLabel = rapidjson::GetValueByPointer(pDocument, "/involvedObject/uid");

    std::lock_guard<std::mutex> lockGuard(mutex);

    string sPodUid;
    if (sReason == "Scheduled" || sReason == "FailedScheduling") {

        auto pPodNameLabel = rapidjson::GetValueByPointer(pDocument, "/involvedObject/name");
        assert(pPodNameLabel != nullptr && pPodNameLabel->IsString());
        auto sPodName = SFromP(pPodNameLabel);

        vector<std::string> v = tars::TC_Common::sepstr<string>(sPodName, "-");
        if (v.size() <= 2) {
            return;  //不符合 tars 服务预期格式
        }

        const auto &sServerApp = v[0];
        const auto &sServerName = v[1];

        strStream.str("");

        // taf服务调度失败，无uid
        if (pPodUidLabel == nullptr || !pPodUidLabel->IsString()) {
            strStream << "insert into t_server_notify (f_app_server, f_pod_name, f_notify_message, f_notify_level, f_notify_thread, f_notify_source) ";
            strStream << "values (" << SqlTr(sServerApp+"."+sServerName) << "," << SqlTr(sPodName) << "," << SqlTr("[alarm] " + sMessage) << "," ;
            strStream << SqlTr("[alarm]") << "," << SqlTr("-1") << "," << SqlTr("tars-tafregistry") << ")";

            try {
                mysql.execute(strStream.str());
            } catch (TC_Mysql_Exception &ex) {
                LOG->error() << ex.what() << std::endl;;
            }
            return;
        }

        sPodUid = SFromP(pPodUidLabel);

        strStream << "insert into t_pod (f_pod_id, f_pod_name, f_present_state, f_present_message, f_resource_version, f_server_app, f_server_name)";
        strStream << " select " << SqlTr(sPodUid) << "," << SqlTr(sPodName) << "," << SqlTr(sReason) << "," << SqlTr(sMessage) << "," << SqlTr(iResourceVersion);
        strStream << ",ts.f_server_app, ts.f_server_name from t_server as ts";
        strStream << " where ts.f_server_app=" << SqlTr(sServerApp) << " and ts.f_server_name=" << SqlTr(sServerName);
        strStream << " on duplicate key update";
        strStream << " f_present_state=(if(f_resource_version<" << SqlTr(iResourceVersion) << "," << SqlTr(sReason) << ",f_present_state))";
        strStream << ",f_present_message=(if(f_resource_version<" << SqlTr(iResourceVersion) << "," << SqlTr(sMessage) << ",f_present_message))";
        strStream << ",f_resource_version=(if(f_resource_version<" << SqlTr(iResourceVersion) << "," << SqlTr(iResourceVersion) << ",f_resource_version))";
        try {
            mysql.execute(strStream.str());
        }
        catch (tars::TC_Mysql_Exception &ex) {
            LOG->error() << ex.what() << std::endl;
        }
        return;
    }

    // 更新逻辑一定有uid
    assert(pPodUidLabel != nullptr && pPodUidLabel->IsString());
    sPodUid = SFromP(pPodUidLabel);

    strStream.str("");
    strStream << "update t_pod set ";
    strStream << "f_setting_state=" << "'" << "Active" << "'";
    strStream << ",f_present_state=" << SqlTr(sReason);
    strStream << ",f_present_message=" << SqlTr(sMessage);
    strStream << ",f_update_time=" << "CURRENT_TIMESTAMP";
    strStream << ",f_resource_version=" << SqlTr(iResourceVersion);
    strStream << " where f_pod_id=" << SqlTr(sPodUid) << " and f_resource_version <" << SqlTr(iResourceVersion);

    try {
        mysql.execute(strStream.str());
    }
    catch (tars::TC_Mysql_Exception &ex) {
        LOG->error() << ex.what() << std::endl;
    }
}

void K8SEndpointInterface::onEventUpdate(const rapidjson::Value &pDocument) {
    onEventAdd(pDocument);
}

void K8SEndpointInterface::onEventDeleted(const rapidjson::Value &pDocument) {
    //do nothing
}

std::shared_ptr<K8SEndpoint> K8SEndpointInterface::parseEndpoint(const rapidjson::Value &pDocument) {
    auto pK8SEndpoint_ = std::make_shared<K8SEndpoint>();

    auto pPorts = rapidjson::GetValueByPointer(pDocument, "/subsets/0/ports");
    if (pPorts == nullptr) {
        return nullptr;
    }

    assert(pPorts->IsArray());
    for (const auto &v :pPorts->GetArray()) {
        constexpr int NodeObjectPortValue = 19385;
        auto pPort = rapidjson::GetValueByPointer(v, "/port");
        int portValue = pPort->GetInt();
        if (portValue == NodeObjectPortValue) {
            continue;
        }

        auto pName = rapidjson::GetValueByPointer(v, "/name");
        assert(pName->IsString());

        auto pProtocol = rapidjson::GetValueByPointer(v, "/protocol");
        assert(pProtocol->IsString());
        bool bIsTcp = strncmp(pProtocol->GetString(), "TCP", pProtocol->GetStringLength()) == 0;
        pK8SEndpoint_->ports.emplace_back(std::make_tuple(SFromP(pName), portValue, bIsTcp));
    }

    if (pK8SEndpoint_->ports.empty()) {
        return nullptr;
    }

    auto pAddresses = rapidjson::GetValueByPointer(pDocument, "/subsets/0/addresses");
    if (pAddresses != nullptr) {
        assert(pAddresses->IsArray());
        for (const auto &v :pAddresses->GetArray()) {
            auto pIp = rapidjson::GetValueByPointer(v, "/ip");
            assert(pIp != nullptr && pIp->IsString());
            auto pHostname = rapidjson::GetValueByPointer(v, "/hostname");
            auto sHostname = (pHostname == nullptr ? "" : SFromP(pHostname));
            pK8SEndpoint_->addresses.emplace_back(std::make_shared<K8SPod>(sHostname, SFromP(pIp)));
        }
    }

    auto pNotReadyAddresses = rapidjson::GetValueByPointer(pDocument, "/subsets/0/notReadyAddresses");
    if (pNotReadyAddresses != nullptr) {
        assert(pNotReadyAddresses->IsArray());
        for (const auto &v :pNotReadyAddresses->GetArray()) {
            auto pIp = rapidjson::GetValueByPointer(v, "/ip");
            assert(pIp != nullptr && pIp->IsString());
            auto pHostname = rapidjson::GetValueByPointer(v, "/hostname");
            auto sHostname = (pHostname == nullptr ? "" : SFromP(pHostname));
            pK8SEndpoint_->notReadyAddresses.emplace_back(std::make_shared<K8SPod>(sHostname, SFromP(pIp)));
        }
    }

    return pK8SEndpoint_;
}

bool K8SEndpointInterface::checkTafIOEndpoint(const rapidjson::Value &pDocument, std::string& sK8SEndpointName) {
    auto pK8SEndpointName = rapidjson::GetValueByPointer(pDocument, "/metadata/name");
    assert(pK8SEndpointName != nullptr && pK8SEndpointName->IsString());

    // tars.io/ServerApp , tars.io/ServerName 标签, 则认为不是 tars 业务服务的 endpoint, 不做处理
    auto pServerApp = rapidjson::GetValueByPointer(pDocument, "/metadata/labels/tars.io~1ServerApp");
    if (pServerApp == nullptr) {
        return false;
    }
    auto pServerName = rapidjson::GetValueByPointer(pDocument, "/metadata/labels/tars.io~1ServerName");
    if (pServerName == nullptr) {
        return false;
    }

    sK8SEndpointName = SFromP(pK8SEndpointName);

    return true;
}
