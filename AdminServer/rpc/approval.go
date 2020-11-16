package rpc

import (
	"base"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"runtime"
)

const _RequestKindApproval = "ServerApproval"

func createApproval(rpcRequest *Request) (Result, Error) {

	type ServerApprovalCreateMetadata struct {
		DeployId       int    `json:"DeployId" valid:"required,matches-DeployId"`
		ApprovalResult bool   `json:"ApprovalResult" valid:"required"`
		ApprovalMark   string `json:"ApprovalMark" valid:"-"`
	}

	var err error

	var metadata ServerApprovalCreateMetadata
	if err = json.Unmarshal(*rpcRequest.Params.Metadata, &metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Format", -1}
	}

	if _, err = govalidator.ValidateStruct(metadata); err != nil {
		return nil, Error{"Bad Schema : Bad Params.Metadata Value", -1}
	}

	const LoadDeployRequestSql = " SELECT f_server_app, f_server_name,f_server_mark, f_request_person,f_server_option,f_server_servant,f_server_k8s FROM t_deploy_request WHERE f_request_id=?"

	var row *sql.Row
	row = tafDb.QueryRow(LoadDeployRequestSql, metadata.DeployId)
	var serverApp string
	var serverName string
	var serverMark string
	var requestPerson string
	var deployRequestServerOption []byte
	var deployRequestServerK8S []byte
	var deployRequestServerServants []byte
	if err = row.Scan(&serverApp, &serverName, &serverMark, &requestPerson, &deployRequestServerOption, &deployRequestServerServants, &deployRequestServerK8S); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	var dbTx *sql.Tx
	allSuccess := false
	defer func() {
		if !allSuccess {
			if dbTx != nil {
				_ = dbTx.Rollback()
			}
			k8sClientImp.DeleteServer(serverApp, serverName)
		}
	}()

	if dbTx, err = tafDb.Begin(); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	if metadata.ApprovalResult {

		var ServerOption base.ServerOption

		if err = json.Unmarshal(deployRequestServerOption, &ServerOption); err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

			return nil, Error{"内部错误", -1}
		}

		if _, err = govalidator.ValidateStruct(ServerOption); err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

			return nil, Error{"内部错误", -1}
		}

		var serverK8S base.ServerK8S
		if err = json.Unmarshal(deployRequestServerK8S, &serverK8S); err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

			return nil, Error{"内部错误", -1}
		}

		validateK8SFun, _ := govalidator.CustomTypeTagMap.Get("matches-ServerK8S")
		if match := validateK8SFun(serverK8S, nil); !match {
			fmt.Printf("内部错误: \n")
			return nil, Error{"内部错误", -1}
		}

		var serverServant base.ServerServant
		if err = json.Unmarshal(deployRequestServerServants, &serverServant); err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

			return nil, Error{"内部错误", -1}
		}

		validateServantFun, _ := govalidator.CustomTypeTagMap.Get("matches-ServerServant")

		if match := validateServantFun(serverServant, nil); !match {
			fmt.Printf("内部错误: \n")
			return nil, Error{"内部错误", -1}
		}

		const InsertServerSql = "INSERT INTO t_server (f_server_app, f_server_name, f_server_mark, f_deploy_person,f_approval_person) VALUES (?,?,?,?,?)"
		var serverId int64
		var lastInsert sql.Result
		if lastInsert, err = dbTx.Exec(InsertServerSql, serverApp, serverName, serverMark, requestPerson, rpcRequest.RequestAccount.Name); err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

			return nil, Error{"内部错误", -1}
		} else if serverId, err = lastInsert.LastInsertId(); err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

			return nil, Error{"内部错误", -1}
		}

		if err = k8sClientImp.CreateServer(serverApp, serverName, serverServant, &serverK8S); err != nil {
			const NotifySql = "INSERT INTO t_server_notify (f_app_server, f_notify_level, f_notify_message, f_notify_thread, f_notify_source) VALUE (?,?,?,?,?)"
			_, _ = tafDb.Exec(NotifySql, serverApp+"."+serverName, "NOTIFYWARN", "Failed : Create Server On K8S Error , "+err.Error(), -1, "tafAdmin")
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

			return nil, Error{"内部错误", -1}
		}

		const InsertServerServantSql = "INSERT INTO t_server_adapter(f_server_id, f_name, f_port,f_threads, f_connections,f_capacity,f_timeout,f_is_taf,f_is_tcp) VALUES(?,?,?,?,?,?,?,?,?)"
		for i := range serverServant {
			if _, err = dbTx.Exec(InsertServerServantSql, serverId, serverServant[i].Name, serverServant[i].Port, serverServant[i].Threads, serverServant[i].Connections, serverServant[i].Capacity, serverServant[i].Timeout, serverServant[i].IsTaf, serverServant[i].IsTcp); err != nil {
				_, file, line, _ := runtime.Caller(0)
				fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

				return nil, Error{"内部错误", -1}
			}
		}

		const InsertServerOptionSql = "INSERT INTO t_server_option(f_server_id, f_server_template, f_server_profile, f_start_script_path, f_stop_script_path, f_monitor_script_path, f_async_thread, f_important_type, f_remote_log_reserve_time, f_remote_log_compress_time, f_remote_log_type) VALUES (?,?,?,?,?,?,?,?,?,?,?)"
		if _, err = dbTx.Exec(InsertServerOptionSql, serverId, ServerOption.ServerTemplate, ServerOption.ServerProfile, ServerOption.StartScript, ServerOption.StopScript, ServerOption.MonitorScript, ServerOption.AsyncThread, ServerOption.ServerImportant, ServerOption.RemoteLogReserveTime, ServerOption.RemoteLogCompressTime, ServerOption.RemoteLogEnable); err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

			return nil, Error{"内部错误", -1}
		}

		const InsertSeverK8SSql = "INSERT INTO t_server_k8s(f_server_id, f_server_app, f_server_name, f_replicas, f_node_selector, f_host_network, f_host_ipc,f_host_port, f_not_stacked) VALUES( ?,?,?,?,?,?,?,?,? )"
		nodeSelectorBytes, _ := json.Marshal(serverK8S.NodeSelector)
		hostPortBytes, _ := json.Marshal(serverK8S.HostPort)
		if _, err = dbTx.Exec(InsertSeverK8SSql, serverId, serverApp, serverName, serverK8S.Replicas, nodeSelectorBytes, serverK8S.HostNetwork, serverK8S.HostIpc, hostPortBytes, serverK8S.NotStacked); err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

			return nil, Error{"内部错误", -1}
		}
	}

	const InsertApprovalSql = "INSERT INTO t_deploy_approval (f_request_time, f_request_person, f_server_app, f_server_name, f_server_mark,f_server_servant,f_server_option,f_server_k8s, f_approval_person, f_approval_result, f_approval_mark) select f_request_time, f_request_person, f_server_app, f_server_name, f_server_mark,f_server_servant,f_server_option,f_server_k8s, ?, ?, ? from t_deploy_request where f_request_id = ?"
	if _, err = dbTx.Exec(InsertApprovalSql, rpcRequest.RequestAccount.Name, metadata.ApprovalResult, metadata.ApprovalMark, metadata.DeployId); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	const DeleteDeployRequestSql = "delete from t_deploy_request  where f_request_id=?"
	if _, err = dbTx.Exec(DeleteDeployRequestSql, metadata.DeployId); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	const InsertSuccessNotify = "INSERT INTO t_server_notify (f_app_server, f_notify_level, f_notify_message, f_notify_thread, f_notify_source) VALUE (?,?,?,?,?)"
	_, _ = tafDb.Exec(InsertSuccessNotify, serverApp+"."+serverName, "NOTIFYNORMAL", "Success : Create Server On K8S Success ", -1, "tafAdmin")

	if err = dbTx.Commit(); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}
	allSuccess = true
	return Success, Error{"", Success}
}

func selectApproval(rpcRequest *Request) (Result, Error) {
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"ApprovalTime": SqlColumn{
			ColumnName: "f_approval_time",
			ColumnType: "string",
		},
		"ApprovalPerson": SqlColumn{
			ColumnName: "f_approval_person",
			ColumnType: "string",
		},
		"ApprovalResult": SqlColumn{
			ColumnName: "f_approval_result",
			ColumnType: "bool",
		},
		"ApprovalMark": SqlColumn{
			ColumnName: "f_approval_mark",
			ColumnType: "string",
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
		"ServerOption": SqlColumn{
			ColumnName: "f_server_option",
			ColumnType: "json",
		},
		"ServerServant": SqlColumn{
			ColumnName: "f_server_servant",
			ColumnType: "json",
		},
	}
	const from string = "t_deploy_approval"

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
	registryHandle(RequestMethodCreate+_RequestKindApproval, createApproval)
	registryHandle(RequestMethodSelect+_RequestKindApproval, selectApproval)
}
