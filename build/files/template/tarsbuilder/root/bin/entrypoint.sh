#!/usr/bin/env bash

cd /tars-src/TarsFramework || exit 255

sh ../bootstrap.sh || exit 255

cp build_release/servers/tarslog/bin/tarslog /tars-k8s-binary || exit 255

cp build_release/servers/tarsstat/bin/tarsstat /tars-k8s-binary || exit 255

cp build_release/servers/tarsquerystat/bin/tarsquerystat /tars-k8s-binary || exit 255

cp build_release/servers/tarsproperty/bin/tarsproperty /tars-k8s-binary || exit 255
cp build_release/servers/tarsqueryproperty/bin/tarsqueryproperty /tars-k8s-binary || exit 255

cd ../src || exit 255
mkdir -p build_tmp && rm -rf build_tmp/* || exit 255
cd build_tmp || exit 255
cmake .. || exit 255
make -j4 || exit 255

cp bin/tarsadmin /tars-k8s-binary || exit 255
cp bin/tarsregistry /tars-k8s-binary || exit 255
cp bin/tarsnotify /tars-k8s-binary || exit 255
cp bin/tarsconfig /tars-k8s-binary || exit 255
cp bin/tarscontrol /tars-k8s-binary || exit 255
cp bin/tarsnode /tars-k8s-binary || exit 255
cp bin/tarsagent /tars-k8s-binary || exit 255
cp bin/tarsimage /tars-k8s-binary || exit 255