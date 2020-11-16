
## Create 接口定义

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
 		"ServerType": SqlColumn{
 			ColumnName: "f_server_type",
 			DataType:   "string",
 		},       
		"ServerMark": SqlColumn{
			ColumnName: "f_server_mark",
			DataType:   "string",
		},
		"DeployPerson": SqlColumn{
			ColumnName: "f_deploy_person",
			DataType:   "string",
		},
		"DeployTime": SqlColumn{
			ColumnName: "f_deploy_time",
			DataType:   "string",
		},
	}
```

## Update 接口定义

> update metadata

```gotemplate
	type UpdateServerMetadata struct {
		ServerId int `json:"ServerId"`
	}
```

> update 支持字段

```gotemplate
    []string { "ServerMark","ServerType"} 
```

