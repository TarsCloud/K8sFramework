package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	tafConf "github.com/TarsCloud/TarsGo/tars/util/conf"
	"github.com/asaskevich/govalidator"
	"github.com/elgris/sqrl"
	"runtime"
)

const _RequestKindServerOption = "ServerOption"

func selectServerOption(rpcRequest *Request) (Result, Error) {

	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"ServerId": SqlColumn{
			ColumnName: "f_server_id",
			ColumnType: "int",
		},
		"ServerTemplate": SqlColumn{
			ColumnName: "f_server_template",
			ColumnType: "string",
		},
		"ServerProfile": SqlColumn{
			ColumnName: "f_server_profile",
			ColumnType: "string",
		},
		"StartScript": SqlColumn{
			ColumnName: "f_start_script_path",
			ColumnType: "string",
		},
		"StopScript": SqlColumn{
			ColumnName: "f_stop_script_path",
			ColumnType: "string",
		},
		"MonitorScript": SqlColumn{
			ColumnName: "f_monitor_script_path",
			ColumnType: "string",
		},
		"AsyncThread": SqlColumn{
			ColumnName: "f_async_thread",
			ColumnType: "int",
		},
		"ServerImportant": SqlColumn{
			ColumnName: "f_important_type",
			ColumnType: "int",
		},
		"RemoteLogType": SqlColumn{
			ColumnName: "f_remote_log_type",
			ColumnType: "int",
		},
		"RemoteLogReserve": SqlColumn{
			ColumnName: "f_remote_log_reserve_time",
			ColumnType: "int",
		},
		"RemoteLogCompress": SqlColumn{
			ColumnName: "f_remote_log_compress_time",
			ColumnType: "int",
		},
	}
	const from = "t_server_option"
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

func extractServerTemplate(key string, value *json.RawMessage) (string, interface{}, Error) {
	const columnName = "f_server_template"

	if value == nil {
		return columnName, nil, Error{"", -1}
	}

	var serverTemplate string
	if err := json.Unmarshal(*value, &serverTemplate); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}
	return columnName, serverTemplate, Error{"", -1}
}

func extractServerProfile(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = "f_server_profile"

	if value == nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var serverProfile string
	if err := json.Unmarshal(*value, &serverProfile); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	conf := tafConf.New()

	if err := conf.InitFromString(serverProfile); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	return columnName, serverProfile, Error{"", Success}
}

func extractStartScript(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = "f_start_script_path"

	if value == nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var script string
	if err := json.Unmarshal(*value, &script); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	return columnName, script, Error{"", -1}
}

func extractStopScript(key string, value *json.RawMessage) (string, interface{}, Error) {
	const columnName = "f_stop_script_path"

	if value == nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var script string
	if err := json.Unmarshal(*value, &script); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	return columnName, script, Error{"", -1}
}

func extractMonitorScript(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = "f_monitor_script_path"

	if value == nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var script string
	if err := json.Unmarshal(*value, &script); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	return columnName, script, Error{"", -1}
}

func extractAsyncThread(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = "f_async_thread"

	if value == nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var asyncThread int
	if err := json.Unmarshal(*value, &asyncThread); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	if asyncThread < 1 || asyncThread > 20 {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	return columnName, asyncThread, Error{"", -1}
}

func extractServerImportant(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = "f_important_type"

	if value == nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var important int
	if err := json.Unmarshal(*value, &important); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	if important < 1 || important > 5 {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	return columnName, important, Error{"", -1}
}

func updateServerOption(rpcRequest *Request) (Result, Error) {
	type UpdateServerOptionMetadata struct {
		ServerId int `json:"ServerId"  valid:"required,matches-ServerId"`
	}

	var err error

	if !rpcRequest.Params.Confirmation {
		return nil, Error{"使该操作生效需要重启该服务的所有运行实例,请确认?", -2}
	}

	var metadata UpdateServerOptionMetadata
	if err = json.Unmarshal(*rpcRequest.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err := govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	updateSqlBuilder := sqrl.Update("t_server_option")

	for k, v := range *rpcRequest.Params.Target {
		var extractFunctionTableKey = _RequestKindServerOption + k
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

	updateSqlBuilder.Where(sqrl.Eq{"f_server_id": metadata.ServerId})

	if _, err = updateSqlBuilder.RunWith(tafDb).Exec(); err != nil {
		return nil, Error{"内部错误", -1}
	}
	return Success, Error{"", Success}
}

func doPreviewTemplateContent(rpcRequest *Request) (Result, Error) {

	type PreviewTemplateContentMetadata struct {
		ServerId int `json:"ServerId"  valid:"required,matches-ServerId"`
	}

	var err error
	var metadata PreviewTemplateContentMetadata

	if err = json.Unmarshal(*rpcRequest.Params.Action.Value, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Value", -1}
	}

	const queryProfileSql = "select f_server_profile,f_server_template from t_server_option where f_server_id=?"
	row := tafDb.QueryRow(queryProfileSql, metadata.ServerId)
	profile := make([]byte, 0)
	var templateName string
	if err := row.Scan(&profile, &templateName); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	allTemplateContent := make([][]byte, 0, 10)
	allTemplateContent = append(allTemplateContent, profile)

	queryParentTemplateFun := func(templateName string) (string, string, string, error) {
		const QueryTemplateSql = "SELECT  a.f_template_content , b.f_template_name , b.f_template_content FROM t_template a JOIN t_template b ON a.f_template_parent =b.f_template_name WHERE a.f_template_name=?"
		row := tafDb.QueryRow(QueryTemplateSql, templateName)
		var selfContent string
		var parentName string
		var parentContent string
		if err := row.Scan(&selfContent, &parentName, &parentContent); err != nil {
			return "", "", "", err
		}
		return selfContent, parentName, parentContent, nil
	}

	for {
		templateContent, parentTemplateName, parentContent, err := queryParentTemplateFun(templateName)

		if err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

			return nil, Error{"内部错误", -1}
		}

		if len(allTemplateContent) == 1 {
			allTemplateContent = append(allTemplateContent, []byte(templateContent))
		}

		if parentTemplateName == templateName {
			break
		}

		allTemplateContent = append(allTemplateContent, []byte(parentContent))
		break
	}

	reverseSliceFun := func(s [][]byte) [][]byte {
		for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
			s[i], s[j] = s[j], s[i]
		}
		return s
	}
	allTemplateContent = reverseSliceFun(allTemplateContent)
	conf := tafConf.New()
	afterJoinTemplateContent := bytes.Join(allTemplateContent, nil)

	if err := conf.InitFromBytes(afterJoinTemplateContent); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	return conf.ToString(), Error{"", Success}
}

func init() {
	registryExtract(_RequestKindServerOption+"ServerTemplate", extractServerTemplate)
	registryExtract(_RequestKindServerOption+"ServerProfile", extractServerProfile)
	registryExtract(_RequestKindServerOption+"StartScript", extractStartScript)
	registryExtract(_RequestKindServerOption+"StopScript", extractStopScript)
	registryExtract(_RequestKindServerOption+"MonitorScript", extractMonitorScript)
	registryExtract(_RequestKindServerOption+"AsyncThread", extractAsyncThread)
	registryExtract(_RequestKindServerOption+"ServerImportant", extractServerImportant)
	registryAction(_RequestKindServerOption+"PreviewTemplateContent", doPreviewTemplateContent)
	registryHandle(RequestMethodSelect+_RequestKindServerOption, selectServerOption)
	registryHandle(RequestMethodUpdate+_RequestKindServerOption, updateServerOption)
}
