
## Create 接口定义

### Create Metadata

```gotemplate
	type CreateServicePoolMetadata struct {
		ServerId     int              `json:"ServerId" valid:"required,matches-ServerId"`
		ServiceImage string           `json:"ServiceImage" valid:"required,matches-ServiceImage"`
		ImageDetail  *json.RawMessage `json:"ImageDetail" valid:"-"`
		ServiceMark  string           `json:"ServiceMark" valid:"-"`
	}

```

## Select 接口定义

> select 全部字段

```gotemplate
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"ServiceId": SqlColumn{
			ColumnName: "f_service_id",
			DataType:   "int",
		},
		"ServiceVersion": SqlColumn{
			ColumnName: "f_service_version",
			DataType:   "int",
		},
		"ServiceMark": SqlColumn{
			ColumnName: "f_service_mark",
			DataType:   "string",
		},
		"ServiceImage": SqlColumn{
			ColumnName: "f_service_image",
			DataType:   "string",
		},
		"ImageDetail": SqlColumn{
			ColumnName: "f_image_detail",
			DataType:   "json",
		},
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
```text
    暂时 不支持 attach
```

## Update接口定义

> update metadata
```gotemplate
	type UpdateServicePoolMetadata struct {
		ServiceId int `json:"ServiceId" valid:"required,matches-ServiceId"`
	}    
```

> update 支持字段
```gotemplate
    [] string = {"ServiceMark"}
```

## Delete接口定义

> delete metadata

```gotemplate
	type DeleteServicePoolMetadata struct {
		ServiceId int `json:"ServiceId"  valid:"required,matches-ServiceId"`
	}
```

## do 接口定义
### Action EnableService
> 用于启用某个版本
```gotemplate
	type EnableService struct {
		ServiceId  int    `json:"ServiceId"  valid:"required,matches-ServiceId"`
		EnableMark string `json:"EnableMark" valid:"-"`
	}
```