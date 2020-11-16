
## Create 接口定义

### Create Metadata

```gotemplate
	type CreateAppMetadata struct {
		NodeName  string `json:"NodeName" valid:"required"`
		AppServer string `json:"AppServer"`
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
		"AppServer": SqlColumn{
			ColumnName: "f_app_server",
			DataType:   "string",
		},
	}
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

## do 接口定义
### Action DeleteNodeEnableServer
> 用于批量删除 指定Node 支持的 Server
```gotemplate
    type DeleteNodeEnableServer struct {
    	NodeName  string   `json:"NodeName"`
   		AppServer []string `json:"AppServer"`
    }
```

### Action AddNodeEnableServer
> 用于批量增加 指定Node 支持的 Server
```gotemplate
    type AddNodeEnableServer struct {
    	NodeName  string   `json:"NodeName"`
   		AppServer []string `json:"AppServer"`
    }
```

### Action DeleteServerEnableNode
> 用于批量删除 支持 Server 的 Node
```gotemplate
    type DeleteServerEnableNode struct {
         AppServer  string   `json:"AppServer"`
         NodeName []string `json:"NodeName"`
    }
```

### Action AddServerEnableNode
> 用于批量增加 支持 Server 的 Node
```gotemplate
    type AddServerEnableNode struct {
         AppServer  string   `json:"AppServer"`
         NodeName []string `json:"NodeName"`
    }
```

### Action ListAffinityGroupByNode
> 用于聚合 单个Node的 EnableAppServer

> 请求格式
```gotemplate
	type ListAffinityGroupByNode struct {
		NodeName []string `json:"NodeName" valid:"each-matches-NodeName"`
	}
```

> 格式说明
+ NodeName 字段可以忽略,也可以填写 NodeName 值.如果忽略,则响应全部 Node对应的 AppServer .如果填写了值,则只响应 填写值对应的AppServer

> 响应结果
```json
{
    "result" : [
        {
            "NodeName": "kube.node68",
            "AppServer": [
                "Login",
                "News",
                "Test"
            ]
        },
        {
            "NodeName": "kube.node94",
            "AppServer": [
                "Login",
                "NewS",
                "Test"
            ]
        }
    ]
}
```

### Action ListAffinityGroupByAppServer
> 用于聚合 可以部署 AppServer 的 Node

>请求格式
```gotemplate
	type ListAffinityGroupByAppServer struct {
		AppServer []string `json:"AppServer" valid:"each-matches-AppServer"`
	}
```

> 格式说明
+ AppServer 字段可以忽略,也可以填写 AppServer 值.如果忽略,则响应全部 AppServer 对应的 NodeName .如果填写了值,则只响应 填写值对应的 NodeName

> 响应结果
```json
{
    "result": [
        {
            "AppServer": "Login",
            "NodeName": [
                "kube.node68",
                "kube.node94"
            ]
        },
        {
            "AppServer": "News",
            "NodeName": [
                "kube.node68"
            ]
        },
        {
            "AppServer": "NewS",
            "NodeName": [
                "kube.node94"
            ]
        },
        {
            "AppServer": "Test",
            "NodeName": [
                "kube.node68",
                "kube.node94"
            ]
        }
    ]
}
```