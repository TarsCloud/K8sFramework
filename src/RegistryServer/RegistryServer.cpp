#include "RegistryServer.h"
#include "RegistryImp.h"
#include "QueryImp.h"
#include "K8SListWatchSession.h"
#include "UpChainSession.h"
#include "IOContext.h"
#include "ServerInfoInterface.h"
#include "K8SClient.h"

void RegistryServer::initialize() {

    LOG->debug() << "RegistryServer::initialize..." << endl;

    K8SRuntimeParams::interface().init();

    _ioThread = std::thread([]() {
        auto TEndpointsWatcher = std::make_shared<K8SListWatchSession>(
                IOContext::instance().getIOContext());
        std::string tendpointsWatchUrl =
                std::string("/apis/k8s.taf.io/v1alpha1/namespaces/") + K8SRuntimeParams::interface().bindNamespace() + "/tendpoints";
        TEndpointsWatcher->setResourceUrl(tendpointsWatchUrl);
        TEndpointsWatcher->setCallBack(handleEndpointsEvent);
        TEndpointsWatcher->prepare();

        auto TTemplateWatcher = std::make_shared<K8SListWatchSession>(IOContext::instance().getIOContext());
        std::string podsWatchUrl =
                std::string("/apis/k8s.taf.io/v1alpha1/namespaces/") + K8SRuntimeParams::interface().bindNamespace() + "/ttemplates";
        TTemplateWatcher->setResourceUrl(podsWatchUrl);
        TTemplateWatcher->setCallBack(handleTemplateEvent);
        TTemplateWatcher->prepare();

        auto UpChainLoader = std::make_shared<UpChainLoadSession>(IOContext::instance().getIOContext());
        UpChainLoader->setCallBack([]() {
            ServerInfoInterface::instance().loadUpChainConf();
        });

        UpChainLoader->prepare();

        IOContext::instance().run();
    });

    _ioThread.detach();

    _k8sClientThread = std::thread([]() {
        K8SClient::instance().run();
    });

    _k8sClientThread.detach();

    if (!K8SListWatchSession::WaitForCacheSync()) {
        LOG->error() << " WaitForCacheSync K8S Error";
        exit(-1);
    }

    constexpr char PodNameEnv[] = "PodName";

    std::string podName = ::getenv(PodNameEnv);
    if (podName.empty()) {
        LOG->error() << "Get Empty PodName Value ,Program Will Exit " << endl;
        cerr << "Get Empty PodName Value ,Program Will Exit " << endl;
    }

    std::stringstream strStream;
    strStream.str("");
    strStream << "/api/v1/namespaces/" << K8SRuntimeParams::interface().bindNamespace() << "/pods/" << podName
              << "/status";
    const std::string url = strStream.str();
    strStream.str("");
    strStream << R"({"status":{"conditions":[{"type":"taf.io/active","status":"True","reason":"Active/Active"}]}})";
    K8SClient::instance().postTask(StrategicMergePatch, url, strStream.str());

    try {
        constexpr char FIXED_QUERY_SERVANT[] = "taf.tafregistry.QueryObj";
        constexpr char FIXED_REGISTRY_SERVANT[] = "taf.tafregistry.RegistryObj";
        addServant<QueryImp>(FIXED_QUERY_SERVANT);
        addServant<RegistryImp>(FIXED_REGISTRY_SERVANT);
    }
    catch (TC_Exception &ex) {
        LOG->error() << "RegistryServer initialize exception:" << ex.what() << endl;
        cerr << "RegistryServer initialize exception:" << ex.what() << endl;
        LOG->flush();
        exit(-1);
    }
    LOG->debug() << "RegistryServer::initialize OK!" << endl;
}

void RegistryServer::destroyApp() {
    LOG->error() << "RegistryServer::destroyApp ok" << endl;
}
