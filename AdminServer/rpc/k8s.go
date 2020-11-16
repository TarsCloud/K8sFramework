package rpc

import (
	"base"
	"database/sql"
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/elgris/sqrl"
)

const _RequestKindK8S = "K8S"

func selectK8S(rpcRequest *Request) (Result, Error) {
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"ServerId": SqlColumn{
			ColumnName: "f_server_id",
			ColumnType: "int",
		},
		"ServerApp": SqlColumn{
			ColumnName: "f_server_app",
			ColumnType: "string",
		},
		"ServerName": SqlColumn{
			ColumnName: "f_server_name",
			ColumnType: "string",
		},
		"HostPort": SqlColumn{
			ColumnName: "f_host_port",
			ColumnType: "json",
		},
		"HostIpc": SqlColumn{
			ColumnName: "f_host_ipc",
			ColumnType: "bool",
		},
		"HostNetwork": SqlColumn{
			ColumnName: "f_host_network",
			ColumnType: "bool",
		},
		"Replicas": SqlColumn{
			ColumnName: "f_replicas",
			ColumnType: "int",
		},
		"NodeSelector": SqlColumn{
			ColumnName: "f_node_selector",
			ColumnType: "json",
		},
	}

	const from string = "t_server_k8s"

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

func extractK8SReplicas(key string, value *json.RawMessage) (string, interface{}, Error) {
	const columnName = string(base.Replicas)

	var replicas int32
	if err := json.Unmarshal(*value, &replicas); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	if replicas < 0 || replicas > 20 {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	return columnName, replicas, Error{"", Success}
}

func extractK8SHostNetwork(key string, value *json.RawMessage) (string, interface{}, Error) {
	const columnName = string(base.HostNetwork)

	var hostNetWork bool
	if err := json.Unmarshal(*value, &hostNetWork); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	return columnName, hostNetWork, Error{"", Success}
}

func extractK8SHostIpc(key string, value *json.RawMessage) (string, interface{}, Error) {
	const columnName = string(base.HostIpc)

	var hostIpc bool
	if err := json.Unmarshal(*value, &hostIpc); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	return columnName, hostIpc, Error{"", Success}
}

func extractK8SNotStacked(key string, value *json.RawMessage) (string, interface{}, Error) {
	const columnName = string(base.NotStacked)

	var notStacked bool
	if err := json.Unmarshal(*value, &notStacked); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	return columnName, notStacked, Error{"", Success}
}

func extractK8SHostPort(key string, value *json.RawMessage) (string, interface{}, Error) {
	const columnName = string(base.HostPort)

	var hostPort map[string]int32
	if value == nil {
		return columnName, map[string]int32{}, Error{"", Success}
	}

	if err := json.Unmarshal(*value, &hostPort); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	return columnName, hostPort, Error{"", Success}
}

func extractK8SNodeSelector(key string, value *json.RawMessage) (string, interface{}, Error) {
	const columnName = string(base.NodeSelect)

	var nodeSelect base.NodeSelector
	if err := json.Unmarshal(*value, &nodeSelect); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}
	return columnName, nodeSelect, Error{"", Success}
}

func updateK8S(rpcRequest *Request) (Result, Error) {
	type UpdateK8SMetadata struct {
		ServerId int `json:"ServerId"  valid:"required,matches-ServerId"`
	}

	var err error
	var metadata UpdateK8SMetadata

	if err = json.Unmarshal(*rpcRequest.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Value", -1}
	}

	var sServerApp string
	var sServerName string
	const getServerInfoSql = "select f_server_app,f_server_name from t_server where f_server_id=?"
	row := tafDb.QueryRow(getServerInfoSql, metadata.ServerId)

	if err = row.Scan(&sServerApp, &sServerName); err != nil {
		if err == sql.ErrNoRows {
			return nil, Error{"ServerId 不存在", -1}
		}
		return nil, Error{"内部错误", -1}
	}

	params := make(map[base.UpdateK8SKey]interface{}, len(*rpcRequest.Params.Target))
	for k, v := range *rpcRequest.Params.Target {
		var extractFunctionTableKey = _RequestKindK8S + k
		if fun, ok := extractFunctionTable[extractFunctionTableKey]; ok == false {
			return nil, Error{"Bad Schema : Unsupported Params.Target[" + k + "]", -1}
		} else {
			if columnName, columnValue, err := fun(k, v); err.Code() != Success {
				return nil, err
			} else {
				params[base.UpdateK8SKey(columnName)] = columnValue
			}
		}
	}

	if err = k8sClientImp.UpdateServerK8S(sServerApp, sServerName, params); err != nil {
		return nil, Error{err.Error(), -1}
	}

	updateSqlBuilder := sqrl.Update("t_server_k8s")
	for k, v := range params {
		switch k {
		case base.Replicas:
			updateSqlBuilder.Set("f_replicas", v)
		case base.NodeSelect:
			bs, _ := json.Marshal(v)
			updateSqlBuilder.Set("f_node_selector", bs)
		case base.HostNetwork:
			updateSqlBuilder.Set("f_host_network", v)
		case base.HostIpc:
			updateSqlBuilder.Set("f_host_ipc", v)
		case base.HostPort:
			bs, _ := json.Marshal(v)
			updateSqlBuilder.Set("f_host_port", bs)
		case base.NotStacked:
			updateSqlBuilder.Set("f_not_stacked", v)
			//Image 和 Version 不需要记录到数据库
		}
	}
	updateSqlBuilder.Where(sqrl.Eq{"f_server_id": metadata.ServerId})
	if _, err = updateSqlBuilder.RunWith(tafDb).Exec(); err != nil {
		return nil, Error{"内部错误", -1}
	}

	return Success, Error{"", Success}
}

func init() {
	registryExtract(_RequestKindK8S+"Replicas", extractK8SReplicas)
	registryExtract(_RequestKindK8S+"NodeSelector", extractK8SNodeSelector)
	registryExtract(_RequestKindK8S+"HostNetwork", extractK8SHostNetwork)
	registryExtract(_RequestKindK8S+"HostIpc", extractK8SHostIpc)
	registryExtract(_RequestKindK8S+"HostPort", extractK8SHostPort)
	registryExtract(_RequestKindK8S+"NotStacked", extractK8SNotStacked)
}

func init() {
	registryHandle(RequestMethodSelect+_RequestKindK8S, selectK8S)
	registryHandle(RequestMethodUpdate+_RequestKindK8S, updateK8S)
}
