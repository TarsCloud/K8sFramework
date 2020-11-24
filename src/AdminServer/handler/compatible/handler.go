package compatible

import (
	"encoding/json"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"golang.org/x/net/context"
	k8sCoreV1 "k8s.io/api/core/v1"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sort"
	"tarsadmin/handler/k8s"
	"tarsadmin/handler/mysql"
	"tarsadmin/openapi/models"
	"tarsadmin/openapi/restapi/operations/affinity"
	"tarsadmin/openapi/restapi/operations/node"
)

var nodeColumnsSqlColumnsMap = mysql.RequestColumnSqlColumnMap{
	"NodePublic": mysql.SqlColumn{
		ColumnName: "f_public",
		ColumnType: "bool",
	},
}

type SelectNodeHandler struct {}

func (s *SelectNodeHandler) Handle(params node.SelectNodeParams) middleware.Responder {
	nodeDetails := nodeLabelRecord.ListNodeDetail()

	result := &models.SelectResult{
		Data: models.ArrayMapInterface{},
		Count: models.Count{},
	}
	result.Count["AllCount"] = int32(len(nodeDetails))
	result.Count["FilterCount"] = int32(len(nodeDetails))

	nodeDetailWrapper := NodeDetailWrapper{Detail: nodeDetails, By: func(e1, e2 *NodeRecordDetail) bool {
		return e1.NodeName < e2.NodeName
	}}
	sort.Sort(nodeDetailWrapper)

	for _, detail := range nodeDetailWrapper.Detail {
		result.Data = append(result.Data, map[string]interface{} {
			"NodeName": 	detail.NodeName,
			"NodeAbility": 	detail.Ability,
			"NodeAddress": 	detail.Address,
			"NodInfo": 		detail.Info,
			"NodePublic": 	detail.Public,
		})
	}

	return node.NewSelectNodeOK().WithPayload(result)
}


type DoListClusterNodeHandler struct {}

func (s *DoListClusterNodeHandler) Handle(params node.DoListClusterNodeParams) middleware.Responder {
	clusterNodes := nodeLabelRecord.listNode()
	sort.Strings(clusterNodes)
	return node.NewDoListClusterNodeOK().WithPayload(&node.DoListClusterNodeOKBody{Data: clusterNodes})
}

type DoSetPublicNodeHandler struct {}

func (s *DoSetPublicNodeHandler) Handle(params node.DoSetPublicNodeParams) middleware.Responder {
	nodeNames := params.Params.Metadata.NodeName

	nodeInterface := k8s.K8sOption.K8SClientSet.CoreV1().Nodes()

	var err error
	var k8sNode *k8sCoreV1.Node

	for _, nodeName := range nodeNames {
		if !nodeLabelRecord.hadNode(nodeName) {
			continue
		}

		if k8sNode, err = nodeInterface.Get(context.TODO(), nodeName, k8sMetaV1.GetOptions{}); err != nil {
			return node.NewDoSetPublicNodeInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}

		k8sNodeCopy := k8sNode.DeepCopy()
		if _, ok := k8sNode.Labels[k8s.TafPublicNodeLabel]; !ok {
			k8sNodeCopy.Labels[k8s.TafPublicNodeLabel] = ""
		}

		if _, err := nodeInterface.Update(context.TODO(), k8sNodeCopy, k8sMetaV1.UpdateOptions{}); err != nil {
			return node.NewDoSetPublicNodeInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
	}

	return node.NewDoSetPublicNodeOK().WithPayload(&node.DoSetPublicNodeOKBody{Result: 0})
}


type DoDeletePublicNodeHandler struct {}

func (s *DoDeletePublicNodeHandler) Handle(params node.DoDeletePublicNodeParams) middleware.Responder {
	nodeNames := params.Params.Metadata.NodeName

	nodeInterface := k8s.K8sOption.K8SClientSet.CoreV1().Nodes()

	var err error
	var k8sNode *k8sCoreV1.Node

	for _, nodeName := range nodeNames {
		if !nodeLabelRecord.hadNode(nodeName) {
			continue
		}

		if k8sNode, err = nodeInterface.Get(context.TODO(), nodeName, k8sMetaV1.GetOptions{}); err != nil {
			return node.NewDoDeletePublicNodeInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}

		k8sNodeCopy := k8sNode.DeepCopy()
		if _, ok := k8sNode.Labels[k8s.TafPublicNodeLabel]; ok {
			delete(k8sNodeCopy.Labels, k8s.TafPublicNodeLabel)
		}

		if _, err := nodeInterface.Update(context.TODO(), k8sNodeCopy, k8sMetaV1.UpdateOptions{}); err != nil {
			return node.NewDoDeletePublicNodeInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
	}

	return node.NewDoDeletePublicNodeOK().WithPayload(&node.DoDeletePublicNodeOKBody{Result: 0})
}

type DoDeleteNodeEnableServerHandler struct {}

func (s *DoDeleteNodeEnableServerHandler) Handle(params affinity.DoDeleteNodeEnableServerParams) middleware.Responder {

	metadata := params.Params.Metadata

	if err := deleteNodeAbility(metadata.NodeName, metadata.ServerApp ...) ; err != nil {
		return affinity.NewDoDeleteNodeEnableServerInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return affinity.NewDoDeleteNodeEnableServerOK().WithPayload(&affinity.DoDeleteNodeEnableServerOKBody{Result: 0})
}

type DoAddNodeEnableServerHandler struct {}

func (s *DoAddNodeEnableServerHandler) Handle(params affinity.DoAddNodeEnableServerParams) middleware.Responder {

	metadata := params.Params.Metadata

	if err := addNodeAbility(metadata.NodeName, metadata.ServerApp ...) ; err != nil {
		return affinity.NewDoAddNodeEnableServerInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return affinity.NewDoAddNodeEnableServerOK().WithPayload(&affinity.DoAddNodeEnableServerOKBody{Result: 0})
}


type DoDeleteServerEnableNodeHandler struct {}

func (s *DoDeleteServerEnableNodeHandler) Handle(params affinity.DoDeleteServerEnableNodeParams) middleware.Responder {
	metadata := params.Params.Metadata

	for _, nodeName := range metadata.NodeName {
		if err := deleteNodeAbility(nodeName, metadata.ServerApp); err != nil {
			return affinity.NewDoDeleteServerEnableNodeInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
	}

	return affinity.NewDoDeleteServerEnableNodeOK().WithPayload(&affinity.DoDeleteServerEnableNodeOKBody{Result: 0})
}


type DoAddServerEnableNodeHandler struct {}

func (s *DoAddServerEnableNodeHandler) Handle(params affinity.DoAddServerEnableNodeParams) middleware.Responder {
	metadata := params.Params.Metadata

	for _, nodeName := range metadata.NodeName {
		if err := addNodeAbility(nodeName, metadata.ServerApp); err != nil {
			return affinity.NewDoDeleteServerEnableNodeInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
	}

	return affinity.NewDoDeleteServerEnableNodeOK().WithPayload(&affinity.DoDeleteServerEnableNodeOKBody{Result: 0})
}
func deleteNodeAbility(nodeName string, serverApp ...string) error {
	if !nodeLabelRecord.hadNode(nodeName) {
		return fmt.Errorf("%s Is Not Tars Node. ", nodeName)
	}

	nodeInterface := k8s.K8sOption.K8SClientSet.CoreV1().Nodes()

	var err error
	var k8sNode *k8sCoreV1.Node

	if k8sNode, err = nodeInterface.Get(context.TODO(), nodeName, k8sMetaV1.GetOptions{}); err != nil {
		return err
	}

	deletedAnyLabel := false

	k8sNodeCopy := k8sNode.DeepCopy()
	for _, v := range serverApp {
		abilityLabel := k8s.TafAbilityNodeLabelPrefix + v
		if _, ok := k8sNode.Labels[abilityLabel]; ok {
			deletedAnyLabel = true
			delete(k8sNodeCopy.Labels, abilityLabel)
		}
	}

	if deletedAnyLabel {
		if _, err := nodeInterface.Update(context.TODO(), k8sNodeCopy, k8sMetaV1.UpdateOptions{}); err != nil {
			return err
		}
	}

	return nil
}


type DoListAffinityGroupByNodeHandler struct {}

func (s *DoListAffinityGroupByNodeHandler) Handle(params affinity.DoListAffinityGroupByNodeParams) middleware.Responder {
	var err error

	NodeName := make([]string, 0, 5)
	if params.NodeName != nil && *params.NodeName == "" {
		err = json.Unmarshal([]byte(*params.NodeName), &NodeName)
		if err != nil {
			return affinity.NewDoListAffinityGroupByNodeInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
	}

	nodeAbility := nodeLabelRecord.listNodeAbility(NodeName)

	result := &models.SelectResult{
		Data: models.ArrayMapInterface{},
		Count: models.Count{},
	}

	result.Count["AllCount"] = int32(len(nodeAbility))
	result.Count["FilterCount"] = int32(len(nodeAbility))

	for _, na := range nodeAbility {
		result.Data = append(result.Data, map[string]interface{} {
			"NodeName": na.NodeName,
			"ServerApp": na.ServerApp,
		})
	}

	return affinity.NewDoListAffinityGroupByNodeOK().WithPayload(result)
}

type DoListAffinityGroupByAbilityHandler struct {}

func (s *DoListAffinityGroupByAbilityHandler) Handle(params affinity.DoListAffinityGroupByAbilityParams) middleware.Responder {
	var err error

	ServerApp := make([]string, 0, 5)
	if params.ServerApp != nil && *params.ServerApp == "" {
		err = json.Unmarshal([]byte(*params.ServerApp), &ServerApp)
		if err != nil {
			return affinity.NewDoListAffinityGroupByAbilityInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
	}

	abilityNode := nodeLabelRecord.listAbilityNode(ServerApp)

	result := &models.SelectResult{
		Data: models.ArrayMapInterface{},
		Count: models.Count{},
	}

	result.Count["AllCount"] = int32(len(abilityNode))
	result.Count["FilterCount"] = int32(len(abilityNode))

	for _, an := range abilityNode {
		result.Data = append(result.Data, map[string]interface{} {
			"NodeName": an.NodeName,
			"ServerApp": an.ServerApp,
		})
	}

	return affinity.NewDoListAffinityGroupByNodeOK().WithPayload(result)
}

func addNodeAbility(nodeName string, serverApp ...string) error {
	if !nodeLabelRecord.hadNode(nodeName) {
		return fmt.Errorf("%s Is Not Tars Node. ", nodeName)
	}

	nodeInterface := k8s.K8sOption.K8SClientSet.CoreV1().Nodes()

	var err error
	var k8sNode *k8sCoreV1.Node

	if k8sNode, err = nodeInterface.Get(context.TODO(), nodeName, k8sMetaV1.GetOptions{}); err != nil {
		return err
	}

	addAnyLabel := false

	k8sNodeCopy := k8sNode.DeepCopy()
	for _, v := range serverApp {
		abilityLabel := k8s.TafAbilityNodeLabelPrefix + v
		if _, ok := k8sNode.Labels[abilityLabel]; !ok {
			addAnyLabel = true
			k8sNodeCopy.Labels[abilityLabel] = ""
		}
	}

	if addAnyLabel {
		if _, err := nodeInterface.Update(context.TODO(), k8sNodeCopy, k8sMetaV1.UpdateOptions{}); err != nil {
			return err
		}
	}

	return nil
}
