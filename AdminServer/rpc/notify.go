package rpc

const _RequestKindNotify = "Notify"

func SelectNotify(rpcRequest *Request) (Result, Error) {
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"NotifyId": SqlColumn{
			ColumnName: "f_notify_id",
			ColumnType: "int",
		},
		"AppServer": SqlColumn{
			ColumnName: "f_app_server",
			ColumnType: "string",
		},
		"PodName": SqlColumn{
			ColumnName: "f_pod_name",
			ColumnType: "string",
		},
		"NotifyLevel": SqlColumn{
			ColumnName: "f_notify_level",
			ColumnType: "string",
		},
		"NotifyMessage": SqlColumn{
			ColumnName: "f_notify_message",
			ColumnType: "string",
		},
		"NotifyTime": SqlColumn{
			ColumnName: "f_notify_time",
			ColumnType: "string",
		},
		"NotifyThread": SqlColumn{
			ColumnName: "f_notify_thread",
			ColumnType: "string",
		},
		"NotifySource": SqlColumn{
			ColumnName: "f_notify_source",
			ColumnType: "string",
		},
	}
	const from = "t_server_notify"
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
	registryHandle(RequestMethodSelect+_RequestKindNotify, SelectNotify)
}
