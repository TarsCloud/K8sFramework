#include "NotifyServer.h"
#include <iostream>

NotifyServer g_app;

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


