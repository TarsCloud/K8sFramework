
#pragma once

#include "K8SRuntimeParams.h"
#include "HttpParser.h"
#include "asio/io_context.hpp"
#include "asio/ssl/stream.hpp"
#include "asio/streambuf.hpp"
#include "asio/ip/tcp.hpp"
#include "asio/error.hpp"
#include "rapidjson/document.h"
#include "rapidjson/pointer.h"
#include <string>
#include <stdexcept>
#include <utility>
#include <vector>
#include <map>
#include <memory>
#include <mutex>
#include <functional>
#include <util/tc_file.h>

enum k8SEventTypeEnum {
    K8SEventTypeAdded = 1u,
    K8SEventTypeDeleted = 2u,
    K8SEventTypeModified = 3u,
    K8SEventTypeBookmark = 4u,
    K8SEventTypeError = 5u,
};

enum K8SHandleListEvent {
    AsAdd,
    Discard
};

class K8SListWatchSession : public std::enable_shared_from_this<K8SListWatchSession> {
    struct ResponseState {
        ResponseState() = default;

        bool headComplete_{false};
        bool messageComplete_{false};
        bool bodyArrive_{false};
        unsigned int code_{0};
        std::string bodyBuffer_{};
    };

public:
    explicit K8SListWatchSession(asio::io_context &ioc);

    ~K8SListWatchSession() = default;


    void setCallBack(std::function<void(k8SEventTypeEnum, const rapidjson::Value &)> callBack) {
        _callBack = std::move(callBack);
    }

    void setResourceUrl(std::string url) {
        resourceUrl_ = std::move(url);
    }

    void setHandleListType(K8SHandleListEvent howHandleListEvent) {
        handleListResult_ = howHandleListEvent;
    }

    void prepare();

    static bool WaitForCacheSync() {
        constexpr size_t MAX_WAIT_CACHE_SYNC_TIME = 3 * 60 * 1000 * 1000; //  3 min;
        constexpr size_t WAIT_CACHE_SYNC_INTERVAL = 500 * 1000;     // 500 ms
        static_assert(WAIT_CACHE_SYNC_INTERVAL && MAX_WAIT_CACHE_SYNC_TIME % WAIT_CACHE_SYNC_INTERVAL == 0, "");

        for (size_t i = 0; i < MAX_WAIT_CACHE_SYNC_TIME; i += WAIT_CACHE_SYNC_INTERVAL) {
            usleep(WAIT_CACHE_SYNC_INTERVAL);
            if (_waitCacheSyncSize < 0) {
                // means some error;
                return false;
            }
            if (_waitCacheSyncSize == 0) {
                return true;
            }
        }
        return false;
    };

private:

    static void fail(asio::error_code ec, char const *what);

    void clearResponseState();

    void afterConnect(asio::error_code ec);

    void afterHandshake(asio::error_code ec);

    void doListRequest();

    void afterListRequest(asio::error_code ec, std::size_t bytes_transferred);

    void doReadListResponse();

    void afterReadListResponse(asio::error_code ec, size_t bytes_transferred);

    void doWatchRequest();

    void afterWatchRequest(asio::error_code ec, std::size_t bytes_transferred);

    void doReadWatchResponse();

    void afterReadWatchResponse(asio::error_code ec, size_t bytes_transferred);

    bool handlerK8SWatchEvent(const rapidjson::Document &document);

private:
    K8SHandleListEvent handleListResult_{AsAdd};
    unsigned long long handledBiggestResourceVersion_{0};
    asio::ssl::context &sslContext_;
    asio::ssl::stream<asio::ip::tcp::socket> stream_;
    asio::streambuf responseBuffer_{};
    std::string requestContext_{};
    std::string resourceUrl_;
    std::string limitContinue_{};
    http_parser responseParser_{};
    http_parser_settings responseParserSetting_{};
    ResponseState responseState_{};
    std::function<void(k8SEventTypeEnum, const rapidjson::Value &)> _callBack;

    static std::atomic_int _waitCacheSyncSize;
};

