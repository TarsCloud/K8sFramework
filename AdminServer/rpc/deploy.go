package rpc

import (
	"base"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/elgris/sqrl"
	"runtime"
	"strings"
)

const _RequestKindDeploy = "ServerDeploy"

func createDeploy(rpcRequest *Request) (Result, Error) {
	type CreateDeployMetadata struct {
		ServerApp     string             `json:"ServerApp"  valid:"required,matches-ServerApp"`
		ServerName    string             `json:"ServerName" valid:"required,matches-ServerName"`
		ServerMark    string             `json:"ServerMark" valid:"-"`
		ServerOption  base.ServerOption  `json:"ServerOption" valid:"required,matches-ServerOption"`
		ServerK8S     base.ServerK8S     `json:"ServerK8S" valid:"required,matches-ServerK8S"`
		ServerServant base.ServerServant `json:"ServerServant" valid:"required,matches-ServerServant"`
	}

	var err error

	var metadata CreateDeployMetadata
	if err = json.Unmarshal(*rpcRequest.Params.Metadata, &metadata); err != nil {
		fmt.Println(err.Error())
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	newServerServant := make(base.ServerServant, len(metadata.ServerServant)) //客户端传入进来的 Servant.key,可能是大小写混合体,,需要统一成小写
	for _, v := range metadata.ServerServant {
		newServerServant[strings.ToLower(v.Name)] = v
	}
	metadata.ServerServant = newServerServant

	//部署时 Replicas的值只能为0 ,因为此时没有镜像服务镜像
	metadata.ServerK8S.Replicas = 0

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	if metadata.ServerK8S.HostNetwork {
		metadata.ServerK8S.HostPort = map[string]int32{}
	} else {
		newServeK8SHostPort := make(map[string]int32, len(metadata.ServerK8S.HostPort)) //客户端传入进来的 ServerK8S.HostPort.key 可能是大小写混合体,需要统一成小写
		for k, v := range metadata.ServerK8S.HostPort {
			newServeK8SHostPort[strings.ToLower(k)] = v
		}
		metadata.ServerK8S.HostPort = newServeK8SHostPort
	}

	var row *sql.Row
	const CheckAppExistSql string = "SELECT true FROM t_app where f_app_name=?"
	appExist := false
	row = tafDb.QueryRow(CheckAppExistSql, metadata.ServerApp)
	if err = row.Scan(&appExist); err != nil {
		if err == sql.ErrNoRows {
			return nil, Error{"App 不存在", -1}
		}
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	const CheckServerExistSql = "SELECT true FROM t_deploy_request WHERE f_server_app=? and f_server_name=? union SELECT true FROM t_server where f_server_app=? and f_server_name=?"
	row = tafDb.QueryRow(CheckServerExistSql, metadata.ServerApp, metadata.ServerName, metadata.ServerApp, metadata.ServerName)

	hadSameServer := false
	err = row.Scan(&hadSameServer)
	if hadSameServer == true {
		return nil, Error{"已存在同名服务", -1}
	} else if err != nil && err != sql.ErrNoRows {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	serverServerOptionByte, _ := json.Marshal(metadata.ServerOption)
	serverK8SByte, _ := json.Marshal(metadata.ServerK8S)
	serverServantByte, _ := json.Marshal(metadata.ServerServant)

	const InsertDeploySql = "INSERT INTO t_deploy_request(f_server_app, f_server_name,f_server_mark,f_server_option,f_server_k8s,f_server_servant,f_request_person) VALUES(?,?,?,?,?,?,?)"
	if _, err := tafDb.Exec(InsertDeploySql, metadata.ServerApp, metadata.ServerName, metadata.ServerMark, serverServerOptionByte, serverK8SByte, serverServantByte, rpcRequest.RequestAccount.Name); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}
	return Success, Error{"", Success}
}

func selectDeploy(rpcRequest *Request) (Result, Error) {
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"DeployId": SqlColumn{
			ColumnName: "f_request_id",
			ColumnType: "int",
		},
		"RequestTime": SqlColumn{
			ColumnName: "f_request_time",
			ColumnType: "string",
		},
		"RequestPerson": SqlColumn{
			ColumnName: "f_request_person",
			ColumnType: "string",
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

		"ServerK8S": SqlColumn{
			ColumnName: "f_server_k8s",
			ColumnType: "json",
		},

		"ServerServant": SqlColumn{
			ColumnName: "f_server_servant",
			ColumnType: "json",
		},

		"ServerOption": SqlColumn{
			ColumnName: "f_server_option",
			ColumnType: "json",
		},
	}
	const from = "t_deploy_request"
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

func extractServerServant(key string, value *json.RawMessage) (string, interface{}, Error) {
	const columnName = "f_server_servant"

	if value == nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var serverServant base.ServerServant
	if err := json.Unmarshal(*value, &serverServant); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	validateServantFun, _ := govalidator.CustomTypeTagMap.Get("matches-ServerServant")

	if match := validateServantFun(serverServant, nil); !match {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	return columnName, value, Error{"", Success}
}

func extractServerK8S(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = "f_server_k8s"

	if value == nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var serverK8S base.ServerK8S
	if err := json.Unmarshal(*value, &serverK8S); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}
	if _, err := govalidator.ValidateStruct(serverK8S); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	return "f_server_k8s", value, Error{"", Success}
}

func extractServerOption(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = "f_server_option"

	if value == nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var serverOption base.ServerOption
	if err := json.Unmarshal(*value, &serverOption); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	if _, err := govalidator.ValidateStruct(serverOption); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}
	return columnName, value, Error{"", Success}
}

func updateDeploy(rpcRequest *Request) (Result, Error) {

	type UpdateDeployMetadata struct {
		DeployId int `json:"DeployId" valid:"required,matches-DeployId"`
	}

	var err error

	var metadata UpdateDeployMetadata
	if err = json.Unmarshal(*rpcRequest.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err := govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	updateSqlBuilder := sqrl.Update("t_deploy_request")

	for k, v := range *rpcRequest.Params.Target {
		var extractFunctionTableKey = _RequestKindDeploy + k
		if fun, ok := extractFunctionTable[extractFunctionTableKey]; ok == false {
			return nil, Error{"Bad Schema : Unsupported Params.Target[" + k + "]", -1}
		} else {
			if columnName, columnValue, err := fun(k, v); err.Code() != Success {
				return nil, err
			} else {
				updateSqlBuilder.Set(columnName, columnValue)
			}
		}
	}
	updateSqlBuilder.Where(sqrl.Eq{"f_request_id": metadata.DeployId})
	if _, err = updateSqlBuilder.RunWith(tafDb).Exec(); err != nil {
		return nil, Error{"内部错误", -1}
	}

	return Success, Error{"", Success}
}

func deleteDeploy(rpcRequest *Request) (Result, Error) {

	type DeleteDeployMetadata struct {
		DeployId int `json:"DeployId"`
	}

	var err error
	var metadata DeleteDeployMetadata
	if err = json.Unmarshal(*rpcRequest.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err := govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	DeleteDeploySql := "DELETE FROM t_deploy_request WHERE f_request_id =?"
	if _, err = tafDb.Exec(DeleteDeploySql, metadata.DeployId); err != nil {
		return nil, Error{"内部错误", -1}
	}

	return Success, Error{"", Success}
}

func init() {
	registryExtract(_RequestKindDeploy+"ServerServant", extractServerServant)
	registryExtract(_RequestKindDeploy+"ServerOption", extractServerOption)
	registryExtract(_RequestKindDeploy+"ServerK8S", extractServerK8S)
}

func init() {
	registryHandle(RequestMethodCreate+_RequestKindDeploy, createDeploy)
	registryHandle(RequestMethodSelect+_RequestKindDeploy, selectDeploy)
	registryHandle(RequestMethodUpdate+_RequestKindDeploy, updateDeploy)
	registryHandle(RequestMethodDelete+_RequestKindDeploy, deleteDeploy)
}
