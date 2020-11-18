package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)
const DefaultConfigMapPath = "/etc/registry-env/"
const DefaultUploadDir = "/uploadDir/"
const DefaultBuildWorkPath = "/buildDir/"
const DefaultBuildStatusFileName = "BuildStatus"
const DefaultServerDetailFileName = "ServerDetail"
const BuildBashShell = "/bin/build.sh"
const FixedBaseTag = "a"


type BuildRequest struct {
	ServerApp    string `json:"ServerApp" valid:"required"`
	ServerName   string `json:"ServerName" valid:"required"`
	ServerType   string `json:"ServerType" valid:"required"`
	ServerTGZ    string `json:"ServerTGZ" valid:"required"`
	CommitPerson string `json:"CommitPerson"`
}

type BuildStatus struct {
	BuildId      string `json:"BuildId"`
	BuildStatus  string `json:"BuildStatus"`
	BuildMessage string `json:"BuildMessage"`
	//ServerApp    string `json:"ServerApp"`
	//ServerName   string `json:"ServerName"`
	BuildImage string `json:"BuildImage,omm"`
}

type BuildTask struct {
	imageTag string
	buildDir string
	process  *exec.Cmd
}

type BuildTaskManager struct {
	lock    sync.Mutex
	builder map[string]*BuildTask
}

func (bm *BuildTaskManager) init() {
	bm.builder = make(map[string]*BuildTask, 5)
}

func (bm *BuildTaskManager) CreateBuildTask(request BuildRequest) *BuildStatus {
	bm.lock.Lock()
	defer bm.lock.Unlock()

	imageTag := strings.ToLower(request.ServerApp) + "." + strings.ToLower(request.ServerName) + ":" + FixedBaseTag + strconv.FormatInt(time.Now().UnixNano(), 10)
	buildId := base64.RawStdEncoding.EncodeToString([]byte(imageTag))
	absoluteBuildDir := DefaultBuildWorkPath + strings.Replace(imageTag, ":", ".", 1)
	absoluteServerFile := DefaultUploadDir + request.ServerTGZ

	bs, err := ioutil.ReadFile(DefaultConfigMapPath+"DockerRegistryUrl")
	if err != nil {
		BuildStatus := &BuildStatus{
			BuildId:      buildId,
			BuildStatus:  "error",
			BuildMessage: "内部错误: " + err.Error(),
			BuildImage:   "Unknown/" + imageTag,
		}
		return BuildStatus
	}

	DockerRegistryUrl := string(bs)

	if err := os.MkdirAll(absoluteBuildDir, 0777); err != nil {
		BuildStatus := &BuildStatus{
			BuildId:      buildId,
			BuildStatus:  "error",
			BuildMessage: "内部错误: " + err.Error(),
			BuildImage:   DockerRegistryUrl + "/" + imageTag,
		}
		return BuildStatus
	}

	bytes := []byte("#!/bin/bash" + "\nServerApp=" + request.ServerApp + "\nServerName=" + request.ServerName + "\nServerType=" + request.ServerType + "\nImageTag=" + imageTag)
	serverDetailFile := absoluteBuildDir + "/" + DefaultServerDetailFileName
	if err := ioutil.WriteFile(serverDetailFile, bytes, 0644); err != nil {
		BuildStatus := &BuildStatus{
			BuildId:      buildId,
			BuildStatus:  "error",
			BuildMessage: "内部错误: " + err.Error(),
			BuildImage:   DockerRegistryUrl + "/" + imageTag,
		}
		return BuildStatus
	}

	cmd := exec.Command(BuildBashShell, buildId, absoluteBuildDir, absoluteServerFile, DefaultConfigMapPath)
	if err := cmd.Start(); err != nil {
		BuildStatus := &BuildStatus{
			BuildId:      buildId,
			BuildStatus:  "error",
			BuildMessage: "内部错误: " + err.Error(),
			BuildImage:   DockerRegistryUrl + "/" + imageTag,
		}
		return BuildStatus
	}

	buildTask := &BuildTask{
		imageTag: DockerRegistryUrl + "/" + imageTag,
		buildDir: absoluteBuildDir,
		process:  cmd,
	}

	go func() {
		_ = cmd.Wait()
		time.Sleep(time.Minute * 10)
		_ = os.RemoveAll(absoluteBuildDir)
		bm.DeleteTask(buildId)
	}()

	bm.builder[buildId] = buildTask
	BuildStatus := &BuildStatus{
		BuildId:      buildId,
		BuildStatus:  "working",
		BuildMessage: "正在启动中",
		BuildImage:   DockerRegistryUrl + "/" + imageTag,
	}
	return BuildStatus
}

func (bm *BuildTaskManager) GetStatus(buildId string) *BuildStatus {
	bm.lock.Lock()
	defer bm.lock.Unlock()

	var buildStatus *BuildStatus

	for {
		task, ok := bm.builder[buildId]
		if !ok {
			buildStatus = &BuildStatus{
				BuildId:      buildId,
				BuildStatus:  "error",
				BuildMessage: "任务不存在或已删除",
			}
			break
		}

		buildStatusFile := task.buildDir + "/" + DefaultBuildStatusFileName
		statusData, err := ioutil.ReadFile(buildStatusFile)
		if err != nil {
			if syscall.Kill(task.process.Process.Pid, 0) == nil {
				buildStatus = &BuildStatus{
					BuildId:      buildId,
					BuildStatus:  "working",
					BuildMessage: "正在启动中",
					BuildImage:   task.imageTag,
				}
				break
			}

			buildStatus = &BuildStatus{
				BuildId:      buildId,
				BuildStatus:  "error",
				BuildMessage: "内部错误: " + err.Error(),
				BuildImage:   task.imageTag,
			}
			break
		}

		var buildStatusFromStatusFile BuildStatus
		if err := json.Unmarshal(statusData, &buildStatusFromStatusFile); err != nil {
			buildStatus = &BuildStatus{
				BuildId:      buildId,
				BuildStatus:  "error",
				BuildMessage: "内部错误: " + err.Error(),
				BuildImage:   task.imageTag,
			}
			break
		}

		buildStatus = &buildStatusFromStatusFile
		break
	}
	return buildStatus
}

func (bm *BuildTaskManager) DeleteTask(buildId string) {
	bm.lock.Lock()
	defer bm.lock.Unlock()
	delete(bm.builder, buildId)
}
