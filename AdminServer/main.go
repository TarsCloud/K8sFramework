package main

import (
	"base"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io/ioutil"
	k8sInformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	k8sClient "k8s/client"
	k8sWatch "k8s/watch"
	"net/http"
	"os"
	"rpc"
	"tafadmin/shell"
)

var rpcImp base.RPCInterface
var k8sWatchImp base.K8SWatchInterface
var k8sClientImp base.K8SClientInterface

var dev  = flag.Bool("dev", false, "bool类型参数: 本地启动")
var config = flag.String("conf", "/root/.kube/config", "string类型参数：本地启动时，配置文件路径")
var namespace = flag.String("ns", "taf", "string类型参数：本地启动时，namespace")

func loadEnv() (*sql.DB, *rest.Config, *kubernetes.Clientset, k8sInformers.SharedInformerFactory, string, error) {
	fmt.Printf("run tafadmin. dev: %t, conf: %s\n", *dev, *config)

	var tafDb *sql.DB
	var k8sConfig *rest.Config
	var k8sClientSet *kubernetes.Clientset
	var k8sInformerFactor k8sInformers.SharedInformerFactory
	var k8sNamespace string
	var err error

	if !*dev {
		tafDb, err = loadTafDB()
		if err != nil {
			return nil, nil, nil, nil, "", fmt.Errorf("Open TafDb Error: %s\n", err.Error())
		}
		k8sConfig, k8sClientSet, k8sInformerFactor, k8sNamespace, err = loadK8S()
		if err != nil {
			return nil, nil, nil, nil, "", fmt.Errorf("Load K8S Error: %s\n", err.Error())
		}
	} else {
		tafDb, err = loadTafDBDev()
		if err != nil {
			return nil, nil, nil, nil, "", fmt.Errorf("Open TafDbDev Error: %s\n", err.Error())
		}
		k8sConfig, k8sClientSet, k8sInformerFactor, k8sNamespace, err = loadK8SDev(*config, *namespace)
		if err != nil {
			return nil, nil, nil, nil, "", fmt.Errorf("Load K8SDev Error: %s\n", err.Error())
		}
	}

	return tafDb, k8sConfig, k8sClientSet, k8sInformerFactor, k8sNamespace, nil
}

func main() {
	flag.Parse()

	tafDb, k8sConfig, k8sClientSet, k8sInformerFactor, k8sNamespace, err := loadEnv()
	if err != nil {
		fmt.Println(fmt.Sprintf("loadEnv error: %s\n", err))
		return
	}

	k8sWatchImp = k8sWatch.K8SWatchImp{}
	k8sWatchImp.SetTafDb(tafDb)
	k8sWatchImp.SetInformerFactor(k8sInformerFactor)
	k8sWatchImp.StartWatch()

	k8sClientImp = k8sClient.K8SClientImp{}
	k8sClientImp.SetK8SClient(k8sClientSet)
	k8sClientImp.SetWorkNamespace(k8sNamespace)
	k8sClientImp.SetK8SWatchImp(k8sWatchImp)

	rpcImp = rpc.Imp{}
	rpcImp.SetK8SClientImp(k8sClientImp)
	rpcImp.SetK8SWatchImp(k8sWatchImp)
	rpcImp.SetTafDb(tafDb)

	podShellImp := shell.PodShellImp{}
	podShellImp.SetK8SClient(k8sClientSet)
	podShellImp.SetK8SConfig(k8sConfig)
	podShellImp.SetK8SNamespace(k8sNamespace)
	podShellImp.SetK8SWatchImp(k8sWatchImp)

	router := mux.NewRouter()

	bindAdminHandler(rpcImp, router)
	bindAgentHandler(k8sWatchImp, router)
	bindShellHandler(&podShellImp, router)

	if err = http.ListenAndServe(":80", router); err != nil {
		fmt.Println(err)
	}
}

func loadTafDBDev() (*sql.DB, error) {
	dbHost := "172.16.8.229"
	dbName := "taf_db"
	dbPort := "3306"
	dbPass := "8788"
	dbUser := "root"
	dbSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8\n", dbUser, dbPass, dbHost, dbPort, dbName)
	return sql.Open("mysql", dbSourceName)
}

func loadK8SDev(confPath, namespace string) (*rest.Config, *kubernetes.Clientset, k8sInformers.SharedInformerFactory, string, error) {
	var k8sNamespace = namespace

	k8sConfig, err := clientcmd.BuildConfigFromFlags("", confPath)
	if err != nil {
		return nil, nil, nil, "", errors.New("Get K8SConfig Value Error , Did You Run Program In K8S ? ")
	}

	k8sClientSet, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, nil, nil, "", errors.New("Get K8SClientSet Value Error , Did You Run Program In K8S ? ")
	}

	k8sInformerFactory := k8sInformers.NewSharedInformerFactoryWithOptions(k8sClientSet, 0, k8sInformers.WithNamespace(k8sNamespace))

	return k8sConfig, k8sClientSet, k8sInformerFactory, k8sNamespace, nil
}

func loadTafDB() (*sql.DB, error) {
	dbHost := os.Getenv("_DB_HOST_")
	dbName := os.Getenv("_DB_NAME_")
	dbPort := os.Getenv("_DB_PORT_")
	dbPass := os.Getenv("_DB_PASSWORD_")
	dbUser := os.Getenv("_DB_USER_")

	dbSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8\n", dbUser, dbPass, dbHost, dbPort, dbName)
	return sql.Open("mysql", dbSourceName)
}

func loadK8S() (*rest.Config, *kubernetes.Clientset, k8sInformers.SharedInformerFactory, string, error) {
	const namespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

	var k8sNamespace string

	if bs, err := ioutil.ReadFile(namespaceFile); err != nil {
		return nil, nil, nil, "", errors.New("Get K8SNamespace Value Error , Did You Run Program In K8S ? ")
	} else {
		k8sNamespace = string(bs)
	}

	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, nil, nil, "", errors.New("Get K8SConfig Value Error , Did You Run Program In K8S ? ")
	}

	k8sClientSet, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, nil, nil, "", errors.New("Get K8SClientSet Value Error , Did You Run Program In K8S ? ")
	}

	k8sInformerFactory := k8sInformers.NewSharedInformerFactoryWithOptions(k8sClientSet, 0, k8sInformers.WithNamespace(k8sNamespace))

	return k8sConfig, k8sClientSet, k8sInformerFactory, k8sNamespace, nil
}

func bindAdminHandler(rpcImp base.RPCInterface, router *mux.Router) {
	router.HandleFunc("/admin", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		writer.Header().Add("Connection", "keep-alive")
		result := rpcImp.Handler(request.Body)
		_, _ = writer.Write(result)
	})
}

func bindShellHandler(shellImp base.ShellInterface, router *mux.Router) {
	router.HandleFunc("/shell", func(writer http.ResponseWriter, request *http.Request) {
		pty, err := shell.NewTerminalSession(writer, request, nil)
		if err != nil {
			http.Error(writer,
				fmt.Sprintf("%s terminal get failed: %s\n", request.URL.Query().Encode(), err),
				http.StatusInternalServerError)
			return
		}
		defer func() {
			_ = pty.Close()
		}()

		shellImp.Handler(request.URL.Query(), pty)
	})
}

func bindAgentHandler(watchImp base.K8SWatchInterface, router *mux.Router)  {
	agentRouter := router.PathPrefix("/agent").Subrouter()
	agentRouter.HandleFunc("/{node}/cpu", func(writer http.ResponseWriter, request *http.Request) {
		daemonProxy("cpu", watchImp, writer, request)
	})
	agentRouter.HandleFunc("/{node}/memory", func(writer http.ResponseWriter, request *http.Request) {
		daemonProxy("memory", watchImp, writer, request)
	})
	agentRouter.HandleFunc("/{node}/disk", func(writer http.ResponseWriter, request *http.Request) {
		daemonProxy("disk", watchImp, writer, request)
	})
	agentRouter.HandleFunc("/{node}/host", func(writer http.ResponseWriter, request *http.Request) {
		daemonProxy("host", watchImp, writer, request)
	})
	agentRouter.HandleFunc("/{node}/net", func(writer http.ResponseWriter, request *http.Request) {
		daemonProxy("net", watchImp, writer, request)
	})
	agentRouter.HandleFunc("/{node}/port", func(writer http.ResponseWriter, request *http.Request) {
		daemonProxy("port",  watchImp, writer, request)
	})
}

func daemonProxy(path string, watchImp base.K8SWatchInterface, writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	node := vars["node"]

	daemon := watchImp.GetDaemonSetPodByName(node)
	if daemon == nil {
		http.Error(writer, "host ip error\n", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("http://%s:%d/%s", daemon.HostIP, daemon.HostPort, path)

	query := request.URL.Query().Encode()
	if query != "" {
		url = fmt.Sprintf("%s?host=%s&%s", url, daemon.HostIP, query)
	}

	rsp, err := http.Get(url)
	if err != nil {
		http.Error(writer, fmt.Sprintf("proxy daemon: %s:%s, error: %s\n",
			daemon.HostIP, daemon.HostPort, err), http.StatusInternalServerError)
		return
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("read rsp from daemon: %s:%s, error: %s\n",
			daemon.HostIP, daemon.HostPort, err), http.StatusInternalServerError)
		return
	}
	_, err = fmt.Fprintf(writer, string(body))
	if err != nil {
		http.Error(writer, fmt.Sprintf("send rsp from web: %s:%s, error: %s\n",
			daemon.HostIP, daemon.HostPort, err), http.StatusInternalServerError)
		return
	}
}
