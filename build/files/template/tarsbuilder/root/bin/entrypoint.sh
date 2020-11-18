#!/usr/bin/env bash

cd /tars-k8s-src/TafCpp || exit 255

sh ./bootstrap.sh || exit 255

cp build_release/tars-cpp-framework/servers/tarslog/bin/tarslog /tars-k8s-binary || exit 255

cp build_release/tars-cpp-framework/servers/tarsstat/bin/tarsstat /tars-k8s-binary || exit 255

cp build_release/tars-cpp-framework/servers/tarsquerystat/bin/tarsquerystat /tars-k8s-binary || exit 255

cp build_release/tars-cpp-framework/servers/tarsproperty/bin/tarsproperty /tars-k8s-binary || exit 255
cp build_release/tars-cpp-framework/servers/tarsqueryproperty/bin/tarsqueryproperty /tars-k8s-binary || exit 255

cd /tars-k8s-src || exit 255
cmake . || exit 255
make -j4 || exit 255

cp bin/tarsadmin /tars-k8s-binary || exit 255
cp bin/tarsregistry /tars-k8s-binary || exit 255
cp bin/tarsnotify /tars-k8s-binary || exit 255
cp bin/tarsconfig /tars-k8s-binary || exit 255
cp bin/tarscontrol /tars-k8s-binary || exit 255
cp bin/tarsnode /tars-k8s-binary || exit 255
cp bin/tarsagent /tars-k8s-binary || exit 255
cp bin/tarsimage /tars-k8s-binary || exit 255