package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const FixedBaseTag = "a"

type ConfigMap struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Data struct {
		Namespace			   string `yaml:"Namespace"`
		CppImageBase           string `yaml:"CppImageBase"`
		DockerRegistryPassword string `yaml:"DockerRegistryPassword"`
		DockerRegistryURL      string `yaml:"DockerRegistryUrl"`
		DockerRegistryUser     string `yaml:"DockerRegistryUser"`
		JavaImageBase          string `yaml:"JavaImageBase"`
		Node8ImageBase         string `yaml:"Node8ImageBase"`
		Node10ImageBase        string `yaml:"Node10ImageBase"`
		NodeImageBase          string `yaml:"NodeImageBase"`
	} `yaml:"data"`
}

type BuildRequest struct {
	ServerApp    string `json:"ServerApp" valid:"required"`
	ServerName   string `json:"ServerName" valid:"required"`
	ServerType   string `json:"ServerType" valid:"required"`
	ServerDir 	 string `json:"ServerDir" valid:"required"`

	TReleaseTemplatePath 	string
	TConfigTemplatePath 	string
	TServerTemplatePath 	string
}

func NewBuildRequest() BuildRequest {
	return BuildRequest{
		TReleaseTemplatePath: treleaseYamlPath,
		TConfigTemplatePath: tconfigYamlPath,
		TServerTemplatePath: tserverYamlPath}
}

var BashShell = "/bin/bash"
var syncGroup sync.WaitGroup
var totalServerNum int32
var completeServerNum int32

func BuildPatchImage() {
	// docker登陆
	sh := fmt.Sprintf("docker login -u %s -p %s %s", DockerRegistryUser, DockerRegistryPassword, DockerRegistryUrl)
	if _, err = exec.Command(BashShell, "-c", sh).CombinedOutput(); err != nil {
		panic(fmt.Sprintf("run docker login err: %s", err))
	}

	// 遍历应用目录，逐个生成服务镜像
	dir, err := ioutil.ReadDir(AppBaseDir)
	if err != nil {
		panic(fmt.Sprintf("read directory %s err: %s\n", AppBaseDir, err))
	}

	// 构建任务集合
	var requests = make([]BuildRequest, 0, 0)
	for _, tgz := range dir {
		request := NewBuildRequest()

		fileName := tgz.Name()
		if tgz.IsDir() {
			fmt.Println(fmt.Sprintf("%s/%s is a directory, what happened? skip it.", AppBaseDir, fileName))
			continue
		}
		ok := strings.HasSuffix(fileName, ".tar.gz") || strings.HasSuffix(fileName, ".tgz")
		if !ok {
			fmt.Println(fmt.Sprintf("%s/%s is invalid name, what happened? skip it.", AppBaseDir, fileName))
			continue
		}

		// 解压判断服务类型
		fileFields := strings.Split(fileName, ".")
		request.ServerApp = fileFields[0]
		request.ServerName= fileFields[1]
		request.ServerType, request.ServerDir = getServerType(AppBaseDir, fileName)
		if !Init {
			request.TReleaseTemplatePath = fmt.Sprintf("%s/release/%s-%s.yaml", AppBaseDir, strings.ToLower(request.ServerApp), strings.ToLower(request.ServerName))
		}
		requests = append(requests, request)
	}

	// 创建go协程任务
	totalServerNum := len(requests)
	syncGroup.Add(totalServerNum)
	for i := 0; i < totalServerNum; i++ {
		go createBuildTask(requests[i])
	}
	syncGroup.Wait()

	fmt.Println("complete image build.")
}

func createBuildTask(request BuildRequest) {
	begTime := time.Now()

	imageTag := strings.ToLower(request.ServerApp) + "." + strings.ToLower(request.ServerName) + ":" + FixedBaseTag + strconv.FormatInt(time.Now().UnixNano(), 10)

	// 创建镜像目录
	createImageDir(request)
	fmt.Println(fmt.Sprintf("%s: 正在启动构建.", imageTag))

	// 创建etc/detail和Dockerfile
	writeImageConfig(request)
	fmt.Println(fmt.Sprintf("%s: 预处理完成.", imageTag))

	// 拉取基镜像
	baseImage := imageBaseMap[request.ServerType]
	sh := fmt.Sprintf("docker pull %s", baseImage)
	if _, err := exec.Command(BashShell, "-c", sh).CombinedOutput(); err != nil {
		panic(fmt.Sprintf("run pull image %s, err: %s", baseImage, err))
	}
	fmt.Println(fmt.Sprintf("%s: 基础镜像拉取完成.", imageTag))

	// 构建镜像
	dockerImage := fmt.Sprintf("%s/%s", DockerRegistryUrl, imageTag)

	sh = fmt.Sprintf("docker build -t %s %s", dockerImage, request.ServerDir)
	if _, err := exec.Command(BashShell, "-c", sh).CombinedOutput(); err != nil {
		panic(fmt.Sprintf("run build image %s, err: %s", imageTag, err))
	}
	fmt.Println(fmt.Sprintf("%s: 服务镜像编译完成.", imageTag))

	sh = fmt.Sprintf("docker push %s", dockerImage)
	cmd := exec.Command(BashShell, "-c", sh)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(fmt.Sprintf("%s: 服务镜像上传stdout重定向错误.", err))
	}
	_ = cmd.Start()
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		if err != nil {
			break
		}
		fmt.Println(string(tmp))
	}
	_ = cmd.Wait()
	fmt.Println(fmt.Sprintf("%s: 服务镜像上传完成.", imageTag))

	sh = fmt.Sprintf("docker rmi %s", dockerImage)
	if _, err := exec.Command(BashShell, "-c", sh).CombinedOutput(); err != nil {
		panic(fmt.Sprintf("run rui image %s, err: %s", imageTag, err))
	}

	// 删除构建目录
	_ = os.RemoveAll(request.ServerDir)

	// 写TRelease.yaml和TServer.yaml文件
	writeTServerFile(request, writeTReleaseFile(request, dockerImage))

	atomic.AddInt32(&completeServerNum, 1)
	syncGroup.Done()

	endTime := time.Now()

	fmt.Println(fmt.Sprintf("%s: 服务镜像构建完成. 已完成: %d/%d, 本任务耗时：%d(s)",
		imageTag, completeServerNum, totalServerNum, endTime.Unix()-begTime.Unix()))

	time.Sleep(time.Duration(1)*time.Second)
}

func unTarGz(srcFilePath string, destDir string) string {
	fr, err := os.Open(srcFilePath)
	if err != nil {
		panic(fmt.Sprintf("untar file: %s, err: %s", srcFilePath, err))
	}
	defer fr.Close()

	gr, err := gzip.NewReader(fr)
	tr := tar.NewReader(gr)

	unTarDir := destDir
	for {
		hdr, err := tr.Next()
		if hdr == nil {
			break
		}
		if err == io.EOF {
			break
		}
		if hdr.Typeflag != tar.TypeDir {
			_ = os.MkdirAll(destDir+"/"+path.Dir(hdr.Name), os.ModePerm)
			fw, _ := os.Create(destDir+"/"+hdr.Name)
			_, err = io.Copy(fw, tr)
		}
		if unTarDir == destDir {
			unTarDir = fmt.Sprintf("%s/%s", unTarDir, strings.Split(path.Dir(hdr.Name), "/")[0])
		}
	}

	return unTarDir
}

func getServerType(dir, file string) (string, string) {
	tarDir := fmt.Sprintf("%s/%s", dir, file)
	unTarDir := unTarGz(tarDir, dir)

	server, err := ioutil.ReadDir(unTarDir)
	if err != nil {
		panic(fmt.Sprintf("read directory %s err: %s\n", unTarDir, err))
	}

	sh := fmt.Sprintf("chmod -R 777 %s", unTarDir)
	_, err = exec.Command(BashShell, "-c", sh).CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("chmod directory %s err: %s\n", unTarDir, err))
	}

	serverType := "taf_cpp"
	for _, file := range server {
		if file.Name() == "app.js" || file.Name() == "package.json" {
			serverType = "taf_node10"
		}
	}

	return serverType, unTarDir
}

func createImageDir(request BuildRequest) {
	_ = os.MkdirAll(fmt.Sprintf("%s/root", request.ServerDir), os.ModePerm)
	_ = os.MkdirAll(fmt.Sprintf("%s/root/etc", request.ServerDir), os.ModePerm)
	_ = os.MkdirAll(fmt.Sprintf("%s/root/usr/lib", request.ServerDir), os.ModePerm)
	_ = os.MkdirAll(fmt.Sprintf("%s/root/usr/local/server/bin", request.ServerDir), os.ModePerm)

	// 遍历所有解压后的文件，并移动到bin目录下
	dir, err := ioutil.ReadDir(request.ServerDir)
	if err != nil {
		panic(fmt.Sprintf("read directory %s err: %s\n", request.ServerDir, err))
	}
	for _, file := range dir {
		_ = os.Rename(fmt.Sprintf("%s/%s", request.ServerDir, file.Name()),
			fmt.Sprintf("%s/root/usr/local/server/bin/%s", request.ServerDir, file.Name()))
	}
}

func writeImageConfig(request BuildRequest) {
	baseImage := imageBaseMap[request.ServerType]

	// 创建etc文件，tafnode使用
	bytes := []byte(fmt.Sprintf("#!/bin/bash\nexport ServerName=%s\nexport ServerType=%s\nexport BuildPerson=%s\nexport BuildTime=%s\n",
		request.ServerName, request.ServerType, "BatchScript", time.Now().Format("2006-01-02 15:04:05")))
	detailFile := fmt.Sprintf("%s/root/etc/detail", request.ServerDir)
	if err := ioutil.WriteFile(detailFile, bytes, os.ModePerm); err != nil {
		panic(fmt.Sprintf("%s.%s write detail err: %s", request.ServerApp, request.ServerName, err))
	}

	// 创建dockerfile
	bytes = []byte(fmt.Sprintf("From %s \nCopy /root /\n", baseImage))
	dockerFile := fmt.Sprintf("%s/Dockerfile", request.ServerDir)
	if err := ioutil.WriteFile(dockerFile, bytes, os.ModePerm); err != nil {
		panic(fmt.Sprintf("%s.%s write dockerfile err: %s", request.ServerApp, request.ServerName, err))
	}
}

func writeTReleaseFile(request BuildRequest, dockerImage string) ReleaseImageItem {
	// 加载TRelease模板
	var treleaseTemplate = TRelease{}
	release, err := ioutil.ReadFile(request.TReleaseTemplatePath)
	if err != nil {
		panic(fmt.Sprintf("read from %s err: %s\n", "trelease.yaml", err))
	}
	err = yaml.Unmarshal(release, &treleaseTemplate)
	if err != nil {
		panic(fmt.Sprintf("unmarshal from %s err: %s\n", "trelease.yaml", err))
	}

	treleaseTemplate.Metadata.Name = fmt.Sprintf("%s-%s", strings.ToLower(request.ServerApp), strings.ToLower(request.ServerName))
	treleaseTemplate.Metadata.Namespace = Namespace

	var item ReleaseImageItem

	if Init {
		if len(treleaseTemplate.Spec.List) != 1 {
			panic(fmt.Sprintf("unmarshal from %s list size: %d > 1.\n", "trelease.yaml", len(treleaseTemplate.Spec.List)))
		}
		item = treleaseTemplate.Spec.List[0]
		item.Image 	= dockerImage
		item.Tag 	= "10000"
		item.ServerType = request.ServerType
		treleaseTemplate.Spec.List[0] = item
	} else {
		if len(treleaseTemplate.Spec.List) <= 0 {
			panic(fmt.Sprintf("unmarshal from %s list size: %d <= 0.\n", "trelease.yaml", len(treleaseTemplate.Spec.List)))
		}
		item.Image = dockerImage
		item.ServerType = request.ServerType

		oldElem := treleaseTemplate.Spec.List[len(treleaseTemplate.Spec.List)-1]
		item.ImagePullSecret = oldElem.ImagePullSecret
		tag, _ := strconv.Atoi(oldElem.Tag)

		item.Tag = strconv.Itoa(tag+1)
		treleaseTemplate.Spec.List = append(treleaseTemplate.Spec.List, item)
	}

	output, err := yaml.Marshal(&treleaseTemplate)
	if err != nil {
		panic(fmt.Sprintf("marshal from %v err: %s\n", treleaseTemplate, err))
	}
	_ = ioutil.WriteFile(fmt.Sprintf("%s/%s.yaml", AppReleaseDir, treleaseTemplate.Metadata.Name), output, os.ModePerm)

	return item
}

func writeTServerFile(request BuildRequest, release ReleaseImageItem) {
	server, err := ioutil.ReadFile(request.TServerTemplatePath)
	if err != nil {
		panic(fmt.Sprintf("read from %s err: %s\n", "tserver.yaml", err))
	}
	tafserver := &TServer{}
	err = yaml.Unmarshal(server, &tafserver)
	if err != nil {
		panic(fmt.Sprintf("unmarshal from %s err: %s\n", "tserver.yaml", err))
	}

	if FromK8SDB {
		ok := AdapterK8SDBTServerData(tafserver, request, release)
		if !ok {
			fmt.Println(fmt.Sprintf("cannot find tserver: %s.%s in k8s db.", request.ServerApp, request.ServerName))
		}
	} else {
		ok := AdapterTafDBTServerData(tafserver, request, release)
		if !ok {
			fmt.Println(fmt.Sprintf("cannot find tserver: %s.%s in taf db.", request.ServerApp, request.ServerName))
		}
	}
}

