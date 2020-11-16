#!/bin/bash

# PodName   From K8S Downward API
# 如果 Pod的网络模式为 HostNetwork ,则 PodName与 HostName 会不一致。为避免错误,PodName由 K8S Downward API 注入
_K8S_POD_NAME_=${PodName}

if [ -z "$_K8S_POD_NAME_" ]; then
  echo "未获取到PodName"
  exit 255
fi

_K8S_POD_IP_=$(awk '{if ($2~"'"$HOSTNAME"'") print $1}' /etc/hosts)

if [ -z "$_K8S_POD_IP_" ]; then
  echo "未获取到PodIP"
  exit 255
fi

export HOSTNAME=${_K8S_POD_NAME_}
# 重设 $HOSTNAME 值是因为tars框架会使用 ::getenv("HOSTNAME") 来获取PodName 值

export PodIp="${_K8S_POD_IP_}"
# 设置 PodIp  值是因为 tarsnode 程序 会使用 ::getenv("PodIp") 来获取 PodIp 值

TAFNODE_EXECUTION_FILE="/usr/local/app/tars/tarsnode/bin/tarsnode"

TAFNODE_CONFIG_FILE="/usr/local/app/tars/tarsnode/conf/tarsnode.conf"

TAFNODE_DATA_DIR="/usr/local/app/tars/tarsnode/data"

IMAGE_BIND_SERVER_DIR="/usr/local/server/bin"

IMAGE_BIND_LOG_DIR="/usr/local/app/tars/app_log"

declare -a ReplaceKeyList=(
  "_K8S_POD_IP_"
)

declare -a ReplaceFileList=(
  "${TAFNODE_CONFIG_FILE}"
)

for KEY in "${ReplaceKeyList[@]}"; do
  for FILE in "${ReplaceFileList[@]}"; do
    if ! sed -i "s#${KEY}#${!KEY}#g" "${FILE}"; then
      exit 255
    fi
  done
done

# ServerApp From K8S Downward API
_TAF_SERVER_APP_=${ServerApp}
if [ -z "$_TAF_SERVER_APP_" ]; then
  echo "未获取到SererApp"
  exit 255
fi

# ServerName From /etc/detail
source "/etc/detail"
_TAF_SERVER_NAME_=${ServerName}
if [ -z "$_TAF_SERVER_NAME_" ]; then
  echo "未获取到ServerName"
  exit 255
fi

ServerBaseDir=${TAFNODE_DATA_DIR}/"${_TAF_SERVER_APP_}"."${_TAF_SERVER_NAME_}"
export ServerBinDir=${ServerBaseDir}/bin
export ServerDataDir=${ServerBaseDir}/data
export ServerConfDir=${ServerBaseDir}/conf
export ServerConfFile=${ServerConfDir}/${_TAF_SERVER_APP_}"."${_TAF_SERVER_NAME_}.config.conf
export ServerLauncherFile="${ServerBinDir}/${_TAF_SERVER_NAME_}"
export ServerLauncherArgv="${_TAF_SERVER_NAME_} --config=${ServerConfFile}"
export ServerLogDir=${IMAGE_BIND_LOG_DIR}

mkdir -p /host-log-path/"${_K8S_POD_NAME_}"
ln -sf /host-log-path/"${_K8S_POD_NAME_}/" ${IMAGE_BIND_LOG_DIR}

mkdir -p "${ServerBaseDir}"
ln -s ${IMAGE_BIND_SERVER_DIR} "${ServerBinDir}"
mkdir -p "${ServerDataDir}"
mkdir -p "${ServerConfDir}"

export LD_LIBRARY_PATH=${ServerBinDir}:${ServerBinDir}/lib

ldconfig
exec ${TAFNODE_EXECUTION_FILE} --config=${TAFNODE_CONFIG_FILE}
#exec /bin/guard