

## Select 接口定义

> select 全部字段

```gotemplate
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"ServerId": SqlColumn{
			ColumnName: "a.f_server_id",
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
		"ServiceId": SqlColumn{
			ColumnName: "f_service_id",
			DataType:   "int",
		},
		"ServiceVersion": SqlColumn{
			ColumnName: "f_service_version",
			DataType:   "int",
		},
		"ServiceImage": SqlColumn{
			ColumnName: "f_service_image",
			DataType:   "string",
		},
		"ImageDetail": SqlColumn{
			ColumnName: "f_image_detail",
			DataType:   "string",
		},
		"EnableTime": SqlColumn{
			ColumnName: "f_enable_time",
			DataType:   "string",
		},
		"EnablePerson": SqlColumn{
			ColumnName: "f_enable_person",
			DataType:   "string",
		},
		"EnableMark": SqlColumn{
			ColumnName: "f_enable_mark",
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
```

> select attach 支持
```text
    暂时 不支持 attach
```