#!/usr/bin/env bash

echo "-----------------------begin make tars----------------------"

cmake -E make_directory build_release
cmake -E chdir build_release cmake .. -DCMAKE_BUILD_TYPE=Release -DExt=""
cmake -E chdir build_release cmake --build . --config release --target install -- -j 8

echo "-----------------------make sdk success----------------------"