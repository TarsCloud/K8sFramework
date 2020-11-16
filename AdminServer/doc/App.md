
## Create 接口定义

### Create Metadata

```gotemplate
	type CreateAppMetadata struct {
		AppName      string `json:"AppName"  valid:"required,matches-ServerApp"`
		AppMark      string `json:"AppMark"  valid:"-"`
		BusinessName string `json:"BusinessName" valid:"matches-BusinessName"`
	}
```

## Select 接口定义

> select 全部字段

```gotemplate
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"AppName": SqlColumn{
			ColumnName: "f_app_name",
			DataType:   "string",
		},
		"AppMark": SqlColumn{
			ColumnName: "f_app_mark",
			DataType:   "string",
		},
		"BusinessName": SqlColumn{
			ColumnName: "f_business_name",
			DataType:   "string",
		},
		"CreateTime": SqlColumn{
			ColumnName: "f_create_time",
			DataType:   "string",
		},
		"CreatePerson": SqlColumn{
			ColumnName: "f_create_person",
			DataType:   "string",
		},
	}
```

> select attach 支持
+ AllCount , 用与查询所有的 app 数量

## Update 接口定义

> update metadata

```gotemplate
	type UpdateAppMetadata struct {
		AppName string `json:"AppName"  valid:"required,matches-ServerApp"`
	}

```

> update 支持字段

```gotemplate
    []string { "AppMark","BusinessName"}
```

## Delete接口定义

> delete metadata

```gotemplate
	type DeleteAppResourceMetadata struct {
		AppName string `json:"AppName"  valid:"required,matches-ServerApp"`
	}
```

