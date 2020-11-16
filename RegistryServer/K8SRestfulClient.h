
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

enum K8SRestfulTaskType {
    Patch,
    Post,
    Delete,
    StrategicMergePatch
};

class K8SRestClientSession;

class K8SRestfulClient {

public:

    static K8SRestfulClient &interface() {
        static K8SRestfulClient k8SRestClient;
        return k8SRestClient;
    }

    void run();

    void postTask(K8SRestfulTaskType type, const std::string &url, const std::string &body);

private:
    static std::string buildPostContent(const std::string &url, const std::string &body);

    static std::string buildPatchContent(const std::string &url, const std::string &body);

    static std::string buildSMPatchContent(const std::string &url, const std::string &body);

    static std::string buildDeleteContent(const std::string &url);

    K8SRestfulClient() = default;

private:
    std::shared_ptr<K8SRestClientSession> sessionPtr_{};
    asio::io_context ioContext_{1};
};