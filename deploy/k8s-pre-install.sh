#!/bin/bash

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

MIRROR=http://mirrors.cloud.tencent.com
NODE_VERSION="v12.13.0"
CURRENT_NODE_SUCC=`node -e "console.log('succ')"`
CURRENT_NODE_VERSION=`node --version`

export NVM_NODEJS_ORG_MIRROR=${MIRROR}/nodejs-release/

if [[ "${CURRENT_NODE_SUCC}" != "succ" ]]; then

  rm -rf v0.35.1.zip
  #centos8 need chmod a+x
  chmod a+x /usr/bin/unzip
  wget https://github.com/nvm-sh/nvm/archive/v0.35.1.zip --no-check-certificate;/usr/bin/unzip v0.35.1.zip

  NVM_HOME=$HOME

  rm -rf $NVM_HOME/.nvm; rm -rf $NVM_HOME/.npm; cp -rf nvm-0.35.1 $NVM_HOME/.nvm; rm -rf nvm-0.35.1;

  NVM_DIR=$NVM_HOME/.nvm;
  echo "export NVM_DIR=$NVM_DIR; [ -s $NVM_DIR/nvm.sh ] && \. $NVM_DIR/nvm.sh; [ -s $NVM_DIR/bash_completion ] && \. $NVM_DIR/bash_completion;" >> /etc/profile

  source /etc/profile

  nvm install ${NODE_VERSION};
fi

################################################################################
#check node version
CURRENT_NODE_VERSION=`node --version`

if [[ "${CURRENT_NODE_VERSION}" < "${NODE_VERSION}" ]]; then
    echo "node is not valid, must be after version:${NODE_VERSION}, please remove your node first."
    exit 1
fi

echo "install node success! Version is ${CURRENT_NODE_VERSION}"

source /etc/profile

rm -rf taf-node 
git clone http://gitlab.whup.com/taf/taf-node.git
cd taf-node/src/taf-node-agent; npm install --registry=http://172.16.8.152:7001
cd ../../..

rm -rf taf-web-nodejs-k8s
git clone http://gitlab.whup.com/taf/taf-web-nodejs-k8s.git 
cd taf-web-nodejs-k8s; npm install --registry=http://172.16.8.152:7001
cd ..

yum install -y mysql