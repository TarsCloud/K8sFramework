#!/usr/bin/env bash

# PodName   From K8S Downward API
_K8S_POD_NAME_=${PodName}

if [ -z "$_K8S_POD_NAME_" ]; then
  echo "未获取到PodName"
  exit 255
fi

# PodIP   From K8S Downward API
_K8S_POD_IP_=${PodIP}

if [ -z "$_K8S_POD_IP_" ]; then
  echo "未获取到PodIP"
  exit 255
fi

export HOSTNAME=${_K8S_POD_NAME_}
# 重设 $HOSTNAME 值是因为tars框架会使用 ::getenv("HOSTNAME") 来获取PodName 值

TARSNODE_EXECUTION_FILE="/usr/local/app/tars/tarsnode/bin/tarsnode"

TARSNODE_CONFIG_FILE="/usr/local/app/tars/tarsnode/conf/tarsnode.conf"

TARSNODE_DATA_DIR="/usr/local/app/tars/tarsnode/data"

IMAGE_BIND_SERVER_DIR="/usr/local/server/bin"

IMAGE_BIND_LOG_DIR="/usr/local/app/tars/app_log"

declare -a ReplaceKeyList=(
  "_K8S_POD_IP_"
)

declare -a ReplaceFileList=(
  "${TARSNODE_CONFIG_FILE}"
)

for KEY in "${ReplaceKeyList[@]}"; do
  for FILE in "${ReplaceFileList[@]}"; do
    if ! sed -i "s#${KEY}#${!KEY}#g" "${FILE}"; then
      exit 255
    fi
  done
done

# ServerApp From K8S Downward API
_TARS_SERVER_APP_=${ServerApp}
if [ -z "$_TARS_SERVER_APP_" ]; then
  echo "未获取到SererApp"
  exit 255
fi

# ServerName From /etc/detail
source "/etc/detail"
_TARS_SERVER_NAME_=${ServerName}
if [ -z "$_TARS_SERVER_NAME_" ]; then
  echo "未获取到ServerName"
  exit 255
fi

_TARS_SERVER_TYPE_=${ServerType}
if [ -z "$_TARS_SERVER_TYPE_" ]; then
  echo "未获取到ServerType"
  exit 255
fi

ServerBaseDir=${TARSNODE_DATA_DIR}/"${_TARS_SERVER_APP_}"."${_TARS_SERVER_NAME_}"
export ServerBinDir=${ServerBaseDir}/bin
export ServerDataDir=${ServerBaseDir}/data
export ServerConfDir=${ServerBaseDir}/conf
export ServerConfFile=${ServerConfDir}/${_TARS_SERVER_APP_}"."${_TARS_SERVER_NAME_}.config.conf
export ServerLogDir=${IMAGE_BIND_LOG_DIR}

mkdir -p /host-log-path/"${_K8S_POD_NAME_}"
ln -sf /host-log-path/"${_K8S_POD_NAME_}/" ${IMAGE_BIND_LOG_DIR}

mkdir -p "${ServerBaseDir}"
ln -s ${IMAGE_BIND_SERVER_DIR} "${ServerBinDir}"
mkdir -p "${ServerDataDir}"
mkdir -p "${ServerConfDir}"

case ${ServerType} in
"tars_cpp" | "tars_node_pkg" | "tars_go")
  export LD_LIBRARY_PATH=${ServerBinDir}:${ServerBinDir}/lib
  export ServerLauncherFile="${ServerBinDir}/${_TARS_SERVER_NAME_}"
  export ServerLauncherArgv="${_TARS_SERVER_NAME_} --config=${ServerConfFile}"
  ;;
"tars_node" | "tars_node8" | "tars_node10" | "tars_node12")
  NODE_EXECUTION_FILE="/usr/local/bin/node"
  NODE_AGENT_BIN="/usr/local/app/tars/tars-node-agent/bin/tars-node-agent"
  export ServerLauncherFile=${NODE_EXECUTION_FILE}
  export ServerLauncherArgv="node ${NODE_AGENT_BIN} ${ServerBinDir}/ -c ${ServerConfFile}"
  ;;
"tars_java_war")
  JAVA_EXECUTION_FILE="/usr/local/openjdk-8/bin/java"
  export ServerLauncherFile=${JAVA_EXECUTION_FILE}
  export ServerLauncherArgv="java -Dconfig=${ServerConfFile} #{jvmparams} -cp #{classpath} #{mainclass}"
  ;;
"tars_java_jar")
  JAVA_EXECUTION_FILE="/usr/local/openjdk-8/bin/java"
  export ServerLauncherFile=${JAVA_EXECUTION_FILE}
  export ServerLauncherArgv="java -Dconfig=${ServerConfFile} #{jvmparams} -jar ${ServerBinDir}/${ServerName}.jar"
  ;;
esac

exec ${TARSNODE_EXECUTION_FILE} --config=${TARSNODE_CONFIG_FILE}
#exec /bin/guard
