
## Create 接口定义

### Create Metadata

```gotemplate
	type CreateHttpRouterMetadata struct {
		DomainValue        string
		MatchValue         string
		MatchType          string
		BackendServerId    int
		BackendAdapterId   string
	}
```

## Select 接口定义

> select 全部字段

```gotemplate
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"RouterId": SqlColumn{
			ColumnName: "f_router_id",
			DataType:   "int",
		},
		"DomainValue": SqlColumn{
			ColumnName: "f_domain_value",
			DataType:   "string",
		},
		"MatchValue": SqlColumn{
			ColumnName: "f_match_value",
			DataType:   "string",
		},
		"MatchType": SqlColumn{
			ColumnName: "f_match_value",
			DataType:   "string",
		},
		"BackendServerId": SqlColumn{
			ColumnName: "f_backend_server_id",
			DataType:   "int",
		},
		"BackendServerApp": SqlColumn{
			ColumnName: "f_backend_server_app",
			DataType:   "int",
		},
		"BackendServerName": SqlColumn{
			ColumnName: "f_backend_server_name",
			DataType:   "int",
		},
		"BackendAdapterId": SqlColumn{
			ColumnName: "f_backend_adapter_id",
			DataType:   "int",
		},
		"BackendAdapterName": SqlColumn{
			ColumnName: "f_backend_adapter_name",
			DataType:   "string",
		},
	}
```

> select attach 支持
+ AllCount , 用与查询所有的 BusinessName 数量


## Delete接口定义

> delete metadata

```gotemplate
	type DeleteHttpRouterMetadata struct {
		RouterId int
	}
```

## do 接口定义

### Action UpdateHttpRouter

> 用于更新一条 HttpRouter 
  ```gotemplate
	type UpdateHttpRouterMetadata struct {
		RouterId         int    `json:"RouterId" valid:"required,matches-RouterId"`
		DomainValue      string `json:"DomainValue" valid:"required,matches-Domain"`
		MatchValue       string `json:"MatchValue" valid:"required,matches-HttpRouterMatchValue"`
		MatchType        string `json:"MatchType" valid:"required,matches-HttpRouterMatchType"`
		BackendServerId  int    `json:"BackendServerId" valid:"required,matches-ServerId"`
		BackendAdapterId string `json:"BackendAdapterId" valid:"required,matches-AdapterId"`
	}
  ```