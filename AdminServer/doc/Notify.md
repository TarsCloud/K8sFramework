
## Select 接口定义

> select 全部字段

```gotemplate
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"NotifyId": SqlColumn{
			ColumnName: "f_notify_id",
			DataType:   "int",
		},
		"AppServer": SqlColumn{
			ColumnName: "f_app_server",
			DataType:   "string",
		},
		"PodName": SqlColumn{
			ColumnName: "f_pod_name",
			DataType:   "string",
		},
		"NotifyLevel": SqlColumn{
			ColumnName: "f_notify_level",
			DataType:   "string",
		},
		"NotifyMessage": SqlColumn{
			ColumnName: "f_notify_message",
			DataType:   "string",
		},
		"NotifyTime": SqlColumn{
			ColumnName: "f_notify_time",
			DataType:   "string",
		},
		"NotifyThread": SqlColumn{
			ColumnName: "f_notify_thread",
			DataType:   "string",
		},
		"NotifySource": SqlColumn{
			ColumnName: "f_notify_source",
			DataType:   "string",
		},
	}
```
