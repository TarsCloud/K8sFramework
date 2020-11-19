#!/bin/bash

mkdir -p openapi
docker run --rm -it -e GOPATH=$GOPATH:/go -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger generate server --exclude-main -A tarsadmin/openapi -t openapi -f ./doc/Admin.yaml
cd openapi
go mod init tarsadmin/openapi