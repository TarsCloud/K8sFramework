package rpc

const _RequestKindPodAlive = "PodAlive"
const _RequestKindPodPerished = "PodPerished"

func SelectPodAlive(rpcRequest *Request) (Result, Error) {

	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"PodId": SqlColumn{
			ColumnName: "f_pod_id",
			ColumnType: "string",
		},
		"PodName": SqlColumn{
			ColumnName: "f_pod_name",
			ColumnType: "string",
		},
		"PodIp": SqlColumn{
			ColumnName: "f_pod_ip",
			ColumnType: "string",
		},
		"NodeIp": SqlColumn{
			ColumnName: "f_node_ip",
			ColumnType: "string",
		},
		"CreateTime": SqlColumn{
			ColumnName: "f_create_time",
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
		"ServiceVersion": SqlColumn{
			ColumnName: "f_service_version",
			ColumnType: "int",
		},
		"SettingState": SqlColumn{
			ColumnName: "f_setting_state",
			ColumnType: "string",
		},
		"PresentState": SqlColumn{
			ColumnName: "f_present_state",
			ColumnType: "string",
		},
		"PresentMessage": {
			ColumnName: "f_present_message",
			ColumnType: "string",
		},
		"UpdateTime": SqlColumn{
			ColumnName: "f_update_time",
			ColumnType: "string",
		},
	}

	const from = "t_pod"

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

func SelectPodPerished(rpcRequest *Request) (Result, Error) {
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"PodId": SqlColumn{
			ColumnName: "f_pod_id",
			ColumnType: "string",
		},
		"PodName": SqlColumn{
			ColumnName: "f_pod_name",
			ColumnType: "string",
		},
		"PodIp": SqlColumn{
			ColumnName: "f_pod_ip",
			ColumnType: "string",
		},
		"NodeIp": SqlColumn{
			ColumnName: "f_node_ip",
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
		"ServiceVersion": SqlColumn{
			ColumnName: "f_service_version",
			ColumnType: "int",
		},
		"CreateTime": SqlColumn{
			ColumnName: "f_create_time",
			ColumnType: "string",
		},
		"DeleteTime": SqlColumn{
			ColumnName: "f_delete_time",
			ColumnType: "string",
		},
	}
	const from = "t_pod_history"

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
	registryHandle(RequestMethodSelect+_RequestKindPodAlive, SelectPodAlive)
	registryHandle(RequestMethodSelect+_RequestKindPodPerished, SelectPodPerished)
}
