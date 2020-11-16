package watch

import (
	"base"
	"encoding/json"
	"fmt"
	k8sAppsV1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/tools/cache"
	"runtime"
	"strings"
	"sync"
)

func loadServerK8SFromStatefulSet(statefulSet *k8sAppsV1.StatefulSet) *base.ServerK8S {
	serverK8s := &base.ServerK8S{
		Replicas:    *statefulSet.Spec.Replicas,
		HostIpc:     statefulSet.Spec.Template.Spec.HostIPC,
		HostNetwork: statefulSet.Spec.Template.Spec.HostNetwork,
		Image:       statefulSet.Spec.Template.Spec.Containers[0].Image,
		Version:     statefulSet.Spec.Template.Labels[base.TafServerVersionLabel],
		NodeSelector: base.NodeSelector{
			Kind: base.AbilityPool,
		},
		HostPort: map[string]int32{},
	}

	nodeSelectorAnnotations, ok := statefulSet.Annotations[base.TafNodeSelectorLabel]
	if !ok {
		//没有 [base.TafNodeSelectorLabel] 注解说明不是符合预期的 statefulSet 格式
		return nil
	}

	if err := json.Unmarshal([]byte(nodeSelectorAnnotations), &serverK8s.NodeSelector); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, fmt.Sprintf("Unexpected NodeSelector Fromat")))
		return nil
	}

	notStackedAnnotations, ok := statefulSet.Annotations[base.TafNotStackedLabel]
	if !ok {
		//没有 [base.TafNotStackedLabel] 注解说明不是符合预期的 statefulSet 格式
		return nil
	}
	if err := json.Unmarshal([]byte(notStackedAnnotations), &serverK8s.NotStacked); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, fmt.Sprintf("Unexpected NotStacked Fromat")))
		return nil
	}

	statefulSetPorts := statefulSet.Spec.Template.Spec.Containers[0].Ports
	if statefulSetPorts != nil {
		for i := range statefulSetPorts {
			if statefulSetPorts[i].HostPort != 0 {
				serverK8s.HostPort[statefulSetPorts[i].Name] = statefulSetPorts[i].HostPort
			}
		}
	}

	return serverK8s
}

type tStatefulSetRecord struct {
	statefulSets map[string]*base.ServerK8S
	mutex        sync.RWMutex
}

func (c *tStatefulSetRecord) onAddStatefulSet(statefulSet *k8sAppsV1.StatefulSet) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	serverK8S := loadServerK8SFromStatefulSet(statefulSet)
	if serverK8S != nil {
		c.statefulSets[statefulSet.Name] = serverK8S
	}
}

func (c *tStatefulSetRecord) onUpdateStatefulSet(oldStatefulSet, newStatefulSet *k8sAppsV1.StatefulSet) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	serverK8S := loadServerK8SFromStatefulSet(newStatefulSet)
	if serverK8S != nil {
		c.statefulSets[newStatefulSet.Name] = serverK8S
	}
}

func (c *tStatefulSetRecord) onDeleteStatefulSet(statefulSet *k8sAppsV1.StatefulSet) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.statefulSets, statefulSet.Name)
}

func (c *tStatefulSetRecord) getServerK8S(app string, name string) *base.ServerK8S {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	k8sResourceName := strings.ToLower(app + "-" + name)
	return c.statefulSets[k8sResourceName]
}

var statefulSetRecord tStatefulSetRecord

func prepareStatefulSetWatch() {
	svcInformer := informerFactory.Apps().V1().StatefulSets().Informer()
	svcInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			statefulSet := obj.(*k8sAppsV1.StatefulSet)
			statefulSetRecord.onAddStatefulSet(statefulSet)
		},
		DeleteFunc: func(obj interface{}) {
			statefulSet := obj.(*k8sAppsV1.StatefulSet)
			statefulSetRecord.onDeleteStatefulSet(statefulSet)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldStatefulSet := oldObj.(*k8sAppsV1.StatefulSet)
			newStatefulSet := newObj.(*k8sAppsV1.StatefulSet)
			statefulSetRecord.onUpdateStatefulSet(oldStatefulSet, newStatefulSet)
		},
	})
}

func init() {
	statefulSetRecord.statefulSets = make(map[string]*base.ServerK8S, 30)
	registryK8SWatchPrepare(prepareStatefulSetWatch)
}
