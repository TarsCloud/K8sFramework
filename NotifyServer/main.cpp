#include "NotifyServer.h"
#include <iostream>

using namespace taf;

// TC_Config *g_pconf;

NotifyServer g_app;

int main(int argc, char *argv[]) {
    try {

        g_app.main(argc, argv);

        // g_pconf = &NotifyServer::getConfig();
        g_app.waitForShutdown();
    }
    catch (exception &ex) {
        cerr << ex.what() << endl;
    }

    return 0;
}


