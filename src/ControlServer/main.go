package main

import (
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

var stopCh chan struct{}

func notifyStop() {
	close(stopCh)
}

func main() {
	flag.Parse()

	stopCh = make(chan struct{})

	k8SOption, err := LoadK8SOption()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	watcher := NewWatcher(k8SOption)

	reconciles := []TReconcile{
		NewTDeployReconcile(1, k8SOption, watcher),
		NewServiceReconcile(2, k8SOption, watcher),
		NewStatefulSetReconcile(2, k8SOption, watcher),
		NewDaemonSetReconcile(1, k8SOption, watcher),
		NewTEndpointReconcile(4, k8SOption, watcher),
		NewTExitedPodReconcile(2, k8SOption, watcher),
		NewTTreeReconcile(1, k8SOption, watcher),
	}

	watcher.Start(stopCh)
	if !watcher.WaitForCacheSync(stopCh) {
		fmt.Println("WaitForCacheSync Error")
		return
	}

	for _, reconcile := range reconciles {
		reconcile.Start(stopCh)
	}

	webhook := NewWebhook(k8SOption, watcher)
	webhook.Start(stopCh)

	wait.Until(func() {
		time.Sleep(time.Second * 1)
	}, 5, stopCh)
}
