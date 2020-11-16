package rpc

import (
	"base"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/elgris/sqrl"
	"runtime"
	"strings"
)

const _RequestKindServerAdapter = "ServerAdapter"

func createServerAdapter(rpcRequest *Request) (Result, Error) {

	type CreateAdapterMetadata struct {
		ServerId int                `json:"ServerId" valid:"required,matches-ServerId"`
		Servant  base.ServerServant `json:"Servant" valid:"required,matches-ServerServant"`
	}

	var err error
	var metadata CreateAdapterMetadata

	if err = json.Unmarshal(*rpcRequest.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	newServerServant := make(base.ServerServant, len(metadata.Servant)) //客户端传入进来的 Servant.key,可能是大小写混合体,,需要统一成小写
	for _, v := range metadata.Servant {
		newServerServant[strings.ToLower(v.Name)] = v
	}
	metadata.Servant = newServerServant

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	const loadInfoSql = "select f_server_app,f_server_name from t_server where f_server_id=?"
	row := tafDb.QueryRow(loadInfoSql, metadata.ServerId)

	var serverApp string
	var serverName string
	if err := row.Scan(&serverApp, &serverName); err != nil {
		if err == sql.ErrNoRows {
			return nil, Error{"ServerId 不存在", -1}
		}
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))
		return nil, Error{"内部错误", -1}
	}

	if err = k8sClientImp.AppendServant(serverApp, serverName, metadata.Servant); err != nil {
		return nil, Error{err.Error(), -1}
	}

	const insertServantSql = "INSERT INTO t_server_adapter(f_server_id, f_name, f_port,f_threads, f_connections,f_capacity,f_timeout,f_is_taf,f_is_tcp) VALUES(?,?,?,?,?,?,?,?,?)"
	for _, v := range metadata.Servant {
		if _, err = tafDb.Exec(insertServantSql, metadata.ServerId, v.Name, v.Port, v.Threads, v.Connections, v.Capacity, v.Timeout, v.IsTaf, v.IsTcp); err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))
			return nil, Error{"内部错误", -1}
		}
	}

	return Success, Error{"", Success}
}

func loadAdapterInfo(adapterId int) (string, string, string, error) {
	const LoadAdapterInfoSql = "select f_server_app,f_server_name ,f_name from t_server_adapter a left join t_server b using (f_server_id) where f_adapter_id=?"
	row := tafDb.QueryRow(LoadAdapterInfoSql, adapterId)

	var serverApp string
	var serverName string
	var adapterName string
	if err := row.Scan(&serverApp, &serverName, &adapterName); err != nil {
		if err == sql.ErrNoRows {
			return "", "", "", errors.New("AdapterId 不存在")
		}
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))
		return "", "", "", errors.New("内部错误")
	}
	return serverApp, serverName, adapterName, nil
}

func deleteServerAdapter(rpcRequest *Request) (Result, Error) {

	type DeleteAdapterMetadata struct {
		AdapterId int `json:"AdapterId" valid:"required,matches-AdapterId"`
	}

	var err error
	var metadata DeleteAdapterMetadata

	if err = json.Unmarshal(*rpcRequest.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	var serverApp string
	var serverName string
	var adapterName string

	if serverApp, serverName, adapterName, err = loadAdapterInfo(metadata.AdapterId); err != nil {
		return nil, Error{err.Error(), -1}
	}

	if err = k8sClientImp.EraseServant(serverApp, serverName, adapterName); err != nil {
		return nil, Error{err.Error(), -1}
	}

	const deleteServantSql = "delete from t_server_adapter where f_adapter_id=?"
	_, _ = tafDb.Exec(deleteServantSql, metadata.AdapterId)

	return Success, Error{"", Success}
}

func selectServerAdapter(rpcRequest *Request) (Result, Error) {

	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"AdapterId": SqlColumn{
			ColumnName: "f_adapter_id",
			ColumnType: "int",
		},
		"ServerId": SqlColumn{
			ColumnName: "f_server_id",
			ColumnType: "int",
		},
		"Name": SqlColumn{
			ColumnName: "f_name",
			ColumnType: "string",
		},
		"Threads": SqlColumn{
			ColumnName: "f_threads",
			ColumnType: "int",
		},
		"Connections": SqlColumn{
			ColumnName: "f_connections",
			ColumnType: "int",
		},
		"Port": SqlColumn{
			ColumnName: "f_port",
			ColumnType: "int",
		},
		"Capacity": SqlColumn{
			ColumnName: "f_capacity",
			ColumnType: "int",
		},

		"Timeout": SqlColumn{
			ColumnName: "f_timeout",
			ColumnType: "int",
		},

		"IsTaf": SqlColumn{
			ColumnName: "f_is_taf",
			ColumnType: "bool",
		},
		"IsTcp": SqlColumn{
			ColumnName: "f_is_tcp",
			ColumnType: "bool",
		},
	}

	const from = "t_server_adapter"
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

func extractAdapterName(key string, value *json.RawMessage) (string, interface{}, Error) {
	const columnName = string(base.ServantName)
	if value == nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var adapterName string
	if err := json.Unmarshal(*value, &adapterName); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	validateServantFun, _ := govalidator.TagMap["matches-ServantName"]

	if match := validateServantFun(adapterName); !match {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	return columnName, adapterName, Error{"", Success}
}

func extractAdapterPort(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = string(base.ServantPort)

	if value == nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var port int
	if err := json.Unmarshal(*value, &port); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	if !govalidator.TagMap["matches-ServantPort"](govalidator.ToString(port)) {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	return columnName, port, Error{"", Success}
}

func extractAdapterThread(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = string(base.ServantThreads)

	if value == nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var threads int
	if err := json.Unmarshal(*value, &threads); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	if !govalidator.TagMap["matches-ServantThreads"](govalidator.ToString(threads)) {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	return columnName, threads, Error{"", Success}
}

func extractAdapterConnection(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = string(base.ServantConnections)

	if value == nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var connections int
	if err := json.Unmarshal(*value, &connections); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	if !govalidator.TagMap["matches-ServantConnections"](govalidator.ToString(connections)) {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	return columnName, connections, Error{"", Success}
}

func extractAdapterCapacity(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = string(base.ServantCapacity)

	if value == nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var capacity int
	if err := json.Unmarshal(*value, &capacity); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	if !govalidator.TagMap["matches-ServantCapacity"](govalidator.ToString(capacity)) {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	return columnName, capacity, Error{"", Success}
}

func extractAdapterTimeout(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = string(base.ServantTimeout)

	if value == nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var timeout int
	if err := json.Unmarshal(*value, &timeout); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	if !govalidator.TagMap["matches-ServantTimeout"](govalidator.ToString(timeout)) {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	return columnName, timeout, Error{"", Success}
}

func extractAdapterIsTaf(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = string(base.ServantIsTaf)

	if value == nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var isTaf bool
	if err := json.Unmarshal(*value, &isTaf); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	return columnName, isTaf, Error{"", Success}
}

func extractAdapterIsTcp(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = string(base.ServantIsTcp)

	if value == nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var isTcp bool
	if err := json.Unmarshal(*value, &isTcp); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	return columnName, isTcp, Error{"", Success}
}

func updateServerAdapter(rpcRequest *Request) (Result, Error) {

	if !rpcRequest.Params.Confirmation {
		return nil, Error{"使该操作生效需要重启服务的所有运行实例,请确认?", -2}
	}

	var err error
	type UpdateAdapterMetadata struct {
		AdapterId int `json:"AdapterId"  valid:"required,matches-AdapterId"`
	}

	var metadata UpdateAdapterMetadata

	if err = json.Unmarshal(*rpcRequest.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	var serverApp string
	var serverName string
	var adapterName string

	if serverApp, serverName, adapterName, err = loadAdapterInfo(metadata.AdapterId); err != nil {
		return nil, Error{err.Error(), -1}
	}

	params := make(map[base.UpdateServantKey]interface{}, len(*rpcRequest.Params.Target))
	for k, v := range *rpcRequest.Params.Target {
		var extractFunctionTableKey = _RequestKindServerAdapter + k
		if fun, ok := extractFunctionTable[extractFunctionTableKey]; ok == false {
			return nil, Error{"Bad Schema : Unsupported Params.Target[" + k + "]", BadParamsSchema}
		} else {
			if columnName, columnValue, err := fun(k, v); err.Code() != 0 {
				return nil, err
			} else {
				params[base.UpdateServantKey(columnName)] = columnValue
			}
		}
	}

	if err = k8sClientImp.UpdateServant(serverApp, serverName, adapterName, params); err != nil {
		return nil, Error{err.Error(), -1}
	}

	updateSqlBuilder := sqrl.Update("t_server_adapter")
	for k, v := range params {
		switch k {
		case base.ServantName:
			updateSqlBuilder.Set("f_name", v)
		case base.ServantPort:
			updateSqlBuilder.Set("f_port", v)
		case base.ServantIsTcp:
			updateSqlBuilder.Set("f_is_tcp", v)
		case base.ServantCapacity:
			updateSqlBuilder.Set("f_capacity", v)
		case base.ServantConnections:
			updateSqlBuilder.Set("f_connections", v)
		case base.ServantIsTaf:
			updateSqlBuilder.Set("f_is_taf", v)
		case base.ServantThreads:
			updateSqlBuilder.Set("f_threads", v)
		case base.ServantTimeout:
			updateSqlBuilder.Set("f_timeout", v)
		}
	}
	updateSqlBuilder.Where(sqrl.Eq{"f_adapter_id": metadata.AdapterId})
	if _, err = updateSqlBuilder.RunWith(tafDb).Exec(); err != nil {
		return nil, Error{"内部错误", -1}
	}
	return Success, Error{"", Success}
}

func init() {
	registryExtract(_RequestKindServerAdapter+"Name", extractAdapterName)
	registryExtract(_RequestKindServerAdapter+"Port", extractAdapterPort)
	registryExtract(_RequestKindServerAdapter+"Threads", extractAdapterThread)
	registryExtract(_RequestKindServerAdapter+"Connections", extractAdapterConnection)
	registryExtract(_RequestKindServerAdapter+"Capacity", extractAdapterCapacity)
	registryExtract(_RequestKindServerAdapter+"Timeout", extractAdapterTimeout)
	registryExtract(_RequestKindServerAdapter+"IsTaf", extractAdapterIsTaf)
	registryExtract(_RequestKindServerAdapter+"IsTcp", extractAdapterIsTcp)
}

func init() {
	registryHandle(RequestMethodSelect+_RequestKindServerAdapter, selectServerAdapter)
	registryHandle(RequestMethodCreate+_RequestKindServerAdapter, createServerAdapter)
	registryHandle(RequestMethodUpdate+_RequestKindServerAdapter, updateServerAdapter)
	registryHandle(RequestMethodDelete+_RequestKindServerAdapter, deleteServerAdapter)
}
