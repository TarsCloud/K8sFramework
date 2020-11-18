#!/usr/bin/env bash

###### LOG 函数
function LOG_ERROR() {
  msg=$(date +%Y-%m-%d" "%H:%M:%S)
  msg="${msg} $*"
  echo -e "\033[31m $msg \033[0m"
}

function LOG_WARNING() {
  msg=$(date +%Y-%m-%d" "%H:%M:%S)
  msg="${msg} $*"
  echo -e "\033[33m $msg \033[0m"
}

function LOG_INFO() {
  echo -e "\033[31m $* \033[0m"
}

function LOG_DEBUG() {
  msg=$(date +%Y-%m-%d" "%H:%M:%S)
  msg="${msg} $*"
  echo -e "\033[40;37m $msg \033[0m"
}
###### LOG 函数

if (($# < 3)); then
  LOG_INFO "Usage:  push.sh <DOCKER_REGISTRY_URL> <DOCKER_REGISTRY_USER> <DOCKER_REGISTRY_PASSWORD>"
  LOG_INFO "Example: push.sh dockerhub.com/tafk8s tafk8s tafk8s@image"
  exit 255
fi

#仓库信息
_DOCKER_REGISTRY_URL_=$1
_DOCKER_REGISTRY_USER_=$2
_DOCKER_REGISTRY_PASSWORD_=$3
#

#### 构建基础镜像
declare -a BaseImages=(
  taf.base
  taf.cppbase
  taf.javabase
  taf.node6base
  taf.node8base
  taf.node10base
  helm.wait
)

for KEY in "${BaseImages[@]}"; do
  if ! docker build -t "${KEY}" -f build/"${KEY}.Dockerfile" build; then
    LOG_ERROR "Build ${KEY} image failed"
    exit 255
  fi
done
#### 构建基础镜像

#### 构建基础服务镜像
declare -a FrameworkImages=(
  taf.tafcontrol
  taf.tafregistry
  taf.tafagent
  taf.tafimage
  taf.tafadmin
  taf.tafweb
)

for KEY in "${FrameworkImages[@]}"; do
  if ! docker build -t "${KEY}" -f build/"${KEY}.Dockerfile" build; then
    LOG_ERROR "Build ${KEY} image failed"
    exit 255
  fi
done
#### 构建基础服务镜像

#### 构建基础服务镜像
declare -a ServerImages=(
  taflog
  tafconfig
  tafnotify
  tafstat
  tafquerystat
  tafproperty
  tafqueryproperty
)

for KEY in "${ServerImages[@]}"; do
  mkdir -p build/files/template/taf."${KEY}"
  mkdir -p build/files/template/taf."${KEY}"/root/etc
  mkdir -p build/files/template/taf."${KEY}"/root/usr/local/server/bin

  if ! cp build/files/binary/"${KEY}" build/files/template/taf."${KEY}"/root/usr/local/server/bin/"${KEY}"; then
    LOG_ERROR "copy ${KEY} failed, please check ${KEY} is in directory: build/files/binary"
    exit 255
  fi

  echo "FROM taf.cppbase
COPY /root /
" >build/files/template/taf."${KEY}"/Dockerfile

  echo "#!/bin/bash
export ServerName=\"${KEY}\"
export ServerType=\"taf_cpp\"
export BuildPerson=\"admin\"
export BuildTime=\"$(date)\"
" >build/files/template/taf."${KEY}"/root/etc/detail

  if ! docker build -t taf."${KEY}" build/files/template/taf."${KEY}"; then
    LOG_ERROR "Build ${KEY} image failed"
    exit 255
  fi
done

#### 构建基础服务镜像

LOG_INFO "Build All Images Ok"

declare -a LocalImages=(
  taf.base
  taf.cppbase
  taf.javabase
  taf.node6base
  taf.node8base
  taf.node10base
  taf.tafcontrol
  taf.tafregistry
  taf.tafagent
  taf.tafimage
  taf.tafadmin
  taf.tafweb
  taf.taflog
  taf.tafconfig
  taf.tafnotify
  taf.tafstat
  taf.tafquerystat
  taf.tafproperty
  taf.tafqueryproperty
  helm.wait
)

# 登陆
if ! docker login -u "${_DOCKER_REGISTRY_USER_}" -p "${_DOCKER_REGISTRY_PASSWORD_}" "${_DOCKER_REGISTRY_URL_}"; then
  LOG_ERROR "docker login to ${_DOCKER_REGISTRY_URL_} failed!"
  exit 255
fi

MYSQL_IMAGE="mysql:5.7.24"

if ! docker pull "${MYSQL_IMAGE}"; then
  LOG_ERROR "docker pull \"${MYSQL_IMAGE}\" failed!"
  exit 255
fi

if ! docker tag "${MYSQL_IMAGE}" "${_DOCKER_REGISTRY_URL_}"/mysqlclient; then
  LOG_ERROR "docker pull \"${MYSQL_IMAGE}\" failed!"
  exit 255
fi

if ! docker push "${_DOCKER_REGISTRY_URL_}"/mysqlclient; then
  LOG_ERROR "docker pull \"${_DOCKER_REGISTRY_URL_}/mysqlclient\" failed!"
  exit 255
fi

for KEY in "${LocalImages[@]}"; do
  RemoteImagesTag="${_DOCKER_REGISTRY_URL_}"/"${KEY}":10000
  if ! docker tag "${KEY}" "${RemoteImagesTag}"; then
    LOG_ERROR "Tag ${KEY} image failed"
    exit 255
  fi

  if ! docker push "${RemoteImagesTag}"; then
    LOG_ERROR "Push ${RemoteImagesTag} image failed"
    exit 255
  fi
done

LOG_INFO "Push All Images Ok"
