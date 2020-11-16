
## Create 接口定义

### Create Metadata

```gotemplate
	type CreateNodeMetadata struct {
		NodeName string `json:"NodeName"  valid:"required,matches-NodeName"`
		NodeMark string `json:"NodeMark"`
	}
```

## Select 接口定义

> select 全部字段

```gotemplate
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"NodeName": SqlColumn{
			ColumnName: "f_node_name",
			DataType:   "string",
		},
		"NodeMark": SqlColumn{
			ColumnName: "f_node_mark",
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

## Update 接口定义

> update metadata

```gotemplate
	type UpdateNodeMetadata struct {
		NodeName string `json:"NodeName"`
	}
```

> update 支持字段

```gotemplate
    []string { "NodeMark" } 
```

## do 接口定义

### Action ListClusterNode
> 用于查询 当前 k8s 集群所有的 节点
```gotemplate
   type ListClusterNode struct {
   }
```
> 响应结果
```json
{
    "result": [
        "kube.master",
        "kube.node67",
        "kube.node68",
        "kube.node94"
    ]
}
```
