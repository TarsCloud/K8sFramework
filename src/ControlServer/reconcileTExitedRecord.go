package main

import (
	"context"
	"encoding/json"
	"fmt"
	k8sCoreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	patchTypes "k8s.io/apimachinery/pkg/types"
	utilRuntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
	crdV1Alpha1 "k8s.taf.io/crd/v1alpha1"
	"strings"
	"time"
)

type TExitedPodReconcile struct {
	k8sOption *K8SOption
	watcher   *Watcher
	threads   int
	workQueue workqueue.RateLimitingInterface
}

func NewTExitedPodReconcile(threads int, runtimeOption *K8SOption, watcher *Watcher) *TExitedPodReconcile {
	reconcile := &TExitedPodReconcile{
		k8sOption: runtimeOption,
		watcher:   watcher,
		threads:   threads,
		workQueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), ""),
	}
	watcher.Registry(reconcile)
	return reconcile
}

func (r *TExitedPodReconcile) EnqueueObj(obj interface{}) {
	switch obj.(type) {
	case *crdV1Alpha1.TServer:
		tserver := obj.(*crdV1Alpha1.TServer)
		key := fmt.Sprintf("tserver/%s", tserver.Name)
		r.workQueue.Add(key)
	case *crdV1Alpha1.TExitedRecord:
		tserver := obj.(*crdV1Alpha1.TExitedRecord)
		key := fmt.Sprintf("tserver/%s", tserver.Name)
		r.workQueue.Add(key)
	case *k8sCoreV1.Pod:
		pod := obj.(*k8sCoreV1.Pod)
		if pod.DeletionTimestamp != nil && pod.Labels != nil {
			app, appExist := pod.Labels[TServerAppLabel]
			server, serverExist := pod.Labels[TServerNameLabel]
			if appExist && serverExist {
				tExitedEvent := &crdV1Alpha1.TExitedRecord{
					App:    app,
					Server: server,
					Pods: []crdV1Alpha1.TExitedPod{
						{
							UID:        string(pod.UID),
							Name:       pod.Name,
							Tag:        pod.Labels[TServerTagLabel],
							NodeIP:     pod.Status.HostIP,
							PodIP:      pod.Status.PodIP,
							CreateTime: pod.CreationTimestamp,
							DeleteTime: *pod.DeletionTimestamp,
						},
					},
				}
				bs, _ := json.Marshal(tExitedEvent)
				key := fmt.Sprintf("event/%s", bs)
				r.workQueue.Add(key)
				return
			}
		}
	default:
		return
	}
}

func (r *TExitedPodReconcile) splitKey(key string) []string {
	return strings.Split(key, "/")
}

func (r *TExitedPodReconcile) processItem() bool {

	obj, shutdown := r.workQueue.Get()

	if shutdown {
		return false
	}

	defer r.workQueue.Done(obj)

	key, ok := obj.(string)
	if !ok {
		utilRuntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
		r.workQueue.Forget(obj)
		return true
	}

	var res ReconcileResult
	v := r.splitKey(key)
	if len(v) != 2 {
		//todo log error
		res = AllOk
	} else {
		switch v[0] {
		case "tserver":
			res = r.reconcileBaseTServer(v[1])
		case "event":
			res = r.reconcileBasePod(v[1])
		default:
			res = AllOk
		}
	}
	switch res {
	case AllOk:
		r.workQueue.Forget(obj)
		return true
	case RateLimit:
		r.workQueue.AddRateLimited(obj)
		return true
	case FatalError:
		r.workQueue.ShutDown()
		return false
	default:
		//code should not reach here
		utilRuntime.HandleError(fmt.Errorf("should not reach place"))
		return false
	}
}

func (r *TExitedPodReconcile) Start(stopCh chan struct{}) {
	for i := 0; i < r.threads; i++ {
		workFun := func() {
			for r.processItem() {
			}
			r.workQueue.ShutDown()
		}
		go wait.Until(workFun, time.Second, stopCh)
	}
}

func (r *TExitedPodReconcile) reconcileBaseTServer(name string) ReconcileResult {
	namespace := r.k8sOption.namespace
	tserver, err := r.watcher.tServerLister.TServers(namespace).Get(name)
	if err != nil {
		if !errors.IsNotFound(err) {
			utilRuntime.HandleError(fmt.Errorf(ResourceGetError, "tserver", namespace, name, err.Error()))
			return RateLimit
		}
		err = r.k8sOption.crdClientSet.CrdV1alpha1().TExitedRecords(namespace).Delete(context.TODO(), name, k8sMetaV1.DeleteOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				return AllOk
			}
			utilRuntime.HandleError(fmt.Errorf(ResourceDeleteError, "texitedrecord", namespace, name, err.Error()))
			return RateLimit
		}
		return AllOk
	}

	if tserver.DeletionTimestamp != nil {
		err = r.k8sOption.crdClientSet.CrdV1alpha1().TExitedRecords(namespace).Delete(context.TODO(), name, k8sMetaV1.DeleteOptions{})
		if err != nil {
			utilRuntime.HandleError(fmt.Errorf(ResourceDeleteError, "texitedrecord", namespace, name, err.Error()))
			return RateLimit
		}
		return AllOk
	}

	tExitedRecord, err := r.watcher.tExitedRecordLister.TExitedRecords(namespace).Get(name)
	if err != nil {
		if !errors.IsNotFound(err) {
			utilRuntime.HandleError(fmt.Errorf(ResourceGetError, "texitedrecord", namespace, name, err.Error()))
			return RateLimit
		}
		tExitedRecord = buildTExitedRecord(tserver)
		tExitedPodInterface := r.k8sOption.crdClientSet.CrdV1alpha1().TExitedRecords(namespace)
		if _, err = tExitedPodInterface.Create(context.TODO(), tExitedRecord, k8sMetaV1.CreateOptions{}); err != nil {
			utilRuntime.HandleError(fmt.Errorf(ResourceCreateError, "texitedrecord", namespace, name, err.Error()))
			return RateLimit
		}
		return AllOk
	}
	return AllOk
}

func (r *TExitedPodReconcile) reconcileBasePod(tExitedPodSpecString string) ReconcileResult {
	namespace := r.k8sOption.namespace
	var tExitedEvent crdV1Alpha1.TExitedRecord
	_ = json.Unmarshal([]byte(tExitedPodSpecString), &tExitedEvent)

	tExitedRecordName := fmt.Sprintf("%s-%s", strings.ToLower(tExitedEvent.App), strings.ToLower(tExitedEvent.Server))
	tExitedRecord, err := r.watcher.tExitedRecordLister.TExitedRecords(namespace).Get(tExitedRecordName)
	if err != nil {
		if !errors.IsNotFound(err) {
			utilRuntime.HandleError(fmt.Errorf(ResourceGetError, "texitedrecord", namespace, tExitedRecordName, err.Error()))
			return RateLimit
		}
		return AllOk
	}

	recordedPodsLen := len(tExitedRecord.Pods)

	const MaxCheckLen = 5
	var sentinelValue int

	if recordedPodsLen <= MaxCheckLen {
		sentinelValue = 0
	} else {
		sentinelValue = recordedPodsLen - MaxCheckLen
	}

	for i := recordedPodsLen - 1; i >= sentinelValue; i-- {
		if tExitedRecord.Pods[i].UID == tExitedEvent.Pods[0].UID {
			// means exited events had recorded
			return AllOk
		}
	}

	const MaxRecordLen = 80
	var patchString string

	exitedPodString, _ := json.Marshal(tExitedEvent.Pods[0])

	if recordedPodsLen < MaxRecordLen {
		patchString = fmt.Sprintf("[{\"op\":\"add\",\"path\":\"/pods/-\",\"value\":%s}]", exitedPodString)
	} else {
		patchString = fmt.Sprintf("[{\"op\":\"remove\",\"path\":\"/pods/0\"},{\"op\":\"add\",\"path\":\"/pods/-\",\"value\":%s}]", exitedPodString)
	}

	_, err = r.k8sOption.crdClientSet.CrdV1alpha1().TExitedRecords(namespace).Patch(context.TODO(), tExitedRecordName, patchTypes.JSONPatchType, []byte(patchString), k8sMetaV1.PatchOptions{})

	if err != nil {
		utilRuntime.HandleError(fmt.Errorf(ResourcePatchError, "texitedrecord", namespace, tExitedRecordName, err.Error()))
		return RateLimit
	}

	return AllOk
}
