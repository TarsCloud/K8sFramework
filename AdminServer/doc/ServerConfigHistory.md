

## Select 接口定义

> select 全部字段

```gotemplate
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"HistoryId": SqlColumn{
			ColumnName: "f_history_id",
			DataType:   "int",
		},
		"ConfigId": SqlColumn{
			ColumnName: "f_config_id",
			DataType:   "int",
		},
		"ConfigName": SqlColumn{
			ColumnName: "f_config_name",
			DataType:   "string",
		},
		"ConfigVersion": SqlColumn{
			ColumnName: "f_config_version",
			DataType:   "int",
		},
		"ConfigContent": SqlColumn{
			ColumnName: "f_config_content",
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
		"CreateMark": SqlColumn{
			ColumnName: "f_create_mark",
			DataType:   "string",
		},
		"AppServer": SqlColumn{
			ColumnName: "f_app_server",
			DataType:   "string",
		},
	}
```

## Delete接口定义

> delete metadata

```gotemplate
	type DeleteServerConfigHistoryMetadata struct {
		HistoryId int `json:"HistoryId"  valid:"required,matches-HistoryConfigId"`
	}
```

## do  接口

### Action ActiveHistoryConfig
> 用于 激活(回滚历史版本的配置)

> 请求格式

```gotemplate
	type ChangeVersionMetadata struct {
		HistoryId int `json:"HistoryId"`
	}
```
