package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"hash/crc32"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type MigrateServer struct {
	Nodes      []string `json:"Nodes"`
	ServerApp  string   `json:"ServerApp",valid:"required,matches-ServerApp"`
	ServerName []string `json:"ServerNames",valid:"required,each-matches-ServerName"`
}

type ConfigMigrateServer struct {
	MigrateServer      []MigrateServer `json:"MigrateServer",valid:"required,matches-MigrateServer"`
}

type ServerList struct {
	Info   *InfoFromTars `json:"info"`
	result string
}

var globalMigrateConfig ConfigMigrateServer
var globalServerList []ServerList


func LoadTarsDBServerData() bool {
	if bs, err := ioutil.ReadFile(migrateConfigPath); err != nil {
		fmt.Println(fmt.Sprintf("Open Config File Error: %s", err))
		return false
	} else {
		err := json.Unmarshal(bs, &globalMigrateConfig)
		if err != nil {
			fmt.Print("Config File Had Bad Format")
			return false
		}
	}

	for _, migrateServer := range globalMigrateConfig.MigrateServer {
		for _, serverName := range migrateServer.ServerName {
			info, err := LoadInfoFromTars(migrateServer.ServerApp, serverName)
			if err != nil {
				fmt.Printf("Load Info From TarsDB Error , App = %s ,ServerName = %s, Error: %s\n", migrateServer.ServerApp, serverName, err.Error())
				return false
			}
			if globalServerList == nil {
				globalServerList = make([]ServerList, 0, 100)
			}
			globalServerList = append(globalServerList, ServerList{info, serverName})
		}
	}

	return true
}

func LoadDumpFromHZServerData(dumpedServerDataPath string) bool {
	if bs, err := ioutil.ReadFile(dumpedServerDataPath); err != nil {
		fmt.Println(fmt.Sprintf("Open Dump File Error: %s", err))
		return false
	} else {
		err := json.Unmarshal(bs, &globalServerList)
		if err != nil {
			fmt.Print("Config File Had Bad Format")
			return false
		}
	}

	return true
}

func AdapterTarsDBTServerData(tarsserver *TServer, request BuildRequest, release ReleaseImageItem) bool {
	for _, val := range globalServerList {
		if request.ServerApp == val.Info.ServerApp && request.ServerName == val.Info.ServerName {
			enableTServer(tarsserver, val.Info, release)
			enableConfig(val.Info, request)
			return true
		}
	}
	return false
}

func enableTServer(tarsserver *TServer, info *InfoFromTars, release ReleaseImageItem)  {
	name := fmt.Sprintf("%s-%s", strings.ToLower(info.ServerApp), strings.ToLower(info.ServerName))
	tarsserver.Metadata.Name = name
	tarsserver.Metadata.Namespace = Namespace

	tarsserver.Spec.App = info.ServerApp
	tarsserver.Spec.Server = info.ServerName

	tarsserver.Spec.Release.Source = name
	tarsserver.Spec.Release.ServerType = release.ServerType
	tarsserver.Spec.Release.Tag = release.Tag
	tarsserver.Spec.Release.Image = release.Image

	if strings.Contains(info.ServerOption.ServerTemplate, "nodejs") {
		tarsserver.Spec.Tars.Template = "tars.nodejs"
	} else {
		tarsserver.Spec.Tars.Template = info.ServerOption.ServerTemplate
	}
	if len(info.ServerOption.ServerProfile) > 5 {
		tarsserver.Spec.Tars.Profile = info.ServerOption.ServerProfile
	} else {
		tarsserver.Spec.Tars.Profile = ""
	}
	tarsserver.Spec.Tars.Servants = make([]Servant, 0, 0)
	for _, servant := range info.ServerAdapters {
		objSeps := strings.Split(servant.Name, ".")
		tarsserver.Spec.Tars.Servants = append(tarsserver.Spec.Tars.Servants,
			Servant{Name: objSeps[len(objSeps)-1], Port: servant.Port, IsTars: servant.IsTars})
	}

	tarsserver.Spec.K8S.Replicas = 1
	tarsserver.Spec.K8S.HostPorts = nil

	delete(tarsserver.Spec.K8S.NodeSelector, "nodeBind")
	tarsserver.Spec.K8S.NodeSelector["abilityPool"] = NodeSelectorValues{Values: make([]string, 0, 1)}

	output, err := yaml.Marshal(&tarsserver)
	if err != nil {
		panic(fmt.Sprintf("marshal from %v err: %s\n", tarsserver, err))
	}
	_ = ioutil.WriteFile(fmt.Sprintf("%s/%s.yaml", AppServerDir, name), output, os.ModePerm)
}

func enableConfig(info *InfoFromTars, request BuildRequest) {
	config, err := ioutil.ReadFile(request.TConfigTemplatePath)
	if err != nil {
		panic(fmt.Sprintf("read from %s err: %s\n", "tconfig.yaml", err))
	}
	tarsconfig := &TConfig{}
	err = yaml.Unmarshal(config, &tarsconfig)
	if err != nil {
		panic(fmt.Sprintf("unmarshal from %s err: %s\n", "tconfig.yaml", err))
	}

	for _, config := range info.ServerConfig {
		crc := crc32.ChecksumIEEE([]byte(config.CurrentConfigFile))
		name := strings.ToLower(fmt.Sprintf("%s-%s-%s", info.ServerApp, info.ServerName, strconv.FormatUint(uint64(crc), 10)))

		tarsconfig.Metadata.Name = name
		tarsconfig.ServerConfig.App = info.ServerApp
		tarsconfig.ServerConfig.Server = info.ServerName
		tarsconfig.ServerConfig.ConfigName = config.CurrentConfigFile
		tarsconfig.ServerConfig.ConfigContent = config.CurrentConfigFileContent

		output, err := yaml.Marshal(&tarsconfig)
		if err != nil {
			panic(fmt.Sprintf("marshal from %v err: %s\n", tarsconfig, err))
		}
		_ = ioutil.WriteFile(fmt.Sprintf("%s/%s.yaml", AppConfigDir, name), output, os.ModePerm)
	}
}

