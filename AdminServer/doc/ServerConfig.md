
## Create 接口定义

### Create Metadata

```gotemplate
	type CreateServerConfigMetadata struct {
		AppServer     string `json:"AppServer" valid:"required,matches-AppServer"`
		ConfigName    string `json:"ConfigName" valid:"required,matches-ConfigName"`
		ConfigContent string `json:"ConfigContent" valid:"required"`
		CreateMark    string `json:"CreateMark" valid:"-"`
		PodSeq        int    `json:"PodSeq" valid:"required,matches-ConfigPodSeq"`
	}
```

## Select 接口定义

> select 全部字段

```gotemplate
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
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
		"PodSeq": SqlColumn{
			ColumnName: "f_pod_seq",
			DataType:   "int",
		},
```

> select attach 支持
```text
    暂时 不支持 attach
```

## Delete接口定义

> delete metadata

```gotemplate
	type DeleteAffinityMetadata struct {
		NodeName  string   `json:"NodeName"`
		AppServer string `json:"AppServer"`
	}
```


## Update接口定义

> update metadata

```gotemplate
	type UpdateServerConfigMetadata struct {
		ConfigId int `json:"ConfigId"  valid:"required,matches-ConfigId"`
	}
```

> update 支持字段
```gotemplate
    []string { ConfigContent,CreateMark}
``` 


## do 接口

### Action PreviewConfigContent

> 用于 预览节点配置与主配置合并后的内容

```gotemplate
	type SeekServerConfigMetadata struct {
		AppServer  string `json:"AppServer" valid:"required,matches-AppServer"`
		ConfigName string `json:"ConfigName" valid:"required,matches-ConfigName"`
		PodSeq     int    `json:"PodSeq" valid:"matches-ConfigPodSeq"`
	}
```

> 响应格式

```json
{
    "result": "xxxxx"
}
```