#!/bin/bash

############################################

# 填写 docker 镜像仓库信息
_DOCKER_REGISTRY_URL_=
_DOCKER_REGISTRY_USER_=
_DOCKER_REGISTRY_PASSWORD_=
#

# 填写 taf 框架 数据库信息
_DB_TAF_HOST_=
_DB_TAF_PORT_=
_DB_TAF_DATABASE_=
_DB_TAF_USER_=
_DB_TAF_PASSWORD_=

# 填写 存储 stat 数据库信息 
_DB_TAF_STAT_HOST_=
_DB_TAF_STAT_DATABASE_=taf_stat  #发现部分代码对此数据库名耦合,建议保留
_DB_TAF_STAT_USER_=
_DB_TAF_STAT_PASSWORD_=
_DB_TAF_STAT_PORT_=

# 填写 存储 property 数据库信息
_DB_TAF_PROPERTY_HOST_=
_DB_TAF_PROPERTY_DATABASE_=taf_property #发现部分代码对此数据库名耦合,建议保留
_DB_TAF_PROPERTY_USER_=
_DB_TAF_PROPERTY_PASSWORD_=
_DB_TAF_PROPERTY_PORT_=


# 填写 部署 taf 程序的 k8s node ,建议两台或以上
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
  _DB_TAF_HOST_
  _DB_TAF_PORT_
  _DB_TAF_DATABASE_
  _DB_TAF_USER_
  _DB_TAF_PASSWORD_

  #
  _DB_TAF_STAT_HOST_
  _DB_TAF_STAT_DATABASE_
  _DB_TAF_STAT_USER_
  _DB_TAF_STAT_PASSWORD_
  _DB_TAF_STAT_PORT_

  #
  _DB_TAF_PROPERTY_HOST_
  _DB_TAF_PROPERTY_DATABASE_
  _DB_TAF_PROPERTY_USER_
  _DB_TAF_PROPERTY_PASSWORD_
  _DB_TAF_PROPERTY_PORT_

)

declare -a AllSqlFile=(
  # MySql/db_property.sql
  # MySql/db_stat.sql
  MySql/db_taf.sql
)

declare -a AllYamlFile=(
  Yaml/tafadmin.yaml
  Yaml/tafconfig.yaml
  Yaml/tafcommon.yaml
  Yaml/tafimage.yaml
  Yaml/taflog.yaml
  Yaml/tafnotify.yaml
  Yaml/tafproperty.yaml
  Yaml/tafqueryproperty.yaml
  Yaml/tafquerystat.yaml
  Yaml/tafregistry.yaml
  Yaml/tafstat.yaml
  Yaml/tafweb.yaml
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
if ! mysql -h${_DB_TAF_HOST_} -P${_DB_TAF_PORT_} -u${_DB_TAF_USER_} -p${_DB_TAF_PASSWORD_} <MySql/db_taf.sql; then
  Record_Error_Info "备份数据到 db_taf 失败, 请检测 db_taf 相关参数是否设置正确,或者网络是否正常"
  exit 255
fi

# 生成镜像
if ! docker login -u ${_DOCKER_REGISTRY_USER_} -p ${_DOCKER_REGISTRY_PASSWORD_} ${_DOCKER_REGISTRY_URL_}; then
  Record_Error_Info "登录镜像仓库失败,请检测 Docker仓库相关参数是否设置正确,或者网络是否正常"
  exit 255
fi

declare -a BaseImage=(
  taf.base
  taf.cppbase
  taf.javabase
  taf.nodebase
  taf.node8base
  taf.node10base
  taf.tafimage
  taf.tafadmin
  taf.tafagent
  taf.tafregistry
  taf.tafweb
)

if ! cp Program/tafnode ImageBase/taf.cppbase/root/usr/local/app/taf/tafnode/bin/tafnode; then
  Record_Error_Info "拷贝程序 tafnode 失败， 请检测 tafnode 是否放置在 Program 目录"
  exit 255
fi
chmod a+x ImageBase/taf.cppbase/root/usr/local/app/taf/tafnode/bin/tafnode

if ! cp Program/tafregistry ImageBase/taf.tafregistry/root/usr/local/app/taf/tafregistry/bin; then
  Record_Error_Info "拷贝程序 tafregistry 失败， 请检测 tafregistry 是否放置在 Program 目录"
  chmod +x ImageBase/taf.tafregistry/root/usr/local/app/taf/tafregistry/bin/tafregistry
  exit 255
fi

if ! cp Program/tafimage ImageBase/taf.tafimage/root/usr/local/app/taf/tafimage/bin; then
  Record_Error_Info "拷贝程序 tafimage 失败， 请检测 tafimage 是否放置在 Program 目录"
  exit 255
fi
chmod a+x ImageBase/taf.tafimage/root/usr/local/app/taf/tafimage/bin/tafimage

if ! cp Program/tafadmin ImageBase/taf.tafadmin/root/usr/local/app/taf/tafadmin/bin; then
  Record_Error_Info "拷贝程序 tafadmin 失败， 请检测 tafadmin 是否放置在 Program 目录"
  exit 255
fi
chmod a+x ImageBase/taf.tafadmin/root/usr/local/app/taf/tafadmin/bin/tafadmin

if ! cp Program/tafagent ImageBase/taf.tafagent/root/usr/local/app/taf/tafagent/bin; then
  Record_Error_Info "拷贝程序 tafagent 失败， 请检测 tafagent 是否放置在 Program 目录"
  exit 255
fi
chmod a+x ImageBase/taf.tafagent/root/usr/local/app/taf/tafagent/bin/tafagent

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
  tafnotify
  taflog
  tafconfig
  tafstat
  tafquerystat
  tafproperty
  tafqueryproperty
)

for KEY in "${ServerImage[@]}"; do
  mkdir -p ImageBase/taf."${KEY}"
  mkdir -p ImageBase/taf."${KEY}"/root/etc
  mkdir -p ImageBase/taf."${KEY}"/root/usr/local/server/bin
  if ! cp Program/"${KEY}" ImageBase/taf."${KEY}"/root/usr/local/server/bin; then
    Record_Error_Info "拷贝程序 ${KEY} 失败， 请检测 ${KEY} 是否放置在 Program 目录"
    exit 255
  fi

  echo "FROM taf.cppbase
COPY /root /
CMD [\"/bin/entrypoint\"]
" >ImageBase/taf."${KEY}"/Dockerfile

  echo "#!/bin/bash
export ServerName=\"${KEY}\"
export ServerType=\"taf_cpp\"
export BuildPerson=\"admin\"
export BuildTime=\"$(date)\"
" >ImageBase/taf."${KEY}"/root/etc/detail

  if ! docker build -t ${_DOCKER_REGISTRY_URL_}/taf."${KEY}":10000 ImageBase/taf."${KEY}"; then
    Record_Error_Info "生成 ${KEY} 镜像失败"
    exit 255
  fi
  if ! docker push ${_DOCKER_REGISTRY_URL_}/taf."${KEY}":10000; then
    Record_Error_Info "推送 ${KEY} 镜像失败"
    exit 255
  fi
done
# 生成镜像
#
# 部署到 K8S

kubectl create namespace taf

kubectl create secret -n taf docker-registry taf-image-secret --docker-server=${_DOCKER_REGISTRY_URL_} --docker-username=${_DOCKER_REGISTRY_USER_} --docker-password=${_DOCKER_REGISTRY_PASSWORD_}

for K8SNode in "${DeployNode[@]}"; do
  if ! kubectl label node "${K8SNode}" 'taf.io/node=' --overwrite ; then
    Record_Error_Info "在 K8S Add Label Failed"
    exit 255
  fi

  if ! kubectl label node "${K8SNode}" 'taf.io/ability.taf=' --overwrite; then
    Record_Error_Info "在 K8S Add Label Failed"
    exit 255
  fi
done

declare -a FirstYamlFile=(
  Yaml/tafcommon.yaml
  Yaml/tafregistry.yaml
  Yaml/tafadmin.yaml
)

for YamlFile in "${FirstYamlFile[@]}"; do
  if ! kubectl apply -f "${YamlFile}"; then
    Record_Error_Info "在 K8S Apply ${YamlFile} Failed"
    exit 255
  fi
done

# 等待 tafregistry 启动
sleep 120

declare -a SecondYamlFile=(
  Yaml/tafconfig.yaml
  Yaml/tafimage.yaml
  Yaml/taflog.yaml
  Yaml/tafnotify.yaml
  Yaml/tafproperty.yaml
  Yaml/tafqueryproperty.yaml
  Yaml/tafquerystat.yaml
  Yaml/tafstat.yaml
  Yaml/tafweb.yaml
)

for YamlFile in "${SecondYamlFile[@]}"; do
  if ! kubectl apply -f "${YamlFile}"; then
    Record_Error_Info "在 K8S Apply ${YamlFile} Failed"
    exit 255
  fi
done
