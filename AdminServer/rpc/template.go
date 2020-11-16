package rpc

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/elgris/sqrl"
	"runtime"
)

const _RequestKindTemplate = "Template"

func createTemplate(rpcRequest *Request) (Result, Error) {
	type CreateTemplateMetadata struct {
		TemplateName    string `json:"TemplateName"    valid:"required,matches-TemplateName"`
		TemplateParent  string `json:"TemplateParent"  valid:"required,matches-TemplateName"`
		TemplateContent string `json:"TemplateContent" valid:"-"`
		CreateMark      string `json:"CreateMark"      valid:"-"`
	}

	var err error
	var dbTx *sql.Tx
	allSuccess := false
	defer func() {
		if dbTx != nil && !allSuccess {
			_ = dbTx.Rollback()
		}
	}()

	var metadata CreateTemplateMetadata

	if err = json.Unmarshal(*rpcRequest.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	if dbTx, err = tafDb.Begin(); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{ErrorMessage: "内部错误", ErrorCode: -1}
	}

	if metadata.TemplateName != metadata.TemplateParent {
		CreateTemplateSql1 := "SELECT true FROM t_template WHERE f_template_name=?"
		row := dbTx.QueryRow(CreateTemplateSql1, metadata.TemplateParent)
		var templateExist bool
		if err = row.Scan(&templateExist); err != nil {
			if err == sql.ErrNoRows {
				return nil, Error{ErrorMessage: "TemplateParent Not Exist", ErrorCode: -1}
			}
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

			return nil, Error{ErrorMessage: "内部错误", ErrorCode: -1}
		}
	}

	CreateTemplateSql2 := "INSERT INTO t_template (f_template_name,f_template_parent,f_template_content, f_create_person, f_create_mark) VALUES (?,?,?,?,?)"
	if _, err = dbTx.Exec(CreateTemplateSql2, metadata.TemplateName, metadata.TemplateParent, metadata.TemplateContent, rpcRequest.RequestAccount.Name, metadata.CreateMark); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{ErrorMessage: "内部错误", ErrorCode: -1}
	}

	if err = dbTx.Commit(); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}
	return Success, Error{ErrorCode: Success}
}

func selectTemplate(rpcRequest *Request) (Result, Error) {
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"TemplateId": SqlColumn{
			ColumnName: "f_template_id",
			ColumnType: "int",
		},
		"TemplateName": SqlColumn{
			ColumnName: "f_template_name",
			ColumnType: "string",
		},
		"TemplateParent": SqlColumn{
			ColumnName: "f_template_parent",
			ColumnType: "string",
		},

		"TemplateContent": SqlColumn{
			ColumnName: "f_template_content",
			ColumnType: "string",
		},
		"CreatePerson": SqlColumn{
			ColumnName: "f_create_person",
			ColumnType: "string",
		},
		"CreateTime": SqlColumn{
			ColumnName: "f_create_time",
			ColumnType: "string",
		},

		"CreateMark": SqlColumn{
			ColumnName: "f_create_mark",
			ColumnType: "string",
		},

		"UpdatePerson": SqlColumn{
			ColumnName: "f_update_person",
			ColumnType: "string",
		},

		"UpdateMark": SqlColumn{
			ColumnName: "f_update_mark",
			ColumnType: "string",
		},

		"UpdateTime": SqlColumn{
			ColumnName: "f_update_time",
			ColumnType: "string",
		},
	}
	const from = "t_template"
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

func extractTemplateContent(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = "f_template_content"

	if value == nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var templateContent string
	if err := json.Unmarshal(*value, &templateContent); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	return columnName, templateContent, Error{"", Success}
}

func extractTemplateParent(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = "f_template_parent"

	if value == nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var templateParent string
	if err := json.Unmarshal(*value, &templateParent); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	if !govalidator.TagMap["matches-TemplateName"](templateParent) {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	CreateTemplateSql1 := "SELECT true FROM t_template WHERE f_template_name=?"
	row := tafDb.QueryRow(CreateTemplateSql1, templateParent)
	var templateExist bool
	if err := row.Scan(&templateExist); err != nil {
		if err == sql.ErrNoRows {
			return columnName, nil, Error{"TemplateParent Not Exist", -1}
		}
		return columnName, nil, Error{"内部错误", -1}
	}
	return columnName, templateParent, Error{"", Success}
}

func extractCreateMark(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = "f_create_mark"

	if value == nil {
		return columnName, nil, Error{"", Success}
	}

	var appMark string
	if err := json.Unmarshal(*value, &appMark); err != nil {
		return columnName, nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}
	return columnName, appMark, Error{"", Success}
}

func updateTemplate(rpcRequest *Request) (Result, Error) {

	type UpdateTemplateMetadata struct {
		TemplateId int `json:"TemplateId" valid:"required,matches-TemplateId"`
	}

	var err error

	var metadata UpdateTemplateMetadata

	if err = json.Unmarshal(*rpcRequest.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err := govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	updateSqlBuilder := sqrl.Update("t_template")

	for k, v := range *rpcRequest.Params.Target {
		var extractFunctionTableKey = _RequestKindTemplate + k
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

	updateSqlBuilder.Set("f_update_person", rpcRequest.RequestAccount.Name)
	updateSqlBuilder.Where(sqrl.Eq{"f_template_id": metadata.TemplateId})
	if _, err = updateSqlBuilder.RunWith(tafDb).Exec(); err != nil {
		return nil, Error{"内部错误", -1}
	}

	return Success, Error{"", Success}
}

func deleteTemplate(rpcRequest *Request) (Result, Error) {
	type DeleteTemplateMetadata struct {
		TemplateId int `json:"TemplateId"`
	}

	var err error
	var metadata DeleteTemplateMetadata
	if err = json.Unmarshal(*rpcRequest.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err := govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	DeleteTemplateSql := "DELETE FROM t_template WHERE f_template_id =?"
	if _, err = tafDb.Exec(DeleteTemplateSql, metadata.TemplateId); err != nil {
		return nil, Error{"内部错误", -1}
	}
	return Success, Error{"", Success}
}

func init() {
	registryExtract(_RequestKindTemplate+"TemplateParent", extractTemplateParent)
	registryExtract(_RequestKindTemplate+"TemplateContent", extractTemplateContent)
	registryExtract(_RequestKindTemplate+"CreateMark", extractCreateMark)
}

func init() {
	registryHandle(RequestMethodCreate+_RequestKindTemplate, createTemplate)
	registryHandle(RequestMethodSelect+_RequestKindTemplate, selectTemplate)
	registryHandle(RequestMethodDelete+_RequestKindTemplate, deleteTemplate)
	registryHandle(RequestMethodUpdate+_RequestKindTemplate, updateTemplate)
}
