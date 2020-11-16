
## Select 接口定义

> select 全部字段

```gotemplate
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"ServerId": SqlColumn{
			ColumnName: "f_server_id",
			DataType:   "int",
		},
		"ServerApp": SqlColumn{
			ColumnName: "f_server_app",
			DataType:   "string",
		},
		"ServerName": SqlColumn{
			ColumnName: "f_server_name",
			DataType:   "string",
		},
		"Replicas": SqlColumn{
			ColumnName: "f_replicas",
			DataType:   "int",
		},
		"NodeSelector": SqlColumn{
			ColumnName: "f_node_selector",
			DataType:   "json",
		},
		"HostIpc": SqlColumn{
			ColumnName: "f_host_ipc",
			DataType:   "bool",
		},
		"HostNetwork": SqlColumn{
			ColumnName: "f_host_network",
			DataType:   "bool",
		},
		"HostPort": SqlColumn{
			ColumnName: "f_host_port",
			DataType:   "json",
		},
	}
```

## Update 接口定义
> update metadata
```gotemplate
	type UpdateK8SMetadata struct {
		ServerId int `json:"ServerId" valid:"required,matches-ServerId"`
	}
```

> update 支持字段

```gotemplate
    [] string {"Replicas", "NodeSelector","HostNetwork","HostIpc","HostPort"}
```