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
    echo "$0 DOCKER_REGISTRY_URL DOCKER_REGISTRY_USER DOCKER_REGISTRY_PASSWORD DB_TARS_HOST DB_TARS_PORT DB_TARS_USER DB_TARS_PASSWORD NODEIP";
    echo "for example $0 docker.tars.com/k8s tars 12345 xx.xx.xx.xx 3306 root 12345 \"10.211.55.9 10.211.55.10\""
    exit 1
fi

#仓库信息
_DOCKER_REGISTRY_URL_=$1
_DOCKER_REGISTRY_USER_=$2
_DOCKER_REGISTRY_PASSWORD_=$3
#

# 填写 tars 框架 数据库信息
_DB_TARS_HOST_=$4
_DB_TARS_PORT_=$5
_DB_TARS_DATABASE_=k8s_tars_db
_DB_TARS_USER_=$6
_DB_TARS_PASSWORD_=$7

# 填写 存储 stat 数据库信息 
_DB_TARS_STAT_HOST_=$4
_DB_TARS_STAT_PORT_=$5
_DB_TARS_STAT_DATABASE_=k8s_tars_stat  #发现部分代码对此数据库名耦合,建议保留
_DB_TARS_STAT_USER_=$6
_DB_TARS_STAT_PASSWORD_=$7

# 填写 存储 property 数据库信息
_DB_TARS_PROPERTY_HOST_=$4
_DB_TARS_PROPERTY_PORT_=$5
_DB_TARS_PROPERTY_DATABASE_=k8s_tars_property #发现部分代码对此数据库名耦合,建议保留
_DB_TARS_PROPERTY_USER_=$6
_DB_TARS_PROPERTY_PASSWORD_=$7

#节点IP
_NODEIP_=$8


LOG_INFO "====================================================================";
LOG_INFO "===**********************tars-k8s-install*****************************===";
LOG_INFO "====================================================================";

#输出配置信息
LOG_DEBUG "===>print config info >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>";
LOG_DEBUG "PARAMS:                     "$*
LOG_DEBUG "DOCKER_REGISTRY_URL:        "$_DOCKER_REGISTRY_URL_
LOG_DEBUG "DOCKER_REGISTRY_USER:       "$_DOCKER_REGISTRY_USER_
LOG_DEBUG "DOCKER_REGISTRY_PASSWORD:   "$_DOCKER_REGISTRY_PASSWORD_
LOG_DEBUG "DB_TARS_HOST:               "$_DB_TARS_HOST_
LOG_DEBUG "DB_TARS_PORT:               "$_DB_TARS_PORT_
LOG_DEBUG "DB_TARS_USER:               "$_DB_TARS_USER_
LOG_DEBUG "DB_TARS_PASSWORD:           "$_DB_TARS_PASSWORD_
LOG_DEBUG "NODEIP:                     "$_NODEIP_
LOG_DEBUG "===<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< print config info finish.\n";

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

# cp -rf tars-node/src/tars-node-agent k8s-template/ImageBase/tars.nodebase
# cp -rf tars-node/src/tars-node-agent k8s-template/ImageBase/tars.nodebase
# cp -rf tars-node/src/tars-node-agent k8s-template/ImageBase/tars.node8base
# cp -rf tars-node/src/tars-node-agent k8s-template/ImageBase/tars.node10base

cp -rf k8s-web k8s-template/ImageBase/tars.tarsweb

declare -a OriginServerImage=(
  tarslog
  tarsstat
  tarsquerystat
  tarsproperty
  tarsqueryproperty
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
if ! mysql -h${_DB_TARS_HOST_} -P${_DB_TARS_PORT_} -u${_DB_TARS_USER_} -p${_DB_TARS_PASSWORD_} <MySql/db_tars.sql; then
  LOG_ERROR "backup db_tars failed"
  exit 255
fi

# 生成镜像
if ! docker login -u ${_DOCKER_REGISTRY_USER_} -p ${_DOCKER_REGISTRY_PASSWORD_} ${_DOCKER_REGISTRY_URL_}; then
  LOG_ERROR "docker login to ${_DOCKER_REGISTRY_URL_} failed!"
  exit 255
fi

declare -a BaseImage=(
  tars.base
  tars.cppbase
  tars.javabase
  tars.nodejsbase
  tars.tarsimage
  tars.tarsadmin
  tars.tarsagent
  tars.tarsregistry
  tars.tarsweb
)

if ! cp Program/tarsnode ImageBase/tars.cppbase/root/usr/local/app/tars/tarsnode/bin/tarsnode; then
  LOG_ERROR "copy tarsnode failed please check  tarsnode is in directory Program"
  exit 255
fi
chmod a+x ImageBase/tars.cppbase/root/usr/local/app/tars/tarsnode/bin/tarsnode

if ! cp Program/tarsregistry ImageBase/tars.tarsregistry/root/usr/local/app/tars/tarsregistry/bin; then
  LOG_ERROR "copy tarsregistry failed please check  tarsregistry is in directory Program"
  chmod +x ImageBase/tars.tarsregistry/root/usr/local/app/tars/tarsregistry/bin/tarsregistry
  exit 255
fi

if ! cp Program/tarsimage ImageBase/tars.tarsimage/root/usr/local/app/tars/tarsimage/bin; then
  LOG_ERROR "copy tarsimage failed please check  tarsimage is in directory Program"
  exit 255
fi
chmod a+x ImageBase/tars.tarsimage/root/usr/local/app/tars/tarsimage/bin/tarsimage

if ! cp Program/tarsadmin ImageBase/tars.tarsadmin/root/usr/local/app/tars/tarsadmin/bin; then
  LOG_ERROR "copy tarsadmin failed, please check tarsadmin is in directory Program"
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
  tarsnotify
  tarslog
  tarsconfig
  tarsstat
  tarsquerystat
  tarsproperty
  tarsqueryproperty
)

for KEY in "${ServerImage[@]}"; do
  mkdir -p ImageBase/tars.${KEY}
  mkdir -p ImageBase/tars.${KEY}/root/etc
  mkdir -p ImageBase/tars.${KEY}/root/usr/local/server/bin

  if ! cp Program/${KEY} ImageBase/tars.${KEY}/root/usr/local/server/bin; then
    LOG_ERROR "copy ${KEY} failed, please check ${KEY} is in directory: Program/"
    exit 255
  fi

  echo "FROM tars.cppbase
COPY /root /
CMD [\"/bin/entrypoint.sh\"]
" >ImageBase/tars.${KEY}/Dockerfile

  echo "#!/bin/bash
export ServerName=\"${KEY}\"
export ServerType=\"tars_cpp\"
export BuildPerson=\"admin\"
export BuildTime=\"$(date)\"
" >ImageBase/tars.${KEY}/root/etc/detail

  if ! docker build -t ${_DOCKER_REGISTRY_URL_}/tars.${KEY}:10000 ImageBase/tars.${KEY}; then
    LOG_ERROR "docker build ${KEY} image failed"
    exit 255
  fi
  if ! docker push ${_DOCKER_REGISTRY_URL_}/tars.${KEY}:10000; then
    LOG_ERROR "push ${KEY} to docker registry failed"
    exit 255
  fi
done
# 生成镜像
#
# 部署到 K8S

kubectl create namespace tars

kubectl create secret -n tars docker-registry tars-image-secret --docker-server=${_DOCKER_REGISTRY_URL_} --docker-username=${_DOCKER_REGISTRY_USER_} --docker-password=${_DOCKER_REGISTRY_PASSWORD_}

for K8SNode in $_NODEIP_; do
  if ! kubectl label node ${K8SNode} 'tars.io/node=' --overwrite ; then
    LOG_ERROR "K8S Add Label Failed"
    exit 255
  fi

  if ! kubectl label node ${K8SNode} 'tars.io/ability.tars=' --overwrite; then
    LOG_ERROR "K8S Add Label Failed"
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
    LOG_ERROR "K8S Apply ${YamlFile} Failed"
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
    LOG_ERROR "K8S Apply ${YamlFile} Failed"
    exit 255
  fi
done

cd ..
rm -rf k8s-install.tmp