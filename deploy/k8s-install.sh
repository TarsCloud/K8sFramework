#!/bin/bash

############################################


#公共函数
function LOG_ERROR()
{
	local msg=$(date +%Y-%m-%d" "%H:%M:%S);

    msg="${msg} $@";

	echo -e "\033[31m $msg \033[0m";	
}

function LOG_WARNING()
{
	local msg=$(date +%Y-%m-%d" "%H:%M:%S);

    msg="${msg} $@";

	echo -e "\033[33m $msg \033[0m";	
}

function LOG_DEBUG()
{
	local msg=$(date +%Y-%m-%d" "%H:%M:%S);

    msg="${msg} $@";

 	echo -e "\033[40;37m $msg \033[0m";	
}

function LOG_INFO()
{
	local msg=$(date +%Y-%m-%d" "%H:%M:%S);
	
	for p in $@
	do
		msg=${msg}" "${p};
	done
	
	echo -e "\033[32m $msg \033[0m"  	
}

if (( $# < 8 ))
then
    echo $#
    echo "$0 DOCKER_REGISTRY_URL DOCKER_REGISTRY_USER DOCKER_REGISTRY_PASSWORD DB_TAF_HOST DB_TAF_PORT DB_TAF_USER DB_TAF_PASSWORD NODEIP";
    echo "for example $0 docker.taf.com/k8s taf 12345 xx.xx.xx.xx 3306 root 12345 \"10.211.55.9 10.211.55.10\""
    exit 1
fi

#仓库信息
_DOCKER_REGISTRY_URL_=$1
_DOCKER_REGISTRY_USER_=$2
_DOCKER_REGISTRY_PASSWORD_=$3
#

# 填写 taf 框架 数据库信息
_DB_TAF_HOST_=$4
_DB_TAF_PORT_=$5
_DB_TAF_DATABASE_=taf_db
_DB_TAF_USER_=$6
_DB_TAF_PASSWORD_=$7

# 填写 存储 stat 数据库信息 
_DB_TAF_STAT_HOST_=$4
_DB_TAF_STAT_PORT_=$5
_DB_TAF_STAT_DATABASE_=taf_stat  #发现部分代码对此数据库名耦合,建议保留
_DB_TAF_STAT_USER_=$6
_DB_TAF_STAT_PASSWORD_=$7

# 填写 存储 property 数据库信息
_DB_TAF_PROPERTY_HOST_=$4
_DB_TAF_PROPERTY_PORT_=$5
_DB_TAF_PROPERTY_DATABASE_=taf_property #发现部分代码对此数据库名耦合,建议保留
_DB_TAF_PROPERTY_USER_=$6
_DB_TAF_PROPERTY_PASSWORD_=$7

#节点IP
_NODEIP_=$8


LOG_INFO "====================================================================";
LOG_INFO "===**********************taf-k8s-install*****************************===";
LOG_INFO "====================================================================";

#输出配置信息
LOG_DEBUG "===>print config info >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>";
LOG_DEBUG "PARAMS:                     "$*
LOG_DEBUG "DOCKER_REGISTRY_URL:        "$_DOCKER_REGISTRY_URL_
LOG_DEBUG "DOCKER_REGISTRY_USER:       "$_DOCKER_REGISTRY_USER_
LOG_DEBUG "DOCKER_REGISTRY_PASSWORD:   "$_DOCKER_REGISTRY_PASSWORD_
LOG_DEBUG "DB_TAF_HOST:                "$_DB_TAF_HOST_
LOG_DEBUG "DB_TAF_PORT:                "$_DB_TAF_PORT_
LOG_DEBUG "DB_TAF_USER:                "$_DB_TAF_USER_
LOG_DEBUG "DB_TAF_PASSWORD:            "$_DB_TAF_PASSWORD_
LOG_DEBUG "NODEIP:                     "$_NODEIP_
LOG_DEBUG "===<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< print config info finish.\n";

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

###########################################################################
OSNAME=`uname`
OS=1

if [[ "$OSNAME" == "Darwin" ]]; then
    OS=2
elif [[ "$OSNAME" == "Windows_NT" ]]; then
    OS=3
else
    OS=1
fi

cp -rf taf-node/src/taf-node-agent k8s-template/ImageBase/taf.nodebase
cp -rf taf-node/src/taf-node-agent k8s-template/ImageBase/taf.node8base
cp -rf taf-node/src/taf-node-agent k8s-template/ImageBase/taf.node10base

cp -rf taf-web-nodejs-k8s k8s-template/ImageBase/taf.tafweb

declare -a OriginServerImage=(
  taflog
  tafstat
  tafquerystat
  tafproperty
  tafqueryproperty
)

for KEY in "${OriginServerImage[@]}"; do
  echo "cp -rf framework/servers/$KEY/bin/$KEY k8s-template/Program/"

  cp -rf framework/servers/$KEY/bin/$KEY k8s-template/Program/
done

rm -rf k8s-install.tmp && cp -r k8s-template k8s-install.tmp

WORKDIR=$(cd $(dirname $0); pwd)

cd k8s-install.tmp

echo ${WORKDIR}

for KEY in "${AllConfig[@]}"; do

  # sed -i时，第一个参数mac必须，linux可选，加上后兼容2个平台
  if [ $OS == 2 ]; then

    for SqlFile in "${AllSqlFile[@]}"; do
      if ! sed -i"" "s#${KEY}#${!KEY}#g" "${SqlFile}"; then
        exit 255
      fi
    done

    for YamlFile in "${AllYamlFile[@]}"; do
      if ! sed -i "" "s#${KEY}#${!KEY}#g" "${YamlFile}"; then
        exit 255
      fi
    done

  else 
  
    for SqlFile in "${AllSqlFile[@]}"; do
      if ! sed -i "s#${KEY}#${!KEY}#g" "${SqlFile}"; then
        exit 255
      fi
    done

    for YamlFile in "${AllYamlFile[@]}"; do
      if ! sed -i "s#${KEY}#${!KEY}#g" "${YamlFile}"; then
        exit 255
      fi

    done
  fi
done

function CheckCommandExist() {
  for i in "$@"; do
    type "$i" >/dev/null 2>&1 || {
      LOG_ERROR "not check $i, please install $i first"
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
  LOG_ERROR "backup db_taf failed"
  exit 255
fi

# 生成镜像
if ! docker login -u ${_DOCKER_REGISTRY_USER_} -p ${_DOCKER_REGISTRY_PASSWORD_} ${_DOCKER_REGISTRY_URL_}; then
  LOG_ERROR "docker login to ${_DOCKER_REGISTRY_URL_} failed!"
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
  LOG_ERROR "copy tafnode failed please check  tafnode is in directory Program"
  exit 255
fi
chmod a+x ImageBase/taf.cppbase/root/usr/local/app/taf/tafnode/bin/tafnode

if ! cp Program/tafregistry ImageBase/taf.tafregistry/root/usr/local/app/taf/tafregistry/bin; then
  LOG_ERROR "copy tafregistry failed please check  tafregistry is in directory Program"
  chmod +x ImageBase/taf.tafregistry/root/usr/local/app/taf/tafregistry/bin/tafregistry
  exit 255
fi

if ! cp Program/tafimage ImageBase/taf.tafimage/root/usr/local/app/taf/tafimage/bin; then
  LOG_ERROR "copy tafimage failed please check  tafimage is in directory Program"
  exit 255
fi
chmod a+x ImageBase/taf.tafimage/root/usr/local/app/taf/tafimage/bin/tafimage

if ! cp Program/tafadmin ImageBase/taf.tafadmin/root/usr/local/app/taf/tafadmin/bin; then
  LOG_ERROR "copy tafadmin failed, please check tafadmin is in directory Program"
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
    LOG_ERROR "Build ${KEY} image failed"
    exit 255
  fi

  if ! docker tag "${KEY}" ${_DOCKER_REGISTRY_URL_}/"${KEY}":10000; then
    LOG_ERROR "Tag ${KEY} image failed" 
    exit 255
  fi

  if ! docker push ${_DOCKER_REGISTRY_URL_}/"${KEY}":10000; then
    LOG_ERROR "Push ${KEY} image failed"
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
  mkdir -p ImageBase/taf.${KEY}
  mkdir -p ImageBase/taf.${KEY}/root/etc
  mkdir -p ImageBase/taf.${KEY}/root/usr/local/server/bin

  if ! cp Program/${KEY} ImageBase/taf.${KEY}/root/usr/local/server/bin; then
    LOG_ERROR "copy ${KEY} failed, please check ${KEY} is in directory: Program/"
    exit 255
  fi

  echo "FROM taf.cppbase
COPY /root /
CMD [\"/bin/entrypoint.sh\"]
" >ImageBase/taf.${KEY}/Dockerfile

  echo "#!/bin/bash
export ServerName=\"${KEY}\"
export ServerType=\"taf_cpp\"
export BuildPerson=\"admin\"
export BuildTime=\"$(date)\"
" >ImageBase/taf.${KEY}/root/etc/detail

  if ! docker build -t ${_DOCKER_REGISTRY_URL_}/taf.${KEY}:10000 ImageBase/taf.${KEY}; then
    LOG_ERROR "docker build ${KEY} image failed"
    exit 255
  fi
  if ! docker push ${_DOCKER_REGISTRY_URL_}/taf.${KEY}:10000; then
    LOG_ERROR "push ${KEY} to docker registry failed"
    exit 255
  fi
done
# 生成镜像
#
# 部署到 K8S

kubectl create namespace taf

kubectl create secret -n taf docker-registry taf-image-secret --docker-server=${_DOCKER_REGISTRY_URL_} --docker-username=${_DOCKER_REGISTRY_USER_} --docker-password=${_DOCKER_REGISTRY_PASSWORD_}

for K8SNode in $_NODEIP_; do
  if ! kubectl label node ${K8SNode} 'taf.io/node=' --overwrite ; then
    LOG_ERROR "K8S Add Label Failed"
    exit 255
  fi

  if ! kubectl label node ${K8SNode} 'taf.io/ability.taf=' --overwrite; then
    LOG_ERROR "K8S Add Label Failed"
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
    LOG_ERROR "K8S Apply ${YamlFile} Failed"
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
    LOG_ERROR "K8S Apply ${YamlFile} Failed"
    exit 255
  fi
done

cd ..
rm -rf k8s-install.tmp