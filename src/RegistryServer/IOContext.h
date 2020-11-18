
#pragma once

#include <asio/io_context.hpp>
#include <memory>

class IOContext {
public:
    static IOContext &instance() {
        static IOContext watcher{};
        return watcher;
    }

    inline void run() {
        asio::io_context::work work(ioContext);
        ioContext.run();
    };

    inline asio::io_context &getIOContext() { return ioContext; }

private:
    IOContext() = default;

private:
    asio::io_context ioContext{1};
};