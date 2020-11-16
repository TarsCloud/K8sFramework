

## Select 接口定义

> select 全部字段

```gotemplate
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"ServerId": SqlColumn{
			ColumnName: "f_server_id",
			DataType:   "int",
		},
		"ServerTemplate": SqlColumn{
			ColumnName: "f_server_template",
			DataType:   "string",
		},
		"ServerProfile": SqlColumn{
			ColumnName: "f_server_profile",
			DataType:   "string",
		},
		"StartScript": SqlColumn{
			ColumnName: "f_start_script_path",
			DataType:   "string",
		},
		"StopScript": SqlColumn{
			ColumnName: "f_stop_script_path",
			DataType:   "string",
		},
		"MonitorScript": SqlColumn{
			ColumnName: "f_monitor_script_path",
			DataType:   "string",
		},
		"AsyncThread": SqlColumn{
			ColumnName: "f_async_thread",
			DataType:   "int",
		},
		"ServerImportant": SqlColumn{
			ColumnName: "f_important_type",
			DataType:   "int",
		},
		"RemoteLogType": SqlColumn{
			ColumnName: "f_remote_log_type",
			DataType:   "int",
		},
		"RemoteLogReserve": SqlColumn{
			ColumnName: "f_remote_log_reserve_time",
			DataType:   "int",
		},
		"RemoteLogCompress": SqlColumn{
			ColumnName: "f_remote_log_compress_time",
			DataType:   "int",
		},
	}
```
## Update接口定义

> update metadata

```gotemplate
	type UpdateServerOptionMetadata struct {
		ServerId int `json:"ServerId"  valid:"required,matches-ServerId"`
	}
```

> update 支持字段
```gotemplate
    []string { ServerTemplate,ServerProfile,StartScript,StopScript,MonitorScript,AsyncThread,ServerImportant}
``` 

## do 接口定义

### Action  PreviewTemplateContent
> 用于预览查看服务的模板内容(合并 ServerProfile,Template 以及 ParentTemplate之后的内容)
>
```gotemplate
	type PreviewTemplateContentMetadata struct {
		ServerId int `json:"ServerId"  valid:"required,matches-ServerId"`
	}
```

> 响应格式

```json
{
    "result": "xxxxx"
}
```
