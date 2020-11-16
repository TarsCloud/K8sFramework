
## Create 接口定义

```gotemplate
	type CreateAdapterMetadata struct {
		ServerId int                 `json:"ServerId" valid:"required,matched-ServerId"`
		Servant  ServerDetailServant `json:"Servant" valid:"required,matches-ServerServant"`
	}
```

## Select 接口定义

> select 全部字段

```gotemplate
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"AdapterId": SqlColumn{
			ColumnName: "f_adapter_id",
			DataType:   "int",
		},
		"ServerId": SqlColumn{
			ColumnName: "f_server_id",
			DataType:   "int",
		},
		"Name": SqlColumn{
			ColumnName: "f_name",
			DataType:   "string",
		},
		"Threads": SqlColumn{
			ColumnName: "f_threads",
			DataType:   "int",
		},
		"Connections": SqlColumn{
			ColumnName: "f_connections",
			DataType:   "int",
		},
		"Port": SqlColumn{
			ColumnName: "f_port",
			DataType:   "int",
		},

		"Capacity": SqlColumn{
			ColumnName: "f_capacity",
			DataType:   "int",
		},

		"Timeout": SqlColumn{
			ColumnName: "f_timeout",
			DataType:   "int",
		},

		"IsTaf": SqlColumn{
			ColumnName: "f_is_taf",
			DataType:   "bool",
		},
		"IsTcp": SqlColumn{
			ColumnName: "f_is_tcp",
			DataType:   "bool",
		},
	}
```


## Update接口定义

> update metadata

```gotemplate
	type UpdateServerOptionMetadata struct {
		AdapterId int `json:"AdapterId"  valid:"required,matches-AdapterId"`
	}
```

> update 支持字段
```gotemplate
    []string { Name,Threads,Connections,Port,Capacity,Timeout,IsTaf,IsTcp}
``` 

## Delete接口定义

> delete metadata

```gotemplate
    	type DeleteAdapterMetadata struct {
    		AdapterId int `json:"AdapterId" valid:"required,matched-AdapterId"`
    	}
```