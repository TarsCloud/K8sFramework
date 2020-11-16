
## Create 接口定义

### Create Metadata

```gotemplate
	type CreateIngressHostMetadata struct {
		DomainValue string
	}
```

## Select 接口定义

> select 全部字段

```gotemplate
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"DomainId": SqlColumn{
			ColumnName: "f_domain_id",
			DataType:   "int",
		},
		"DomainValue": SqlColumn{
			ColumnName: "f_domain_value",
			DataType:   "string",
		},
		"FirstLevelDomainValue": SqlColumn{
			ColumnName: "f_first_lever",
			DataType:   "string",
		},
		"SecondLevelDomainValue": SqlColumn{
			ColumnName: "f_second_lever",
			DataType:   "string",
		},
		"ThirdLevelDomainValue": SqlColumn{
			ColumnName: "f_third_level",
			DataType:   "string",
		},
		"FourthLevelDomainValue": SqlColumn{
			ColumnName: "f_fourth_level",
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
##  不支持 update


## Delete接口定义

> delete metadata

```gotemplate
	type DeleteIngressHostMetadata struct {
		DomainId int
	}
```