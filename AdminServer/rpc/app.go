package rpc

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/elgris/sqrl"
)

const RequestKindApp = "App"

func createApp(request *Request) (Result, Error) {
	type CreateAppMetadata struct {
		AppName      string  `json:"AppName"  valid:"required,matches-ServerApp"`
		AppMark      string  `json:"AppMark"  valid:"-"`
		BusinessName *string `json:"BusinessName" valid:"matches-BusinessName"`
	}

	var err error
	var metadata CreateAppMetadata
	if err = json.Unmarshal(*request.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	if metadata.BusinessName != nil && *metadata.BusinessName == "" {
		metadata.BusinessName = nil
	}

	const CreateAppResourceSql1 = "INSERT INTO t_app (f_app_name, f_app_mark, f_create_person,f_business_name) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE f_app_name=f_app_name"
	if _, err = tafDb.Exec(CreateAppResourceSql1, metadata.AppName, metadata.AppMark, request.RequestAccount.Name, metadata.BusinessName); err != nil {
		return nil, Error{"内部错误 ", -1}
	}
	return Success, Error{"", 0}
}

func selectApp(request *Request) (Result, Error) {
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"AppName": SqlColumn{
			ColumnName: "f_app_name",
			ColumnType: "string",
		},
		"AppMark": SqlColumn{
			ColumnName: "f_app_mark",
			ColumnType: "string",
		},
		"BusinessName": SqlColumn{
			ColumnName: "f_business_name",
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
	}
	const from = "t_app"
	var err error
	selectResult := SelectResult{
		Data:  nil,
		Count: nil,
	}
	if selectResult.Data, selectResult.Count, err = execSelectSql(tafDb, from, request.Params, requestColumnsSqlColumnsMap, nil); err != nil {
		return nil, Error{err.Error(), -1}
	}

	return selectResult, Error{"", 0}
}

func deleteApp(request *Request) (Result, Error) {
	type DeleteAppResourceMetadata struct {
		AppName string `json:"AppName"  valid:"required,matches-ServerApp"`
	}

	var err error

	var metadata DeleteAppResourceMetadata
	if err = json.Unmarshal(*request.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	//检测 APP下是否存在 Server
	const DeleteAppResourceSql1 = "DELETE FROM t_app where f_app_name=?"
	if _, err = tafDb.Exec(DeleteAppResourceSql1, metadata.AppName); err != nil {
		return nil, Error{"内部错误", -1}
	}

	return Success, Error{"", Success}
}

func extractAppMark(key string, value *json.RawMessage) (string, interface{}, Error) {
	const columnName = "f_app_mark"
	if value == nil {
		return columnName, nil, Error{"", Success}
	}

	var appMark string
	if err := json.Unmarshal(*value, &appMark); err != nil {
		return "", nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", BadParamsSchema}
	}
	return columnName, appMark, Error{"", Success}
}

func extractBusiness(key string, value *json.RawMessage) (string, interface{}, Error) {
	const columnName = "f_business_name"

	if value == nil {
		return columnName, nil, Error{"", Success}
	}

	var businessName string

	if err := json.Unmarshal(*value, &businessName); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", BadTargetSchema}
	}

	if businessName == "" {
		return columnName, nil, Error{"", Success}
	}

	if ok := govalidator.TagMap["matches-BusinessName"](businessName); !ok {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", BadTargetSchema}
	}

	return columnName, businessName, Error{"", Success}
}

func updateApp(request *Request) (Result, Error) {

	type UpdateAppMetadata struct {
		AppName string `json:"AppName"  valid:"required,matches-ServerApp"`
	}

	var err error
	var metadata UpdateAppMetadata

	if _, err := govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", BadMetadataSchema}
	}

	updateSqlBuilder := sqrl.Update("t_app")

	for k, v := range *request.Params.Target {
		var extractFunctionTableKey = RequestKindApp + k
		if fun, ok := extractFunctionTable[extractFunctionTableKey]; ok == false {
			return nil, Error{"Bad Schema : Unsupported Params.Target[" + k + "]", BadParamsSchema}
		} else {
			if columnName, columnValue, err := fun(k, v); err.Code() != 0 {
				return nil, err
			} else {
				updateSqlBuilder.Set(columnName, columnValue)
			}
		}
	}

	updateSqlBuilder.Where(sqrl.Eq{"f_app_name": metadata.AppName})

	if _, err = updateSqlBuilder.RunWith(tafDb).Exec(); err != nil {
		return nil, Error{"内部错误", BadParamsSchema}
	}
	return 0, Error{"", Success}
}

func init() {
	registryExtract(RequestKindApp+"AppMark", extractAppMark)
	registryExtract(RequestKindApp+"BusinessName", extractBusiness)
}

func init() {
	registryHandle(RequestMethodCreate+RequestKindApp, createApp)
	registryHandle(RequestMethodSelect+RequestKindApp, selectApp)
	registryHandle(RequestMethodUpdate+RequestKindApp, updateApp)
	registryHandle(RequestMethodDelete+RequestKindApp, deleteApp)
}
