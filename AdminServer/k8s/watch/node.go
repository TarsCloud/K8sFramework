package watch

import (
	"base"
	"encoding/json"
	k8sApiV1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"strconv"
	"sync"
)

type tNodeLabelRecord struct {
	nodes     map[string]map[string]interface{} // map[node][ability]nil
	abilities map[string]map[string]interface{} // map[ability][node]nil
	mutex     sync.RWMutex
}

func (c *tNodeLabelRecord) onAddNode(node *k8sApiV1.Node) {
	nodeLabels := node.Labels
	if _, ok := nodeLabels[base.TafNodeLabel]; !ok {
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
			if base.IsAbilityLabel(label) {
				abilityMap[label] = nil
				appArray = append(appArray, label[len(base.TafAbilityNodeLabelPrefix):])
				break
			}

			isPublicLabel := base.IsPublicNodeLabel(label)
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

	go func() {
		const insertNodeInfoSql = "insert into t_node(f_node_name, f_ability, f_public, f_address, f_info, f_resource_version) values (?,?,?,?,?,?) on duplicate key update f_node_name=f_node_name"
		abilityArrayBs, _ := json.Marshal(appArray)
		addressBs, _ := json.Marshal(node.Status.Addresses)
		infoBs, _ := json.Marshal(node.Status.NodeInfo)
		resourceVersion, _ := strconv.ParseInt(node.ResourceVersion, 10, 64)
		_, _ = tafDb.Exec(insertNodeInfoSql, node.Name, abilityArrayBs, isPublicNode, addressBs, infoBs, resourceVersion)
	}()
}

func (c *tNodeLabelRecord) onUpdateNode(oldNode, newNode *k8sApiV1.Node) {

	oldNodeLabels := oldNode.Labels
	newNodeLabels := newNode.Labels

	_, oldNodeIsTafNode := oldNodeLabels[base.TafNodeLabel]
	_, newNodeIsTafNode := newNodeLabels[base.TafNodeLabel]

	if !oldNodeIsTafNode {
		if !newNodeIsTafNode {
			return
		}
		c.onAddNode(newNode)
		return
	}

	if !newNodeIsTafNode {
		c.onDeleteNode(newNode)
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	nodeName := newNode.Name

	for label := range oldNodeLabels {
		if base.IsAbilityLabel(label) || base.IsPublicNodeLabel(label) {
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
			if base.IsAbilityLabel(label) {
				abilityMap[label] = nil
				appArray = append(appArray, label[len(base.TafAbilityNodeLabelPrefix):])
				break
			}

			isPublicLabel := base.IsPublicNodeLabel(label)
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

	go func() {
		const updateNodeInfoSql = "update t_node set f_ability =?,f_public=?,f_address=?, f_info=?,f_resource_version=? where f_node_name=? and f_resource_version<?"
		abilityArrayBs, _ := json.Marshal(appArray)
		addressBs, _ := json.Marshal(newNode.Status.Addresses)
		infoBs, _ := json.Marshal(newNode.Status.NodeInfo)
		resourceVersion, _ := strconv.ParseInt(newNode.ResourceVersion, 10, 64)
		_, _ = tafDb.Exec(updateNodeInfoSql, abilityArrayBs, isPublicNode, addressBs, infoBs, resourceVersion, nodeName, resourceVersion)
	}()
}

func (c *tNodeLabelRecord) onDeleteNode(node *k8sApiV1.Node) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.nodes, node.Name)

	for _, v := range c.abilities {
		delete(v, node.Name)
	}

	go func() {
		const deleteNodeSql = "delete from t_node where f_node_name=?"
		_, _ = tafDb.Exec(deleteNodeSql, node.Name)
	}()
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

func (c *tNodeLabelRecord) listAbilityNode(apps []string) []base.AbilityNode {

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var result []base.AbilityNode

	if apps == nil || len(apps) == 0 {
		result = make([]base.AbilityNode, 0, len(c.abilities))
		for abilityLabel, nodeMap := range c.abilities {
			if !base.IsPublicNodeLabel(abilityLabel) {
				nodeArray := make([]string, 0, len(nodeMap))
				for node := range nodeMap {
					nodeArray = append(nodeArray, node)
				}
				result = append(result, base.AbilityNode{
					ServerApp: abilityLabel[len(base.TafAbilityNodeLabelPrefix):],
					NodeName:  nodeArray,
				})
			}
		}
	} else {
		result = make([]base.AbilityNode, 0, len(apps))
		for _, app := range apps {
			abilityLabel := base.TafAbilityNodeLabelPrefix + app
			nodeMap, ok := c.abilities[abilityLabel]
			if !ok {
				result = append(result, base.AbilityNode{
					ServerApp: app,
					NodeName:  []string{},
				})
				continue
			}

			nodeArray := make([]string, 0, len(nodeMap))
			for node := range nodeMap {
				nodeArray = append(nodeArray, node)
			}
			result = append(result, base.AbilityNode{
				ServerApp: app,
				NodeName:  nodeArray,
			})
		}
	}
	return result
}

func (c *tNodeLabelRecord) listNodeAbility(nodes []string) []base.NodeAbility {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var result []base.NodeAbility

	if nodes == nil || len(nodes) == 0 {
		result = make([]base.NodeAbility, 0, len(c.nodes))
		for node, addNodeAbilities := range c.nodes {
			appArray := make([]string, 0, len(addNodeAbilities))
			for ability := range addNodeAbilities {
				if base.IsAbilityLabel(ability) {
					appArray = append(appArray, ability[len(base.TafAbilityNodeLabelPrefix):])
				}
			}
			result = append(result, base.NodeAbility{
				NodeName:  node,
				ServerApp: appArray,
			})
		}
		return result
	} else {
		result = make([]base.NodeAbility, 0, len(nodes))
		for _, node := range nodes {
			abilityMap, abilityMapOk := c.nodes[node]
			if !abilityMapOk {
				result = append(result, base.NodeAbility{
					NodeName:  node,
					ServerApp: []string{},
				})
				continue
			}

			appArray := make([]string, 0, len(abilityMap))
			for ability := range abilityMap {
				if base.IsAbilityLabel(ability) {
					appArray = append(appArray, ability[len(base.TafAbilityNodeLabelPrefix):])
				}
			}
			result = append(result, base.NodeAbility{
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
		if _, ok := abilityMap[base.TafPublicNodeLabel]; ok {
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

var nodeLabelRecord tNodeLabelRecord

func prepareNodeWatch() {
	svcInformer := informerFactory.Core().V1().Nodes().Informer()
	svcInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			node := obj.(*k8sApiV1.Node)
			nodeLabelRecord.onAddNode(node)
		},
		DeleteFunc: func(obj interface{}) {
			node := obj.(*k8sApiV1.Node)
			nodeLabelRecord.onDeleteNode(node)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldNode := oldObj.(*k8sApiV1.Node)
			newNode := newObj.(*k8sApiV1.Node)
			nodeLabelRecord.onUpdateNode(oldNode, newNode)
		},
	})
}

func init() {
	nodeLabelRecord.nodes = make(map[string]map[string]interface{}, 30)
	nodeLabelRecord.abilities = make(map[string]map[string]interface{}, 30)
	nodeLabelRecord.mutex = sync.RWMutex{}
	registryK8SWatchPrepare(prepareNodeWatch)
}
