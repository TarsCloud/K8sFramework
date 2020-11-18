#!/bin/bash

############################################

# 填写 docker 镜像仓库信息
_DOCKER_REGISTRY_URL_=
_DOCKER_REGISTRY_USER_=
_DOCKER_REGISTRY_PASSWORD_=
#

# 填写 tars 框架 数据库信息
_DB_TARS_HOST_=
_DB_TARS_PORT_=
_DB_TARS_DATABASE_=
_DB_TARS_USER_=
_DB_TARS_PASSWORD_=

# 填写 存储 stat 数据库信息 
_DB_TARS_STAT_HOST_=
_DB_TARS_STAT_DATABASE_= k8s_tars_stat  #发现部分代码对此数据库名耦合,建议保留
_DB_TARS_STAT_USER_=
_DB_TARS_STAT_PASSWORD_=
_DB_TARS_STAT_PORT_=

# 填写 存储 property 数据库信息
_DB_TARS_PROPERTY_HOST_=
_DB_TARS_PROPERTY_DATABASE_= k8s_tars_property #发现部分代码对此数据库名耦合,建议保留
_DB_TARS_PROPERTY_USER_=
_DB_TARS_PROPERTY_PASSWORD_=
_DB_TARS_PROPERTY_PORT_=


# 填写 部署 tars 程序的 k8s node ,建议两台或以上
declare -a DeployNode=(
  #kube.node117
  #kube.node118
  #kube.node119
)

##############################################

# 以下不用填写
declare -a AllConfig=(
  _DOCKER_REGISTRY_URL_
  _DOCKER_REGISTRY_USER_
  _DOCKER_REGISTRY_PASSWORD_

  #
  _DB_TARS_HOST_
  _DB_TARS_PORT_
  _DB_TARS_DATABASE_
  _DB_TARS_USER_
  _DB_TARS_PASSWORD_

  #
  _DB_TARS_STAT_HOST_
  _DB_TARS_STAT_DATABASE_
  _DB_TARS_STAT_USER_
  _DB_TARS_STAT_PASSWORD_
  _DB_TARS_STAT_PORT_

  #
  _DB_TARS_PROPERTY_HOST_
  _DB_TARS_PROPERTY_DATABASE_
  _DB_TARS_PROPERTY_USER_
  _DB_TARS_PROPERTY_PASSWORD_
  _DB_TARS_PROPERTY_PORT_

)

declare -a AllSqlFile=(
  # MySql/db_property.sql
  # MySql/db_stat.sql
  MySql/db_tars.sql
)

declare -a AllYamlFile=(
  Yaml/tarsadmin.yaml
  Yaml/tarsconfig.yaml
  Yaml/tarscommon.yaml
  Yaml/tarsimage.yaml
  Yaml/tarslog.yaml
  Yaml/tarsnotify.yaml
  Yaml/tarsproperty.yaml
  Yaml/tarsqueryproperty.yaml
  Yaml/tarsquerystat.yaml
  Yaml/tarsregistry.yaml
  Yaml/tarsstat.yaml
  Yaml/tarsweb.yaml
)

rm -rf Install && cp -r Template Install && cd Install || exit

for KEY in "${AllConfig[@]}"; do

  # sed -i时，第一个参数mac必须，linux可选，加上后兼容2个平台

  for SqlFile in "${AllSqlFile[@]}"; do
    if ! sed -i "" "s#${KEY}#${!KEY}#g" "${SqlFile}"; then
      exit 255
    fi
  done

  for YamlFile in "${AllYamlFile[@]}"; do
    if ! sed -i "" "s#${KEY}#${!KEY}#g" "${YamlFile}"; then
      exit 255
    fi
  done
done

function Record_Error_Info() {
  echo -e "\e[31m $* \e[0m"
}

function CheckCommandExist() {
  for i in "$@"; do
    type "$i" >/dev/null 2>&1 || {
      Record_Error_Info "未检测到 $i 程序, 请安装好 $i 程序之后重试"
      exit 255
    }
  done
}

## 检测参数
#
## 检测环境
CheckCommandExist kubectl docker mysql
#
#初始化 MySQL
if ! mysql -h${_DB_TARS_HOST_} -P${_DB_TARS_PORT_} -u${_DB_TARS_USER_} -p${_DB_TARS_PASSWORD_} <MySql/db_tars.sql; then
  Record_Error_Info "备份数据到 db_tars 失败, 请检测 db_tars 相关参数是否设置正确,或者网络是否正常"
  exit 255
fi

# 生成镜像
if ! docker login -u ${_DOCKER_REGISTRY_USER_} -p ${_DOCKER_REGISTRY_PASSWORD_} ${_DOCKER_REGISTRY_URL_}; then
  Record_Error_Info "登录镜像仓库失败,请检测 Docker仓库相关参数是否设置正确,或者网络是否正常"
  exit 255
fi

declare -a BaseImage=(
  tars.base
  tars.cppbase
  tars.javabase
  tars.nodebase
  tars.node8base
  tars.node10base
  tars.tarsimage
  tars.tarsadmin
  tars.tarsagent
  tars.tarsregistry
  tars.tarsweb
)

if ! cp Program/tarsnode ImageBase/tars.cppbase/root/usr/local/app/tars/tarsnode/bin/tarsnode; then
  Record_Error_Info "拷贝程序 tarsnode 失败， 请检测 tarsnode 是否放置在 Program 目录"
  exit 255
fi
chmod a+x ImageBase/tars.cppbase/root/usr/local/app/tars/tarsnode/bin/tarsnode

if ! cp Program/tarsregistry ImageBase/tars.tarsregistry/root/usr/local/app/tars/tarsregistry/bin; then
  Record_Error_Info "拷贝程序 tarsregistry 失败， 请检测 tarsregistry 是否放置在 Program 目录"
  chmod +x ImageBase/tars.tarsregistry/root/usr/local/app/tars/tarsregistry/bin/tarsregistry
  exit 255
fi

if ! cp Program/tarsimage ImageBase/tars.tarsimage/root/usr/local/app/tars/tarsimage/bin; then
  Record_Error_Info "拷贝程序 tarsimage 失败， 请检测 tarsimage 是否放置在 Program 目录"
  exit 255
fi
chmod a+x ImageBase/tars.tarsimage/root/usr/local/app/tars/tarsimage/bin/tarsimage

if ! cp Program/tarsadmin ImageBase/tars.tarsadmin/root/usr/local/app/tars/tarsadmin/bin; then
  Record_Error_Info "拷贝程序 tarsadmin 失败， 请检测 tarsadmin 是否放置在 Program 目录"
  exit 255
fi
chmod a+x ImageBase/tars.tarsadmin/root/usr/local/app/tars/tarsadmin/bin/tarsadmin

if ! cp Program/tarsagent ImageBase/tars.tarsagent/root/usr/local/app/tars/tarsagent/bin; then
  Record_Error_Info "拷贝程序 tarsagent 失败， 请检测 tarsagent 是否放置在 Program 目录"
  exit 255
fi
chmod a+x ImageBase/tars.tarsagent/root/usr/local/app/tars/tarsagent/bin/tarsagent

for KEY in "${BaseImage[@]}"; do
  if ! docker build -t "${KEY}" ImageBase/"${KEY}"; then
    Record_Error_Info "Build ${KEY} 镜像失败"
    exit 255
  fi

  if ! docker tag "${KEY}" ${_DOCKER_REGISTRY_URL_}/"${KEY}":10000; then
    Record_Error_Info "Tag ${KEY} 镜像失败"
    exit 255
  fi

  if ! docker push ${_DOCKER_REGISTRY_URL_}/"${KEY}":10000; then
    Record_Error_Info "Push ${KEY} 镜像失败"
    exit 255
  fi
done

declare -a ServerImage=(
  tarsnotify
  tarslog
  tarsconfig
  tarsstat
  tarsquerystat
  tarsproperty
  tarsqueryproperty
)

for KEY in "${ServerImage[@]}"; do
  mkdir -p ImageBase/tars."${KEY}"
  mkdir -p ImageBase/tars."${KEY}"/root/etc
  mkdir -p ImageBase/tars."${KEY}"/root/usr/local/server/bin
  if ! cp Program/"${KEY}" ImageBase/tars."${KEY}"/root/usr/local/server/bin; then
    Record_Error_Info "拷贝程序 ${KEY} 失败， 请检测 ${KEY} 是否放置在 Program 目录"
    exit 255
  fi

  echo "FROM tars.cppbase
COPY /root /
CMD [\"/bin/entrypoint\"]
" >ImageBase/tars."${KEY}"/Dockerfile

  echo "#!/bin/bash
export ServerName=\"${KEY}\"
export ServerType=\"tars_cpp\"
export BuildPerson=\"admin\"
export BuildTime=\"$(date)\"
" >ImageBase/tars."${KEY}"/root/etc/detail

  if ! docker build -t ${_DOCKER_REGISTRY_URL_}/tars."${KEY}":10000 ImageBase/tars."${KEY}"; then
    Record_Error_Info "生成 ${KEY} 镜像失败"
    exit 255
  fi
  if ! docker push ${_DOCKER_REGISTRY_URL_}/tars."${KEY}":10000; then
    Record_Error_Info "推送 ${KEY} 镜像失败"
    exit 255
  fi
done
# 生成镜像
#
# 部署到 K8S

kubectl create namespace tars

kubectl create secret -n tars docker-registry tars-image-secret --docker-server=${_DOCKER_REGISTRY_URL_} --docker-username=${_DOCKER_REGISTRY_USER_} --docker-password=${_DOCKER_REGISTRY_PASSWORD_}

for K8SNode in "${DeployNode[@]}"; do
  if ! kubectl label node "${K8SNode}" 'tars.io/node=' --overwrite ; then
    Record_Error_Info "在 K8S Add Label Failed"
    exit 255
  fi

  if ! kubectl label node "${K8SNode}" 'tars.io/ability.tars=' --overwrite; then
    Record_Error_Info "在 K8S Add Label Failed"
    exit 255
  fi
done

declare -a FirstYamlFile=(
  Yaml/tarscommon.yaml
  Yaml/tarsregistry.yaml
  Yaml/tarsadmin.yaml
)

for YamlFile in "${FirstYamlFile[@]}"; do
  if ! kubectl apply -f "${YamlFile}"; then
    Record_Error_Info "在 K8S Apply ${YamlFile} Failed"
    exit 255
  fi
done

# 等待 tarsregistry 启动
sleep 120

declare -a SecondYamlFile=(
  Yaml/tarsconfig.yaml
  Yaml/tarsimage.yaml
  Yaml/tarslog.yaml
  Yaml/tarsnotify.yaml
  Yaml/tarsproperty.yaml
  Yaml/tarsqueryproperty.yaml
  Yaml/tarsquerystat.yaml
  Yaml/tarsstat.yaml
  Yaml/tarsweb.yaml
)

for YamlFile in "${SecondYamlFile[@]}"; do
  if ! kubectl apply -f "${YamlFile}"; then
    Record_Error_Info "在 K8S Apply ${YamlFile} Failed"
    exit 255
  fi
done
