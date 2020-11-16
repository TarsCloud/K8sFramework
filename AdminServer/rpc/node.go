package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"runtime"
)

const _RequestKindNode = "Node"

func selectNode(rpcRequest *Request) (Result, Error) {
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"NodeName": SqlColumn{
			ColumnName: "f_node_name",
			ColumnType: "string",
		},
		"NodeAbility": SqlColumn{
			ColumnName: "f_ability",
			ColumnType: "json",
		},
		"NodePublic": SqlColumn{
			ColumnName: "f_public",
			ColumnType: "bool",
		},
		"NodeAddress": SqlColumn{
			ColumnName: "f_address",
			ColumnType: "json",
		},
		"NodInfo": SqlColumn{
			ColumnName: "f_info",
			ColumnType: "json",
		},
	}
	const from = "t_node"
	var err error
	selectResult := SelectResult{
		Data:  nil,
		Count: nil,
	}

	if selectResult.Data, selectResult.Count, err = execSelectSql(tafDb, from, rpcRequest.Params, requestColumnsSqlColumnsMap, nil); err != nil {
		return nil, Error{err.Error(), -1}
	}

	return selectResult, Error{"", Success}
}

func doListClusterNode(rpcRequest *Request) (Result, Error) {
	clusterNodes := k8sWatchImp.ListNode()
	return clusterNodes, Error{"", Success}
}

func doListPublicNode(rpcRequest *Request) (Result, Error) {
	publicNodes := k8sWatchImp.ListPublicNode()
	return publicNodes, Error{"", Success}
}

func doSetPublicNode(rpcRequest *Request) (Result, Error) {
	type SetNodePublicMetadata struct {
		NodeName []string `json:"NodeName"  valid:"required, each-matches-NodeName"`
	}

	var err error
	var metadata SetNodePublicMetadata

	if err = json.Unmarshal(*rpcRequest.Params.Action.Value, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Value", -1}
	}

	if err = k8sClientImp.SetPublicNode(metadata.NodeName...); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}
	return Success, Error{"", Success}
}

func doDeletePublicNode(rpcRequest *Request) (Result, Error) {
	type DeleteNodePublicMetadata struct {
		NodeName []string `json:"NodeName"  valid:"required, each-matches-NodeName"`
	}

	var err error
	var metadata DeleteNodePublicMetadata

	if err = json.Unmarshal(*rpcRequest.Params.Action.Value, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Value", -1}
	}

	if err = k8sClientImp.DeletePublicNode(metadata.NodeName...); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}
	return Success, Error{"", Success}
}

func init() {
	registryAction(_RequestKindNode+"ListClusterNode", doListClusterNode)
	registryAction(_RequestKindNode+"ListPublicNode", doListPublicNode)
	registryAction(_RequestKindNode+"SetPublicNode", doSetPublicNode)
	registryAction(_RequestKindNode+"DeletePublicNode", doDeletePublicNode)
	registryHandle(RequestMethodSelect+_RequestKindNode, selectNode)
}

const _RequestKindAffinity = "Affinity"

func doListAffinityGroupByNode(rpcRequest *Request) (Result, Error) {

	type ListAffinityGroupByNode struct {
		NodeName []string `json:"NodeName" valid:"each-matches-NodeName"`
	}

	var err error
	var metadata ListAffinityGroupByNode
	if err = json.Unmarshal(*rpcRequest.Params.Action.Value, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Value", -1}
	}

	result := k8sWatchImp.ListAbilityNode(metadata.NodeName)
	return result, Error{"", Success}
}

func doListAffinityGroupByAbility(rpcRequest *Request) (Result, Error) {

	type ListAffinityGroupByAbility struct {
		ServerApp []string `json:"ServerApp" valid:"each-matches-ServerApp"`
	}

	var err error
	var metadata ListAffinityGroupByAbility
	if err = json.Unmarshal(*rpcRequest.Params.Action.Value, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Value", -1}
	}

	result := k8sWatchImp.ListAbilityNode(metadata.ServerApp)
	return result, Error{"", Success}
}

func doDeleteNodeEnableSever(rpcRequest *Request) (Result, Error) {
	type DeleteNodeEnableServer struct {
		NodeName  string   `json:"NodeName" valid:"required,matches-NodeName"`
		ServerApp []string `json:"ServerApp" valid:"required,each-matches-ServerApps"`
	}

	var err error
	var metadata DeleteNodeEnableServer

	if err = json.Unmarshal(*rpcRequest.Params.Action.Value, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Value", -1}
	}

	if err = k8sClientImp.DeleteNodeAbility(metadata.NodeName, metadata.ServerApp...); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	return Success, Error{"", Success}
}

func doAddNodeEnableServer(rpcRequest *Request) (Result, Error) {
	type AddNodeEnableServer struct {
		NodeName  string   `json:"NodeName" valid:"required,matches-NodeName"`
		ServerApp []string `json:"ServerApp" valid:"required , each-matches-ServerApps"`
	}

	var err error
	var metadata AddNodeEnableServer

	if err = json.Unmarshal(*rpcRequest.Params.Action.Value, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Value", -1}
	}

	if err = k8sClientImp.AddNodeAbility(metadata.NodeName, metadata.ServerApp...); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}
	return Success, Error{"", Success}
}

func doDeleteEnableServerNode(rpcRequest *Request) (Result, Error) {
	type DeleteEnableServerNode struct {
		ServerApp string   `json:"ServerApp" valid:"required, matches-ServerApps"`
		NodeName  []string `json:"NodeName"  valid:"required, each-matches-NodeName"`
	}

	var err error
	var metadata DeleteEnableServerNode

	if err = json.Unmarshal(*rpcRequest.Params.Action.Value, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Value", -1}
	}

	for i := range metadata.NodeName {
		if err = k8sClientImp.DeleteNodeAbility(metadata.NodeName[i], metadata.ServerApp); err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

			return nil, Error{"内部错误", -1}
		}
	}

	return Success, Error{"", Success}
}

func doAddEnableServerNode(rpcRequest *Request) (Result, Error) {
	type AddEnableServerNode struct {
		ServerApp string   `json:"ServerApp" valid:"required, matches-ServerApps"`
		NodeName  []string `json:"NodeName"  valid:"required, each-matches-NodeName"`
	}

	var err error
	var metadata AddEnableServerNode

	if err = json.Unmarshal(*rpcRequest.Params.Action.Value, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Value", -1}
	}

	for i := range metadata.NodeName {
		if err = k8sClientImp.AddNodeAbility(metadata.NodeName[i], metadata.ServerApp); err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

			return nil, Error{"内部错误", -1}
		}
	}
	return Success, Error{"", Success}
}

func init() {
	registryAction(_RequestKindAffinity+"DeleteNodeEnableServer", doDeleteNodeEnableSever)
	registryAction(_RequestKindAffinity+"AddNodeEnableServer", doAddNodeEnableServer)
	registryAction(_RequestKindAffinity+"DeleteServerEnableNode", doDeleteEnableServerNode)
	registryAction(_RequestKindAffinity+"AddServerEnableNode", doAddEnableServerNode)
	registryAction(_RequestKindAffinity+"ListAffinityGroupByNode", doListAffinityGroupByNode)
	registryAction(_RequestKindAffinity+"ListAffinityGroupByAbility", doListAffinityGroupByAbility)
}
