package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	configName    string
	configContent string
}

type serverK8S struct {
	// server
	id int
	name string
	app string
	server string
	// image
	image string
	tag string
	// option
	template string
	profile string
	// adapter
	objs []Servant
	// k8s
	replicas int
	nodeSelc NodeSelector
	hostPorts []Hostport
	// config
	configs []Config
}

type templateTars struct {
	parent string
	content string
}

var err error
var serverK8SCache = make(map[int]serverK8S)
var templateTarsCache = make(map[string]templateTars)

var _DOCKER_REGISTRY_URL_ = "registry.cn-shenzhen.aliyuncs.com/tars-k8s"

func LoadK8SDBServerData() {
	server, err := selectServer(globalK8SDb)
	if err.ErrorCode != 0 {
		panic(err)
	}
	for i := 0; i < len(server.Data); i++ {
		id := server.Data[i]["f_server_id"].(int)
		_, ok := serverK8SCache[id]
		if ok {
			panic(fmt.Sprintf("t_server f_server_id: %d has duplicated number.", id))
		}

		one := serverK8S{}

		one.id = id
		one.name = strings.ToLower(server.Data[i]["f_server_app"].(string))+"-"+strings.ToLower(server.Data[i]["f_server_name"].(string))
		one.app = server.Data[i]["f_server_app"].(string)
		one.server = server.Data[i]["f_server_name"].(string)

		serverK8SCache[id] = one
	}

	service, err 	:= selectServicePool(globalK8SDb)
	if err.ErrorCode != 0 {
		panic(err)
	}
	for i := 0; i < len(service.Data); i++ {
		id := service.Data[i]["f_server_id"].(int)
		item, ok := serverK8SCache[id]
		if !ok {
			panic(fmt.Sprintf("t_service_pool f_server_id: %d has not existed in t_server.", id))
		}

		image := service.Data[i]["f_service_image"].(string)
		item.image = strings.Replace(image, _DOCKER_REGISTRY_URL_, "_DOCKER_REGISTRY_URL_", 1)
		item.tag = "10000"

		serverK8SCache[id] = item
	}

	option, err 	:= selectServerOption(globalK8SDb)
	if err.ErrorCode != 0 {
		panic(err)
	}
	for i := 0; i < len(option.Data); i++ {
		id := option.Data[i]["f_server_id"].(int)
		item, ok := serverK8SCache[id]
		if !ok {
			panic(fmt.Sprintf("t_server_option f_server_id: %d has not existed in t_server.", id))
		}

		item.template = option.Data[i]["f_server_template"].(string)
		item.profile = option.Data[i]["f_server_profile"].(string)

		serverK8SCache[id] = item
	}

	adapter, err 	:= selectServerAdapter(globalK8SDb)
	if err.ErrorCode != 0 {
		panic(err)
	}
	for i := 0; i < len(adapter.Data); i++ {
		id := adapter.Data[i]["f_server_id"].(int)
		item, ok := serverK8SCache[id]
		if !ok {
			panic(fmt.Sprintf("t_server_adapter f_server_id: %d has not existed in t_server.", id))
		}

		if item.objs == nil {
			item.objs = make([]Servant, 0)
		}

		obj := Servant{}
		obj.Name = adapter.Data[i]["f_name"].(string)
		obj.Port = adapter.Data[i]["f_port"].(int)
		obj.IsTars = adapter.Data[i]["f_is_tars"].(bool)
		item.objs = append(item.objs, obj)

		serverK8SCache[id] = item
	}

	k8s, err := selectK8S(globalK8SDb)
	if err.ErrorCode != 0 {
		panic(err)
	}
	for i := 0; i < len(k8s.Data); i++ {
		id := k8s.Data[i]["f_server_id"].(int)
		item, ok := serverK8SCache[id]
		if !ok {
			panic(fmt.Sprintf("t_server_k8s f_server_id: %d has not existed in t_server.", id))
		}

		item.replicas = k8s.Data[i]["f_replicas"].(int)

		if k8s.Data[i]["f_node_selector"] != nil {
			var selector = struct {
				Kind string `json:"Kind"`
				Value []string `json:"Value"`
			}{}
			if err := json.Unmarshal(k8s.Data[i]["f_node_selector"].(json.RawMessage), &selector); err != nil {
				panic(fmt.Sprintf("t_server_k8s f_server_id: %d unmarshal node_selector json err: %s.", id, err))
			}
			if selector.Kind == "AbilityPool" {
				item.nodeSelc.AbilityPool.Values = selector.Value
				item.nodeSelc.NodeBind.Values = nil
			} else if selector.Kind == "NodeBind" {
				item.nodeSelc.AbilityPool.Values = nil
				item.nodeSelc.NodeBind.Values = selector.Value
			} else {
				panic(fmt.Sprintf("t_server_k8s f_server_id: %d unmarshal unknow node_selector type: %s.", id, selector.Kind))
			}
		}

		if k8s.Data[i]["f_host_port"] != nil {
			var hostPorts = make(map[string]int)
			if err := json.Unmarshal(k8s.Data[i]["f_host_port"].(json.RawMessage), &hostPorts); err != nil {
				panic(fmt.Sprintf("t_server_k8s f_server_id: %d unmarshal host_port json err: %s.", id, err))
			}
			for k, v := range hostPorts {
				obj := Hostport{NameRef: k, Port: v}
				item.hostPorts = append(item.hostPorts, obj)
			}
		}

		serverK8SCache[id] = item
	}

	config, err := selectConfig(globalK8SDb)
	if err.ErrorCode != 0 {
		panic(err)
	}
	for i := 0; i < len(config.Data); i++ {
		appServer := config.Data[i]["f_app_server"].(string)
		res := strings.Split(appServer, ".")
		if len(res) != 2 {
			fmt.Println(fmt.Sprintf("here is an unexpected config: %s", appServer))
			continue
		}
		id := -1
		for k, v := range serverK8SCache {
			if res[0] == v.app && res[1] == v.server {
				id = k
				break
			}
		}
		if id == -1 {
			panic(fmt.Sprintf("config %s is not in server.?", appServer))
		}

		item, _ := serverK8SCache[id]
		if item.configs == nil {
			item.configs = make([]Config, 0, 5)
		}
		item.configs = append(item.configs, Config{
			configName: config.Data[i]["f_config_name"].(string),
			configContent: config.Data[i]["f_config_content"].(string),
			})
		serverK8SCache[id] = item
	}
}
func AdapterK8SDBTServerData(tarsserver *TServer, request BuildRequest, release ReleaseImageItem) bool {
	for _, val := range serverK8SCache {
		if request.ServerApp == val.app && request.ServerName == val.server {
			TServerAdapter(tarsserver, val, release)
			TConfigAdapter(val, request)
			return true
		}
	}
	return false
}

func LoadTarsDBTemplateData() {
	template, err 	:= selectTemplate(globalK8SDb)
	if err.ErrorCode != 0 {
		panic(err)
	}
	for i := 0; i < len(template.Data); i++ {
		name := template.Data[i]["f_template_name"].(string)
		_, ok := templateTarsCache[name]
		if ok {
			panic(fmt.Sprintf("t_template f_template_name: %s has duplicated number.", name))
		}

		one := templateTars{}

		one.parent = template.Data[i]["f_template_parent"].(string)
		if template.Data[i]["f_template_content"] != nil {
			one.content = template.Data[i]["f_template_content"].(string)
		}

		templateTarsCache[name] = one
	}
}
func DumpTTempateData()  {
	template, err := ioutil.ReadFile(ttemplateYamlPath)
	if err != nil {
		panic(fmt.Sprintf("read from %s err: %s\n", "ttemplate.yaml", err))
	}
	tarstemplate := &TTemplate{}
	err = yaml.Unmarshal(template, &tarstemplate)
	if err != nil {
		panic(fmt.Sprintf("unmarshal from %s err: %s\n", "ttemplate.yaml", err))
	}
	for key, val := range templateTarsCache {
		tarstemplate.Metadata.Name = key
		tarstemplate.Spec.Parent = val.parent
		tarstemplate.Spec.Content = val.content

		output, err := yaml.Marshal(&tarstemplate)
		if err != nil {
			panic(fmt.Sprintf("marshal from %v err: %s\n", tarstemplate, err))
		}
		_ = ioutil.WriteFile(fmt.Sprintf("%s/%s.yaml", AppTemplateDir, key), output, os.ModePerm)
	}
}

func TServerAdapter(tarsserver *TServer, val serverK8S, release ReleaseImageItem) {
	name := val.name
	tarsserver.Metadata.Name = val.name
	tarsserver.Metadata.Namespace = Namespace

	tarsserver.Spec.App = val.app
	tarsserver.Spec.Server = val.server

	tarsserver.Spec.Release.Source = name
	tarsserver.Spec.Release.Image = release.Image
	tarsserver.Spec.Release.Tag = release.Tag
	tarsserver.Spec.Release.ServerType = release.ServerType

	if strings.Contains(val.template, "nodejs") {
		tarsserver.Spec.Tars.Template = "tars.nodejs"
	} else {
		tarsserver.Spec.Tars.Template = val.template
	}
	tarsserver.Spec.Tars.Servants = val.objs
	if len(val.profile) > 5 {
		tarsserver.Spec.Tars.Profile = val.profile
	} else {
		tarsserver.Spec.Tars.Profile = ""
	}

	tarsserver.Spec.K8S.Replicas = val.replicas
	tarsserver.Spec.K8S.HostPorts = val.hostPorts
	if val.nodeSelc.AbilityPool.Values == nil {
		delete(tarsserver.Spec.K8S.NodeSelector, "abilityPool")
		tarsserver.Spec.K8S.NodeSelector["nodeBind"] = val.nodeSelc.NodeBind
	} else {
		delete(tarsserver.Spec.K8S.NodeSelector, "nodeBind")
		tarsserver.Spec.K8S.NodeSelector["abilityPool"] = val.nodeSelc.AbilityPool
	}

	output, err := yaml.Marshal(&tarsserver)
	if err != nil {
		panic(fmt.Sprintf("marshal from %v err: %s\n", tarsserver, err))
	}
	_ = ioutil.WriteFile(fmt.Sprintf("%s/%s.yaml", AppServerDir, name), output, os.ModePerm)
}

func TConfigAdapter(val serverK8S, request BuildRequest) {
	config, err := ioutil.ReadFile(request.TConfigTemplatePath)
	if err != nil {
		panic(fmt.Sprintf("read from %s err: %s\n", "tconfig.yaml", err))
	}
	tarsconfig := &TConfig{}
	err = yaml.Unmarshal(config, &tarsconfig)
	if err != nil {
		panic(fmt.Sprintf("unmarshal from %s err: %s\n", "tconfig.yaml", err))
	}

	for _, config := range val.configs {
		name := strings.ToLower(fmt.Sprintf("%s-%s-%s", val.app, val.server, config.configName))

		tarsconfig.Metadata.Name = name
		tarsconfig.ServerConfig.App = val.app
		tarsconfig.ServerConfig.Server = val.server
		tarsconfig.ServerConfig.ConfigName = config.configName
		tarsconfig.ServerConfig.ConfigContent = config.configContent

		output, err := yaml.Marshal(&tarsconfig)
		if err != nil {
			panic(fmt.Sprintf("marshal from %v err: %s\n", tarsconfig, err))
		}
		_ = ioutil.WriteFile(fmt.Sprintf("%s/%s.yaml", AppConfigDir, name), output, os.ModePerm)
	}
}
