## Create 接口定义

### Create Metadata

```gotemplate
type ServerApprovalCreateMetadata struct {
	DeployId       int    `json:"DeployId" valid:"required"`
	ApprovalResult bool   `json:"ApprovalResult"`
	ApprovalMark   string `json:"ApprovalMark"`
}
```

## Select 接口定义

> select 全部字段

```gotemplate
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"ApprovalTime": SqlColumn{
			ColumnName: "f_approval_time",
			DataType:   "string",
		},
		"ApprovalPerson": SqlColumn{
			ColumnName: "f_approval_person",
			DataType:   "string",
		},
		"ApprovalResult": SqlColumn{
			ColumnName: "f_approval_result",
			DataType:   "bool",
		},
		"ApprovalMark": SqlColumn{
			ColumnName: "f_approval_mark",
			DataType:   "string",
		},
		"RequestTime": SqlColumn{
			ColumnName: "f_request_time",
			DataType:   "string",
		},
		"RequestPerson": SqlColumn{
			ColumnName: "f_request_person",
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
		"ServerMark": SqlColumn{
			ColumnName: "f_server_mark",
			DataType:   "string",
		},
		"ServerK8S": SqlColumn{
			ColumnName: "f_server_k8s",
			DataType:   "json",
		},
		"ServerOption": SqlColumn{
			ColumnName: "f_server_option",
			DataType:   "json",
		},
		"ServerServant": SqlColumn{
			ColumnName: "f_server_servant",
			DataType:   "json",
		},
	}
```