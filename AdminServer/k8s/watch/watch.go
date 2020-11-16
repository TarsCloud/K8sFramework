package watch

import (
	"base"
	"database/sql"
	"k8s.io/client-go/informers"
	"sync"
)

var tafDb *sql.DB

var informerFactory informers.SharedInformerFactory

var prepareFunc []func()

func registryK8SWatchPrepare(f func()) {
	if prepareFunc == nil {
		prepareFunc = make([]func(), 0, 5)
	}
	prepareFunc = append(prepareFunc, f)
}

type K8SWatchImp struct {
}

func (i K8SWatchImp) GetServerServant(app string, name string) base.ServerServant {
	return serviceRecord.GetServerServant(app, name)
}

func (i K8SWatchImp) GetServerK8S(serverApp, serverName string) *base.ServerK8S {
	return statefulSetRecord.getServerK8S(serverApp, serverName)
}

func (i K8SWatchImp) SetInformerFactor(factor informers.SharedInformerFactory) {
	informerFactory = factor
}

func (i K8SWatchImp) SetTafDb(db *sql.DB) {
	tafDb = db
}

var once sync.Once

func (i K8SWatchImp) StartWatch() {
	once.Do(
		func() {
			for _, f := range prepareFunc {
				f()
			}
			informerFactory.Start(nil)
		})
}

func (i K8SWatchImp) ListNode() []string {
	return nodeLabelRecord.listNode()
}

func (i K8SWatchImp) ListAbilityNode(apps []string) []base.AbilityNode {
	return nodeLabelRecord.listAbilityNode(apps)
}

func (i K8SWatchImp) ListNodeAbility(nodes []string) []base.NodeAbility {
	return nodeLabelRecord.listNodeAbility(nodes)
}

func (i K8SWatchImp) ListPublicNode() []string {
	return nodeLabelRecord.listPublicNode()
}

func (i K8SWatchImp) IsClusterHadNode(node string) bool {
	return nodeLabelRecord.hadNode(node)
}

func (i K8SWatchImp) GetDaemonSetPodByName(nodeName string) *base.DaemonPodK8S {
	 daemon, ok := daemonPodRecord.getDamonPodByName(nodeName)
	 if !ok {
	 	return nil
	 }
	 return &daemon
}

func (i K8SWatchImp) GetDaemonSetPodByIP(nodeIP string) *base.DaemonPodK8S {
	daemon, ok := daemonPodRecord.getDamonPodByIP(nodeIP)
	if !ok {
		return nil
	}
	return &daemon
}
