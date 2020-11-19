package handler

import (
	"database/sql"
	"fmt"
	"github.com/go-openapi/loads"
	"k8s.io/client-go/rest"
	"tarsadmin/handler/compatible"
	"tarsadmin/handler/k8s"
	"tarsadmin/handler/mysql"
	"tarsadmin/openapi/restapi"
	"tarsadmin/openapi/restapi/operations"
)

func loadSwagger() (*restapi.Server, *operations.TarsadminOpenapiAPI) {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		fmt.Println(fmt.Sprintf("swagger loads error: %s\n", err))
	}

	api := operations.NewTarsadminOpenapiAPI(swaggerSpec)
	server := restapi.NewServer(api)

	return server, api
}

func StartServer(namespace string, config *rest.Config, tarsDb *sql.DB, port int) error {
	var err error

	mysql.TarsDb = tarsDb

	// start common watcher
	if k8s.K8sOption, k8s.K8sWatcher, err = k8s.StartWatcher(namespace, config); err != nil {
		return fmt.Errorf("start swatcher err: %v", err)
	}

	// start node watcher: 暂时无TNode的crd资源，故复用旧版Admin的逻辑，在内存中构建缓存操作
	compatible.StartNodeWatch()

	// start restful server
	server, api := loadSwagger()
	defer server.Shutdown()

	server.Host = "0.0.0.0"
	server.Port = port

	adminHandler := tarsAdminHandler{}
	server.SetHandler(adminHandler.ConfigureAPI(api))
	return server.Serve()
}

