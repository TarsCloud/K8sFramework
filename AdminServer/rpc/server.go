package rpc

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/elgris/sqrl"
	"runtime"
)

const _RequestKindServer = "Server"

func selectServer(rpcRequest *Request) (Result, Error) {
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
		"ServerMark": SqlColumn{
			ColumnName: "f_server_mark",
			ColumnType: "string",
		},
		"ServerType": SqlColumn{
			ColumnName: "f_server_type",
			ColumnType: "string",
		},
		"DeployPerson": SqlColumn{
			ColumnName: "f_deploy_person",
			ColumnType: "string",
		},
		"DeployTime": SqlColumn{
			ColumnName: "f_deploy_time",
			ColumnType: "string",
		},
	}
	const from = "t_server"

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

func extractServerMark(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = "f_server_mark"

	if value == nil {
		return columnName, nil, Error{"", Success}
	}

	var serverMark string
	if err := json.Unmarshal(*value, &serverMark); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	return columnName, serverMark, Error{"", Success}
}

func extractServerType(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = "f_server_type"

	if value == nil {
		return "", nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var serverType string
	if err := json.Unmarshal(*value, &serverType); err != nil {
		return "", nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	if !govalidator.TagMap["matches-ServerType"](serverType) {
		return "", nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	return columnName, serverType, Error{"", Success}
}

func updateServer(rpcRequest *Request) (Result, Error) {

	type UpdateServerMetadata struct {
		ServerId int `json:"ServerId" valid:"required,matches-ServerId"`
	}

	var err error

	var metadata UpdateServerMetadata
	if err = json.Unmarshal(*rpcRequest.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err := govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	updateSqlBuilder := sqrl.Update("t_server")

	for k, v := range *rpcRequest.Params.Target {
		var extractFunctionTableKey = _RequestKindServer + k
		if fun, ok := extractFunctionTable[extractFunctionTableKey]; ok == false {
			return nil, Error{"Bad Schema : Unsupported Params.Target[" + k + "]", BadParamsSchema}
		} else {
			if columnName, columnValue, err := fun(k, v); err.Code() != Success {
				return nil, err
			} else {
				updateSqlBuilder.Set(columnName, columnValue)
			}
		}
	}
	updateSqlBuilder.Where(sqrl.Eq{"f_server_id": metadata.ServerId})

	if _, err = updateSqlBuilder.RunWith(tafDb).Exec(); err != nil {
		return nil, Error{"内部错误", BadParamsSchema}
	}
	return Success, Error{"", Success}
}

func deleteServer(rpcRequest *Request) (Result, Error) {
	type DeleteServerMetadata struct {
		ServerId []int `json:"ServerId" valid:"required,each-matches-ServerId"`
	}

	var err error

	var metadata DeleteServerMetadata
	if err = json.Unmarshal(*rpcRequest.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err := govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	getServerInfoSql := "select f_server_app,f_server_name,f_server_kind from t_server where f_server_id=?"

	deleteSql := []string{
		"delete from t_server_k8s where f_server_id=?",
		"delete from t_server_option where f_server_id=?",
		"delete from t_server_adapter where f_server_id=?",
		"delete from t_service_enabled where f_server_id=?",
		"delete from t_service_pool where f_server_id=?",
		"delete from t_server where f_server_id=?",
	}

	result := make([]string, 0, len(metadata.ServerId))

	for _, serverId := range metadata.ServerId {

		row := tafDb.QueryRow(getServerInfoSql, serverId)

		var sServerApp string
		var sServerName string
		var sServerKind string

		if err = row.Scan(&sServerApp, &sServerName, &sServerKind); err != nil {
			if err == sql.ErrNoRows {
				result = append(result, "ServerId 不存在")
				continue
			}

			result = append(result, "内部错误")
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

			continue
		}

		if sServerKind == "System" {
			result = append(result, "不允许下线System服务")
			continue
		}

		k8sClientImp.DeleteServer(sServerApp, sServerName)

		deleteSqlRunOk := true
		for _, v := range deleteSql {
			_, err = tafDb.Exec(v, serverId)
			if err != nil {
				result = append(result, "内部错误")
				_, file, line, _ := runtime.Caller(0)
				fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

				deleteSqlRunOk = false
				break
			}
		}

		if deleteSqlRunOk {
			result = append(result, "下线"+sServerApp+"."+sServerName+"成功")
		}
	}
	return result, Error{"", Success}
}

func init() {

	registryExtract(_RequestKindServer+"ServerMark", extractServerMark)
	registryExtract(_RequestKindServer+"ServerType", extractServerType)

	registryHandle(RequestMethodSelect+_RequestKindServer, selectServer)
	registryHandle(RequestMethodUpdate+_RequestKindServer, updateServer)
	registryHandle(RequestMethodDelete+_RequestKindServer, deleteServer)
}
