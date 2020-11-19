#include "RegistryImp.h"
#include "RegistryServer.h"
#include "util/tc_mysql.h"
#include "ServerInfoInterface.h"
#include "K8SClient.h"

void RegistryImp::initialize() {
    LOG->debug() << "RegistryImp init ok." << endl;
}

void RegistryImp::updateServerState(const string &podName, const string &settingState, const string &presentState, tars::CurrentPtr current) {
    std::stringstream strStream;
    strStream.str("");
    strStream << "/api/v1/namespaces/" << K8SRuntimeParams::interface().bindNamespace() << "/pods/" << podName
              << "/status";
    const std::string url = strStream.str();
    strStream.str("");
    strStream << R"({"status":{"conditions":[{"type":"tars.io/active")" << ","
              << R"("status":")" << ((settingState == "Active" && presentState == "Active") ? "True" : "False") << R"(",)"
              << R"("reason":")" << settingState << "/" << presentState << R"("}]}})";
    K8SClient::instance().postTask(StrategicMergePatch, url, strStream.str());
}

tars::Int32 RegistryImp::getServerDescriptor(const string &serverApp, const string &serverName, ServerDescriptor &serverDescriptor, tars::CurrentPtr current) {
    return ServerInfoInterface::instance().getServerDescriptor(serverApp, serverName, serverDescriptor);
}

