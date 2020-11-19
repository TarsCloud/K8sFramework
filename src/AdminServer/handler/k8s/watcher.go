package k8s

import (
	"fmt"
	k8sCoreV1 "k8s.io/api/core/v1"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilRuntime "k8s.io/apimachinery/pkg/util/runtime"
	k8sInformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	k8sSchema "k8s.io/client-go/kubernetes/scheme"
	k8sCoreV1Typed "k8s.io/client-go/kubernetes/typed/core/v1"
	k8sAppsLister "k8s.io/client-go/listers/apps/v1"
	k8sCoreLister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"
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
	Namespace    string
	CrdClientSet crdClientSet.Interface
	K8SClientSet kubernetes.Interface
	recorder     record.EventRecorder
	config       *rest.Config
}

type Watcher struct {
	k8sInformerFactory k8sInformers.SharedInformerFactory
	crdInformerFactory crdInformers.SharedInformerFactory

	serviceLister     k8sCoreLister.ServiceLister
	podLister         k8sCoreLister.PodLister
	statefulSetLister k8sAppsLister.StatefulSetLister
	daemonSetLister   k8sAppsLister.DaemonSetLister

	tServerLister   crdLister.TServerLister
	tEndpointLister crdLister.TEndpointLister
	tTemplateLister crdLister.TTemplateLister
	tReleaseLister  crdLister.TReleaseLister
	tTreeLister     crdLister.TTreeLister
	tTExitedPod     crdLister.TExitedRecordLister
	tDeployLister   crdLister.TDeployLister

	synced []cache.InformerSynced
}

func StartWatcher(namespace string, config *rest.Config) (*K8SOption, *Watcher, error) {
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	crdClient, err := crdClientSet.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	utilRuntime.Must(crdScheme.AddToScheme(k8sSchema.Scheme))

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&k8sCoreV1Typed.EventSinkImpl{Interface: k8sClient.CoreV1().Events(namespace)})

	controllerAgentName := fmt.Sprintf("tars-operator")
	recorder := eventBroadcaster.NewRecorder(k8sSchema.Scheme, k8sCoreV1.EventSource{Component: controllerAgentName})

	k8sOption := &K8SOption{Namespace: namespace, config: config,
		K8SClientSet: k8sClient, CrdClientSet: crdClient, recorder: recorder}

	stopCh := make(chan struct{})

	watcher := newWatcher(k8sOption)

	watcher.start(stopCh)
	if !watcher.waitForCacheSync(stopCh) {
		return nil, nil, fmt.Errorf("handler watcher waitForCacheSync Error")
	}

	return k8sOption, watcher, nil
}

func newWatcher(r *K8SOption) *Watcher {

	k8sInformerFactory := k8sInformers.NewSharedInformerFactoryWithOptions(r.K8SClientSet, time.Minute*30, k8sInformers.WithNamespace(r.Namespace))
	crdInformerFactory := crdInformers.NewSharedInformerFactoryWithOptions(r.CrdClientSet, time.Minute*15, crdInformers.WithNamespace(r.Namespace))

	serviceInformer := k8sInformerFactory.Core().V1().Services()
	daemonSetInformer := k8sInformerFactory.Apps().V1().DaemonSets()
	podInformer := k8sInformerFactory.Core().V1().Pods()
	statefulSetInformer := k8sInformerFactory.Apps().V1().StatefulSets()
	tServerInformer := crdInformerFactory.Crd().V1alpha1().TServers()
	tEndpointInformer := crdInformerFactory.Crd().V1alpha1().TEndpoints()
	tTemplateInformer := crdInformerFactory.Crd().V1alpha1().TTemplates()
	tReleaseInformer := crdInformerFactory.Crd().V1alpha1().TReleases()
	tTreeInformer := crdInformerFactory.Crd().V1alpha1().TTrees()
	tExitedPodInformer := crdInformerFactory.Crd().V1alpha1().TExitedRecords()
	tDeployInformer := crdInformerFactory.Crd().V1alpha1().TDeploys()

	watcher := &Watcher{
		k8sInformerFactory: k8sInformerFactory,
		crdInformerFactory: crdInformerFactory,
		serviceLister:      serviceInformer.Lister(),
		daemonSetLister:    daemonSetInformer.Lister(),
		statefulSetLister:  statefulSetInformer.Lister(),
		podLister:          podInformer.Lister(),
		tServerLister:      tServerInformer.Lister(),
		tEndpointLister:    tEndpointInformer.Lister(),
		tTemplateLister:    tTemplateInformer.Lister(),
		tReleaseLister:     tReleaseInformer.Lister(),
		tTreeLister:        tTreeInformer.Lister(),
		tTExitedPod:        tExitedPodInformer.Lister(),
		tDeployLister: 		tDeployInformer.Lister(),

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
			tExitedPodInformer.Informer().HasSynced,
			tDeployInformer.Informer().HasSynced,
		},
	}

	eventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldMeta := oldObj.(k8sMetaV1.Object)
			newMeta := newObj.(k8sMetaV1.Object)
			if newMeta.GetResourceVersion() != oldMeta.GetResourceVersion() {
			}
		},
		DeleteFunc: func(obj interface{}) {
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
	tExitedPodInformer.Informer().AddEventHandler(eventHandler)
	tDeployInformer.Informer().AddEventHandler(eventHandler)

	return watcher
}

func (w *Watcher) start(stopCh chan struct{}) {
	w.k8sInformerFactory.Start(stopCh)
	w.crdInformerFactory.Start(stopCh)
}

func (w *Watcher) waitForCacheSync(stopCh chan struct{}) bool {
	return cache.WaitForCacheSync(stopCh, w.synced...)
}
