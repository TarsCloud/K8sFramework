package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilRuntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
	crdV1Alpha1 "k8s.tars.io/api/crd/v1alpha1"
)

type TDeployReconcile struct {
	k8sOption *K8SOption
	watcher   *Watcher
	threads   int
	workQueue workqueue.RateLimitingInterface
}

func NewTDeployReconcile(threads int, runtimeOption *K8SOption, watcher *Watcher) *TDeployReconcile {
	reconcile := &TDeployReconcile{
		k8sOption: runtimeOption,
		watcher:   watcher,
		threads:   threads,
		workQueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), ""),
	}
	watcher.Registry(reconcile)
	return reconcile
}

func (r *TDeployReconcile) processItem() bool {

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

func (r *TDeployReconcile) EnqueueObj(obj interface{}) {
	switch obj.(type) {
	case *crdV1Alpha1.TDeploy:
		tdeploy := obj.(*crdV1Alpha1.TDeploy)
		if tdeploy.Deployed == nil || !*tdeploy.Deployed {
			if tdeploy.Approve != nil && tdeploy.Approve.Result {
				key := tdeploy.Name
				r.workQueue.Add(key)
			}
		}
	default:
		return
	}
}

func (r *TDeployReconcile) Start(stopCh chan struct{}) {
	for i := 0; i < r.threads; i++ {
		workFun := func() {
			for r.processItem() {
			}
			r.workQueue.ShutDown()
		}
		go wait.Until(workFun, time.Second, stopCh)
	}
}

func (r *TDeployReconcile) reconcile(name string) ReconcileResult {
	namespace := r.k8sOption.namespace
	tdeploy, err := r.watcher.tDeployLister.TDeploys(namespace).Get(name)
	if err != nil {
		if !errors.IsNotFound(err) {
			utilRuntime.HandleError(fmt.Errorf(ResourceGetError, "tdeploy", namespace, name, err.Error()))
			return RateLimit
		}
		return AllOk
	}

	if tdeploy.Approve == nil || !tdeploy.Approve.Result {
		return AllOk
	}

	if tdeploy.Deployed != nil && *tdeploy.Deployed {
		return AllOk
	}

	tserverName := fmt.Sprintf("%s-%s", strings.ToLower(tdeploy.Apply.App), strings.ToLower(tdeploy.Apply.Server))

	newTServer := &crdV1Alpha1.TServer{
		TypeMeta: k8sMetaV1.TypeMeta{},
		ObjectMeta: k8sMetaV1.ObjectMeta{
			Name:      tserverName,
			Namespace: namespace,
			Labels: map[string]string{
				TServerAppLabel:  tdeploy.Apply.App,
				TServerNameLabel: tdeploy.Apply.Server,
				TSubTypeLabel:    string(tdeploy.Apply.SubType),
			},
		},
		Spec: tdeploy.Apply,
	}

	if _, err = r.k8sOption.crdClientSet.CrdV1alpha1().TServers(namespace).Create(context.TODO(), newTServer, k8sMetaV1.CreateOptions{}); err != nil && !errors.IsAlreadyExists(err) {
		utilRuntime.HandleError(fmt.Errorf(ResourceCreateError, "newTServer", namespace, newTServer.Name, err.Error()))
		return RateLimit
	}

	deployed := true
	tdeployCopy := tdeploy.DeepCopy()
	tdeployCopy.Deployed = &deployed
	if _, err := r.k8sOption.crdClientSet.CrdV1alpha1().TDeploys(tdeploy.Namespace).Update(context.TODO(), tdeployCopy, k8sMetaV1.UpdateOptions{}); err != nil {
		utilRuntime.HandleError(fmt.Errorf(ResourceUpdateError, "tdeploy", namespace, name, err.Error()))
		return RateLimit
	}

	return AllOk

}
