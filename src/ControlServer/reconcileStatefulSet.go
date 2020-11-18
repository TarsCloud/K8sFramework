package main

import (
	"context"
	"fmt"
	k8sAppsV1 "k8s.io/api/apps/v1"
	k8sCoreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilRuntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
	crdV1Alpha1 "k8s.taf.io/crd/v1alpha1"
	"time"
)

type StatefulSetReconcile struct {
	k8sOption *K8SOption
	watcher   *Watcher
	threads   int
	workQueue workqueue.RateLimitingInterface
}

func NewStatefulSetReconcile(threads int, runtimeOption *K8SOption, watcher *Watcher) *StatefulSetReconcile {
	reconcile := &StatefulSetReconcile{
		k8sOption: runtimeOption,
		watcher:   watcher,
		threads:   threads,
		workQueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), ""),
	}
	watcher.Registry(reconcile)
	return reconcile
}

func (r *StatefulSetReconcile) processItem() bool {

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

func (r *StatefulSetReconcile) EnqueueObj(obj interface{}) {
	switch obj.(type) {
	case *crdV1Alpha1.TServer:
		tserver := obj.(*crdV1Alpha1.TServer)
		key := fmt.Sprintf("%s", tserver.Name)
		r.workQueue.Add(key)
	case *k8sAppsV1.StatefulSet:
		statefulSet := obj.(*k8sAppsV1.StatefulSet)
		if ownerRef := statefulSet.GetOwnerReferences(); ownerRef != nil {
			for i := range ownerRef {
				if ownerRef[i].Kind == TServerKind && ownerRef[i].APIVersion == TServerAPIVersion {
					key := fmt.Sprintf("%s", ownerRef[i].Name)
					r.workQueue.Add(key)
					return
				}
			}
		}
		return
	}
}

func (r *StatefulSetReconcile) Start(stopCh chan struct{}) {
	for i := 0; i < r.threads; i++ {
		workFun := func() {
			for r.processItem() {
			}
			r.workQueue.ShutDown()
		}
		go wait.Until(workFun, time.Second, stopCh)
	}
}

func (r *StatefulSetReconcile) reconcile(name string) ReconcileResult {
	namespace := r.k8sOption.namespace

	tserver, err := r.watcher.tServerLister.TServers(namespace).Get(name)
	if err != nil {
		if !errors.IsNotFound(err) {
			utilRuntime.HandleError(fmt.Errorf(ResourceGetError, "tserver", namespace, name, err.Error()))
			return RateLimit
		}
		err = r.k8sOption.k8sClientSet.AppsV1().StatefulSets(namespace).Delete(context.TODO(), name, k8sMetaV1.DeleteOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				return AllOk
			}
			utilRuntime.HandleError(fmt.Errorf(ResourceDeleteError, "statefulset", namespace, name, err.Error()))
			return RateLimit
		}
		return AllOk
	}

	if tserver.DeletionTimestamp != nil || tserver.Spec.K8S.NodeSelector.DaemonSet != nil {
		err = r.k8sOption.k8sClientSet.AppsV1().StatefulSets(namespace).Delete(context.TODO(), name, k8sMetaV1.DeleteOptions{})
		if err != nil {
			utilRuntime.HandleError(fmt.Errorf(ResourceDeleteError, "statefulset", namespace, name, err.Error()))
			return RateLimit
		}
		return AllOk
	}

	statefulSet, err := r.watcher.statefulSetLister.StatefulSets(namespace).Get(name)
	if err != nil {
		if !errors.IsNotFound(err) {
			utilRuntime.HandleError(fmt.Errorf(ResourceGetError, "statefulset", namespace, name, err.Error()))
			return RateLimit
		}
		if tserver.Spec.K8S.NodeSelector.DaemonSet == nil {
			statefulSet = buildStatefulSet(tserver)
			statefulSetInterface := r.k8sOption.k8sClientSet.AppsV1().StatefulSets(namespace)
			if _, err = statefulSetInterface.Create(context.TODO(), statefulSet, k8sMetaV1.CreateOptions{}); err != nil {
				utilRuntime.HandleError(fmt.Errorf(ResourceCreateError, "statefulset", namespace, name, err.Error()))
				return RateLimit
			}
		}
		return AllOk
	}

	if tserver.DeletionTimestamp != nil {
		err = r.k8sOption.k8sClientSet.AppsV1().StatefulSets(namespace).Delete(context.TODO(), name, k8sMetaV1.DeleteOptions{})
		if err != nil {
			utilRuntime.HandleError(fmt.Errorf(ResourceDeleteError, "statefulset", namespace, name, err.Error()))
			return RateLimit
		}
		return AllOk
	}

	if !isOwnerByTServer(tserver, statefulSet) {
		// 此处意味着出现了非由 controller 管理的同名 statefulSet, 需要警告和重试
		msg := fmt.Sprintf(ResourceOutControlError, "statefulset", namespace, statefulSet.Name, namespace, name)
		r.k8sOption.recorder.Event(tserver, k8sCoreV1.EventTypeWarning, ResourceOutControlReason, msg)
		return RateLimit
	}

	anyChanged := !equalTServerAndStatefulSet(tserver, statefulSet)

	if anyChanged {
		statefulSetCopy := statefulSet.DeepCopy()
		syncStatefulSet(tserver, statefulSetCopy)
		statefulSetInterface := r.k8sOption.k8sClientSet.AppsV1().StatefulSets(namespace)
		if _, err := statefulSetInterface.Update(context.TODO(), statefulSetCopy, k8sMetaV1.UpdateOptions{}); err != nil {
			utilRuntime.HandleError(fmt.Errorf(ResourceUpdateError, "statefulset", namespace, name, err.Error()))
			return RateLimit
		}
	}
	return AllOk
}