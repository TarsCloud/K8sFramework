package main

import (
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"tafadmin/handler"
)

var dev  = flag.Bool("dev", false, "bool类型参数: 本地启动")
var config = flag.String("config", "/root/.kube/config", "string类型参数：本地启动时，配置文件路径")
var namespace = flag.String("namespace", "tao", "string类型参数：本地启动时，K8S命名空间")
var port = flag.Int("port", 80, "int类型参数：本地启动时，监听资源端口")

func main() {
	flag.Parse()

	tafDb, k8sNamespace, k8sConfig, err := LoadEnv()
	if err != nil {
		fmt.Println(fmt.Sprintf("LoadEnv error: %s\n", err))
		return
	}

	if err := handler.StartServer(k8sNamespace, k8sConfig, tafDb, *port); err != nil {
		fmt.Println(fmt.Sprintf("StartServer error: %s\n", err))
		return
	}
}


