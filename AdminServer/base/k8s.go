package base

import (
	"database/sql"
	k8sInformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"strings"
)

/*定义如下几个概念:

1. TafK8S 节点池
> TafK8S 在K8S集群范围内能使用的节点就是 TasK8s 节点池.  在K8S Node上增加 $TafNodeLabel 来作为TafK8S节点标识

2.  Ability 节点池
> TafK8S 提供服务Pod与节点的绑定功能, 即服务的Pod只能调度到固定的某个或某群节点上.
在此场景下, 固定的某个或某群节点就是 $ServerApp 的节点池.
> 一个节点可以既属于 App A 的专用节点池, 又属于 App B 的专用节点池
> 在节点上增加 $TafNodeAbilityPrefix+${ServerApp} 来标识该节点属于哪个 Ability 节点池

3. 公用节点池
> 不使用 服务Pod与节点绑定功能时,可以使用公共节点池.
> 在节点上增加 $PublicTafNodeFlag 标签来标识此节点是公用节点池的一部分
> 一个节点可以既属于 Ability 节点池, 又属于 公用用节点池
*/

const TafNodeLabel = "taf.io/node"                  // 此标签表示 该节点可以被 taf 使用
const TafAbilityNodeLabelPrefix = "taf.io/ability." // 此标签表示 该节点可以被 taf 当做 App节点池使用
const TafPublicNodeLabel = "taf.io/public"          // 此标签表示 该节点可以被 taf 当做 公用节点池使用

const TafServerAppLabel = "taf.io/ServerApp"
const TafServerNameLabel = "taf.io/ServerName"
const TafServerVersionLabel = "taf.io/ServerVersion"
const TafServantLabel = "taf.io/Servant"

const TafNodeSelectorLabel = "taf.io/NodeSelector"
const TafNotStackedLabel = "taf.io/NotStacked"

func IsAbilityLabel(label string) bool {
	return strings.HasPrefix(label, TafAbilityNodeLabelPrefix)
}

func IsPublicNodeLabel(label string) bool {
	return label == TafPublicNodeLabel
}

type NodeAbility struct {
	NodeName  string   `json:"NodeName"`
	ServerApp []string `json:"ServerApp"`
}

type AbilityNode struct {
	ServerApp string   `json:"ServerApp"`
	NodeName  []string `json:"NodeName"`
}

type K8SWatchInterface interface {
	//枚举所有节点
	ListNode() []string

	//枚举可以部署指定apps的 node
	ListAbilityNode(apps []string) []AbilityNode

	//枚举node可以部署哪些apps
	ListNodeAbility(nodes []string) []NodeAbility

	//枚举公共节点
	ListPublicNode() []string

	IsClusterHadNode(nodes string) bool

	SetTafDb(db *sql.DB)
	SetInformerFactor(factor k8sInformers.SharedInformerFactory)

	GetServerK8S(serverApp, serverName string) *ServerK8S
	GetServerServant(serverApp string, serverName string) ServerServant

	GetDaemonSetPodByName(nodeName string) *DaemonPodK8S
	GetDaemonSetPodByIP(nodeIP string) *DaemonPodK8S

	StartWatch()
}

type UpdateK8SKey string

const (
	Replicas    UpdateK8SKey = "Replicas"
	Version     UpdateK8SKey = "Version"
	Image       UpdateK8SKey = "Image"
	NodeSelect  UpdateK8SKey = "NodeSelector"
	HostIpc     UpdateK8SKey = "HostIpc"
	HostPort    UpdateK8SKey = "HostPort"
	HostNetwork UpdateK8SKey = "HostNetwork"
	NotStacked  UpdateK8SKey = "NotStacked"
)

type K8SClientInterface interface {
	SetK8SClient(clientSet *kubernetes.Clientset)

	SetWorkNamespace(namespace string)

	SetK8SWatchImp(k8sWatchImp K8SWatchInterface)

	//增加某个节点的 ability
	AddNodeAbility(node string, apps ...string) error

	//取消某个节点的 ability
	DeleteNodeAbility(node string, apps ...string) error

	//将某些 node 设置为 公用节点
	SetPublicNode(nodes ...string) error

	//将某些 node 设置为 非公用节点
	DeletePublicNode(nodes ...string) error

	CreateServer(serverApp string, serverName string, serverServant ServerServant, serverK8S *ServerK8S) error

	DeleteServer(serverApp string, serverName string)

	UpdateServerK8S(serverApp string, serverName string, params map[UpdateK8SKey]interface{}) error

	AppendServant(serverApp string, serverName string, serverServant ServerServant) error

	EraseServant(serverApp string, serverName string, adapterName string) error

	UpdateServant(serverApp string, serverName string, adapterName string, params map[UpdateServantKey]interface{}) error
}
