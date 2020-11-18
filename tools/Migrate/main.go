package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
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


var (
	AppBaseDir 		string
	AppServerDir  	string
	AppReleaseDir 	string
	AppTemplateDir 	string
	AppConfigDir 	string
)
var (
	Init 		bool
	FromK8SDB 	bool
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


var packageDir  = flag.String("packageDir", "", "string: 程序包的路径，服务名必须是AppName.ServerName.(tar.gz/tgz)")
var genTemplate = flag.Bool("genTemplate", false, "bool类型参数：是否生成TTemplate文件")
var initRelease = flag.Bool("initRelease", true, "bool类型参数：是否初始化TRelease文件")

var fromK8SDB = flag.Bool("fromK8SDB", false, "bool类型参数：是否从K8S DB生成TServer+TRelease")
var fromTAFDB = flag.Bool("fromTAFDB", true, "bool类型参数：是否从TAF DB生成TServer+TRelease")
var dumpedTafFile = flag.String("dumpedTafFile", "", "string类型参数：当fromK8SDB和fromTAFDB都是false，则从该文件生成TServer+TRelease")

// 迁移配置文件
const migrateConfigPath = "./Config/migrate.json"
const imageConfigPath = "./Config/image.yaml"

// 迁移模板文件
const ttemplateYamlPath = "./Template/ttemplate.yaml"

const treleaseYamlPath = "./Template/trelease.yaml"
const tconfigYamlPath = "./Template/tconfig.yaml"
const tserverYamlPath = "./Template/tserver.yaml"

const hemlTemplateDir = "./Template/helm"

// 目前一次只支持一个应用
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
	FromK8SDB = *fromK8SDB

	// 创建输出yaml的目标目录
	AppServerDir = fmt.Sprintf("%s/%s", AppBaseDir, "server")
	AppReleaseDir = fmt.Sprintf("%s/%s", AppBaseDir, "release")
	AppTemplateDir = fmt.Sprintf("%s/%s", AppBaseDir, "template")
	AppConfigDir = fmt.Sprintf("%s/%s", AppBaseDir, "config")

	_ = os.MkdirAll(AppServerDir, os.ModePerm)
	_ = os.MkdirAll(AppReleaseDir, os.ModePerm)
	_ = os.Mkdir(AppTemplateDir, os.ModePerm)
	_ = os.MkdirAll(AppConfigDir, os.ModePerm)

	// 加载全局DB配置
	LoadDBConfig()

	// 是否生成TTemplate
	if *genTemplate {
		LoadTafDBTemplateData()
		// yaml模板
		DumpTTempateData()
	}

	// TServer和TRelease的数据源
	if FromK8SDB {
		LoadK8SDBServerData()
	} else {
		if *fromTAFDB {
			LoadTafDBServerData()
		} else {
			LoadDumpFromHZServerData(*dumpedTafFile)
		}
	}

	// 加载全局Image配置
	LoadImageConfig()

	// 上传镜像，生成TRelease和TServer
	BuildPatchImage()

	// 生成helm安装包
	GenHelmInstallPackage()
}

func LoadDBConfig() {
	if bs, err := ioutil.ReadFile(migrateConfigPath); err != nil {
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

func LoadImageConfig() {
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

func GenHelmInstallPackage() {
	sh := fmt.Sprintf("cp -r %s %s", hemlTemplateDir, AppBaseDir)
	if _, err = exec.Command(BashShell, "-c", sh).CombinedOutput(); err != nil {
		panic(fmt.Sprintf("cp helm dir err: %s", err))
	}


	sh = fmt.Sprintf("mv %s %s/helm/templates", AppServerDir, AppBaseDir)
	if _, err = exec.Command(BashShell, "-c", sh).CombinedOutput(); err != nil {
		panic(fmt.Sprintf("mv tserver err: %s", err))
	}

	sh = fmt.Sprintf("mv %s %s/helm/charts/config/templates", AppConfigDir, AppBaseDir)
	if _, err = exec.Command(BashShell, "-c", sh).CombinedOutput(); err != nil {
		panic(fmt.Sprintf("mv tconfig err: %s", err))
	}

	sh = fmt.Sprintf("mv %s %s/helm/charts/release/templates", AppReleaseDir, AppBaseDir)
	if _, err = exec.Command(BashShell, "-c", sh).CombinedOutput(); err != nil {
		panic(fmt.Sprintf("mv trelease err: %s", err))
	}

	fmt.Println("succ. to generate helm dir")
}

