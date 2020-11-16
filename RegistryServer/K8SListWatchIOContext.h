
#pragma once

#include <asio/io_context.hpp>
#include <memory>

class K8SListWatchIOContext {
public:
    static K8SListWatchIOContext &instance() {
        static K8SListWatchIOContext watcher{};
        return watcher;
    }

    inline void run() {
        asio::io_context::work work(ioContext);
        ioContext.run();
    };

    inline asio::io_context &getIOContext() { return ioContext; }

private:
    K8SListWatchIOContext() = default;

private:
    asio::io_context ioContext{1};
};