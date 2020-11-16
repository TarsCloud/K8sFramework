#include "RegistryServer.h"
#include <iostream>

RegistryServer g_app;

int main(int argc, char *argv[]) {
    try {
        g_app.main(argc, argv);
        g_app.waitForShutdown();
    }
    catch (exception &ex) {
        cerr << ex.what() << endl;
    }
    return 0;
}

//todo 将数据库操作移动到独立线程