package k8s

import (
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	k8sCoreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	k8sSchema "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
	"tarsadmin/handler/websocket"
	"tarsadmin/openapi/models"
	"tarsadmin/openapi/restapi/operations/shell"
)

type SSHPodShellHandler struct {}

func (s *SSHPodShellHandler) Handle(params shell.SSHPodShellParams) middleware.Responder {
	return middleware.ResponderFunc(func(writer http.ResponseWriter, producer runtime.Producer) {
		pty, err := websocket.NewTerminalSession(writer, params.HTTPRequest, nil)
		if err != nil {
			msg, _ := (&models.Error{Code: -1, Message: err.Error()}).MarshalBinary()
			http.Error(writer, string(msg), http.StatusInternalServerError)
			return
		}
		defer func() {
			pty.Done()
			_ = pty.Close()
		}()

		podName := *params.PodName

		sh := fmt.Sprintf("#!/bin/sh\n\ndir=/usr/local/app/tars/app_log/%s/%s\n\nif [ -d $dir ]; then\n  cd $dir\nfi\n\nif [ ! -f /bin/bash ]; then\n  sh\nelse\n  bash\nfi\n",
			*params.AppName, *params.ServerName)

		// 历史pod只能通过tars-agent容器进入查看日志
		if *params.History {
			sh	= fmt.Sprintf("#!/bin/sh\n\ndir=/usr/local/app/tars/app_log/%s/%s/%s\n\nif [ -d $dir ]; then\n  cd $dir\nfi\n\nif [ ! -f /bin/bash ]; then\n  sh\nelse\n  bash\nfi\n",
				podName, *params.AppName, *params.ServerName)

			pod, ok	:= getDaemonPodByIp(*params.NodeIP)
			if !ok {
				errorProcess(pty, fmt.Sprintf("Can not get daemon pod from nodeIP: %s\n", *params.NodeIP))
				return
			}
			podName = pod.Name
		}

		execCmd(sh, podName, pty)
	})
}

func execCmd(sh, podName string, pty websocket.PtyHandler) {
	cmd := []string{"/bin/sh", "-c", sh}

	req := K8sOption.K8SClientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(K8sOption.Namespace).
		SubResource("exec")

	req.VersionedParams(&k8sCoreV1.PodExecOptions{
		Command:   cmd,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, k8sSchema.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(K8sOption.config, "POST", req.URL())
	if err != nil {
		errorProcess(pty, fmt.Sprintf("new SPDY executor error! err: %v\n", err))
		return
	}

	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		TerminalSizeQueue: pty,
		Tty:               true,
	})
	if err != nil {
		errorProcess(pty, fmt.Sprintf("Exec to pod error! err: %v\n", err))
		return
	}
}

func errorProcess(pty websocket.PtyHandler, msg string) {
	_, _ = pty.Write([]byte(msg))
	fmt.Printf("wsID:%d err, msg:%s\n",pty.GetPtyID(), msg)
}

func getDaemonPodByName(nodeName string) (*k8sCoreV1.Pod, bool) {
	return getDaemonPodByField(func(pod *k8sCoreV1.Pod) bool {
		return pod.Spec.NodeName == nodeName})
}

func getDaemonPodByIp(nodeIp string) (*k8sCoreV1.Pod, bool) {
	return getDaemonPodByField(func(pod *k8sCoreV1.Pod) bool {
		return pod.Status.HostIP == nodeIp})
}

func getDaemonPodByField(fun func(pod *k8sCoreV1.Pod) bool) (*k8sCoreV1.Pod, bool) {
	requirement, err := labels.NewRequirement("app", selection.DoubleEquals, []string{TarsAgentDaemonSetName})
	if err != nil {
		return nil, false
	}

	agents, err := K8sWatcher.podLister.Pods(K8sOption.Namespace).List(labels.NewSelector().Add([]labels.Requirement{*requirement} ...))
	if err != nil {
		return nil, false
	}

	index := -1
	for i, agent := range agents {
		if fun(agent) && agent.Status.Phase == k8sCoreV1.PodRunning {
			index = i
			break
		}
	}

	if index == -1 {
		return nil, false
	} else {
		return agents[index], true
	}
}
