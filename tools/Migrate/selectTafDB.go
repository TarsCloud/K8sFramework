package main

import (
	"errors"
	"fmt"
	"strings"
)

type ServerOption struct {
	ServerImportant       int    `json:"ServerImportant" valid:"required,matches-ServerImportant"`
	StartScript           string `json:"StartScript"     valid:"matches-ServerStartScript"`
	StopScript            string `json:"StopScript"      valid:"matches-ServerStopScript"`
	MonitorScript         string `json:"MonitorScript"   valid:"matches-ServerMonitorScript"`
	AsyncThread           int    `json:"AsyncThread"     valid:"required,matches-ServerAsyncThread"`
	ServerTemplate        string `json:"ServerTemplate"  valid:"required,matches-TemplateName"`
	ServerProfile         string `json:"ServerProfile"   valid:"-"`
	RemoteLogEnable       bool   `json:"RemoteLogEnable"`
	RemoteLogReserveTime  string `json:"RemoteLogReserveTime"`
	RemoteLogCompressTime string `json:"RemoteLogCompressTime"`
}

type ConfigHistory struct {
}

type ServerConfigFile struct {
	CurrentConfigFile             string
	CurrentConfigFileContent      string
	CurrentConfigFileUpdatePerson string
	CurrentConfigFileUpdateTime   string
	CurrentConfigFileSelectReason string
	ConfigHistory                 []ConfigHistory
}

type ServerAdapter struct {
	Name        string `json:"Name" valid:"required,matches-ServantName"`
	Port        int    `json:"Port" valid:"required,matches-ServantPort"`
	Threads     int    `json:"Threads" valid:"required,matches-ServantThreads"`
	Connections int    `json:"Connections" valid:"required,matches-ServantConnections"`
	Capacity    int    `json:"Capacity" valid:"required,matches-ServantCapacity"`
	Timeout     int    `json:"Timeout" valid:"required,matches-ServantTimeout"`
	IsTars       bool   `json:"IsTars" valid:"-"`
}

type InfoFromTars struct {
	ServerApp      string                   `json:"ServerApp"`
	ServerName     string                   `json:"ServerName"`
	ServerOption   *ServerOption            `json:"ServerOption"`
	ServerAdapters map[string]ServerAdapter `json:"ServerAdapter"`
	ServerConfig   []ServerConfigFile       `json:"ServerConfigFile"`
}

func LoadOption(serverApp string, serverName string) (*ServerOption, error) {
	var loadSql = "SELECT template_name,server_type,profile,async_thread_num ,server_important_type ,remote_log_reserve_time,remote_log_compress_time,remote_log_type FROM t_server_conf where application=? and server_name=?"
	rowsFromTAdapterConf, err := globalTarsDb.Query(loadSql, serverApp, serverName)
	defer func() {
		if rowsFromTAdapterConf != nil {
			_ = rowsFromTAdapterConf.Close()
		}
	}()

	if err != nil {
		fmt.Print(err.Error())
		panic(err.Error())
	}

	var templateName string
	var server_type string
	var profile string
	var asyncThread int
	var server_important_type int
	var remote_log_reserve_time string
	var remote_log_compress_time string
	var remote_log_type bool

	getOption := false

	for rowsFromTAdapterConf.Next() {
		err := rowsFromTAdapterConf.Scan(&templateName, &server_type, &profile, &asyncThread, &server_important_type, &remote_log_reserve_time, &remote_log_compress_time, &remote_log_type)
		if err != nil {
			return nil, err
		}
		getOption = true
	}

	if !getOption {
		return nil, errors.New("no server option found in TarsDB")
	}

	return &ServerOption{
		StartScript:           "",
		StopScript:            "",
		MonitorScript:         "",
		AsyncThread:           asyncThread,
		ServerTemplate:        templateName,
		ServerProfile:         profile,
		ServerImportant:       server_important_type,
		RemoteLogReserveTime:  remote_log_reserve_time,
		RemoteLogCompressTime: remote_log_compress_time,
		RemoteLogEnable:       remote_log_type,
	}, nil
}

func LoadAdapters(serverApp string, serverName string) map[string]ServerAdapter {
	var loadSql = "SELECT distinct servant,queuecap,queuetimeout ,protocol,thread_num,max_connections FROM t_adapter_conf where application=? and server_name=?"
	rowsFromTAdapterConf, err := globalTarsDb.Query(loadSql, serverApp, serverName)
	defer func() {
		if rowsFromTAdapterConf != nil {
			_ = rowsFromTAdapterConf.Close()
		}
	}()

	if err != nil {
		fmt.Print(err.Error())
		panic(err.Error())
	}

	basePort := 9999
	var servant string
	var queuecap int
	var queuetimeout int
	var thread_num int
	var max_connections int
	var protocol string

	var adapters = make(map[string]ServerAdapter, 0)
	for rowsFromTAdapterConf.Next() {
		err := rowsFromTAdapterConf.Scan(&servant, &queuecap, &queuetimeout, &protocol, &thread_num, &max_connections)
		if err != nil {
			panic(err.Error())
		}
		basePort += 1
		adapters[strings.ToLower(servant)] = ServerAdapter{
			Name:        servant,
			Port:        basePort,
			Threads:     thread_num,
			Connections: max_connections,
			Capacity:    queuecap,
			Timeout:     queuetimeout,
			IsTars:       protocol == "tars",
		}
	}
	return adapters
}

func LoadConfigFileWithHistory(configId int64) ServerConfigFile {

	var loadSql = "select a.filename ,b.content,b.posttime from t_config_files a left join t_config_history_files b on a.id=b.configid where configid=? order by b.id desc limit 3"
	rows, err := globalTarsDb.Query(loadSql, configId)
	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()

	if err != nil {
		fmt.Print(err.Error())
		panic(err.Error())
	}

	var serverConfigFile ServerConfigFile

	var configFileName string
	var configFileContent string
	var configPostTime string

	firstColumns := true
	for rows.Next() {
		err = rows.Scan(&configFileName, &configFileContent, &configPostTime)
		if err != nil {
			fmt.Println(err.Error())
			panic(err.Error())
		}

		if firstColumns {
			serverConfigFile.CurrentConfigFile = configFileName
			serverConfigFile.CurrentConfigFileContent = configFileContent
			serverConfigFile.CurrentConfigFileUpdateTime = configPostTime
			firstColumns = false
			continue
		}
		if serverConfigFile.ConfigHistory == nil {
			serverConfigFile.ConfigHistory = make([]ConfigHistory, 0)
		}
		serverConfigFile.ConfigHistory = append(serverConfigFile.ConfigHistory, ConfigHistory{})
	}
	return serverConfigFile
}

func LoadConfigFile(serverApp string, serverName string) []ServerConfigFile {

	var loadSql = "select id from t_config_files where server_name=? and host =''"
	rows, err := globalTarsDb.Query(loadSql, serverApp+"."+serverName)
	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()

	if err != nil {
		fmt.Print(err.Error())
		panic(err.Error())
	}

	serverConfigFiles := make([]ServerConfigFile, 0)

	var configId int64
	for rows.Next() {
		err = rows.Scan(&configId)
		if err != nil {
			panic(err.Error())
		}
		serverConfigFiles = append(serverConfigFiles, LoadConfigFileWithHistory(configId))
	}
	return serverConfigFiles
}

func LoadInfoFromTars(serverApp string, serverName string) (*InfoFromTars, error) {
	var info = new(InfoFromTars)
	info.ServerApp = serverApp
	info.ServerName = serverName

	var err error
	info.ServerOption, err = LoadOption(serverApp, serverName)
	if err != nil {
		return nil, err
	}

	info.ServerAdapters = LoadAdapters(serverApp, serverName)
	info.ServerConfig = LoadConfigFile(serverApp, serverName)

	return info, nil
}

