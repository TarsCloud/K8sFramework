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

function LOG_DEBUG() {
  msg=$(date +%Y-%m-%d" "%H:%M:%S)
  msg="${msg} $*"
  echo -e "\033[40;37m $msg \033[0m"
}
###### LOG 函数

# ### 检出指定branch和tag的git submodule分支
# if ! git submodule init; then
#   LOG_ERROR "Git Update Submodule Error"
#   exit 255
#   end
# fi

# checkOutCmd="git submodule foreach '
#   branch=\$(git config -f \$toplevel/.gitmodules submodule.\$name.branch);
#   if [[ -n \branch ]]; then
#     git checkout \$branch;
#     git pull;
#   else
#     exit 255;
#   fi;

#   if [[ \$name == tars-cpp ]]; then
#     tag=\$(git config -f \$toplevel/.gitmodules submodule.\$name.tag);
#     if [[ -n \$tag ]]; then
#       git checkout \$tag;
#     else
#       exit 255;
#     fi;
#   fi;
# '"

# eval ${checkOutCmd}
### 检出指定branch和tag的git submodule分支

### 构建 Docker 镜像,并在镜像中编译代码
if ! docker build -t tars.builder -f build/tars.builder.Dockerfile build; then
  LOG_ERROR "Build \"tars.builder\" image error"
  exit 255
  end
fi

if ! docker run -i -v "${PWD}"/:/tars-src -v "${PWD}"/build/files/tars.cpp.bootstrap.sh:/tars-src/bootstrap.sh -v "${PWD}"/build/files/binary:/tars-k8s-binary tars.builder; then
  LOG_ERROR "Build Source Error"
  exit 255
  end
fi
### 构建 Docker 镜像,并在镜像中编译代码


# docker run -i -v "${PWD}"/src:/tars-k8s-src -v "${PWD}"/build/files/tars.cpp.bootstrap.sh:/tars-k8s-src/TarsFramework/bootstrap.sh -v "${PWD}"/build/files/binary:/tars-k8s-binary tars.builder