package k8s

import (
	"encoding/json"
	"github.com/go-openapi/runtime/middleware"
	"io/ioutil"
	"os"
	"tafadmin/openapi/models"
	"tafadmin/openapi/restapi/operations/default_operations"
)

// 默认Selector类型
var K8SNodeSelectorKind = models.ArrayString{"AbilityPool", "PublicPool", "NodeBind"}
// 默认Server类型
var ServerTypeOptional = models.ArrayString{"taf_cpp", "taf_java_war", "taf_java_jar", "taf_node", "taf_node8", "taf_node10", "taf_node_pkg", "not_taf"}
// 默认K8S参数
var ServerK8S = models.ServerK8S{
	HostIpc: false, HostNetwork: false, Replicas: 0,
	HostPort: make([]*models.HostPortElem, 0, 1),
	NodeSelector: &models.NodeSelector{AbilityPool: &models.NodeSelectorElem{
	Value: make([]string, 0, 1),
	}},
}
// 默认私有模板参数
var asyncThread int32 = 3
var serverImportant int32 = 5
var serverSubType = "taf"
var ServerOption = models.ServerOption{AsyncThread: &asyncThread, ServerImportant: &serverImportant, ServerSubType: &serverSubType, ServerTemplate: "taf.default", ServerProfile: ""}
// 默认Servant参数
var isTrue = true
var ServerServantElem = models.ServerServantElem{IsTaf: &isTrue, IsTCP: &isTrue, Threads: 3, Port: 10000, Timeout: 60000, Capacity: 10000, Connections: 10000}

// Handler处理
const DefaultConfigMapPath = "/etc/default-env/"

type SelectDefaultValueHandler struct {}

func (s *SelectDefaultValueHandler) Handle(params default_operations.SelectDefaultValueParams) middleware.Responder {
	// 如果Config挂载存在，加载挂载的ConfigMap，否则使用硬编码默认
	_, err := os.Stat(DefaultConfigMapPath)
	if err == nil {
		if bs, err := ioutil.ReadFile(DefaultConfigMapPath +"K8SNodeSelectorKind"); err == nil {
			_ = json.Unmarshal(bs, &K8SNodeSelectorKind)
		}
		if bs, err := ioutil.ReadFile(DefaultConfigMapPath +"ServerTypeOptional"); err == nil {
			_ = json.Unmarshal(bs, &ServerTypeOptional)
		}
		if bs, err := ioutil.ReadFile(DefaultConfigMapPath +"ServerK8S"); err == nil {
			_ = json.Unmarshal(bs, &ServerK8S)
		}
		if bs, err := ioutil.ReadFile(DefaultConfigMapPath +"ServerOption"); err == nil {
			_ = json.Unmarshal(bs, &ServerOption)
		}
		if bs, err := ioutil.ReadFile(DefaultConfigMapPath +"ServerServantElem"); err == nil {
			_ = json.Unmarshal(bs, &ServerServantElem)
		}
	}

	result := default_operations.SelectDefaultValueOKBodyResult{
		K8SNodeSelectorKind: K8SNodeSelectorKind,
		ServerTypeOptional:  ServerTypeOptional,
		ServerK8S:           &ServerK8S,
		ServerOption:        &ServerOption,
		ServerServantElem:   &ServerServantElem}

	return default_operations.NewSelectDefaultValueOK().WithPayload(&default_operations.SelectDefaultValueOKBody{Result: &result})
}
