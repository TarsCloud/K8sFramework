#!/bin/bash

BuildID="$1"
BuildDir="$2"
ServerTGZ="$3"
RegistryEnv="$4"

DockerRegistry=$(cat ${RegistryEnv}"DockerRegistryUrl")
DockerUser=$(cat ${RegistryEnv}"DockerRegistryUser")
DockerPassword=$(cat ${RegistryEnv}"DockerRegistryPassword")

BuildDetailFile="ServerDetail"
BuildStatusFile="BuildStatus"
BuildMissLibFile="MissLib"
BuildPidFile="BuildPid"
BuildDockerFile="Dockerfile"
BuildDockerRootDir="root"

cd "${BuildDir}" || exit
source "${BuildDetailFile}"

ImageEtcDir="/etc"
BuildEtcDir="root${ImageEtcDir}"

ImageUserLidDir="/usr/lib"
BuildUsrLibDir="root${ImageUserLidDir}"
mkdir -p "${BuildUsrLibDir}"

ImageServerDetail="/etc/detail"
BuildServerDetailFile="root${ImageServerDetail}"

ImageServerBinDir="/usr/local/server/bin"
BuildServerBinDir="root${ImageServerBinDir}"
mkdir -p "${BuildServerBinDir}"

ImageServerLibDir="${ImageServerBinDir}/lib"
BuildServerLibDir="root${ImageServerLibDir}"

DOCKER_IMAGE_LOCAL_TAG="${ImageTag}"
DOCKER_IMAGE_REGISTRY_TAG=${DockerRegistry}/${DOCKER_IMAGE_LOCAL_TAG}

function Record_Task_Error_Status() {
  echo \{\"BuildId\":\""${BuildID}"\",\"BuildStatus\":\"error\",\"BuildMessage\":\""$*"\"\,\"BuildImage\":\""${DOCKER_IMAGE_REGISTRY_TAG}"\"\} >"${BuildStatusFile}"
}

function Record_Task_Working_Status() {
  echo \{\"BuildId\":\""${BuildID}"\",\"BuildStatus\":\"working\",\"BuildMessage\":\""$*"\"\,\"BuildImage\":\""${DOCKER_IMAGE_REGISTRY_TAG}"\"\} >"${BuildStatusFile}"
}

function Record_Task_Done_Status() {
  echo \{\"BuildId\":\""${BuildID}"\",\"BuildStatus\":\"done\",\"BuildMessage\":\""$*"\"\,\"BuildImage\":\""${DOCKER_IMAGE_REGISTRY_TAG}"\"\} >"${BuildStatusFile}"
}

function Decompression_Server_TGZ() {
  Record_Task_Working_Status "正在解压文件"
  if ! tar zxvf "${ServerTGZ}" --strip-components=1 -C "${BuildServerBinDir}"; then
    Record_Task_Error_Status "解压文件失败"
    exit 255
  fi
}

function Copy_Cpp_Lib() {
  Record_Task_Working_Status "正在检测程序依赖文件"
  rm -rf ${BuildMissLibFile}

  BuildServerExeFile="${BuildServerBinDir}/${ServerName}"
  export LD_LIBRARY_PATH=${PWD}/${BuildServerBinDir}:${PWD}/${BuildServerLibDir}:${LD_LIBRARY_PATH}

  ldd "${BuildServerExeFile}" | awk \
    {'
        {
            if (match($3,".so")){
                print "found "$3
                system("cp " $3 " '"$BuildUsrLibDir"'")
            }
            else if ($3=="not" && $4=="found") {
                system("echo -n " $1 ",  >> '"$BuildMissLibFile"'")
                print "not found " $1
            }
        }
        system("echo -n "" >> '"$BuildMissLibFile"'")
    '}

  if [[ ! -f "${BuildMissLibFile}" ]]; then
    Record_Task_Error_Status "检测程序依赖失败"
    exit 255
  fi

  MissingLibs=$(cat ${BuildMissLibFile})
  if [[ -n "$MissingLibs" ]]; then
    Record_Task_Error_Status "缺失依赖文件" "${MissingLibs}"
    exit 255
  fi

  ls /lib64/ | xargs -t -I{} rm -rf "${BuildUsrLibDir}"/{}
  ls ${BuildUsrLibDir} | xargs -t -I{} rm -rf "${BuildServerLibDir}"/{}
}

function Make_Image_Detail() {
  mkdir -p ${BuildEtcDir}
  echo -e "
#!/bin/bash
export ServerName=\"${ServerName}\"
export ServerType=\"${ServerType}\"
export BuildPerson=\"${BuildPerson}\"
export BuildTime=\"$(date)\"" >${BuildServerDetailFile}
}

function Build_and_Push_Image() {
  Record_Task_Working_Status "正在生成镜像"
  cd "${BuildDir}" || exit
  if ! docker build -t "${DOCKER_IMAGE_REGISTRY_TAG}" .; then
    Record_Task_Error_Status 镜像生成失败
    exit 255
  fi

  Record_Task_Working_Status "正在上传镜像"

  if ! docker push "${DOCKER_IMAGE_REGISTRY_TAG}"; then
    Record_Task_Error_Status "镜像上传失败"
    exit 255
  fi

  docker rmi "${DOCKER_IMAGE_REGISTRY_TAG}"
  return
}

echo -n $$ >${BuildPidFile}

if ! docker login -u "${DockerUser}" -p "${DockerPassword}" "${DockerRegistry}"; then
  Record_Task_Error_Status "登录镜像仓库失败"
  exit 255
fi

Decompression_Server_TGZ

Node12ImageBase

case ${ServerType} in
"tars_cpp" | "tars_node_pkg" | "tars_go")
  BuildImageBase=$(cat ${RegistryEnv}"CppImageBase")
  ;;
"tars_java_war" | "tars_java_jar")
  BuildImageBase=$(cat ${RegistryEnv}"JavaImageBase")
  ;;
"tars_node")
  BuildImageBase=$(cat ${RegistryEnv}"NodeImageBase")
  ;;
"tars_node8")
  BuildImageBase=$(cat ${RegistryEnv}"Node8ImageBase")
  ;;
"tars_node10")
  BuildImageBase=$(cat ${RegistryEnv}"Node10ImageBase")
  ;;
*)
  Record_Task_Error_Status "不识别的服务类型"
  exit 255
  ;;
esac

Record_Task_Working_Status "正在下载基础镜像"

if ! docker pull "${BuildImageBase}"; then
  Record_Task_Error_Status "下载基础镜像失败"
  exit 255
fi

echo -e "
FROM ${BuildImageBase}
COPY ${BuildDockerRootDir} /" >"${BuildDockerFile}"

Make_Image_Detail
Build_and_Push_Image
Record_Task_Done_Status "任务完成"