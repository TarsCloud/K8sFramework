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
	Info   *InfoFromTaf `json:"info"`
	result string
}

var globalMigrateConfig ConfigMigrateServer
var globalServerList []ServerList


func LoadTafDBServerData() bool {
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
			info, err := LoadInfoFromTaf(migrateServer.ServerApp, serverName)
			if err != nil {
				fmt.Printf("Load Info From TafDB Error , App = %s ,ServerName = %s, Error: %s\n", migrateServer.ServerApp, serverName, err.Error())
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

func AdapterTafDBTServerData(tafserver *TServer, request BuildRequest, release ReleaseImageItem) bool {
	for _, val := range globalServerList {
		if request.ServerApp == val.Info.ServerApp && request.ServerName == val.Info.ServerName {
			enableTServer(tafserver, val.Info, release)
			enableConfig(val.Info, request)
			return true
		}
	}
	return false
}

func enableTServer(tafserver *TServer, info *InfoFromTaf, release ReleaseImageItem)  {
	name := fmt.Sprintf("%s-%s", strings.ToLower(info.ServerApp), strings.ToLower(info.ServerName))
	tafserver.Metadata.Name = name
	tafserver.Metadata.Namespace = Namespace

	tafserver.Spec.App = info.ServerApp
	tafserver.Spec.Server = info.ServerName

	tafserver.Spec.Release.Source = name
	tafserver.Spec.Release.ServerType = release.ServerType
	tafserver.Spec.Release.Tag = release.Tag
	tafserver.Spec.Release.Image = release.Image

	if strings.Contains(info.ServerOption.ServerTemplate, "nodejs") {
		tafserver.Spec.Taf.Template = "taf.nodejs"
	} else {
		tafserver.Spec.Taf.Template = info.ServerOption.ServerTemplate
	}
	if len(info.ServerOption.ServerProfile) > 5 {
		tafserver.Spec.Taf.Profile = info.ServerOption.ServerProfile
	} else {
		tafserver.Spec.Taf.Profile = ""
	}
	tafserver.Spec.Taf.Servants = make([]Servant, 0, 0)
	for _, servant := range info.ServerAdapters {
		objSeps := strings.Split(servant.Name, ".")
		tafserver.Spec.Taf.Servants = append(tafserver.Spec.Taf.Servants,
			Servant{Name: objSeps[len(objSeps)-1], Port: servant.Port, IsTaf: servant.IsTaf})
	}

	tafserver.Spec.K8S.Replicas = 1
	tafserver.Spec.K8S.HostPorts = nil

	delete(tafserver.Spec.K8S.NodeSelector, "nodeBind")
	tafserver.Spec.K8S.NodeSelector["abilityPool"] = NodeSelectorValues{Values: make([]string, 0, 1)}

	output, err := yaml.Marshal(&tafserver)
	if err != nil {
		panic(fmt.Sprintf("marshal from %v err: %s\n", tafserver, err))
	}
	_ = ioutil.WriteFile(fmt.Sprintf("%s/%s.yaml", AppServerDir, name), output, os.ModePerm)
}

func enableConfig(info *InfoFromTaf, request BuildRequest) {
	config, err := ioutil.ReadFile(request.TConfigTemplatePath)
	if err != nil {
		panic(fmt.Sprintf("read from %s err: %s\n", "tconfig.yaml", err))
	}
	tafconfig := &TConfig{}
	err = yaml.Unmarshal(config, &tafconfig)
	if err != nil {
		panic(fmt.Sprintf("unmarshal from %s err: %s\n", "tconfig.yaml", err))
	}

	for _, config := range info.ServerConfig {
		crc := crc32.ChecksumIEEE([]byte(config.CurrentConfigFile))
		name := strings.ToLower(fmt.Sprintf("%s-%s-%s", info.ServerApp, info.ServerName, strconv.FormatUint(uint64(crc), 10)))

		tafconfig.Metadata.Name = name
		tafconfig.ServerConfig.App = info.ServerApp
		tafconfig.ServerConfig.Server = info.ServerName
		tafconfig.ServerConfig.ConfigName = config.CurrentConfigFile
		tafconfig.ServerConfig.ConfigContent = config.CurrentConfigFileContent

		output, err := yaml.Marshal(&tafconfig)
		if err != nil {
			panic(fmt.Sprintf("marshal from %v err: %s\n", tafconfig, err))
		}
		_ = ioutil.WriteFile(fmt.Sprintf("%s/%s.yaml", AppConfigDir, name), output, os.ModePerm)
	}
}

