
## Create 接口定义

### Create Metadata

```gotemplate
	type CreateTemplateMetadata struct {
		TemplateName    string `json:"TemplateName"`
		TemplateParent  string `json:"TemplateParent"`
		TemplateContent string `json:"TemplateContent"`
		CreateMark      string `json:"CreateMark"`
	}

```

## Select 接口定义

> select 全部字段

```gotemplate
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"TemplateId": SqlColumn{
			ColumnName: "f_template_id",
			DataType:   "int",
		},
		"TemplateName": SqlColumn{
			ColumnName: "f_template_name",
			DataType:   "string",
		},
		"TemplateParent": SqlColumn{
			ColumnName: "f_template_parent",
			DataType:   "string",
		},

		"TemplateContent": SqlColumn{
			ColumnName: "f_template_content",
			DataType:   "string",
		},
		"CreatePerson": SqlColumn{
			ColumnName: "f_create_person",
			DataType:   "string",
		},
		"CreateTime": SqlColumn{
			ColumnName: "f_create_time",
			DataType:   "string",
		},

		"CreateMark": SqlColumn{
			ColumnName: "f_create_mark",
			DataType:   "string",
		},

		"UpdatePerson": SqlColumn{
			ColumnName: "f_update_person",
			DataType:   "string",
		},

		"UpdateMark": SqlColumn{
			ColumnName: "f_update_mark",
			DataType:   "string",
		},

		"UpdateTime": SqlColumn{
			ColumnName: "f_update_time",
			DataType:   "string",
		},
```

## Update 接口定义

> Update metadata

```gotemplate
	type UpdateTemplateMetadata struct {
		TemplateId int `json:"TemplateId" valid:"required,matches-TemplateId"`
	}
```

> update 支持字段

```gotemplate
    []string { "TemplateContent" ,"TemplateParent","CreateMark"} 
```

## Delete接口定义

> delete metadata

```gotemplate
	type DeleteTemplateMetadata struct {
		TemplateId int `json:"TemplateId"`
	}
```
