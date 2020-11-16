## Select 接口定义

> select 全部字段

```gotemplate
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"PodId": SqlColumn{
			ColumnName: "f_pod_id",
			DataType:   "string",
		},
		"PodName": SqlColumn{
			ColumnName: "f_pod_name",
			DataType:   "string",
		},
		"PodIp": SqlColumn{
			ColumnName: "f_pod_ip",
			DataType:   "string",
		},
		"NodeIp": SqlColumn{
			ColumnName: "f_node_ip",
			DataType:   "string",
		},
		"ServerApp": SqlColumn{
			ColumnName: "f_server_app",
			DataType:   "string",
		},
		"ServerName": SqlColumn{
			ColumnName: "f_server_name",
			DataType:   "string",
		},
		"ServiceVersion": SqlColumn{
			ColumnName: "f_service_version",
			DataType:   "int",
		},
		"CreateTime": SqlColumn{
			ColumnName: "f_create_time",
			DataType:   "string",
		},
		"DeleteTime": SqlColumn{
			ColumnName: "f_delete_time",
			DataType:   "string",
		},
	}
```