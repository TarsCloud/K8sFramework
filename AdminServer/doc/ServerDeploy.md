
## 特殊结构定义

```gotemplate
type ServerServantElem struct {
	Name        string `json:"Name" valid:"required,matches-ServantName"`
	Port        int    `json:"Port" valid:"required,matches-ServantPort"`
	HostPort    int    `json:"HostPort" valid:"required,matches-ServantPort"`
	Threads     int    `json:"Threads" valid:"required,matches-ServantThreads"`
	Connections int    `json:"Connections" valid:"required,matches-ServantConnections"`
	Capacity    int    `json:"Capacity" valid:"required,matches-ServantCapacity"`
	Timeout     int    `json:"Timeout" valid:"required,matches-ServantTimeout"`
	IsTaf      bool   `json:"IsTaf" valid:"-"`
	IsTcp       bool   `json:"IsTcp" valid:"-"`
}

type ServerServant map[string]ServerServantElem

type ServerK8S struct {
	Replicas int32 `json:"Replicas" valid:"required"`
	NodeSelect struct {
		Kind string `json:"Kind"`
		// 值选项有 NodeBind -> 需要提供 Node 选择窗口 , AbilityPool , PublicPool
		Value []string `json:"Value"`
		// 如果 Kind 为 NodeBind  则值为 NodeName 数组
		// 如果 Kind 为 AbilityPool 或   PublicPool ,则此值为空
	}

	// 是否允许堆叠, 允许堆叠则一台 Node 可以部署同一个服务的多个Pod, 不允许堆叠则一台Node只能部署同一个服务的一个Pod
	NotStacked bool `json:"NotStacked"`

	// NodeSelect.Kind 值为 NodeBind 时才可以选择
	HostIpc bool `json:"HostIpc"`

	// NodeSelect.Kind 值为 NodeBind 时才可以选择，  与 HostPort 互斥
	HostNetwork bool `json:"HostNetwork"`

	HostPort map[string]int

	Image string `json:"Image"`

	Version string `json:"Version"`
}

type ServerOption struct {
	ServerImportant       int    `json:"ServerImportant" valid:"matches-ServerImportant"`
	StartScript           string `json:"StartScript"     valid:"matches-ServerStartScript"`
	StopScript            string `json:"StopScript"      valid:"matches-ServerStopScript"`
	MonitorScript         string `json:"MonitorScript"   valid:"matches-ServerMonitorScript"`
	AsyncThread           int    `json:"AsyncThread"     valid:"required,matches-ServerAsyncThread"`
	ServerTemplate        string `json:"ServerTemplate"  valid:"required,matches-TemplateName"`
	ServerProfile         string `json:"ServerProfile"   valid:"-"`
	RemoteLogEnable       bool   `json:"RemoteLogEnable"`
	RemoteLogReserveTime  int    `json:"RemoteLogReserveTime"`
	RemoteLogCompressTime int    `json:"RemoteLogCompressTime"`
}
```


## Create 接口定义

### Create Metadata

```gotemplate
	type CreateDeployMetadata struct {
		ServerApp      string           `json:"ServerApp"  valid:"required,alphanum"`
		ServerName     string           `json:"ServerName" valid:"required,alphanum"`
		ServerMark     string           `json:"ServerMark"`
		ServerK8S      ServerDetailK8S `json:"ServerK8S" valid:"required"`
		ServerOption   ServerDetailOption `json:"ServerOption" valid:"required"`
		ServerServant  []ServerDetailServantElem `json:"ServerServant" valid:"required"`
	}
```

## Select 接口定义

> select 全部字段

```gotemplate
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"DeployId": SqlColumn{
			ColumnName: "f_request_id",
			DataType:   "int",
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

		"ServerServant": SqlColumn{
			ColumnName: "f_server_servant",
			DataType:   "json",
		},

		"ServerOption": SqlColumn{
			ColumnName: "f_server_option",
			DataType:   "json",
		},
	}
```

## Update 接口定义

> update metadata

```gotemplate
	type UpdateDeployMetadata struct {
		DeployId int `json:"DeployId"`
	}
```

> update 支持字段

```gotemplate
    []string { "ServerServant" ,"ServerOption" ,"ServerK8S"} 
```

## Delete接口定义

> delete metadata

```gotemplate
	type DeleteDeployMetadata struct {
		DeployId int `json:"DeployId"`
	}
```

