package main

import (
	"context"
	"fmt"
	k8sCoreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	utilRuntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/workqueue"
	crdV1Alpha1 "k8s.taf.io/crd/v1alpha1"
	"strings"
	"time"
)

type TEndpointReconcile struct {
	k8sOption *K8SOption
	watcher   *Watcher
	threads   int
	workQueue workqueue.RateLimitingInterface
}

func NewTEndpointReconcile(threads int, runtimeOption *K8SOption, watcher *Watcher) *TEndpointReconcile {
	reconcile := &TEndpointReconcile{
		k8sOption: runtimeOption,
		watcher:   watcher,
		threads:   threads,
		workQueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), ""),
	}
	watcher.Registry(reconcile)
	return reconcile
}

func (r *TEndpointReconcile) processItem() bool {

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
	res = r.reconcile(key)

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

func (r *TEndpointReconcile) EnqueueObj(obj interface{}) {
	switch obj.(type) {
	case *crdV1Alpha1.TServer:
		tserver := obj.(*crdV1Alpha1.TServer)
		key := tserver.Name
		r.workQueue.Add(key)
	case *crdV1Alpha1.TEndpoint:
		tendpoint := obj.(*crdV1Alpha1.TEndpoint)
		key := tendpoint.Name
		r.workQueue.Add(key)
	case *k8sCoreV1.Pod:
		pod := obj.(*k8sCoreV1.Pod)
		if pod.Labels != nil {
			app, appExist := pod.Labels[TServerAppLabel]
			server, serverExist := pod.Labels[TServerNameLabel]
			if appExist && serverExist {
				key := fmt.Sprintf("%s-%s", strings.ToLower(app), strings.ToLower(server))
				r.workQueue.Add(key)
				return
			}
		}
	default:
		return
	}
}

func (r *TEndpointReconcile) Start(stopCh chan struct{}) {
	for i := 0; i < r.threads; i++ {
		workFun := func() {
			for r.processItem() {
			}
			r.workQueue.ShutDown()
		}
		go wait.Until(workFun, time.Second, stopCh)
	}
}

func (r *TEndpointReconcile) reconcile(name string) ReconcileResult {
	namespace := r.k8sOption.namespace

	tserver, err := r.watcher.tServerLister.TServers(namespace).Get(name)
	if err != nil {
		if !errors.IsNotFound(err) {
			utilRuntime.HandleError(fmt.Errorf(ResourceGetError, "tserver", namespace, name, err.Error()))
			return RateLimit
		}
		err = r.k8sOption.crdClientSet.CrdV1alpha1().TEndpoints(namespace).Delete(context.TODO(), name, k8sMetaV1.DeleteOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				return AllOk
			}
			utilRuntime.HandleError(fmt.Errorf(ResourceDeleteError, "tendpoint", namespace, name, err.Error()))
			return RateLimit
		}
		return AllOk
	}

	if tserver.DeletionTimestamp != nil {
		err = r.k8sOption.crdClientSet.CrdV1alpha1().TEndpoints(namespace).Delete(context.TODO(), name, k8sMetaV1.DeleteOptions{})
		if err != nil {
			utilRuntime.HandleError(fmt.Errorf(ResourceDeleteError, "tendpoint", namespace, name, err.Error()))
			return RateLimit
		}
		return AllOk
	}

	tendpoint, err := r.watcher.tEndpointLister.TEndpoints(namespace).Get(name)
	if err != nil {
		if !errors.IsNotFound(err) {
			utilRuntime.HandleError(fmt.Errorf(ResourceGetError, "tendpoint", namespace, name, err.Error()))
			return RateLimit
		}
		tendpoint = buildTEndpoint(tserver)
		tendpointInterface := r.k8sOption.crdClientSet.CrdV1alpha1().TEndpoints(namespace)
		if _, err = tendpointInterface.Create(context.TODO(), tendpoint, k8sMetaV1.CreateOptions{}); err != nil {
			utilRuntime.HandleError(fmt.Errorf(ResourceCreateError, "tendpoint", namespace, name, err.Error()))
			return RateLimit
		}
		return AllOk
	}

	if !isOwnerByTServer(tserver, tendpoint) {
		// 此处意味着出现了非由 controller 管理的同名 tendpoint, 需要警告和重试
		msg := fmt.Sprintf(ResourceOutControlError, "tendpoint", namespace, tendpoint.Name, namespace, name)
		r.k8sOption.recorder.Event(tserver, k8sCoreV1.EventTypeWarning, ResourceOutControlReason, msg)
		return RateLimit
	}

	anyChanged := !equalTServerAndTEndpoint(tserver, tendpoint)

	if anyChanged {
		tendpointCopy := tendpoint.DeepCopy()
		syncTEndpoint(tserver, tendpointCopy)
		tendpointInterface := r.k8sOption.crdClientSet.CrdV1alpha1().TEndpoints(namespace)
		if _, err := tendpointInterface.Update(context.TODO(), tendpointCopy, k8sMetaV1.UpdateOptions{}); err != nil {
			utilRuntime.HandleError(fmt.Errorf(ResourceUpdateError, "tendpoint", namespace, name, err.Error()))
			return RateLimit
		}
	}
	return r.updateStatus(tendpoint)
}

func (r *TEndpointReconcile) updateStatus(tendpoint *crdV1Alpha1.TEndpoint) ReconcileResult {

	appRequirement, _ := labels.NewRequirement(TServerAppLabel, selection.DoubleEquals, []string{tendpoint.Spec.App})
	serverRequirement, _ := labels.NewRequirement(TServerNameLabel, selection.DoubleEquals, []string{tendpoint.Spec.Server})

	pods, err := r.watcher.podLister.Pods(tendpoint.Namespace).List(labels.NewSelector().Add(*appRequirement).Add(*serverRequirement))
	if err != nil {
		utilRuntime.HandleError(fmt.Errorf(ResourceSelectorError, tendpoint.Namespace, "tendpoint", err.Error()))
		return RateLimit
	}

	tendpointPodStatuses := make([]*crdV1Alpha1.TEndpointPodStatus, 0, len(pods))

	for _, pod := range pods {
		podStatus := &crdV1Alpha1.TEndpointPodStatus{
			UID:               string(pod.UID),
			Name:              pod.Name,
			PodIP:             pod.Status.PodIP,
			HostIP:            pod.Status.HostIP,
			StartTime:         pod.CreationTimestamp,
			ContainerStatuses: pod.Status.ContainerStatuses,
			Tag:               pod.Labels[TServerTagLabel],
		}

		if pod.DeletionTimestamp != nil {
			podStatus.SettingState = "Active"
			podStatus.PresentState = "Terminating"
			podStatus.PresentMessage = fmt.Sprintf("pod/%s is terminating", pod.Name)
			tendpointPodStatuses = append(tendpointPodStatuses, podStatus)
			continue
		}

		if len(pod.Status.Conditions) == 0 {
			podStatus.SettingState = "Active"
			podStatus.PresentState = "UnKnown"
			podStatus.PresentMessage = ""
			tendpointPodStatuses = append(tendpointPodStatuses, podStatus)
			continue
		}

		orderConditions := make([]*k8sCoreV1.PodCondition, 4)
		var readinessGatesCondition *k8sCoreV1.PodCondition

		i := 0
		for i = range pod.Status.Conditions {
			switch pod.Status.Conditions[i].Type {
			case k8sCoreV1.PodScheduled:
				orderConditions[0] = &pod.Status.Conditions[i]
			case k8sCoreV1.PodInitialized:
				orderConditions[1] = &pod.Status.Conditions[i]
			case k8sCoreV1.ContainersReady:
				orderConditions[2] = &pod.Status.Conditions[i]
			case k8sCoreV1.PodReady:
				orderConditions[3] = &pod.Status.Conditions[i]
			case TPodReadinessGate:
				readinessGatesCondition = &pod.Status.Conditions[i]
			}
		}

		if i < 3 {
			podStatus.SettingState = "Active"
			podStatus.PresentState = orderConditions[i].Reason
			podStatus.PresentMessage = orderConditions[i].Message
			tendpointPodStatuses = append(tendpointPodStatuses, podStatus)
			continue
		}

		switch orderConditions[3].Status {
		case k8sCoreV1.ConditionTrue:
			if readinessGatesCondition == nil {
				podStatus.SettingState = "Active"
				podStatus.PresentState = "Started:"
				podStatus.PresentMessage = fmt.Sprintf("pod/%s is started", pod.Name)
				tendpointPodStatuses = append(tendpointPodStatuses, podStatus)
				continue
			}

			v := strings.Split(readinessGatesCondition.Reason, "/")
			podStatus.SettingState = v[0]
			podStatus.PresentState = v[1]
			podStatus.PresentMessage = ""
			tendpointPodStatuses = append(tendpointPodStatuses, podStatus)
			continue
		default:
			if orderConditions[3].Reason != "ReadinessGatesNotReady" {
				podStatus.SettingState = "Active"
				podStatus.PresentState = orderConditions[3].Reason
				podStatus.PresentMessage = orderConditions[3].Message
				tendpointPodStatuses = append(tendpointPodStatuses, podStatus)
				continue
			}

			if readinessGatesCondition != nil {
				v := strings.Split(readinessGatesCondition.Reason, "/")
				podStatus.SettingState = v[0]
				podStatus.PresentState = v[1]
				podStatus.PresentMessage = ""
				tendpointPodStatuses = append(tendpointPodStatuses, podStatus)
				continue
			}

			podStatus.SettingState = "UnKnown"
			podStatus.PresentState = "UnKnown"
			podStatus.PresentMessage = ""
			tendpointPodStatuses = append(tendpointPodStatuses, podStatus)
		}
	}

	tendpointCopy := tendpoint.DeepCopy()
	tendpointCopy.Status.PodStatus = tendpointPodStatuses

	_, err = r.k8sOption.crdClientSet.CrdV1alpha1().TEndpoints(tendpoint.Namespace).UpdateStatus(context.TODO(), tendpointCopy, k8sMetaV1.UpdateOptions{})

	if err != nil {
		utilRuntime.HandleError(fmt.Errorf(ResourceUpdateError, "tendpoint", tendpoint.Namespace, tendpoint.Name, err.Error()))
		return RateLimit
	}

	return AllOk
}
