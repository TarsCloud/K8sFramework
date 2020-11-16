package rpc

import (
	"base"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/elgris/sqrl"
	"runtime"
)

const _RequestKindServicePool = "ServicePool"

func createServicePool(rpcRequest *Request) (Result, Error) {

	type CreateServicePoolMetadata struct {
		ServerId     int              `json:"ServerId" valid:"required,matches-ServerId"`
		ServiceImage string           `json:"ServiceImage" valid:"required,matches-ServiceImage"`
		ImageDetail  *json.RawMessage `json:"ImageDetail" valid:"-"`
		ServiceMark  string           `json:"ServiceMark" valid:"-"`
	}

	var err error
	var dbTx *sql.Tx
	allSuccess := false
	defer func() {
		if dbTx != nil && !allSuccess {
			_ = dbTx.Rollback()
		}
	}()

	var metadata CreateServicePoolMetadata
	if err = json.Unmarshal(*rpcRequest.Params.Metadata, &metadata); err != nil {
		return -1, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return -1, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	if dbTx, err = tafDb.Begin(); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	const CreateResourceSql1 = "SELECT f_server_app,f_server_name FROM t_server WHERE f_server_id=? "
	var serverApp string
	var serverName string

	var row *sql.Row
	row = dbTx.QueryRow(CreateResourceSql1, metadata.ServerId)
	if err = row.Scan(&serverApp, &serverName); err != nil {
		return nil, Error{"ServerId 不存在", -1}
	}

	const CreateResourceSql2 = "SELECT IFNULL(max(f_service_version)+2,10001) FROM t_service_pool where f_server_id=?"
	row = dbTx.QueryRow(CreateResourceSql2, metadata.ServerId)
	var newVersion int
	if err = row.Scan(&newVersion); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	const CreateResourceSql3 = "INSERT INTO  t_service_pool (f_service_version, f_service_image, f_service_mark, f_image_detail, f_server_id, f_server_app, f_server_name,f_create_person) VALUE (?,?,?,?,?,?,?,?)"
	if _, err = dbTx.Exec(CreateResourceSql3, newVersion, metadata.ServiceImage, metadata.ServiceMark, metadata.ImageDetail, metadata.ServerId, serverApp, serverName, rpcRequest.RequestAccount.Name); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	if err = dbTx.Commit(); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	return Success, Error{"", Success}
}

func selectServicePool(rpcRequest *Request) (Result, Error) {
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"ServiceId": SqlColumn{
			ColumnName: "f_service_id",
			ColumnType: "int",
		},
		"ServiceVersion": SqlColumn{
			ColumnName: "f_service_version",
			ColumnType: "int",
		},
		"ServiceMark": SqlColumn{
			ColumnName: "f_service_mark",
			ColumnType: "string",
		},
		"ServiceImage": SqlColumn{
			ColumnName: "f_service_image",
			ColumnType: "string",
		},
		"ImageDetail": SqlColumn{
			ColumnName: "f_image_detail",
			ColumnType: "json",
		},
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
		"CreateTime": SqlColumn{
			ColumnName: "f_create_time",
			ColumnType: "string",
		},
		"CreatePerson": SqlColumn{
			ColumnName: "f_create_person",
			ColumnType: "string",
		},
	}

	const from = "t_service_pool"

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

func extractServiceMark(key string, value *json.RawMessage) (string, interface{}, Error) {

	const columnName = "f_service_mark"

	if value == nil {
		return columnName, "", Error{"", Success}
	}
	var serviceMark string
	if err := json.Unmarshal(*value, &serviceMark); err != nil {
		return columnName, "", Error{"Bad Schema : Bad Params.Target[" + key + "] Format ", -1}
	}
	return columnName, value, Error{"", Success}
}

func updateServicePool(rpcRequest *Request) (Result, Error) {

	type UpdateServicePoolMetadata struct {
		ServiceId int `json:"ServiceId" valid:"required,matches-ServiceId"`
	}

	var err error

	var metadata UpdateServicePoolMetadata
	if err = json.Unmarshal(*rpcRequest.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err := govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	updateSqlBuilder := sqrl.Update("t_service_pool")

	for k, v := range *rpcRequest.Params.Target {
		var extractFunctionTableKey = _RequestKindServicePool + k
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
	updateSqlBuilder.Where(sqrl.Eq{"f_service_id": metadata.ServiceId})
	if _, err = updateSqlBuilder.RunWith(tafDb).Exec(); err != nil {
		return nil, Error{"内部错误", -1}
	}
	return Success, Error{"", Success}
}

func deleteServicePool(rpcRequest *Request) (Result, Error) {
	type DeleteServicePoolMetadata struct {
		ServiceId int `json:"ServiceId"  valid:"required,matches-ServiceId"`
	}

	var err error
	var dbTx *sql.Tx
	allSuccess := false
	defer func() {
		if dbTx != nil && !allSuccess {
			_ = dbTx.Rollback()
		}
	}()

	var metadata DeleteServicePoolMetadata
	if err = json.Unmarshal(*rpcRequest.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	if dbTx, err = tafDb.Begin(); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	const DeleteServicePoolResourceSql1 = "SELECT f_service_id from t_service_enabled where f_service_id=?"
	row := dbTx.QueryRow(DeleteServicePoolResourceSql1, metadata.ServiceId)

	var serviceId int
	if err = row.Scan(&serviceId); err == nil {
		return nil, Error{"无法删除已激活版本", -1}
	}
	if err != sql.ErrNoRows {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	const DeleteServicePoolResourceSql2 = "DELETE FROM t_service_pool where f_service_id=?"
	if _, err = dbTx.Exec(DeleteServicePoolResourceSql2, metadata.ServiceId); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	if err = dbTx.Commit(); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}
	return Success, Error{"", Success}
}

func doEnableService(rpcRequest *Request) (Result, Error) {
	type EnableService struct {
		ServiceId  int    `json:"ServiceId"  valid:"required,matches-ServiceId"`
		Replicas   int    `json:"Replicas"`
		EnableMark string `json:"EnableMark" valid:"-"`
	}

	var err error
	var metadata EnableService
	var dbTx *sql.Tx
	allSuccess := false
	defer func() {
		if dbTx != nil && !allSuccess {
			_ = dbTx.Rollback()
		}
	}()

	if err = json.Unmarshal(*rpcRequest.Params.Action.Value, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Action.Value Value", -1}
	}

	if dbTx, err = tafDb.Begin(); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	const ServicePoolResourceSql1 = "select a.f_server_id, a.f_server_app, a.f_server_name, a.f_service_version, a.f_service_image, ifnull(b.f_service_id,-1) from t_service_pool a left join t_service_enabled b using (f_service_id) where f_service_id = ?"
	row := dbTx.QueryRow(ServicePoolResourceSql1, metadata.ServiceId)

	var serverId int
	var serverApp string
	var serverName string
	var targetServiceVersion string
	var targetServiceImage string
	var enabledServiceId int

	if err = row.Scan(&serverId, &serverApp, &serverName, &targetServiceVersion, &targetServiceImage, &enabledServiceId); err != nil {
		if err == sql.ErrNoRows {
			return nil, Error{"不存在的 ServiceId", -1}
		}
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	if metadata.ServiceId == enabledServiceId {
		return Success, Error{"", Success}
	}

	k8sParams := map[base.UpdateK8SKey]interface{}{
		base.Image:    targetServiceImage,
		base.Version:  targetServiceVersion,
		base.Replicas: int32(metadata.Replicas),
	}

	if err = k8sClientImp.UpdateServerK8S(serverApp, serverName, k8sParams); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	const UpdateK8SSql = "update t_server_k8s set f_replicas=? where f_server_id=?"

	if _, err = dbTx.Exec(UpdateK8SSql, metadata.Replicas, serverId); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))
		return nil, Error{"内部错误", -1}
	}

	const ServicePoolResourceSql3 = "insert into t_service_enabled (f_server_id, f_service_id, f_enable_person, f_enable_mark) VALUE (?, ?, ?, ?) on duplicate key update f_service_id=?, f_enable_person=?,f_enable_mark=?"

	if _, err = dbTx.Exec(ServicePoolResourceSql3, serverId, metadata.ServiceId, rpcRequest.RequestAccount.Name, metadata.EnableMark, metadata.ServiceId, rpcRequest.RequestAccount.Name, metadata.EnableMark); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	if err = dbTx.Commit(); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	allSuccess = true
	return Success, Error{"", Success}
}

func init() {
	registryExtract(_RequestKindServicePool+"ServiceMark", extractServiceMark)
}

func init() {
	registryAction(_RequestKindServicePool+"EnableService", doEnableService)
	registryHandle(RequestMethodCreate+_RequestKindServicePool, createServicePool)
	registryHandle(RequestMethodSelect+_RequestKindServicePool, selectServicePool)
	registryHandle(RequestMethodUpdate+_RequestKindServicePool, updateServicePool)
	registryHandle(RequestMethodDelete+_RequestKindServicePool, deleteServicePool)
}

const _RequestKindServiceEnabled = "ServiceEnabled"

func selectServiceEnabled(rpcRequest *Request) (Result, Error) {
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"ServerId": SqlColumn{
			ColumnName: "a.f_server_id",
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
		"ServiceId": SqlColumn{
			ColumnName: "f_service_id",
			ColumnType: "int",
		},
		"ServiceVersion": SqlColumn{
			ColumnName: "f_service_version",
			ColumnType: "string",
		},
		"ServiceImage": SqlColumn{
			ColumnName: "f_service_image",
			ColumnType: "string",
		},
		"ImageDetail": SqlColumn{
			ColumnName: "f_image_detail",
			ColumnType: "json",
		},
		"EnableTime": SqlColumn{
			ColumnName: "f_enable_time",
			ColumnType: "string",
		},
		"EnablePerson": SqlColumn{
			ColumnName: "f_enable_person",
			ColumnType: "string",
		},
		"EnableMark": SqlColumn{
			ColumnName: "f_enable_mark",
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
	const from = "t_service_enabled a left join t_service_pool using (f_service_id)"
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

func init() {
	registryHandle(RequestMethodSelect+_RequestKindServiceEnabled, selectServiceEnabled)
}
