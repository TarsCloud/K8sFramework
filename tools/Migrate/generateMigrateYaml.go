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


func LoadTafDBServerData(migrateConfigPath string) bool {
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
			enableConfig(val.Info)
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

	tafserver.Spec.Taf.Template = info.ServerOption.ServerTemplate
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
	tafserver.Spec.K8S.NodeSelector["abilityPool"] = make([]string, 0, 0)

	output, err := yaml.Marshal(&tafserver)
	if err != nil {
		panic(fmt.Sprintf("marshal from %v err: %s\n", tafserver, err))
	}
	_ = ioutil.WriteFile(fmt.Sprintf("%s/%s.yaml", AppServerDir, name), output, os.ModePerm)
}

func enableConfig(info *InfoFromTaf) {
	type CreateServerConfigMetadata struct {
		AppServer     string `json:"AppServer" valid:"required,matches-AppServer"`
		ConfigName    string `json:"ConfigName" valid:"required,matches-ConfigName"`
		ConfigContent string `json:"ConfigContent" valid:"required"`
		ConfigMark    string `json:"ConfigMark" valid:"-"`
		PodSeq        int    `json:"PodSeq" valid:"matches-ConfigPodSeq"`
	}

	for _, config := range info.ServerConfig {
		metadata := CreateServerConfigMetadata{
			ConfigName:    config.CurrentConfigFile,
			ConfigContent: config.CurrentConfigFileContent,
			AppServer:     info.ServerApp + "." + info.ServerName,
			ConfigMark:    "",
			PodSeq:        -1,
		}

		const sql = "insert into t_config (f_config_name,f_config_version,f_config_content,f_create_person,f_app_server,f_config_mark,f_pod_seq) value (?,?,?,?,?,?,?) on duplicate key update f_config_version=f_config_version+1"
		_, err := globalK8SDb.Exec(sql, metadata.ConfigName, 10001, metadata.ConfigContent, "migrate", metadata.AppServer, metadata.ConfigMark, metadata.PodSeq)
		if err != nil {
			panic(err.Error())
		}
	}
}

