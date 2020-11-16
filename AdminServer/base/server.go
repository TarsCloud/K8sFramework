package base

type ServerServantElem struct {
	Name        string `json:"Name" valid:"required,matches-ServantName"`
	Port        int    `json:"Port" valid:"required,matches-ServantPort"`
	Threads     int    `json:"Threads" valid:"required,matches-ServantThreads"`
	Connections int    `json:"Connections" valid:"required,matches-ServantConnections"`
	Capacity    int    `json:"Capacity" valid:"required,matches-ServantCapacity"`
	Timeout     int    `json:"Timeout" valid:"required,matches-ServantTimeout"`
	IsTaf      bool   `json:"IsTaf" valid:"-"`
	IsTcp       bool   `json:"IsTcp" valid:"-"`
}

type UpdateServantKey string

const (
	ServantName        UpdateServantKey = "Name"
	ServantPort        UpdateServantKey = "Port"
	ServantThreads     UpdateServantKey = "Threads"
	ServantConnections UpdateServantKey = "Connections"
	ServantCapacity    UpdateServantKey = "Capacity"
	ServantTimeout     UpdateServantKey = "Timeout"
	ServantIsTaf      UpdateServantKey = "IsTaf"
	ServantIsTcp       UpdateServantKey = "IsTcp"
)

type ServerServant map[string]*ServerServantElem

type NodeSelectorKind string

const (
	NodeBind    NodeSelectorKind = "NodeBind"
	AbilityPool NodeSelectorKind = "AbilityPool"
	PublicPool  NodeSelectorKind = "PublicPool"
)

type NodeSelector struct {
	Kind NodeSelectorKind `json:"Kind"`
	// 值选项有 NodeBind , AbilityPool , PublicPool
	Value []string `json:"Value"`
	// 如果 Kind 为 NodeBind  则值为 NodeName 数组
	// 如果 Kind 为 AbilityPool 或   PublicPool ,则此值为空
}

type ServerK8S struct {
	Replicas int32 `json:"Replicas"`

	NodeSelector NodeSelector `json:"NodeSelector"`

	// 是否允许堆叠, 允许堆叠则一台 Node 可以部署同一个服务的多个Pod, 不允许堆叠则一台Node只能部署同一个服务的一个Pod
	NotStacked bool `json:"NotStacked"`

	// NodeSelect.Kind 值为 NodeBind 时才可以选择
	HostIpc bool `json:"HostIpc"`

	// NodeSelect.Kind 值为 NodeBind 时才可以选择，  与 HostPort 互斥
	HostNetwork bool `json:"HostNetwork"`

	// NodeSelect.Kind 值为 NodeBind 时才可以选择
	HostPort map[string]int32 `json:"HostPort"`

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
