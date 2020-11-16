package rpc

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/elgris/sqrl"
	"runtime"
)

const RequestKindBusiness = "Business"

func createBusiness(request *Request) (Result, Error) {
	type CreateBusinessMetadata struct {
		BusinessName  string `json:"BusinessName"  valid:"required,matches-BusinessName"`
		BusinessShow  string `json:"BusinessShow"  valid:"required,matches-BusinessShow"`
		BusinessMark  string `json:"BusinessMark"  valid:"-"`
		BusinessOrder int    `json:"BusinessOrder" valid:"required,matches-BusinessOrder"`
	}

	var err error

	var metadata CreateBusinessMetadata
	if err = json.Unmarshal(*request.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	const CreateBusinessResourceSql1 = "INSERT INTO t_business (f_business_name, f_business_show, f_business_mark,f_business_order, f_create_person) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE f_business_name=f_business_name"

	if _, err = tafDb.Exec(CreateBusinessResourceSql1, metadata.BusinessName, metadata.BusinessShow, metadata.BusinessMark, metadata.BusinessOrder, request.RequestAccount.Name); err != nil {
		return nil, Error{"内部错误", -1}
	}

	return Success, Error{"", Success}
}

func selectBusiness(request *Request) (Result, Error) {
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"BusinessName": SqlColumn{
			ColumnName: "f_business_name",
			ColumnType: "string",
		},
		"BusinessShow": SqlColumn{
			ColumnName: "f_business_show",
			ColumnType: "string",
		},
		"BusinessMark": SqlColumn{
			ColumnName: "f_business_mark",
			ColumnType: "string",
		},
		"BusinessOrder": SqlColumn{
			ColumnName: "f_business_order",
			ColumnType: "int",
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
	const from = "t_business"
	var err error
	selectResult := SelectResult{
		Data:  nil,
		Count: nil,
	}

	if selectResult.Data, selectResult.Count, err = execSelectSql(tafDb, from, request.Params, requestColumnsSqlColumnsMap, nil); err != nil {
		return nil, Error{err.Error(), -1}
	}

	return selectResult, Error{"", Success}
}

func extractBusinessMark(key string, value *json.RawMessage) (string, interface{}, Error) {
	const columnsName = "f_business_mark"

	if value == nil {
		return columnsName, nil, Error{"", -1}
	}

	var BusinessMark string
	if err := json.Unmarshal(*value, &BusinessMark); err != nil {
		return "", nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}
	return columnsName, BusinessMark, Error{"", Success}
}

func extractBusinessShow(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = "f_business_show"

	if value == nil {
		return "", nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var BusinessShow string
	if err := json.Unmarshal(*value, &BusinessShow); err != nil {
		return "", nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	if !govalidator.TagMap["matches-BusinessShow"](BusinessShow) {
		return "", nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	return columnName, BusinessShow, Error{"", Success}
}

func extractBusinessOrder(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = "f_business_order"

	if value == nil {
		return "", nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}

	var BusinessOrder int
	if err := json.Unmarshal(*value, &BusinessOrder); err != nil {
		return "", nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}

	if !govalidator.TagMap["matches-BusinessOrder"](govalidator.ToString(BusinessOrder)) {
		return "", nil, Error{"Bad Schema : Bad Params.Target[" + key + "] Value ", -1}
	}
	return columnName, BusinessOrder, Error{"", -1}
}

func updateBusiness(request *Request) (Result, Error) {

	type UpdateBusinessMetadata struct {
		BusinessName string `json:"BusinessName"  valid:"required,matches-BusinessName"`
	}

	var err error

	var metadata UpdateBusinessMetadata
	if err = json.Unmarshal(*request.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err := govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	updateSqlBuilder := sqrl.Update("t_business")

	for k, v := range *request.Params.Target {
		var extractFunctionTableKey = RequestKindBusiness + k
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
	updateSqlBuilder.Where(sqrl.Eq{"f_business_name": metadata.BusinessName})
	if _, err = updateSqlBuilder.RunWith(tafDb).Exec(); err != nil {
		return nil, Error{"内部错误", BadParamsSchema}
	}
	return 0, Error{"", Success}
}

func deleteBusiness(request *Request) (Result, Error) {
	type DeleteBusinessResourceMetadata struct {
		BusinessName string `json:"BusinessName"  valid:"required,matches-BusinessName"`
	}

	var err error
	var metadata DeleteBusinessResourceMetadata

	if err = json.Unmarshal(*request.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	const DeleteBusinessResourceSql1 = "DELETE FROM t_business where f_business_name=?"
	if _, err = tafDb.Exec(DeleteBusinessResourceSql1, metadata.BusinessName); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}
	return Success, Error{"", Success}
}

func doListBusinessApp(request *Request) (Result, Error) {
	type ListBusinessApp struct {
		BusinessName []string `json:"BusinessName" valid:"each-matches-BusinessName"`
	}

	var err error
	var rows *sql.Rows
	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()

	var metadata ListBusinessApp
	if err = json.Unmarshal(*request.Params.Action.Value, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Value", -1}
	}

	hadBusinessNameSqlBuilder := sqrl.Select("f_business_name", "f_business_show", "ifnull(f_app_name,'')").From("t_business a left join t_app b using (f_business_name)")
	if len(metadata.BusinessName) != 0 {
		hadBusinessNameSqlBuilder.Where(sqrl.Eq{"f_business_name": metadata.BusinessName})
	}
	hadBusinessNameSqlBuilder.OrderBy("f_business_order desc,f_business_name")

	if rows, err = hadBusinessNameSqlBuilder.RunWith(tafDb).Query(); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	type GroupElem struct {
		BusinessName string   `json:"BusinessName"`
		BusinessShow string   `json:"BusinessShow"`
		App          []string `json:"App"`
	}

	result := make([]GroupElem, 0, 15)

	var lastBusinessName string
	var lastBusinessShow string
	var posApp []string
	for rows.Next() {
		var businessName string
		var businessShow string
		var appName string

		if err = rows.Scan(&businessName, &businessShow, &appName); err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

			return nil, Error{"内部错误", -1}
		}

		if businessName != lastBusinessName {
			if posApp != nil {
				result = append(result, GroupElem{
					BusinessName: lastBusinessName,
					BusinessShow: lastBusinessShow,
					App:          posApp,
				})
			}
			posApp = make([]string, 0, 10)
			lastBusinessName = businessName
			lastBusinessShow = businessShow
		}

		if appName != "" {
			posApp = append(posApp, appName)
		}
	}

	result = append(result, GroupElem{
		BusinessName: lastBusinessName,
		BusinessShow: lastBusinessShow,
		App:          posApp,
	})

	if len(metadata.BusinessName) == 0 {
		posApp = make([]string, 0, 10)
		nonBusinessAppSqlBuilder := sqrl.Select("f_app_name").From("t_app").Where(sqrl.Eq{"f_business_name": nil})
		if rows, err = nonBusinessAppSqlBuilder.RunWith(tafDb).Query(); err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

			return nil, Error{"内部错误", -1}
		}
		for rows.Next() {
			var appName string
			if err = rows.Scan(&appName); err != nil {
				_, file, line, _ := runtime.Caller(0)
				fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

				return nil, Error{"内部错误", -1}
			}
			posApp = append(posApp, appName)
		}

		result = append(result, GroupElem{
			BusinessName: "",
			BusinessShow: "",
			App:          posApp,
		})
	}
	return result, Error{"", Success}
}

func doAddBusinessApp(request *Request) (Result, Error) {
	type AddBusinessApp struct {
		BusinessName string   `json:"BusinessName" valid:"required,matches-BusinessName"`
		AppName      []string `json:"AppName" valid:"required,each-matches-ServerApp"`
	}

	var err error

	var metadata AddBusinessApp
	if err = json.Unmarshal(*request.Params.Action.Value, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Value", -1}
	}

	sqlBuilder := sqrl.Update("t_app").Set("f_business_name", metadata.BusinessName).Where(sqrl.Eq{"f_app_name": metadata.AppName})

	if _, err = sqlBuilder.RunWith(tafDb).Exec(); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}
	return Success, Error{"", Success}
}

func doDeleteBusinessApp(request *Request) (Result, Error) {
	type DeleteBusinessApp struct {
		BusinessName string   `json:"BusinessName" valid:"required,matches-BusinessName"`
		AppName      []string `json:"AppName" valid:"required,each-matches-ServerApp"`
	}

	var err error

	var metadata DeleteBusinessApp
	if err = json.Unmarshal(*request.Params.Action.Value, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Value", -1}
	}

	sqlBuilder := sqrl.Update("t_app").Set("f_business_name", nil).Where(sqrl.Eq{"f_business_name": metadata.BusinessName}).Where(sqrl.Eq{"f_app_name": metadata.AppName})

	if _, err = sqlBuilder.RunWith(tafDb).Exec(); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}
	return Success, Error{"", Success}
}

func init() {
	registryExtract(RequestKindBusiness+"BusinessMark", extractBusinessMark)
	registryExtract(RequestKindBusiness+"BusinessShow", extractBusinessShow)
	registryExtract(RequestKindBusiness+"BusinessOrder", extractBusinessOrder)
	registryAction(RequestKindBusiness+"ListBusinessApp", doListBusinessApp)
	registryAction(RequestKindBusiness+"AddBusinessApp", doAddBusinessApp)
	registryAction(RequestKindBusiness+"DeleteBusinessApp", doDeleteBusinessApp)
	registryHandle(RequestMethodCreate+RequestKindBusiness, createBusiness)
	registryHandle(RequestMethodSelect+RequestKindBusiness, selectBusiness)
	registryHandle(RequestMethodUpdate+RequestKindBusiness, updateBusiness)
	registryHandle(RequestMethodDelete+RequestKindBusiness, deleteBusiness)
}
