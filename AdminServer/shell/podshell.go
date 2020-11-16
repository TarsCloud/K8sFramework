package shell

import (
	"base"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"net/url"
	"strings"
)

type PodShellImp struct {
	k8sClientSet 	*clientset.Clientset
	k8sConfig    	*restclient.Config
	k8sWatchImp		base.K8SWatchInterface
	k8sNamespace 	string
}

func (ps *PodShellImp) SetK8SClient(clientSet *clientset.Clientset) {
	ps.k8sClientSet = clientSet
}

func (ps *PodShellImp) SetK8SConfig(config *restclient.Config) {
	ps.k8sConfig = config
}

func (ps *PodShellImp) SetK8SNamespace(namespace string) {
	ps.k8sNamespace = namespace
}

func (ps *PodShellImp) SetK8SWatchImp(watchImp base.K8SWatchInterface) {
	ps.k8sWatchImp = watchImp
}

func (ps *PodShellImp) Handler(query url.Values, pty base.PtyHandler) {
	defer func() {
		pty.Done()
	}()

	online, err := CheckQuery(query)
	if err != nil {
		ps.errorProcess(pty, fmt.Sprintf("check query failed: %s\n", err))
		return
	}

	if online {
		ps.handlerOnline(query, pty)
	} else {
		ps.handlerHistory(query, pty)
	}
}

func (ps *PodShellImp) handlerOnline(query url.Values, pty base.PtyHandler) {
	appName 	:= query.Get("appName")
	serverName 	:= query.Get("serverName")

	// shell.sh -> string
	sh	:= fmt.Sprintf("#!/bin/sh\n\ndir=/usr/local/app/taf/app_log/%s/%s\n\nif [ -d $dir ]; then\n  cd $dir\nfi\n\nif [ ! -f /bin/bash ]; then\n  sh\nelse\n  bash\nfi\n",
		appName, serverName)

	podName 	:= query.Get("podName")
	containerName := strings.ToLower(appName) + "-" + strings.ToLower(serverName)

	ps.execCmd(sh, podName, containerName, pty)
}

func (ps *PodShellImp) handlerHistory(query url.Values, pty base.PtyHandler) {
	appName 	:= query.Get("appName")
	serverName 	:= query.Get("serverName")
	podName 	:= query.Get("podName")

	// shell.sh -> string
	sh	:= fmt.Sprintf("#!/bin/sh\n\ndir=/usr/local/app/taf/app_log/%s/%s/%s\n\nif [ -d $dir ]; then\n  cd $dir\nfi\n\nif [ ! -f /bin/bash ]; then\n  sh\nelse\n  bash\nfi\n",
		podName, appName, serverName)

	nodeIP 		:= query.Get("nodeIP")
	daemonPod	:= ps.k8sWatchImp.GetDaemonSetPodByIP(nodeIP)
	if daemonPod == nil {
		ps.errorProcess(pty, fmt.Sprintf("Can not get daemon pod from nodeIP: %s\n", nodeIP))
		return
	}

	// damonset pod
	podName 		 = daemonPod.PodName
	containerName 	:= daemonPod.ContainerName

	ps.execCmd(sh, podName, containerName, pty)
}

func (ps *PodShellImp) getPod(name string) (*corev1.Pod, error) {
	opt := metav1.GetOptions{}
	return ps.k8sClientSet.CoreV1().Pods(ps.k8sNamespace).Get(name, opt)
}

func (ps *PodShellImp) validatePod(pod *corev1.Pod, containerName string) (bool, error) {
	if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
		return false, fmt.Errorf("cannot exec into a container in a completed pod; current phase is %s", pod.Status.Phase)
	}
	for _, c := range pod.Spec.Containers {
		if containerName == c.Name {
			return true, nil
		}
	}
	return false, fmt.Errorf("pod has no container '%s'", containerName)
}


func (ps *PodShellImp) exec(cmd []string, pty base.PtyHandler, podName, containerName string) error {
	req := ps.k8sClientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(ps.k8sNamespace).
		SubResource("exec")

	req.VersionedParams(&corev1.PodExecOptions{
		Container: containerName,
		Command:   cmd,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(ps.k8sConfig, "POST", req.URL())
	if err != nil {
		return err
	}
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		TerminalSizeQueue: pty,
		Tty:               true,
	})
	return err
}

func (ps *PodShellImp) errorProcess(pty base.PtyHandler, msg string) {
	_, _ = pty.Write([]byte(msg))
	fmt.Printf("wsID:%d err, msg:%s\n",pty.GetPtyID(), msg)
}

func (ps *PodShellImp) execCmd(sh, podName, containerName string, pty base.PtyHandler) {
	pod, err := ps.getPod(podName)
	if err != nil {
		ps.errorProcess(pty, fmt.Sprintf("Get kubernetes client failed: %v\n", err))
		return
	}
	ok, err := ps.validatePod(pod, containerName)
	if !ok {
		ps.errorProcess(pty, fmt.Sprintf("Validate pod error! err: %v\n", err))
		return
	}

	cmd := []string{"/bin/sh", "-c", sh}

	err = ps.exec(cmd, pty, podName, containerName)
	if err != nil {
		ps.errorProcess(pty, fmt.Sprintf("Exec to pod error! err: %v\n", err))
		return
	}
}