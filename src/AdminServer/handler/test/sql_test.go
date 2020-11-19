package test

import (
	"database/sql"
	"fmt"
	"github.com/go-openapi/strfmt"
	_ "github.com/go-sql-driver/mysql"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"tarsadmin/handler/k8s"
	"tarsadmin/handler/mysql"
	"tarsadmin/openapi/models"
	"tarsadmin/openapi/restapi/operations/applications"
	"tarsadmin/openapi/restapi/operations/approval"
	"tarsadmin/openapi/restapi/operations/config"
	"tarsadmin/openapi/restapi/operations/deploy"
	"testing"
)

func loadTarsDBDev() (*sql.DB, error) {
	dbHost := "172.16.8.229"
	dbName := "tars_db"
	dbPort := "3306"
	dbPass := "8788"
	dbUser := "root"
	dbSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8\n", dbUser, dbPass, dbHost, dbPort, dbName)
	return sql.Open("mysql", dbSourceName)
}
func loadK8SDev() (string, *rest.Config, error) {
	k8sConfig, err := clientcmd.BuildConfigFromFlags("", "/Users/zoujiamin/.kube/config")
	if err != nil {
		return "", nil, fmt.Errorf("Get K8SConfig Value Error , Did You Run Program In K8S ? ")
	}
	return "tars", k8sConfig, nil
}

func ConstructSelectParams() (models.SelectRequestFilter, models.SelectRequestLimiter, models.SelectRequestOrder) {
	var filter models.SelectRequestFilter
	filter.Eq = make(models.MapInterface)
	filter.Ne = make(models.MapInterface)
	filter.Like = make(models.MapString)

	var limiter models.SelectRequestLimiter
	offset := int32(0)
	limiter.Offset = &offset

	var order models.SelectRequestOrder
	order = make([]*models.SelectRequestOrderElem, 0, 2)

	return filter, limiter, order
}

func TestSelectAppHandler_Handle(t *testing.T) {
	filter, limiter, _ := ConstructSelectParams()

	filter.Eq = make(models.MapInterface)
	filter.Eq["CreatePerson"] = "admin"
	fb, _ := filter.MarshalBinary()
	fs := string(fb)

	lb, _ := limiter.MarshalBinary()
	ls := string(lb)

	var params applications.SelectAppParams
	params.Filter = &fs
	params.Limiter = &ls

	handler := k8s.SelectAppHandler{}
	response := handler.Handle(params)
	fmt.Println(response)
}

func TestUpdateAppHandler_Handle(t *testing.T) {
	var params applications.UpdateAppParams

	name := "tars"
	params.Params.Metadata = &applications.UpdateAppParamsBodyMetadata{AppName: &name}
	params.Params.Target = &applications.UpdateAppParamsBodyTarget{BusinessName: "FRAMEWORK", AppMark: "jaminTest"}

	handler := k8s.UpdateAppHandler{}
	response := handler.Handle(params)
	fmt.Println(response)
}

func TestSelectServerConfigHandler_Handle(t *testing.T) {
	filter, limiter, _ := ConstructSelectParams()
	filter.Eq["AppServer"] = "Semantics.AnalyserServer"
	filter.Eq["PodSeq"] = -1

	fb, _ := filter.MarshalBinary()
	fs := string(fb)

	lb, _ := limiter.MarshalBinary()
	ls := string(lb)

	var params config.SelectServerConfigParams
	params.Filter = &fs
	params.Limiter = &ls

	handler := k8s.SelectServerConfigHandler{}
	response := handler.Handle(params)
	fmt.Println(response)
}

func CreateDeployMeta() *models.DeployMeta {
	serverApp := "Semantics"
	serverName := "AnalyserServer"

	notStacked := true
	serverK8S := models.ServerK8S{
		Image: "registry.cn-hangzhou.aliyuncs.com/dtool/semantics.analyserserver:a1600911633405313000",
		Version: "10002",
		NotStacked: &notStacked,
		HostIpc: false,
		HostNetwork: false,
		Replicas: 1,
		HostPort: []*models.HostPortElem {
			{
				NameRef: "AnalyserObj",
				Port: 10070,
			},
		},
		NodeSelector: &models.NodeSelector{
			AbilityPool: &models.NodeSelectorElem{
				Value: make([]string, 0, 1),
			},
		},
	}

	asyncThread := int32(3)
	serverImportant := int32(10)
	serverSubType := "tars"

	serverOption := models.ServerOption{
		AsyncThread: &asyncThread,
		ServerImportant: &serverImportant,
		ServerProfile: "",
		ServerSubType: &serverSubType,
		ServerTemplate: "tars.cpp",
	}

	isTars, isTcp := true, true
	serverServant := models.MapServant{
		"AnalyserObj": {
			Capacity: 10000,
			Connections: 10000,
			IsTars: &isTars,
			IsTCP: &isTcp,
			Name: "AnalyserObj",
			Port: 10000,
			Threads: 1,
			Timeout: 60000,
		},
	}

	return &models.DeployMeta{
		RequestPerson: "jaminzou",
		RequestTime: strfmt.NewDateTime(),
		ServerApp: &serverApp,
		ServerK8S: &serverK8S,
		ServerMark: "unit test",
		ServerName: &serverName,
		ServerOption: &serverOption,
		ServerServant: serverServant,
	}
}

func TestCreateDeployHandler_Handle(t *testing.T) {
	var params deploy.CreateDeployParams
	params.Params.Metadata = CreateDeployMeta()

	handler := k8s.CreateDeployHandler{}
	response := handler.Handle(params)
	fmt.Println(response)
}

func TestUpdateDeployHandler_Handle(t *testing.T) {
	deployMeta := CreateDeployMeta()

	deployId := "base-jceproxyserver-y2enbcrjrs-s5jqn"

	var params deploy.UpdateDeployParams
	params.Params.Metadata = &deploy.UpdateDeployParamsBodyMetadata{
		DeployID: &deployId,
	}
	isTars, isTcp := true, true
	deployMeta.ServerServant["TestUpdateObj"] = models.ServerServantElem{
			Capacity: 10000,
			Connections: 10000,
			Name: "TestUpdateObj",
			Port: 10001,
			Threads: 1,
			Timeout: 60000,
			IsTCP: &isTcp,
			IsTars: &isTars,
	}
	params.Params.Target = &deploy.UpdateDeployParamsBodyTarget{
		ServerServant: deployMeta.ServerServant,
		ServerOption: deployMeta.ServerOption,
		ServerK8S: deployMeta.ServerK8S,
	}

	handler := k8s.UpdateDeployHandler{}
	response := handler.Handle(params)
	fmt.Println(response)
}


func TestCreateApprovalHandler_Handle(t *testing.T) {
	var params approval.CreateApprovalParams
	deployId := "base-jceproxyserver-y2enbcrjrs-s5jqn"
	params.Params.Metadata = &approval.CreateApprovalParamsBodyMetadata{
		DeployID: &deployId,
		ApprovalMark: "unit test",
		ApprovalResult: true,
	}

	handler := k8s.CreateApprovalHandler{}
	response := handler.Handle(params)
	fmt.Println(response)
}

func TestSelectDeployHandler_Handle(t *testing.T) {
	filter, limiter, _ := ConstructSelectParams()
	filter.Like["ServerApp"] = "%Semantics%"
	filter.Like["ServerName"] = "%AnalyserServer%"

	fb, _ := filter.MarshalBinary()
	fs := string(fb)

	lb, _ := limiter.MarshalBinary()
	ls := string(lb)

	var params deploy.SelectDeployParams
	params.Filter = &fs
	params.Limiter = &ls

	handler := k8s.SelectDeployHandler{}
	response := handler.Handle(params)
	fmt.Println(response)
}

func init()  {
	var err error
	mysql.TarsDb, err = loadTarsDBDev()
	if err != nil {
		fmt.Println(fmt.Sprintf("load tars_db error: %v", err))
	}

	k8sNamespace, k8sConfig, err := loadK8SDev()
	if err != nil {
		fmt.Println(fmt.Sprintf("load k8s error: %v", err))
	}

	if k8s.K8sOption, k8s.K8sWatcher, err = k8s.StartWatcher(k8sNamespace, k8sConfig); err != nil {
		fmt.Println(fmt.Sprintf("start watcher error: %v", err))
	}
}
