#include "RegistryServer.h"
#include "RegistryImp.h"
#include "QueryImp.h"
#include "K8SListWatchSession.h"
#include "K8SListWatchIOContext.h"
#include "K8SEndpointInterface.h"
#include "K8SRestfulClient.h"

void RegistryServer::initialize() {

    LOG->debug() << "RegistryServer::initialize..." << endl;

    TC_Config &conf = RegistryServer::getConfig();
    K8SEndpointInterface::instance().init(conf);
    K8SRuntimeParams::interface().init();

    _k8sWatchThread = std::thread([]() {
        auto k8sEndpointsWatch = std::make_shared<K8SListWatchSession>(
                K8SListWatchIOContext::instance().getIOContext());
        std::string endpointsWatchUrl =
                std::string("/api/v1/namespaces/") + K8SRuntimeParams::interface().bindNamespace() + "/endpoints";
        k8sEndpointsWatch->setResourceUrl(endpointsWatchUrl);
        k8sEndpointsWatch->setCallBack(handleK8SEndpointsEvent);
        k8sEndpointsWatch->prepare();

        auto k8sPodsWatcher = std::make_shared<K8SListWatchSession>(K8SListWatchIOContext::instance().getIOContext());
        std::string podsWatchUrl =
                std::string("/api/v1/namespaces/") + K8SRuntimeParams::interface().bindNamespace() + "/pods";
        k8sPodsWatcher->setResourceUrl(podsWatchUrl);
        k8sPodsWatcher->setCallBack(handleK8SPodsEvent);
        k8sPodsWatcher->setHandleListType(Discard);
        k8sPodsWatcher->prepare();

        auto k8sEventsWatcher = std::make_shared<K8SListWatchSession>(K8SListWatchIOContext::instance().getIOContext());
        std::string eventsWatchUrl =
                std::string("/api/v1/namespaces/") + K8SRuntimeParams::interface().bindNamespace() + "/events";
        k8sEventsWatcher->setResourceUrl(eventsWatchUrl);
        k8sEventsWatcher->setCallBack(handleK8SEventsEvent);
        k8sEventsWatcher->setHandleListType(Discard);
        k8sEventsWatcher->prepare();

        K8SListWatchIOContext::instance().run();
    });
    _k8sWatchThread.detach();

    _k8sRestfulClientThread = std::thread([]() {
        K8SRestfulClient::interface().run();
    });
    _k8sRestfulClientThread.detach();

    if (!K8SListWatchSession::WaitForCacheSync()) {
        LOG->error() << " WaitForCacheSync K8S Error";
        LOG->flush();
        exit(-1);
    }


    try {
        constexpr char FIXED_QUERY_SERVANT[] = "tars.tafregistry.QueryObj";
        constexpr char FIXED_REGISTRY_SERVANT[] = "tars.tafregistry.RegistryObj";
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
