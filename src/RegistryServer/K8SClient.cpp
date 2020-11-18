
#include "servant/taf_logger.h"
#include "K8SClient.h"

class K8SClientSession {

public:
    explicit K8SClientSession(asio::io_context &ioContext) : sslContext_(
            K8SRuntimeParams::interface().sslContext()), stream_(ioContext, sslContext_) {
        memset(&responseParserSetting_, 0, sizeof(responseParserSetting_));
        responseParser_.data = &responseCompleteFlag_;
        responseParserSetting_.on_message_complete = [](http_parser *p) -> int {
            bool *completer = static_cast<bool *>(p->data);
            *completer = true;
            return 0;
        };
    }

    bool doRestfulRequest(const std::shared_ptr<std::string> &contentPtr) {
        requestContextPtr_ = contentPtr;
        return connectHost() && writeRequest() && clearResponseState() && readResponse();
    }

private:
    bool connectHost() {
        if (connected_) {
            return connected_;
        }
        asio::ip::tcp::endpoint endpoint(asio::ip::address::from_string(K8SRuntimeParams::interface().apiServerHost()),
                                         K8SRuntimeParams::interface().apiServerPort());
        stream_.next_layer().connect(endpoint, errorCode_);
        if (errorCode_) {
            LOG->error() << "Connect K8SApiServer Error : " << errorCode_.message() << std::endl;
            return false;
        }
        stream_.handshake(asio::ssl::stream_base::client, errorCode_);
        if (errorCode_) {
            LOG->error() << "Handshake K8SApiServer Error : " << errorCode_.message() << std::endl;
            return false;
        }
        connected_ = true;
        return connected_;
    }

    bool writeRequest() {
        asio::write(stream_, asio::buffer(*requestContextPtr_), errorCode_);
        if (errorCode_) {
            LOG->error() << "Write K8SApiServer Error" << errorCode_.message() << std::endl;
            return false;
        }
        return true;
    }

    bool readResponse() {
        while (true) {
            size_t bytes_transferred = stream_.read_some(responseBuffer_.prepare(4096), errorCode_);
            if (errorCode_) {
                LOG->error() << "Read K8SApiServer Error" << errorCode_.message() << std::endl;
                return false;
            }
            responseBuffer_.commit(bytes_transferred);
            const char *willParserData =
                    asio::buffer_cast<const char *>(responseBuffer_.data()) + responseParserLength_;
            size_t willParserLength = responseBuffer_.size() - responseParserLength_;
            size_t parserLength = http_parser_execute(&responseParser_, &responseParserSetting_, willParserData,
                                                      willParserLength);
            responseParserLength_ += parserLength;
            if (responseCompleteFlag_) {
                constexpr unsigned int OK_CODE = 200;       //请求得到[成功]响应
                constexpr unsigned int CREATED_CODE = 201;  //请求建立资源得到[成功]响应
                constexpr unsigned int ACCEPTED_CODE = 202;  //请求变更资源得到[成功]响应
                constexpr unsigned int CONFLICT_CODE = 409;  //请求建立资源得到[资源已存在]响应
                if (responseParser_.status_code == OK_CODE ||
                    responseParser_.status_code == CREATED_CODE ||
                    responseParser_.status_code == ACCEPTED_CODE ||
                    responseParser_.status_code == CONFLICT_CODE) {
                    return true;
                }
                const char *responseData = asio::buffer_cast<const char *>(responseBuffer_.data());
                const size_t responseLength = responseBuffer_.size();
                LOG->debug() << "Do Restful Request Ok , But Receive Unexpected Response :\r\n\t Request: "
                             << *requestContextPtr_ << "\r\n\tResponse : " << string(responseData, responseLength)
                             << endl;
                LOG->flush();
                return true;
            }
        }
    }

    bool clearResponseState() {
        responseBuffer_.consume(responseBuffer_.size());
        responseParserLength_ = 0;
        responseCompleteFlag_ = false;
        http_parser_init(&responseParser_, http_parser_type::HTTP_RESPONSE);
        return true;
    }

private:
    asio::ssl::context &sslContext_;
    asio::ssl::stream<asio::ip::tcp::socket> stream_;
    bool connected_{false};
    std::shared_ptr<std::string> requestContextPtr_{};
    asio::streambuf responseBuffer_{};
    http_parser responseParser_{};
    http_parser_settings responseParserSetting_{};
    bool responseCompleteFlag_{false};
    size_t responseParserLength_{};
    asio::error_code errorCode_{};
};

void K8SClient::postTask(K8SClientRequestType type, const std::string &url, const std::string &body) {
    std::string content{};
    switch (type) {
        case Patch: {
            content = buildPatchContent(url, body);
        }
            break;
        case StrategicMergePatch: {
            content = buildSMPatchContent(url, body);
        }
            break;
        case Post: {
            content = buildPostContent(url, body);
        }
            break;
        case Delete: {
            content = buildDeleteContent(url);
        }
            break;
        default:
            return;
    }

    auto contentPtr = std::make_shared<std::string>(std::move(content));
    ioContext_.post(
            [this, contentPtr] {
                bool workSuccess = sessionPtr_->doRestfulRequest(contentPtr);
                if (!workSuccess) {
                    // fixme 如果失败，则重建链接进行尝试. 重试有极低概率的任务堆积过多风险,暂时忽略.
                    constexpr size_t RetryTimes = 8;
                    constexpr size_t RetryIntervalTime = 2;
                    for (size_t i = 0; i < RetryTimes; ++i) {
                        sleep(i * RetryIntervalTime + 1);
                        LOG->error() << "Retry Connect To K8SApiServer , Number Of Retries : " << i << endl;
                        sessionPtr_ = std::make_shared<K8SClientSession>(ioContext_);
                        workSuccess = sessionPtr_->doRestfulRequest(contentPtr);
                        if (workSuccess) {
                            break;
                        }
                    }
                    if (!workSuccess) {
                        LOG->error() << "Retry Connect To K8SApiServer Failed , Program Will Exit " << endl;
                        LOG->flush();
                        exit(-1);
                    }
                }
            }
    );
}

void K8SClient::run() {
    sessionPtr_ = std::make_shared<K8SClientSession>(ioContext_);
    asio::io_context::work work(ioContext_);
    ioContext_.run();
}

std::string K8SClient::buildPostContent(const std::string &url, const std::string &body) {
    std::ostringstream strStream;
    strStream << "POST " << url << " HTTP/1.1\r\n";
    strStream << "Authorization: Bearer " << K8SRuntimeParams::interface().bindToken() << "\r\n";
    strStream << "Host: " << K8SRuntimeParams::interface().apiServerHost() << ":"
              << K8SRuntimeParams::interface().apiServerPort() << "\r\n";
    strStream << "Content-Length: " << body.size() << "\r\n";
    strStream << "Content-Type: application/json\r\n";
    strStream << "Connection: Keep-alive\r\n";
    strStream << "\r\n";
    strStream << body;
    std::string reqContext = strStream.str();
    return reqContext;
}

std::string K8SClient::buildPatchContent(const std::string &url, const std::string &body) {
    std::ostringstream strStream;
    strStream << "PATCH " << url << " HTTP/1.1\r\n";
    strStream << "Authorization: Bearer " << K8SRuntimeParams::interface().bindToken() << "\r\n";
    strStream << "Host: " << K8SRuntimeParams::interface().apiServerHost() << ":"
              << K8SRuntimeParams::interface().apiServerPort() << "\r\n";
    strStream << "Content-Length: " << body.size() << "\r\n";
    strStream << "Content-Type: application/json-patch+json\r\n";
    strStream << "Connection: Keep-alive\r\n";
    strStream << "\r\n";
    strStream << body;
    std::string reqContext = strStream.str();
    return reqContext;
}

std::string K8SClient::buildDeleteContent(const std::string &url) {
    std::ostringstream strStream;
    strStream << "DELETE " << url << " HTTP/1.1\r\n";
    strStream << "Authorization: Bearer " << K8SRuntimeParams::interface().bindToken() << "\r\n";
    strStream << "Host: " << K8SRuntimeParams::interface().apiServerHost() << ":"
              << K8SRuntimeParams::interface().apiServerPort() << "\r\n";
    strStream << "Connection: Keep-alive\r\n";
    strStream << "\r\n";
    std::string reqContext = strStream.str();
    return reqContext;
}

std::string K8SClient::buildSMPatchContent(const string &url, const string &body) {
    std::ostringstream strStream;
    strStream << "PATCH " << url << " HTTP/1.1\r\n";
    strStream << "Authorization: Bearer " << K8SRuntimeParams::interface().bindToken() << "\r\n";
    strStream << "Host: " << K8SRuntimeParams::interface().apiServerHost() << ":"
              << K8SRuntimeParams::interface().apiServerPort() << "\r\n";
    strStream << "Content-Length: " << body.size() << "\r\n";
    strStream << "Content-Type: application/strategic-merge-patch+json\r\n";
    strStream << "Connection: Keep-alive\r\n";
    strStream << "\r\n";
    strStream << body;
    std::string reqContext = strStream.str();
    return reqContext;
}
