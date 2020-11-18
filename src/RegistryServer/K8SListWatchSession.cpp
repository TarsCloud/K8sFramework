
#include "K8SListWatchSession.h"
#include <unordered_map>
#include <servant/taf_logger.h>
#include "rapidjson/document.h"
#include "asio/read.hpp"
#include "rapidjson/pointer.h"
#include "rapidjson/stringbuffer.h"

void K8SListWatchSession::fail(asio::error_code ec, char const *what) {
    std::string message;
    message = message + what + ":" + ec.message();
    _waitCacheSyncSize = -1;
    throw std::runtime_error(message);
}

atomic_int K8SListWatchSession::_waitCacheSyncSize{0};

K8SListWatchSession::K8SListWatchSession(asio::io_context &ioc)
        : sslContext_(K8SRuntimeParams::interface().sslContext()), stream_(ioc, sslContext_) {
    memset(&responseParserSetting_, 0, sizeof(responseParserSetting_));
    responseParser_.data = &responseState_;

    responseParserSetting_.on_headers_complete = [](http_parser *p) -> int {
        auto *state = static_cast<ResponseState * >(p->data);
        state->code_ = p->status_code;
        state->headComplete_ = true;
        return 0;
    };

    responseParserSetting_.on_body = [](http_parser *p, const char *at, size_t length) -> int {
        auto *state = static_cast<ResponseState * >(p->data);
        state->bodyBuffer_.append(at, length);
        state->bodyArrive_ = true;
        return 0;
    };

    responseParserSetting_.on_message_complete = [](http_parser *p) -> int {
        auto *state = static_cast<ResponseState * >(p->data);
        state->messageComplete_ = true;
        return 0;
    };
}


void K8SListWatchSession::prepare() {
    ++_waitCacheSyncSize;
    asio::ip::tcp::endpoint endpoint(asio::ip::address::from_string(K8SRuntimeParams::interface().apiServerHost()),
                                     K8SRuntimeParams::interface().apiServerPort());
    stream_.next_layer().async_connect(endpoint, std::bind(&K8SListWatchSession::afterConnect, shared_from_this(),
                                                           std::placeholders::_1));
}

void K8SListWatchSession::afterConnect(asio::error_code ec) {
    if (ec) {
        return fail(ec, __FUNCTION__);
    }
    stream_.async_handshake(asio::ssl::stream_base::client,
                            std::bind(&K8SListWatchSession::afterHandshake, shared_from_this(), std::placeholders::_1));
}

void K8SListWatchSession::afterHandshake(asio::error_code ec) {
    if (ec) {
        return fail(ec, __FUNCTION__);
    }
    doListRequest();
}

void K8SListWatchSession::clearResponseState() {
    responseState_.headComplete_ = false;
    responseState_.bodyArrive_ = false;
    responseState_.messageComplete_ = false;
    responseState_.bodyBuffer_.clear();
    http_parser_init(&responseParser_, HTTP_RESPONSE);
}

void K8SListWatchSession::doListRequest() {
    std::ostringstream strStream;
    strStream << "GET ";
    strStream << resourceUrl_ << "?limit=30";

    if (!limitContinue_.empty()) {
        strStream << "&continue=" << limitContinue_;
    }
    LOG->debug() << "K8SListUrl : " << strStream.str() << std::endl;
    strStream << " HTTP/1.1\r\n";
    strStream << "Host: " << K8SRuntimeParams::interface().apiServerHost() << ":"
              << K8SRuntimeParams::interface().apiServerPort() << "\r\n";
    strStream << "Authorization: Bearer " << K8SRuntimeParams::interface().bindToken() << "\r\n";
    strStream << "\r\n";

    requestContext_ = strStream.str();
    asio::async_write(stream_, asio::buffer(requestContext_),
                      std::bind(&K8SListWatchSession::afterListRequest, shared_from_this(), std::placeholders::_1,
                                std::placeholders::_2));
}

void K8SListWatchSession::afterListRequest(asio::error_code ec, std::size_t) {
    if (ec) {
        return fail(ec, __FUNCTION__);
    }

    clearResponseState();
    doReadListResponse();
}

void K8SListWatchSession::doReadListResponse() {
    stream_.async_read_some(responseBuffer_.prepare(1024 * 1024 * 2),
                            std::bind(&K8SListWatchSession::afterReadListResponse, shared_from_this(),
                                      std::placeholders::_1, std::placeholders::_2));
}

void K8SListWatchSession::afterReadListResponse(asio::error_code ec, size_t bytes_transferred) {
    if (ec) {
        return fail(ec, __FUNCTION__);
    }

    const char *responseData = asio::buffer_cast<const char *>(responseBuffer_.data());

    size_t parserSize = http_parser_execute(&responseParser_, &responseParserSetting_, responseData, bytes_transferred);
    responseBuffer_.consume(parserSize);

    if (!responseState_.messageComplete_) {
        doReadListResponse();
        return;
    }

    if (responseState_.code_ != HTTP_STATUS_OK) {
        LOG->error() << "K8SListWatch Session Receive Unexpected Response, And Program Will Exit : "
                     << responseState_.bodyBuffer_
                     << std::endl;
        LOG->flush();
        exit(-1);
    }

    rapidjson::Document jsonDocument{};
    assert(!responseState_.bodyBuffer_.empty());
    jsonDocument.Parse(responseState_.bodyBuffer_.data(), responseState_.bodyBuffer_.size());

    if (jsonDocument.HasParseError()) {
        LOG->error() << "K8SListWatch Session Receive Unexpected Response : " << responseState_.bodyBuffer_
                     << std::endl;
        throw std::runtime_error("K8SListWatch Session Receive Unexpected Response");
    }

    if (handleListResult_ == AsAdd) {
        auto &&jsonArray = jsonDocument["items"].GetArray();
        for (auto &&item : jsonArray) {
            if (_callBack) {
                _callBack(K8SEventTypeAdded, item);
            }
        }
    }

    auto pResourceVersion = rapidjson::GetValueByPointer(jsonDocument, "/metadata/resourceVersion");
    assert(pResourceVersion != nullptr);
    assert(pResourceVersion->IsString());
    auto iResourceVersion = std::stoull(pResourceVersion->GetString(), nullptr, 10);
    handledBiggestResourceVersion_ = std::max(handledBiggestResourceVersion_, iResourceVersion);
    auto pLimitContinue = rapidjson::GetValueByPointer(jsonDocument, "/metadata/continue");
    if (pLimitContinue == nullptr) {
        limitContinue_.clear();
        --_waitCacheSyncSize;
        doWatchRequest();
        return;
    }
    assert(pLimitContinue->IsString());
    limitContinue_ = string(pLimitContinue->GetString(), pLimitContinue->GetStringLength());
    if (limitContinue_.empty()) {
        --_waitCacheSyncSize;
        doWatchRequest();
        return;
    };
    doListRequest();
}

void K8SListWatchSession::doWatchRequest() {

    std::ostringstream strStream;
    strStream << "GET ";
    strStream << resourceUrl_ << "?watch=1&allowWatchBookmarks=true&timeoutSeconds=" << (std::rand() % 14 + 45) * 60;
    if (handledBiggestResourceVersion_ != 0) {
        strStream << "&resourceVersion=" << handledBiggestResourceVersion_;
    }
    LOG->debug() << "K8SWatchUrl : " << strStream.str() << std::endl;
    strStream << " HTTP/1.1\r\n";
    strStream << "Host: " << K8SRuntimeParams::interface().apiServerHost() << ":"
              << K8SRuntimeParams::interface().apiServerPort() << "\r\n";
    strStream << "Authorization: Bearer " << K8SRuntimeParams::interface().bindToken() << "\r\n";
    strStream << "\r\n";

    requestContext_ = strStream.str();
    asio::async_write(stream_, asio::buffer(requestContext_),
                      std::bind(&K8SListWatchSession::afterWatchRequest, shared_from_this(), std::placeholders::_1,
                                std::placeholders::_2));
}

void K8SListWatchSession::afterWatchRequest(asio::error_code ec, std::size_t) {
    if (ec) {
        return fail(ec, __FUNCTION__);
    }
    clearResponseState();
    doReadWatchResponse();
}

void K8SListWatchSession::doReadWatchResponse() {
    stream_.async_read_some(responseBuffer_.prepare(8192),
                            std::bind(&K8SListWatchSession::afterReadWatchResponse, shared_from_this(),
                                      std::placeholders::_1, std::placeholders::_2));
}

void K8SListWatchSession::afterReadWatchResponse(asio::error_code ec, std::size_t bytes_transferred) {
    if (ec) {
        return fail(ec, __FUNCTION__);
    }

    const char *responseData = asio::buffer_cast<const char *>(responseBuffer_.data());
    size_t parserSize = http_parser_execute(&responseParser_, &responseParserSetting_, responseData, bytes_transferred);
    assert(parserSize == bytes_transferred);
    responseBuffer_.consume(parserSize);

    if (!responseState_.headComplete_) {
        doReadWatchResponse();
        return;
    }

    assert(responseState_.headComplete_);

    constexpr unsigned int HTTP_OK = 200;

    if (responseState_.code_ != HTTP_OK) {
        if (!responseState_.messageComplete_) {
            doReadWatchResponse();
            return;
        }
        assert(responseState_.messageComplete_);
        LOG->error() << "K8S ListWatch Session Receive Unexpected Response , Program Will Exit : \r\n\t "
                     << responseState_.bodyBuffer_ << std::endl;
        LOG->flush();
        exit(-1);
    }

    const char *bodyData = responseState_.bodyBuffer_.c_str();
    size_t bodyDataLength = responseState_.bodyBuffer_.size();

    if (responseState_.bodyArrive_) {
        rapidjson::StringStream stream(bodyData);
        size_t lastTell = 0;
        while (lastTell < bodyDataLength) {
            rapidjson::Document document{};
            document.ParseStream<rapidjson::kParseStopWhenDoneFlag>(stream);
            if (document.HasParseError()) {
                break;
            }
            stream.Take();
            lastTell = stream.Tell();
            if (!handlerK8SWatchEvent(document)) {
                LOG->error() << "K8S ListWatch Session Receive Unexpected Response , Program Will Exit : \r\n\t "
                             << responseState_.bodyBuffer_ << std::endl;
                throw std::runtime_error("K8S ListWatch Session Receive Unexpected Response");
            }
        }
        assert(lastTell <= bodyDataLength);
        size_t remainingLength = bodyDataLength - lastTell;
        responseState_.bodyBuffer_.replace(0, remainingLength, bodyData + lastTell);
        responseState_.bodyBuffer_.resize(remainingLength);
        responseState_.bodyArrive_ = false;
    }

    responseState_.messageComplete_ ? doWatchRequest() : doReadWatchResponse();
}

bool K8SListWatchSession::handlerK8SWatchEvent(const rapidjson::Document &document) {
    assert(!document.HasParseError());

    auto pType = rapidjson::GetValueByPointer(document, "/type")->GetString();
    assert(pType != nullptr);

    constexpr char ADDEDTypeValue[] = "ADDED";
    constexpr char DELETETypeValue[] = "DELETED";
    constexpr char UPDATETypeValue[] = "MODIFIED";
    constexpr char BOOKMARKTypeValue[] = "BOOKMARK";
    constexpr char ERRORTypeValue[] = "ERROR";

    k8SEventTypeEnum what;

    if (strcmp(ADDEDTypeValue, pType) == 0) {
        what = K8SEventTypeAdded;
    } else if (strcmp(UPDATETypeValue, pType) == 0) {
        what = K8SEventTypeUpdate;
    } else if (strcmp(DELETETypeValue, pType) == 0) {
        what = K8SEventTypeDeleted;
    } else if (strcmp(BOOKMARKTypeValue, pType) == 0) {
        what = K8SEventTypeBookmark;
    } else if (strcmp(ERRORTypeValue, pType) == 0) {
        what = K8SEventTypeError;
    } else {
        throw std::runtime_error(pType);
    }

    if (what == K8SEventTypeError) {
        auto pErrorCode = rapidjson::GetValueByPointer(document, "/object/code");
        unsigned int errorCode = pErrorCode->GetUint();
        constexpr unsigned int HTTP_GONE = 410;
        if (errorCode == HTTP_GONE) {
            auto pMessageJson = rapidjson::GetValueByPointer(document, "/object/message");
            assert(pMessageJson != nullptr);
            auto pMessage = pMessageJson->GetString();
            assert(pMessage != nullptr);
            auto begin = strchr(pMessage, '(');
            if (begin != nullptr) {
                auto v = std::stoull(begin + 1, nullptr);
                handledBiggestResourceVersion_ = std::max(handledBiggestResourceVersion_, v);
            } else {
                handledBiggestResourceVersion_ += 200;
                //fixme 目前采用当前 handledBiggestResourceVersion_+=200 ,来解决 Gone 问题,
                // 但不是足够稳妥. 应该获取到资源对应的 --watch-cache-sizes 后再处理
            }
            return true;
        }
        return false;
    }

    auto pResourceVersion = rapidjson::GetValueByPointer(document, "/object/metadata/resourceVersion");
    assert(pResourceVersion != nullptr);
    assert(pResourceVersion->IsString());
    auto iResourceVersion = std::stoull(pResourceVersion->GetString(), nullptr, 10);

    if (iResourceVersion > handledBiggestResourceVersion_) {
        if (what != K8SEventTypeBookmark) {
            if (_callBack) {
                const auto &item = document["object"];
                _callBack(what, item);
            }
        }
        handledBiggestResourceVersion_ = iResourceVersion;
    }
    return true;
}