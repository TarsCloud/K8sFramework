
#pragma once

#include <memory>
#include <vector>
#include <asio/ssl/stream.hpp>
#include <asio/ip/tcp.hpp>
#include <asio/streambuf.hpp>
#include <asio/write.hpp>
#include <asio/read.hpp>
#include <asio/read_until.hpp>
#include <asio/connect.hpp>
#include <queue>
#include <mutex>
#include <iostream>
#include "HttpParser.h"
#include "K8SRuntimeParams.h"

enum K8SClientRequestType {
    Patch,
    Post,
    Delete,
    StrategicMergePatch
};

class K8SClientSession;

class K8SClient {

public:

    static K8SClient &instance() {
        static K8SClient k8SClient;
        return k8SClient;
    }

    void run();

    void postTask(K8SClientRequestType type, const std::string &url, const std::string &body);

private:
    static std::string buildPostContent(const std::string &url, const std::string &body);

    static std::string buildPatchContent(const std::string &url, const std::string &body);

    static std::string buildSMPatchContent(const std::string &url, const std::string &body);

    static std::string buildDeleteContent(const std::string &url);

    K8SClient() = default;

private:
    std::shared_ptr<K8SClientSession> sessionPtr_{};
    asio::io_context ioContext_{1};
};