package main

import (
	"context"
	"fmt"
	k8sCoreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilRuntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
	crdV1Alpha1 "k8s.taf.io/crd/v1alpha1"
	"time"
)

type DaemonSetReconcile struct {
	k8sOption *K8SOption
	watcher   *Watcher
	threads   int
	workQueue workqueue.RateLimitingInterface
}

func NewDaemonSetReconcile(threads int, runtimeOption *K8SOption, watcher *Watcher) *DaemonSetReconcile {
	reconcile := &DaemonSetReconcile{
		k8sOption: runtimeOption,
		watcher:   watcher,
		threads:   threads,
		workQueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), ""),
	}
	watcher.Registry(reconcile)
	return reconcile
}

func (r *DaemonSetReconcile) processItem() bool {

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

func (r *DaemonSetReconcile) EnqueueObj(obj interface{}) {
	switch obj.(type) {
	case *crdV1Alpha1.TServer:
		tserver := obj.(*crdV1Alpha1.TServer)
		key := fmt.Sprintf("%s", tserver.Name)
		r.workQueue.Add(key)
	default:
		return
	}
}

func (r *DaemonSetReconcile) Start(stopCh chan struct{}) {
	for i := 0; i < r.threads; i++ {
		workFun := func() {
			for r.processItem() {
			}
			r.workQueue.ShutDown()
		}
		go wait.Until(workFun, time.Second, stopCh)
	}
}

func (r *DaemonSetReconcile) reconcile(name string) ReconcileResult {
	namespace := r.k8sOption.namespace

	tserver, err := r.watcher.tServerLister.TServers(namespace).Get(name)
	if err != nil {
		if !errors.IsNotFound(err) {
			utilRuntime.HandleError(fmt.Errorf(ResourceGetError, "tserver", namespace, name, err.Error()))
			return RateLimit
		}
		err = r.k8sOption.k8sClientSet.AppsV1().DaemonSets(namespace).Delete(context.TODO(), name, k8sMetaV1.DeleteOptions{})
		if err != nil {
			utilRuntime.HandleError(fmt.Errorf(ResourceDeleteError, "daemonset", namespace, name, err.Error()))
			return RateLimit
		}
		return AllOk
	}

	daemonSet, err := r.watcher.daemonSetLister.DaemonSets(namespace).Get(name)
	if err != nil {
		if !errors.IsNotFound(err) {
			utilRuntime.HandleError(fmt.Errorf(ResourceGetError, "daemonset", namespace, name, err.Error()))
			return RateLimit
		}
		if tserver.Spec.K8S.NodeSelector.DaemonSet != nil {
			daemonSet = buildDaemonSet(tserver)
			daemonSetInterface := r.k8sOption.k8sClientSet.AppsV1().DaemonSets(namespace)
			if _, err = daemonSetInterface.Create(context.TODO(), daemonSet, k8sMetaV1.CreateOptions{}); err != nil {
				utilRuntime.HandleError(fmt.Errorf(ResourceCreateError, "daemonset", namespace, name, err.Error()))
				return RateLimit
			}
		}
		return AllOk
	}

	if tserver.DeletionTimestamp != nil {
		err = r.k8sOption.k8sClientSet.AppsV1().DaemonSets(namespace).Delete(context.TODO(), name, k8sMetaV1.DeleteOptions{})
		if err != nil {
			utilRuntime.HandleError(fmt.Errorf(ResourceDeleteError, "daemonset", namespace, name, err.Error()))
			return RateLimit
		}
		return AllOk
	}

	if !isOwnerByTServer(tserver, daemonSet) {
		// 此处意味着出现了非由 controller 管理的同名 daemonSet, 需要警告和重试
		msg := fmt.Sprintf(ResourceOutControlError, "daemonset", namespace, daemonSet.Name, namespace, name)
		r.k8sOption.recorder.Event(tserver, k8sCoreV1.EventTypeWarning, ResourceOutControlReason, msg)
		return RateLimit
	}

	anyChanged := !equalTServerAndDaemonSet(tserver, daemonSet)

	if anyChanged {
		daemonSetCopy := daemonSet.DeepCopy()
		syncDaemonSet(tserver, daemonSetCopy)
		daemonSetInterface := r.k8sOption.k8sClientSet.AppsV1().DaemonSets(namespace)
		if _, err := daemonSetInterface.Update(context.TODO(), daemonSetCopy, k8sMetaV1.UpdateOptions{}); err != nil {
			utilRuntime.HandleError(fmt.Errorf(ResourceUpdateError, "daemonset", namespace, name, err.Error()))
			return RateLimit
		}
	}

	return AllOk
}
