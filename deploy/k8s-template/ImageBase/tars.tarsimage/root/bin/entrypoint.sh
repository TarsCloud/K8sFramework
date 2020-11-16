#! /bin/bash

mkdir -p /buildDir
mkdir -p /uploadDir

TAFIMAGE_EXECUTION_FILE="/usr/local/app/tars/tarsimage/bin/tarsimage"
exec ${TAFIMAGE_EXECUTION_FILE}
