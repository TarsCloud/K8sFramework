# TafApp 资源说明
```text
    Business 用于查看,变更 Business 时使用
```

## Create 接口定义

### Create Metadata

```gotemplate
	type CreateBusinessMetadata struct {
		BusinessName  string `json:"BusinessName"  valid:"required"`
		BusinessShow  string `json:"BusinessShow"`
		BusinessMark  string `json:"BusinessMark"`
		BusinessOrder int    `json:"BusinessOrder"`
	}
```

## Select 接口定义

> select 全部字段

```gotemplate
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"BusinessName": SqlColumn{
			ColumnName: "f_business_name",
			DataType:   "string",
		},
		"BusinessShow": SqlColumn{
			ColumnName: "f_business_show",
			DataType:   "string",
		},
		"BusinessMark": SqlColumn{
			ColumnName: "f_business_mark",
			DataType:   "string",
		},
		"BusinessOrder": SqlColumn{
			ColumnName: "f_business_order",
			DataType:   "int",
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
+ AllCount , 用与查询所有的 BusinessName 数量

## Update 接口定义

> update metadata

```gotemplate
	type UpdateBusinessMetadata struct {
		BusinessName string `json:"BusinessName"`
	}
```

> update 支持字段

```gotemplate
    []string { "BusinessMark","BusinessShow", "BusinessOrder"} 
```

## Delete接口定义

> delete metadata

```gotemplate
	type DeleteBusinessResourceMetadata struct {
		BusinessName int `json:"BusinessName"`
	}
```

## do 接口定义

### Action AddBusinessApp
+ 用于批量增加 某个 BusinessName 包含的 App

> 请求格式
```gotemplate
	type AddBusinessApp struct {
		BusinessName string   `json:"BusinessName" valid:"required,matches-BusinessName"`
		AppName      []string `json:"AppName" valid:"required,each-matches-ServerApp"`
	}
```

### Action DeleteBusinessApp
> 用于批量取消 某个 BusinessName 包含的 App

> 请求格式
```gotemplate
	type DeleteBusinessApp struct {
		BusinessName string   `json:"BusinessName" valid:"required,matches-BusinessName"`
		AppName      []string `json:"AppName" valid:"required,each-matches-ServerApp"`
	}
```

### Action ListBusinessApp

> 用于 聚合输出 某个 BusinessName 包含的 App

```gotemplate
    type ListBusinessApp struct {
    	BusinessName []string `json:"BusinessName" valid:"each-matches-BusinessName"`
    }
```

> 请求格式说明
+ BusinessName 字段可以忽略,也可以填写BusinessName .如果忽略.则响应全部的 Business 对应的App .如果有值,则只响应指定的BusinessName 对应App

> 响应格式

```json
{
    "result": [
        {
            "BusinessName": "ZSZQ",
            "BusinessShow": "招商证券",
            "App": [
                "Login",
                "News"
            ]
        },
        {
            "BusinessName": "ZXZQ",
            "BusinessShow": "中信证券",
            "App": [
                "Agent"
            ]
        },
        {
            "BusinessName": "YHZQ",
            "BusinessShow": "银河证券2",
            "App": []
        },
        {
            "BusinessName": "",
            "BusinessShow": "",
            "App": [
                "Push",
                "Test"
            ]
        }
    ]
}
```



