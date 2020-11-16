#!/bin/bash

_K8S_POD_NAME_=$(cat /etc/hostname)

_K8S_POD_IP_=$(awk '{if ($2~"'"$_K8S_POD_NAME_"'") print $1}' /etc/hosts)

REGISTRY_EXECUTION_FILE=/usr/local/app/tars/tarsregistry/bin/tarsregistry

REGISTRY_CONFIG_FILE=/usr/local/app/tars/tarsregistry/conf/tarsregistry.conf

declare -a ReplaceKeyList=(
  _K8S_POD_IP_
  _DB_HOST_
  _DB_PORT_
  _DB_NAME_
  _DB_USER_
  _DB_PASSWORD_
  _EXTERNAL_JCEPROXY_
)

declare -a ReplaceFileList=(
  "${REGISTRY_CONFIG_FILE}"
)

for KEY in "${ReplaceKeyList[@]}"; do
  for FILE in "${ReplaceFileList[@]}"; do
    sed -i "s#${KEY}#${!KEY}#g" "${FILE}"
    if [[ 0 -ne $? ]]; then
      exit 255
    fi
  done
done

ldconfig

exec ${REGISTRY_EXECUTION_FILE} --config=${REGISTRY_CONFIG_FILE}
# exec /bin/guard
