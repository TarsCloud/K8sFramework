package rpc

import (
	"database/sql"
	"encoding/json"
	"fmt"
	tafConf "github.com/TarsCloud/TarsGo/tars/util/conf"
	"github.com/asaskevich/govalidator"
	"github.com/elgris/sqrl"
	"runtime"
)

const _RequestKindServerConfig = "ServerConfig"
const _RequestKindServerConfigHistory = "ServerConfigHistory"

func createServerConfig(rpcRequest *Request) (Result, Error) {

	type CreateServerConfigMetadata struct {
		AppServer     string `json:"AppServer" valid:"required,matches-ServerApps"`
		ConfigName    string `json:"ConfigName" valid:"required,matches-ConfigName"`
		ConfigContent string `json:"ConfigContent" valid:"required"`
		ConfigMark    string `json:"ConfigMark" valid:"-"`
		PodSeq        int    `json:"PodSeq" valid:"matches-ConfigPodSeq"`
	}

	var err error

	var metadata CreateServerConfigMetadata
	if err = json.Unmarshal(*rpcRequest.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	conf := tafConf.New()
	if err = conf.InitFromString(metadata.ConfigContent); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata.ConfigContent Value", -1}
	}

	CreateServerConfigResourceSql1 := "insert into t_config (f_config_version, f_config_name, f_config_content, f_create_person,f_config_mark, f_app_server,f_pod_seq) select IFNULL(max(f_config_version),9999)+2,?,?,?,?,?,? from t_config_history where f_app_server=? and f_config_name=? and f_pod_seq=?"
	if _, err = tafDb.Exec(CreateServerConfigResourceSql1, metadata.ConfigName, metadata.ConfigContent, rpcRequest.RequestAccount.Name, metadata.ConfigMark, metadata.AppServer, metadata.PodSeq, metadata.AppServer, metadata.ConfigName, metadata.PodSeq); err != nil {
		return nil, Error{"已存在相同配置", -1}
	}

	return Success, Error{"", Success}
}

func selectServerConfig(rpcRequest *Request) (Result, Error) {
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"ConfigId": SqlColumn{
			ColumnName: "f_config_id",
			ColumnType: "int",
		},
		"ConfigName": SqlColumn{
			ColumnName: "f_config_name",
			ColumnType: "string",
		},
		"ConfigVersion": SqlColumn{
			ColumnName: "f_config_version",
			ColumnType: "int",
		},
		"ConfigContent": SqlColumn{
			ColumnName: "f_config_content",
			ColumnType: "string",
		},
		"CreateTime": SqlColumn{
			ColumnName: "f_create_time",
			ColumnType: "string",
		},
		"CreatePerson": SqlColumn{
			ColumnName: "f_create_person",
			ColumnType: "string",
		},
		"ConfigMark": SqlColumn{
			ColumnName: "f_config_mark",
			ColumnType: "string",
		},
		"AppServer": SqlColumn{
			ColumnName: "f_app_server",
			ColumnType: "string",
		},
		"PodSeq": SqlColumn{
			ColumnName: "f_pod_seq",
			ColumnType: "int",
		},
	}

	const from = "t_config"

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

func extractConfigMark(key string, value *json.RawMessage) (string, interface{}, Error) {
	const columnsName = "f_config_mark"

	if value == nil {
		return columnsName, nil, Error{"", Success}
	}

	var ConfigMark string
	if err := json.Unmarshal(*value, &ConfigMark); err != nil {
		return columnsName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}
	return columnsName, ConfigMark, Error{"", Success}
}

func extractConfigContent(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = "f_config_content"

	if value == nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var configContent string
	if err := json.Unmarshal(*value, &configContent); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	conf := tafConf.New()
	if err := conf.InitFromString(configContent); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	return columnName, configContent, Error{"", Success}
}

func updateServerConfig(rpcRequest *Request) (Result, Error) {
	type UpdateServerConfigMetadata struct {
		ConfigId int `json:"ConfigId"  valid:"required,matches-ConfigId"`
	}

	var err error

	var metadata UpdateServerConfigMetadata
	if err = json.Unmarshal(*rpcRequest.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	updateSqlBuilder := sqrl.Update("t_config")

	var configContentChanged = false
	for k, v := range *rpcRequest.Params.Target {
		var extractFunctionTableKey = _RequestKindServerConfig + k
		if fun, ok := extractFunctionTable[extractFunctionTableKey]; ok == false {
		} else {
			if columnName, columnValue, err := fun(k, v); err.Code() != Success {
				return nil, err
			} else {
				if columnName == "f_config_content" {
					configContentChanged = true
				}
				updateSqlBuilder.Set(columnName, columnValue)
			}
		}
	}

	if configContentChanged {

		updateConfigSql1 := "insert into t_config_history (f_app_server, f_config_name, f_config_version, f_config_content,f_config_id,f_create_person,f_create_time, f_config_mark,f_pod_seq) select f_app_server,f_config_name,f_config_version,f_config_content,f_config_id,f_create_person,f_create_time,f_config_mark ,f_pod_seq from t_config where f_config_id = ? ON DUPLICATE KEY UPDATE f_history_id=f_history_id"
		if _, err = tafDb.Exec(updateConfigSql1, metadata.ConfigId); err != nil {
			return nil, Error{"内部错误", -1}
		}

		var maxConfigVersion int
		var row = tafDb.QueryRow("select max(f_config_version)+2 from t_config_history where f_config_id=?", metadata.ConfigId)
		if err = row.Scan(&maxConfigVersion); err != nil {
			return nil, Error{"内部错误", -1}
		}
		updateSqlBuilder.Set("f_config_version", maxConfigVersion)
	}

	updateSqlBuilder.Where(sqrl.Eq{"f_config_id": metadata.ConfigId})

	if _, err = updateSqlBuilder.RunWith(tafDb).Exec(); err != nil {
		return nil, Error{"内部错误", -1}
	}

	return Success, Error{"", Success}
}

func deleteServerConfig(rpcRequest *Request) (Result, Error) {

	type DeleteServerConfigMetadata struct {
		ConfigId int `json:"ConfigId"  valid:"required,matches-ConfigId"`
	}

	var err error
	var metadata DeleteServerConfigMetadata

	if err = json.Unmarshal(*rpcRequest.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err := govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	DeleteConfigResourceSql1 := "insert into t_config_history (f_app_server, f_config_name, f_config_version, f_config_content,f_config_id,f_create_person,f_create_time, f_config_mark,f_pod_seq) select f_app_server,f_config_name,f_config_version,f_config_content,f_config_id,f_create_person,f_create_time,f_config_mark ,f_pod_seq from t_config where f_config_id = ? ON DUPLICATE KEY UPDATE f_history_id=f_history_id"
	_, _ = tafDb.Exec(DeleteConfigResourceSql1, metadata.ConfigId)

	DeleteConfigResourceSql2 := "delete from t_config where f_config_id=?"
	if _, err = tafDb.Exec(DeleteConfigResourceSql2, metadata.ConfigId); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	return Success, Error{"", Success}
}

func doPreviewConfigContent(rpcRequest *Request) (Result, Error) {
	type PreviewConfigContentMetadata struct {
		AppServer  string `json:"AppServer" valid:"required,matches-ServerApps"`
		ConfigName string `json:"ConfigName" valid:"required,matches-ConfigName"`
		PodSeq     int    `json:"PodSeq" valid:"matches-ConfigPodSeq"`
	}

	var rows *sql.Rows

	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()

	var err error
	var metadata PreviewConfigContentMetadata
	if err = json.Unmarshal(*rpcRequest.Params.Action.Value, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"内部错误", -1}
	}

	querySql := "select f_config_content from t_config where f_app_server=? and (f_pod_seq=-1 or f_pod_seq=?) order by f_pod_seq"
	if rows, err = tafDb.Query(querySql, metadata.AppServer, metadata.PodSeq); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	var result string
	for rows.Next() {
		var tmp string
		if err = rows.Scan(&tmp); err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

			return nil, Error{"内部错误", -1}
		}
		if result != "" {
			result += "\r\n\r\n"
		}
		result += tmp
	}
	return result, Error{"", Success}
}

func selectServerConfigHistory(rpcRequest *Request) (Result, Error) {
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"HistoryId": SqlColumn{
			ColumnName: "f_history_id",
			ColumnType: "int",
		},
		"ConfigId": SqlColumn{
			ColumnName: "f_config_id",
			ColumnType: "int",
		},
		"ConfigName": SqlColumn{
			ColumnName: "f_config_name",
			ColumnType: "string",
		},
		"ConfigVersion": SqlColumn{
			ColumnName: "f_config_version",
			ColumnType: "int",
		},
		"ConfigContent": SqlColumn{
			ColumnName: "f_config_content",
			ColumnType: "string",
		},
		"CreateTime": SqlColumn{
			ColumnName: "f_create_time",
			ColumnType: "string",
		},
		"CreatePerson": SqlColumn{
			ColumnName: "f_create_person",
			ColumnType: "string",
		},
		"ConfigMark": SqlColumn{
			ColumnName: "f_config_mark",
			ColumnType: "string",
		},
		"ServerApps": SqlColumn{
			ColumnName: "f_app_server",
			ColumnType: "string",
		},
	}
	const from = "t_config_history"

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

func deleteServerConfigHistory(rpcRequest *Request) (Result, Error) {
	type DeleteServerConfigHistoryMetadata struct {
		HistoryId int `json:"HistoryId"  valid:"required,matches-HistoryConfigId"`
	}

	var err error
	var metadata DeleteServerConfigHistoryMetadata

	if err = json.Unmarshal(*rpcRequest.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{err.Error(), -1}
	}

	DeleteConfigHistoryResourceSql1 := "delete from t_config_history where f_history_id=?"
	if _, err = tafDb.Exec(DeleteConfigHistoryResourceSql1, metadata.HistoryId); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	return Success, Error{"", Success}
}

func doActiveHistoryConfig(rpcRequest *Request) (Result, Error) {
	type ChangeVersionMetadata struct {
		HistoryId int `json:"HistoryId"  valid:"required,matches-HistoryConfigId"`
	}

	var err error

	var metadata ChangeVersionMetadata
	if err = json.Unmarshal(*rpcRequest.Params.Action.Value, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Value", -1}
	}

	updateConfigSql1 := "insert into t_config_history (f_app_server, f_config_name, f_config_version, f_config_content,f_config_id,f_create_person,f_create_time, f_config_mark,f_pod_seq) select f_app_server,f_config_name,f_config_version,f_config_content,f_config_id,f_create_person,f_create_time,f_config_mark,f_pod_seq from t_config where f_config_id = (select f_config_id from t_config_history where f_history_id=?) ON DUPLICATE KEY UPDATE f_history_id=f_history_id"
	if _, err = tafDb.Exec(updateConfigSql1, metadata.HistoryId); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	updateConfigSql2 := "update t_config a inner join t_config_history b using (f_config_id) set a.f_config_version=b.f_config_version,a.f_config_content=b.f_config_content,a.f_create_person=b.f_create_person,a.f_config_content=b.f_config_content, a.f_create_time=b.f_create_time,a.f_config_mark=b.f_config_mark,a.f_pod_seq=b.f_pod_seq where b.f_history_id=?"
	if _, err = tafDb.Exec(updateConfigSql2, metadata.HistoryId); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	return Success, Error{"", Success}
}

func init() {
	registryExtract(_RequestKindServerConfig+"ConfigContent", extractConfigContent)
	registryExtract(_RequestKindServerConfig+"ConfigMark", extractConfigMark)
	registryAction(_RequestKindServerConfig+"PreviewConfigContent", doPreviewConfigContent)
	registryHandle(RequestMethodCreate+_RequestKindServerConfig, createServerConfig)
	registryHandle(RequestMethodDelete+_RequestKindServerConfig, deleteServerConfig)
	registryHandle(RequestMethodSelect+_RequestKindServerConfig, selectServerConfig)
	registryHandle(RequestMethodUpdate+_RequestKindServerConfig, updateServerConfig)
}

func init() {
	registryAction(_RequestKindServerConfigHistory+"ActiveHistoryConfig", doActiveHistoryConfig)
	registryHandle(RequestMethodDelete+_RequestKindServerConfigHistory, deleteServerConfigHistory)
	registryHandle(RequestMethodSelect+_RequestKindServerConfigHistory, selectServerConfigHistory)
}
