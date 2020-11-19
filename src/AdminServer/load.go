package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func LoadEnv() (*sql.DB, string, *rest.Config, error) {
	fmt.Printf("run tarsadmin. dev: %t, conf: %s\n", *dev, *config)

	var tarsDb *sql.DB
	var k8sNamespace string
	var k8sConfig *rest.Config
	var err error

	if !*dev {
		tarsDb, err = loadTarsDB()
		if err != nil {
			return nil, "", nil, fmt.Errorf("Open TarsDb Error: %s\n", err.Error())
		}
		k8sNamespace, k8sConfig, err = loadK8S()
		if err != nil {
			return nil, "", nil, fmt.Errorf("Load K8S Error: %s\n", err.Error())
		}
	} else {
		tarsDb, err = loadTarsDBDev()
		if err != nil {
			return nil, "", nil, fmt.Errorf("Open TarsDbDev Error: %s\n", err.Error())
		}
		k8sNamespace, k8sConfig, err = loadK8SDev(*config, *namespace)
		if err != nil {
			return nil, "", nil, fmt.Errorf("Load K8SDev Error: %s\n", err.Error())
		}
	}

	return tarsDb, k8sNamespace, k8sConfig, nil
}

func loadTarsDBDev() (*sql.DB, error) {
	/*
	dbHost := "116.63.36.58"
	dbName := "tars_db"
	dbPort := "3306"
	dbPass := "sT3fg5aQm"
	dbUser := "root"
	 */
	dbHost := "172.16.8.229"
	dbName := "tars_db"
	dbPort := "3306"
	dbPass := "8788"
	dbUser := "root"
	dbSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8\n", dbUser, dbPass, dbHost, dbPort, dbName)
	return sql.Open("mysql", dbSourceName)
}

func loadK8SDev(confPath, namespace string) (string, *rest.Config, error) {
	var k8sNamespace = namespace

	k8sConfig, err := clientcmd.BuildConfigFromFlags("", confPath)
	if err != nil {
		return "", nil, fmt.Errorf("Get K8SConfig Value Error , Did You Run Program In K8S ? ")
	}
	return k8sNamespace, k8sConfig, nil
}

func loadTarsDB() (*sql.DB, error) {
	dbHost := os.Getenv("_DB_HOST_")
	dbName := os.Getenv("_DB_NAME_")
	dbPort := os.Getenv("_DB_PORT_")
	dbPass := os.Getenv("_DB_PASSWORD_")
	dbUser := os.Getenv("_DB_USER_")

	dbSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8\n", dbUser, dbPass, dbHost, dbPort, dbName)
	return sql.Open("mysql", dbSourceName)
}

func loadK8S() (string, *rest.Config, error) {
	const namespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

	var k8sNamespace string

	if bs, err := ioutil.ReadFile(namespaceFile); err != nil {
		return "", nil, fmt.Errorf("Get K8SNamespace Value Error , Did You Run Program In K8S ? ")
	} else {
		k8sNamespace = string(bs)
	}

	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		return "", nil, fmt.Errorf("Get K8SConfig Value Error , Did You Run Program In K8S ? ")
	}

	return k8sNamespace, k8sConfig, nil
}
