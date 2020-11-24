package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	k8sCoreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	patchTypes "k8s.io/apimachinery/pkg/types"
	utilRuntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
	crdV1Alpha1 "k8s.tars.io/api/crd/v1alpha1"
)

type TTreeReconcile struct {
	k8sOption *K8SOption
	watcher   *Watcher
	threads   int
	workQueue workqueue.RateLimitingInterface
}

func NewTTreeReconcile(threads int, runtimeOption *K8SOption, watcher *Watcher) *TTreeReconcile {
	reconcile := &TTreeReconcile{
		k8sOption: runtimeOption,
		watcher:   watcher,
		threads:   threads,
		workQueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), ""),
	}
	watcher.Registry(reconcile)
	return reconcile
}

func (r *TTreeReconcile) processItem() bool {

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

	res := r.reconcile(key)

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

func (r *TTreeReconcile) EnqueueObj(obj interface{}) {
	switch obj.(type) {
	case *crdV1Alpha1.TServer:
		tserver := obj.(*crdV1Alpha1.TServer)
		key := fmt.Sprintf("%s", tserver.Name)
		r.workQueue.Add(key)
	default:
		return
	}
}

func (r *TTreeReconcile) Start(stopCh chan struct{}) {
	for i := 0; i < r.threads; i++ {
		workFun := func() {
			for r.processItem() {
			}
			r.workQueue.ShutDown()
		}
		go wait.Until(workFun, time.Second, stopCh)
	}
}

func (r *TTreeReconcile) reconcile(name string) ReconcileResult {
	namespace := r.k8sOption.namespace
	tserver, err := r.watcher.tServerLister.TServers(namespace).Get(name)
	if err != nil {
		if !errors.IsNotFound(err) {
			utilRuntime.HandleError(fmt.Errorf(ResourceGetError, "tserver", namespace, name, err.Error()))
			return RateLimit
		}
		return AllOk
	}

	if tserver.DeletionTimestamp != nil {
		return AllOk
	}

	ttree, err := r.watcher.tTreeLister.TTrees(namespace).Get(TarsTreeResourceName)
	if err != nil {
		msg := fmt.Sprintf(ResourceGetError, "ttree", namespace, TarsTreeResourceName, err.Error())
		utilRuntime.HandleError(fmt.Errorf(msg))
		r.k8sOption.recorder.Event(tserver, k8sCoreV1.EventTypeWarning, ResourceGetReason, msg)
		return RateLimit
	}

	for i := range ttree.Apps {
		if ttree.Apps[i].Name == tserver.Spec.App {
			return AllOk
		}
	}

	newTressApp := &crdV1Alpha1.TTreeApp{
		Name:         tserver.Spec.App,
		BusinessRef:  "",
		CreatePerson: TarsControlServiceAccount,
		CreateTime:   k8sMetaV1.Now(),
		Mark:         "AddByControl",
	}

	bs, _ := json.Marshal(newTressApp)
	patchContent := fmt.Sprintf("[{\"op\":\"add\",\"path\":\"/apps/-\",\"value\":%s}]", bs)
	_, err = r.k8sOption.crdClientSet.CrdV1alpha1().TTrees(namespace).Patch(context.TODO(), TarsTreeResourceName, patchTypes.JSONPatchType, []byte(patchContent), k8sMetaV1.PatchOptions{})
	if err != nil {
		return RateLimit
	}

	return AllOk
}
