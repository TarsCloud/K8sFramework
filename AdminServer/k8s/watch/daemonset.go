package watch

import (
	"base"
	k8sApiV1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"sync"
)

const ownerKind = "DaemonSet"
const ownerName = "taf-tafagent"

const containerName = "taf-tafagent"
const podStatus 	= "Running"

type tDaemonPodRecord struct {
	daemonPods  map[string]base.DaemonPodK8S // map[PodName]K8SPod
	mutex     	sync.RWMutex
}

func loadDaemonPodK8SFromPod(daemonPod *k8sApiV1.Pod) *base.DaemonPodK8S {
	daemon := base.NewDaemonPodK8S()

	bDaemon := false

	ownerRefers := daemonPod.OwnerReferences
	for _, ownRefer := range ownerRefers {
		if ownRefer.Name == ownerName && ownRefer.Kind == ownerKind {
			bDaemon = true
		}
	}
	if !bDaemon {
		return nil
	}

	if !daemonPod.Spec.HostNetwork {
		return nil
	}

	daemon.PodName 	= daemonPod.Name
	daemon.HostIP	= daemonPod.Status.HostIP
	daemon.NodeName = daemonPod.Spec.NodeName

	containers := daemonPod.Spec.Containers
	for _, container := range containers {
		if container.Name == containerName {
			if len(container.Ports) > 0 {
				daemon.ContainerName = containerName
				daemon.HostPort = container.Ports[0].HostPort
			}
			break
		}
	}

	if daemonPod.Status.Phase == podStatus {
		daemon.Ready = true
	} else {
		daemon.Ready = false
	}

	return daemon
}

func (c *tDaemonPodRecord) getDamonPodByName(nodeName string) (base.DaemonPodK8S, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, v := range c.daemonPods {
		if v.NodeName == nodeName && v.Ready {
			return v, true
		}
	}
	return *base.NewDaemonPodK8S(), false
}

func (c *tDaemonPodRecord) getDamonPodByIP(nodeIP string) (base.DaemonPodK8S, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, v := range c.daemonPods {
		if v.HostIP == nodeIP && v.Ready {
			return v, true
		}
	}
	return *base.NewDaemonPodK8S(), false
}

func (c *tDaemonPodRecord) onAddDaemonPod(daemonPod *k8sApiV1.Pod) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	daemon := loadDaemonPodK8SFromPod(daemonPod)
	if daemon != nil {
		c.daemonPods[daemon.PodName] = *daemon
	}
}

func (c *tDaemonPodRecord) onDeleteDaemonPod(daemonPod *k8sApiV1.Pod) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	daemon := loadDaemonPodK8SFromPod(daemonPod)
	if daemon != nil {
		delete(c.daemonPods, daemon.PodName)
	}
}

func (c *tDaemonPodRecord) onUpdateDaemonPod(oldDaemonPod, newDaemonPod *k8sApiV1.Pod) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	daemon := loadDaemonPodK8SFromPod(newDaemonPod)
	if daemon != nil {
		c.daemonPods[daemon.PodName] = *daemon
	}
}

func prepareDaemonSetWatch() {
	svcInformer := informerFactory.Core().V1().Pods().Informer()
	svcInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			daemonPod := obj.(*k8sApiV1.Pod)
			daemonPodRecord.onAddDaemonPod(daemonPod)
		},
		DeleteFunc: func(obj interface{}) {
			daemonPod := obj.(*k8sApiV1.Pod)
			daemonPodRecord.onDeleteDaemonPod(daemonPod)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldDaemonPod := oldObj.(*k8sApiV1.Pod)
			newDaemonPod := newObj.(*k8sApiV1.Pod)
			daemonPodRecord.onUpdateDaemonPod(oldDaemonPod, newDaemonPod)
		},
	})
}

var daemonPodRecord tDaemonPodRecord

func init() {
	daemonPodRecord.daemonPods = make(map[string]base.DaemonPodK8S, 5)
	nodeLabelRecord.mutex = sync.RWMutex{}
	registryK8SWatchPrepare(prepareDaemonSetWatch)
}
