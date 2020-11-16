#!/bin/bash

_K8S_POD_NAME_=${PodName}

_K8S_POD_IP_=$(awk '{if ($2~"'"$HOSTNAME"'") print $1}' /etc/hosts)

TAFNODE_EXECUTION_FILE="/usr/local/app/tars/tarsnode/bin/tarsnode"

TAFNODE_CONFIG_FILE="/usr/local/app/tars/tarsnode/conf/tarsnode.conf"

TAFNODE_DATA_DIR="/usr/local/app/tars/tarsnode/data"

IMAGE_BIND_SERVER_DIR="/usr/local/server/bin"

IMAGE_BIND_LOG_DIR="/usr/local/app/tars/app_log"

NODE_EXECUTION_FILE="/usr/bin/node"

NODE_AGENT_DIR="/usr/local/app/tars/tars-node-agent/bin/tars-node-agent"

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

export PodIp="${_K8S_POD_IP_}"

source "/etc/detail"

ServerBaseDir=${TAFNODE_DATA_DIR}/"${ServerApp}"."${ServerName}"

export ServerBinDir=${ServerBaseDir}/bin
export ServerDataDir=${ServerBaseDir}/data
export ServerConfDir=${ServerBaseDir}/conf
export ServerConfFile=${ServerConfDir}/${ServerApp}"."${ServerName}.config.conf
export ServerLauncherFile=${NODE_EXECUTION_FILE}
export ServerLauncherArgv="node ${NODE_AGENT_DIR} ${ServerBinDir}/ -c ${ServerConfFile}"
export ServerLogDir=${IMAGE_BIND_LOG_DIR}

mkdir -p /host-log-path/"${_K8S_POD_NAME_}"
ln -sf /host-log-path/"${_K8S_POD_NAME_}/" ${IMAGE_BIND_LOG_DIR}


mkdir -p "${ServerBaseDir}"
ln -s ${IMAGE_BIND_SERVER_DIR} "${ServerBinDir}"
mkdir -p "${ServerDataDir}"
mkdir -p "${ServerConfDir}"

export NODE_PATH="/usr/local/app/tars/tarsnode/lib/node_modules"

exec ${TAFNODE_EXECUTION_FILE} --config=${TAFNODE_CONFIG_FILE}