package main

import (
	"fmt"
	//"io/ioutil"
	k8sCoreV1 "k8s.io/api/core/v1"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilRuntime "k8s.io/apimachinery/pkg/util/runtime"
	k8sInformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	k8sSchema "k8s.io/client-go/kubernetes/scheme"
	k8sCoreV1Typed "k8s.io/client-go/kubernetes/typed/core/v1"
	k8sAppsLister "k8s.io/client-go/listers/apps/v1"
	k8sCoreLister "k8s.io/client-go/listers/core/v1"
	k8sClientCmd "k8s.io/client-go/tools/clientcmd"
	//"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
	crdClientSet "k8s.tars.io/crd/clientset/versioned"
	crdScheme "k8s.tars.io/crd/clientset/versioned/scheme"
	crdInformers "k8s.tars.io/crd/informers/externalversions"
	crdLister "k8s.tars.io/crd/listers/crd/v1alpha1"
	"time"
)

type K8SOption struct {
	namespace string

	k8sClientSet kubernetes.Interface

	crdClientSet crdClientSet.Interface

	recorder record.EventRecorder
}

func LoadK8SOption() (*K8SOption, error) {
	const namespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

	var namespace string

	//if bs, err := ioutil.ReadFile(namespaceFile); err != nil {
	//	return nil, err
	//} else {
	//	namespace = string(bs)
	//}
	//
	//clusterConfig, err := rest.InClusterConfig()
	//if err != nil {
	//	return nil, err
	//}

	namespace = "tars"

	clusterConfig, err := k8sClientCmd.BuildConfigFromFlags("", "/home/adugeek/.kube/config")
	if err != nil {
		return nil, fmt.Errorf("Load K8S Config Error : %s ", err.Error())
	}

	k8sClientSet := kubernetes.NewForConfigOrDie(clusterConfig)

	crdClient := crdClientSet.NewForConfigOrDie(clusterConfig)

	utilRuntime.Must(crdScheme.AddToScheme(k8sSchema.Scheme))

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&k8sCoreV1Typed.EventSinkImpl{Interface: k8sClientSet.CoreV1().Events(namespace)})

	controllerAgentName := fmt.Sprintf("tars-tarsoperator")
	recorder := eventBroadcaster.NewRecorder(k8sSchema.Scheme, k8sCoreV1.EventSource{Component: controllerAgentName})

	option := &K8SOption{
		k8sClientSet: k8sClientSet,
		crdClientSet: crdClient,
		namespace:    namespace,
		recorder:     recorder,
	}

	return option, nil
}

type Watcher struct {
	k8sInformerFactory k8sInformers.SharedInformerFactory
	crdInformerFactory crdInformers.SharedInformerFactory

	serviceLister       k8sCoreLister.ServiceLister
	podLister           k8sCoreLister.PodLister
	statefulSetLister   k8sAppsLister.StatefulSetLister
	daemonSetLister     k8sAppsLister.DaemonSetLister
	tServerLister       crdLister.TServerLister
	tEndpointLister     crdLister.TEndpointLister
	tTemplateLister     crdLister.TTemplateLister
	tReleaseLister      crdLister.TReleaseLister
	tTreeLister         crdLister.TTreeLister
	tExitedRecordLister crdLister.TExitedRecordLister
	tDeployLister       crdLister.TDeployLister

	synced     []cache.InformerSynced
	reconciles []TReconcile
}

func NewWatcher(r *K8SOption) *Watcher {

	k8sInformerFactory := k8sInformers.NewSharedInformerFactoryWithOptions(r.k8sClientSet, time.Minute*30, k8sInformers.WithNamespace(r.namespace))
	crdInformerFactory := crdInformers.NewSharedInformerFactoryWithOptions(r.crdClientSet, time.Minute*15, crdInformers.WithNamespace(r.namespace))

	serviceInformer := k8sInformerFactory.Core().V1().Services()
	daemonSetInformer := k8sInformerFactory.Apps().V1().DaemonSets()
	podInformer := k8sInformerFactory.Core().V1().Pods()
	statefulSetInformer := k8sInformerFactory.Apps().V1().StatefulSets()
	tServerInformer := crdInformerFactory.Crd().V1alpha1().TServers()
	tEndpointInformer := crdInformerFactory.Crd().V1alpha1().TEndpoints()
	tTemplateInformer := crdInformerFactory.Crd().V1alpha1().TTemplates()
	tReleaseInformer := crdInformerFactory.Crd().V1alpha1().TReleases()
	tTreeInformer := crdInformerFactory.Crd().V1alpha1().TTrees()
	tExitedRecordInformer := crdInformerFactory.Crd().V1alpha1().TExitedRecords()
	tDeployInformer := crdInformerFactory.Crd().V1alpha1().TDeploys()

	watcher := &Watcher{
		k8sInformerFactory:  k8sInformerFactory,
		crdInformerFactory:  crdInformerFactory,
		serviceLister:       serviceInformer.Lister(),
		daemonSetLister:     daemonSetInformer.Lister(),
		statefulSetLister:   statefulSetInformer.Lister(),
		podLister:           podInformer.Lister(),
		tServerLister:       tServerInformer.Lister(),
		tEndpointLister:     tEndpointInformer.Lister(),
		tTemplateLister:     tTemplateInformer.Lister(),
		tReleaseLister:      tReleaseInformer.Lister(),
		tExitedRecordLister: tExitedRecordInformer.Lister(),
		tDeployLister:       tDeployInformer.Lister(),

		synced: []cache.InformerSynced{
			serviceInformer.Informer().HasSynced,
			daemonSetInformer.Informer().HasSynced,
			statefulSetInformer.Informer().HasSynced,
			podInformer.Informer().HasSynced,
			tServerInformer.Informer().HasSynced,
			tEndpointInformer.Informer().HasSynced,
			tTemplateInformer.Informer().HasSynced,
			tReleaseInformer.Informer().HasSynced,
			tTreeInformer.Informer().HasSynced,
			tExitedRecordInformer.Informer().HasSynced,
			tDeployInformer.Informer().HasSynced,
		},
	}

	eventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			watcher.EnqueueObj(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldMeta := oldObj.(k8sMetaV1.Object)
			newMeta := newObj.(k8sMetaV1.Object)
			if newMeta.GetResourceVersion() != oldMeta.GetResourceVersion() {
				watcher.EnqueueObj(newObj)
			}
		},
		DeleteFunc: func(obj interface{}) {
			watcher.EnqueueObj(obj)
		},
	}

	serviceInformer.Informer().AddEventHandler(eventHandler)
	daemonSetInformer.Informer().AddEventHandler(eventHandler)
	statefulSetInformer.Informer().AddEventHandler(eventHandler)
	podInformer.Informer().AddEventHandler(eventHandler)
	tServerInformer.Informer().AddEventHandler(eventHandler)
	tEndpointInformer.Informer().AddEventHandler(eventHandler)
	tTemplateInformer.Informer().AddEventHandler(eventHandler)
	tReleaseInformer.Informer().AddEventHandler(eventHandler)
	tTreeInformer.Informer().AddEventHandler(eventHandler)
	tExitedRecordInformer.Informer().AddEventHandler(eventHandler)
	tDeployInformer.Informer().AddEventHandler(eventHandler)

	return watcher
}

func (w *Watcher) EnqueueObj(obj interface{}) {
	if w.reconciles != nil {
		for _, reconcile := range w.reconciles {
			reconcile.EnqueueObj(obj)
		}
	}
}

func (w *Watcher) Registry(reconcile TReconcile) {
	if reconcile == nil {
		return
	}
	if w.reconciles == nil {
		w.reconciles = make([]TReconcile, 0, 5)
	}
	w.reconciles = append(w.reconciles, reconcile)
}

func (w *Watcher) Start(stopCh chan struct{}) {
	w.k8sInformerFactory.Start(stopCh)
	w.crdInformerFactory.Start(stopCh)
}

func (w *Watcher) WaitForCacheSync(stopCh chan struct{}) bool {
	return cache.WaitForCacheSync(stopCh, w.synced...)
}
