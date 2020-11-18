

#pragma once

#include "servant/Application.h"

using namespace taf;

/**
 *  Registry Server
 */
class RegistryServer : public Application {
protected:
    /**
     * 初始化, 只会进程调用一次
     */
    void initialize() override;

    /**
     * 析构, 每个进程都会调用一次
     */
    void destroyApp() override;

private:
    std::thread _ioThread;
    std::thread _k8sClientThread;
};


