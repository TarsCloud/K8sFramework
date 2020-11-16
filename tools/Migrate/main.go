package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type DBConfig struct {
	DbName string `json:"DBName",valid:"required,matches-DBName"`
	DbHost string `json:"DBHost",valid:"required,matches-DBHost"`
	DbUser string `json:"DBUser",valid:"required,matches-DBUser"`
	DbPass string `json:"DBPass",valid:"required,matches-DBPass"`
	DbPort string `json:"DBPort",valid:"required,matches-DBPort"`
}

type ConfigDBContent struct {
	TafDBConfig        DBConfig        `json:"TafDBConfig",valid:"required,matches-DBConfig"`
	K8sDBConfig        DBConfig        `json:"K8SDBConfig",valid:"required,matches-DBConfig"`
}

var globalK8SDb *sql.DB
var globalTafDb *sql.DB
var globalDBConfig ConfigDBContent

func loadDBConfig(dbConfigPath string) {
	if bs, err := ioutil.ReadFile(dbConfigPath); err != nil {
		fmt.Print("Open Config File Error ? ")
		return
	} else {
		err := json.Unmarshal(bs, &globalDBConfig)
		if err != nil {
			fmt.Print("Config File Had Bad Format")
			return
		}
	}

	tafSqlUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", globalDBConfig.TafDBConfig.DbUser, globalDBConfig.TafDBConfig.DbPass, globalDBConfig.TafDBConfig.DbHost, globalDBConfig.TafDBConfig.DbPort, globalDBConfig.TafDBConfig.DbName)

	k8sSqlUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", globalDBConfig.K8sDBConfig.DbUser, globalDBConfig.K8sDBConfig.DbPass, globalDBConfig.K8sDBConfig.DbHost, globalDBConfig.K8sDBConfig.DbPort, globalDBConfig.K8sDBConfig.DbName)

	var err error

	if globalTafDb, err = sql.Open("mysql", tafSqlUrl); err != nil || globalTafDb == nil {
		fmt.Println("Open Taf MySQl Error , Did You Set Right TafDBConfig Params ? ")
		return
	}

	if globalK8SDb, err = sql.Open("mysql", k8sSqlUrl); err != nil || globalK8SDb == nil {
		fmt.Println("Open TafK8S MySQl Error , Did You Set Right K8SDBConfig Params ? ")
		return
	}
}

var (
	AppBaseDir 		string
	AppServerDir  	string
	AppReleaseDir 	string
	AppTemplateDir 	string
)
var (
	Init 		bool
	Semantics 	bool
)
var (
	DockerRegistryUrl 		string
	DockerRegistryUser 		string
	DockerRegistryPassword 	string
)
var (
	Namespace string
)
var imageBaseMap = make(map[string]string)

func loadImageConfig(imageConfigPath string) {
	template, err := ioutil.ReadFile(imageConfigPath)
	if err != nil {
		panic(fmt.Sprintf("read from %s err: %s\n", "configmap.yaml", err))
	}
	configmap := &ConfigMap{}
	err = yaml.Unmarshal(template, &configmap)
	if err != nil {
		panic(fmt.Sprintf("unmarshal from %s err: %s\n", "configmap.yaml", err))
	}

	// 基础镜像
	imageBaseMap["taf_cpp"] = configmap.Data.CppImageBase
	imageBaseMap["taf_node10"] = configmap.Data.Node10ImageBase

	Namespace = configmap.Data.Namespace
	DockerRegistryUrl = configmap.Data.DockerRegistryURL
	DockerRegistryUser = configmap.Data.DockerRegistryUser
	DockerRegistryPassword = configmap.Data.DockerRegistryPassword
}

var packageDir  = flag.String("packageDir", "", "string: 程序包的路径，服务名必须是AppName.ServerName.(tar.gz/tgz)")
var genTemplate = flag.Bool("genTemplate", false, "bool类型参数：是否生成TTemplate文件")
var initRelease = flag.Bool("initRelease", true, "bool类型参数：是否初始化TRelease文件")

var fromK8SDB = flag.Bool("fromK8SDB", false, "bool类型参数：是否从K8S DB生成TServer+TRelease")
var fromTAFDB = flag.Bool("fromTAFDB", true, "bool类型参数：是否从TAF DB生成TServer+TRelease")
var dumpedTafFile = flag.String("dumpedTafFile", "", "string类型参数：当fromK8SDB和fromTAFDB都是false，则从该文件生成TServer+TRelease")


func main()  {
	flag.Parse()

	// 输入参数
	AppBaseDir = *packageDir
	if AppBaseDir == "" {
		fmt.Println("Missing packageDir In Command Line")
		return
	}
	// 是否初始化生成TRelease
	Init = *initRelease
	// 是否从k8s_db拉取数据
	Semantics = *fromK8SDB

	// 创建输出yaml的目标目录
	AppServerDir = fmt.Sprintf("%s/%s", AppBaseDir, "server")
	_ = os.MkdirAll(AppServerDir, os.ModePerm)
	AppReleaseDir = fmt.Sprintf("%s/%s", AppBaseDir, "release")
	_ = os.MkdirAll(AppReleaseDir, os.ModePerm)
	AppTemplateDir = fmt.Sprintf("%s/%s", AppBaseDir, "template")
	_ = os.Mkdir(AppTemplateDir, os.ModePerm)

	// 加载全局DB配置
	migrateConfigPath := "./Config/migrate.json"
	loadDBConfig(migrateConfigPath)

	// 是否生成TTemplate
	if *genTemplate {
		LoadTafDBTemplateData()
		// yaml模板
		ttemplateYamlPath := "./Template/ttemplate.yaml"
		DumpTafDBTempateData(ttemplateYamlPath)
	}

	// TServer和TRelease的数据源
	if Semantics {
		LoadK8SDBServerData()
	} else {
		if *fromTAFDB {
			LoadTafDBServerData(migrateConfigPath)
		} else {
			LoadDumpFromHZServerData(*dumpedTafFile)
		}
	}

	// 加载全局Image配置
	imageConfigPath := "./Config/image.yaml"
	loadImageConfig(imageConfigPath)

	// 上传镜像，生成TRelease和TServer
	BuildPatchImage()
}