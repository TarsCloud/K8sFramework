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

type ServiceReconcile struct {
	k8sOption *K8SOption
	watcher   *Watcher
	threads   int
	workQueue workqueue.RateLimitingInterface
}

func NewServiceReconcile(threads int, runtimeOption *K8SOption, watcher *Watcher) *ServiceReconcile {
	reconcile := &ServiceReconcile{
		k8sOption: runtimeOption,
		watcher:   watcher,
		threads:   threads,
		workQueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), ""),
	}
	watcher.Registry(reconcile)
	return reconcile
}

func (r *ServiceReconcile) processItem() bool {

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

func (r *ServiceReconcile) EnqueueObj(obj interface{}) {
	switch obj.(type) {
	case *crdV1Alpha1.TServer:
		tserver := obj.(*crdV1Alpha1.TServer)
		key := fmt.Sprintf("%s", tserver.Name)
		r.workQueue.Add(key)
	case *k8sCoreV1.Service:
		service := obj.(*k8sCoreV1.Service)
		if ownerRef := service.GetOwnerReferences(); ownerRef != nil {
			for i := range ownerRef {
				if ownerRef[i].Kind == TServerKind && ownerRef[i].APIVersion == TServerAPIVersion {
					key := fmt.Sprintf("%s", ownerRef[i].Name)
					r.workQueue.Add(key)
					return
				}
			}
		}
	default:
		return
	}
}

func (r *ServiceReconcile) Start(stopCh chan struct{}) {
	for i := 0; i < r.threads; i++ {
		workFun := func() {
			for r.processItem() {
			}
			r.workQueue.ShutDown()
		}
		go wait.Until(workFun, time.Second, stopCh)
	}
}

func (r *ServiceReconcile) reconcile(name string) ReconcileResult {
	namespace := r.k8sOption.namespace
	tserver, err := r.watcher.tServerLister.TServers(namespace).Get(name)
	if err != nil {
		if !errors.IsNotFound(err) {
			utilRuntime.HandleError(fmt.Errorf(ResourceGetError, "tserver", namespace, name, err.Error()))
			return RateLimit
		}
		err = r.k8sOption.k8sClientSet.CoreV1().Services(namespace).Delete(context.TODO(), name, k8sMetaV1.DeleteOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				utilRuntime.HandleError(fmt.Errorf(ResourceDeleteError, "service", namespace, name, err.Error()))
			}
			return RateLimit
		}
		return AllOk
	}

	if tserver.DeletionTimestamp != nil {
		err = r.k8sOption.k8sClientSet.CoreV1().Services(namespace).Delete(context.TODO(), name, k8sMetaV1.DeleteOptions{})
		if err != nil {
			utilRuntime.HandleError(fmt.Errorf(ResourceDeleteError, "service", namespace, name, err.Error()))
			return RateLimit
		}
		return AllOk
	}

	service, err := r.watcher.serviceLister.Services(namespace).Get(name)
	if err != nil {
		if !errors.IsNotFound(err) {
			utilRuntime.HandleError(fmt.Errorf(ResourceGetError, "service", namespace, name, err.Error()))
			return RateLimit
		}
		service = buildService(tserver)
		serviceInterface := r.k8sOption.k8sClientSet.CoreV1().Services(namespace)
		if _, err = serviceInterface.Create(context.TODO(), service, k8sMetaV1.CreateOptions{}); err != nil {
			utilRuntime.HandleError(fmt.Errorf(ResourceCreateError, "service", namespace, name, err.Error()))
			return RateLimit
		}
		return AllOk
	}

	if !isOwnerByTServer(tserver, service) {
		// 此处意味着出现了非由 controller 管理的同名 service, 需要警告和重试
		msg := fmt.Sprintf(ResourceOutControlError, "service", namespace, service.Name, namespace, name)
		r.k8sOption.recorder.Event(tserver, k8sCoreV1.EventTypeWarning, ResourceOutControlReason, msg)
		return RateLimit
	}

	anyChanged := !equalTServerAndService(tserver, service)

	if anyChanged {
		serviceCopy := service.DeepCopy()
		syncService(tserver, serviceCopy)
		serviceInterface := r.k8sOption.k8sClientSet.CoreV1().Services(namespace)
		if _, err := serviceInterface.Update(context.TODO(), serviceCopy, k8sMetaV1.UpdateOptions{}); err != nil {
			utilRuntime.HandleError(fmt.Errorf(ResourceUpdateError, "service", namespace, name, err.Error()))
			return RateLimit
		}
	}

	return AllOk
}
