package compatible

import (
	k8sCoreV1 "k8s.io/api/core/v1"
	k8sInformers "k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"strings"
	"sync"
	"tarsadmin/handler/k8s"
	"time"
)

func IsAbilityLabel(label string) bool {
	return strings.HasPrefix(label, k8s.TarsAbilityNodeLabelPrefix)
}

func IsPublicNodeLabel(label string) bool {
	return label == k8s.TarsPublicNodeLabel
}

type NodeAbility struct {
	NodeName  string   `json:"NodeName"`
	ServerApp []string `json:"ServerApp"`
}

type AbilityNode struct {
	ServerApp string   `json:"ServerApp"`
	NodeName  []string `json:"NodeName"`
}

func (c *tNodeLabelRecord) onAddNode(node *k8sCoreV1.Node) {
	// NodeLabel新增了namespace
	TarsNodeLabel := k8s.TarsNodeLabelPrefix + k8s.K8sOption.Namespace

	nodeLabels := node.Labels
	if _, ok := nodeLabels[TarsNodeLabel]; !ok {
		return
	}

	abilityMap, abilityMapOk := c.nodes[node.Name]
	if !abilityMapOk {
		abilityMap = make(map[string]interface{})
		c.nodes[node.Name] = abilityMap
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	appArray := make([]string, 0, len(nodeLabels))
	isPublicNode := false

	for label := range nodeLabels {
		expectedLabel := true
		for {
			if IsAbilityLabel(label) {
				abilityMap[label] = nil
				appArray = append(appArray, label[len(k8s.TarsAbilityNodeLabelPrefix):])
				break
			}

			isPublicLabel := IsPublicNodeLabel(label)
			if isPublicLabel {
				abilityMap[label] = nil
				isPublicNode = true
				break
			}

			expectedLabel = false
			break
		}

		if expectedLabel {
			nodeMap, nodeMapOk := c.abilities[label]
			if !nodeMapOk {
				nodeMap = make(map[string]interface{})
				c.abilities[label] = nodeMap
			}
			nodeMap[node.Name] = nil
		}
	}

	c.details[node.Name] = NodeRecordDetail{NodeName: node.Name, Ability: appArray, Address: node.Status.Addresses,
		Info: node.Status.NodeInfo, Public: isPublicNode, Version: node.ResourceVersion}
}

func (c *tNodeLabelRecord) onUpdateNode(oldNode, newNode *k8sCoreV1.Node) {
	// NodeLabel新增了namespace
	TarsNodeLabel := k8s.TarsNodeLabelPrefix + k8s.K8sOption.Namespace

	oldNodeLabels := oldNode.Labels
	newNodeLabels := newNode.Labels

	_, oldNodeIsTarsNode := oldNodeLabels[TarsNodeLabel]
	_, newNodeIsTarsNode := newNodeLabels[TarsNodeLabel]

	if !oldNodeIsTarsNode {
		if !newNodeIsTarsNode {
			return
		}
		c.onAddNode(newNode)
		return
	}

	if !newNodeIsTarsNode {
		c.onDeleteNode(newNode)
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	nodeName := newNode.Name

	for label := range oldNodeLabels {
		if IsAbilityLabel(label) || IsPublicNodeLabel(label) {
			delete(c.nodes[nodeName], label)
			delete(c.abilities[label], nodeName)
			if len(c.abilities[label]) == 0 {
				delete(c.abilities, nodeName)
			}
		}
	}

	abilityMap, abilityMapOk := c.nodes[nodeName]
	if !abilityMapOk {
		abilityMap = make(map[string]interface{})
		c.nodes[nodeName] = abilityMap
	}

	appArray := make([]string, 0, len(newNodeLabels))
	isPublicNode := false

	for label := range newNodeLabels {
		expectedLabel := true
		for {
			if IsAbilityLabel(label) {
				abilityMap[label] = nil
				appArray = append(appArray, label[len(k8s.TarsAbilityNodeLabelPrefix):])
				break
			}

			isPublicLabel := IsPublicNodeLabel(label)
			if isPublicLabel {
				abilityMap[label] = nil
				isPublicNode = true
				break
			}

			expectedLabel = false
			break
		}

		if expectedLabel {
			nodeMap, nodeMapOk := c.abilities[label]
			if !nodeMapOk {
				nodeMap = make(map[string]interface{})
				c.abilities[label] = nodeMap
			}
			nodeMap[nodeName] = nil
		}
	}

	c.details[nodeName] = NodeRecordDetail{NodeName: nodeName, Ability: appArray, Address: newNode.Status.Addresses,
		Info: newNode.Status.NodeInfo, Public: isPublicNode, Version: newNode.ResourceVersion}
}

func (c *tNodeLabelRecord) onDeleteNode(node *k8sCoreV1.Node) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.nodes, node.Name)

	for _, v := range c.abilities {
		delete(v, node.Name)
	}

	delete(c.details, node.Name)
}

func (c *tNodeLabelRecord) listNode() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	result := make([]string, 0, len(c.nodes))
	for node := range c.nodes {
		result = append(result, node)
	}
	return result
}

func (c *tNodeLabelRecord) listAbilityNode(apps []string) []AbilityNode {

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var result []AbilityNode

	if apps == nil || len(apps) == 0 {
		result = make([]AbilityNode, 0, len(c.abilities))
		for abilityLabel, nodeMap := range c.abilities {
			if !IsPublicNodeLabel(abilityLabel) {
				nodeArray := make([]string, 0, len(nodeMap))
				for node := range nodeMap {
					nodeArray = append(nodeArray, node)
				}
				result = append(result, AbilityNode{
					ServerApp: abilityLabel[len(k8s.TarsAbilityNodeLabelPrefix):],
					NodeName:  nodeArray,
				})
			}
		}
	} else {
		result = make([]AbilityNode, 0, len(apps))
		for _, app := range apps {
			abilityLabel := k8s.TarsAbilityNodeLabelPrefix + app
			nodeMap, ok := c.abilities[abilityLabel]
			if !ok {
				result = append(result, AbilityNode{
					ServerApp: app,
					NodeName:  []string{},
				})
				continue
			}

			nodeArray := make([]string, 0, len(nodeMap))
			for node := range nodeMap {
				nodeArray = append(nodeArray, node)
			}
			result = append(result, AbilityNode{
				ServerApp: app,
				NodeName:  nodeArray,
			})
		}
	}
	return result
}

func (c *tNodeLabelRecord) listNodeAbility(nodes []string) []NodeAbility {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var result []NodeAbility

	if nodes == nil || len(nodes) == 0 {
		result = make([]NodeAbility, 0, len(c.nodes))
		for node, addNodeAbilities := range c.nodes {
			appArray := make([]string, 0, len(addNodeAbilities))
			for ability := range addNodeAbilities {
				if IsAbilityLabel(ability) {
					appArray = append(appArray, ability[len(k8s.TarsAbilityNodeLabelPrefix):])
				}
			}
			result = append(result, NodeAbility{
				NodeName:  node,
				ServerApp: appArray,
			})
		}
		return result
	} else {
		result = make([]NodeAbility, 0, len(nodes))
		for _, node := range nodes {
			abilityMap, abilityMapOk := c.nodes[node]
			if !abilityMapOk {
				result = append(result, NodeAbility{
					NodeName:  node,
					ServerApp: []string{},
				})
				continue
			}

			appArray := make([]string, 0, len(abilityMap))
			for ability := range abilityMap {
				if IsAbilityLabel(ability) {
					appArray = append(appArray, ability[len(k8s.TarsAbilityNodeLabelPrefix):])
				}
			}
			result = append(result, NodeAbility{
				NodeName:  node,
				ServerApp: appArray,
			})
		}
	}
	return result
}

func (c *tNodeLabelRecord) listPublicNode() []string {

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	result := make([]string, 0, 10)
	for node, abilityMap := range c.nodes {
		if _, ok := abilityMap[k8s.TarsPublicNodeLabel]; ok {
			result = append(result, node)
		}
	}
	return result
}

func (c *tNodeLabelRecord) hadNode(node string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	_, ok := c.nodes[node]
	return ok
}

func (c *tNodeLabelRecord) ListNodeDetail() []NodeRecordDetail {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	result := make([]NodeRecordDetail, 0, len(c.details))
	for _, v := range c.details {
		result = append(result, v)
	}
	return result
}

type NodeRecordDetail struct {
	NodeName 	string
	Ability  	[]string
	Public 	 	bool
	Address 	[]k8sCoreV1.NodeAddress
	Info		k8sCoreV1.NodeSystemInfo
	Version 	string
}

type tNodeLabelRecord struct {
	nodes     map[string]map[string]interface{} // map[node][ability]nil
	abilities map[string]map[string]interface{} // map[ability][node]nil
	details	  map[string]NodeRecordDetail
	mutex     sync.RWMutex
}

// 暂时无TNode的crd资源，故复用旧版Admin的逻辑，在内存中构建缓存操作
var nodeLabelRecord tNodeLabelRecord

func StartNodeWatch() bool {
	k8sInformerFactory := k8sInformers.NewSharedInformerFactoryWithOptions(k8s.K8sOption.K8SClientSet, time.Minute*30, k8sInformers.WithNamespace(k8s.K8sOption.Namespace))

	nodeInformer := k8sInformerFactory.Core().V1().Nodes().Informer()
	nodeInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			node := obj.(*k8sCoreV1.Node)
			nodeLabelRecord.onAddNode(node)
		},
		DeleteFunc: func(obj interface{}) {
			node := obj.(*k8sCoreV1.Node)
			nodeLabelRecord.onDeleteNode(node)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldNode := oldObj.(*k8sCoreV1.Node)
			newNode := newObj.(*k8sCoreV1.Node)
			if oldNode.ResourceVersion != newNode.ResourceVersion {
				nodeLabelRecord.onUpdateNode(oldNode, newNode)
			}
		},
	})

	stopCh := make(chan struct{})
	k8sInformerFactory.Start(stopCh)

	synced := []cache.InformerSynced{
		nodeInformer.HasSynced,
	}
	return cache.WaitForCacheSync(stopCh, synced ...)
}

func init() {
	nodeLabelRecord.nodes = make(map[string]map[string]interface{}, 10)
	nodeLabelRecord.abilities = make(map[string]map[string]interface{}, 10)
	nodeLabelRecord.details = make(map[string]NodeRecordDetail, 10)
	nodeLabelRecord.mutex = sync.RWMutex{}
}
