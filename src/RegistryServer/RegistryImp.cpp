#include "RegistryImp.h"
#include "RegistryServer.h"
#include "util/tc_mysql.h"
#include "ServerInfoInterface.h"
#include "K8SClient.h"
#include "K8SParams.h"

void RegistryImp::initialize() {
    LOG->debug() << "RegistryImp init ok." << endl;
}

void RegistryImp::updateServerState(const std::string &podName, const std::string &settingState, const std::string &presentState, taf::CurrentPtr current) {
    std::stringstream strStream;
    strStream.str("");
    strStream << "/api/v1/namespaces/" << K8SParams::instance().bindNamespace() << "/pods/" << podName << "/status";
    const std::string patchUrl = strStream.str();

    strStream.str("");
    strStream << R"({"status":{"conditions":[{"type":"taf.io/active")" << ","
              << R"("status":")" << ((settingState == "Active" && presentState == "Active") ? "True" : "False") << R"(",)"
              << R"("reason":")" << settingState << "/" << presentState << R"("}]}})";

    const std::string patchBody = strStream.str();

    for (auto i = 0; i < 5; ++i) {

        auto patchRequest = K8SClient::instance().postRequest(K8SClientRequestMethod::StrategicMergePatch, patchUrl, patchBody);

        bool finish = patchRequest->waitFinish(std::chrono::milliseconds(20));
        if (!finish) {
            LOG->debug() << "Update Server State Overtime" << endl;
            return;
        }

        if (patchRequest->state() != Done) {
            LOG->debug() << "Update Server State Error: " << patchRequest->stateMessage() << endl;
            continue;
        }

        if (patchRequest->responseCode() != HTTP_STATUS_OK) {
            LOG->debug() << "Update Server State Error: " << patchRequest->stateMessage() << endl;
            continue;
        }

        return;
    }

    LOG->error() << "Update Server State Error "<< endl;
    return;
}

taf::Int32 RegistryImp::getServerDescriptor(const std::string &serverApp, const std::string &serverName, ServerDescriptor &serverDescriptor, taf::CurrentPtr current) {
    return ServerInfoInterface::instance().getServerDescriptor(serverApp, serverName, serverDescriptor);
}
