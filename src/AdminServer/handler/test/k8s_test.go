package test

import (
	"encoding/json"
	"fmt"
	"tarsadmin/handler/compatible"
	"tarsadmin/handler/k8s"
	"tarsadmin/handler/mysql"
	"tarsadmin/openapi/models"
	"tarsadmin/openapi/restapi/operations/affinity"
	"tarsadmin/openapi/restapi/operations/release"
	"tarsadmin/openapi/restapi/operations/server"
	"tarsadmin/openapi/restapi/operations/server_k8s"
	"tarsadmin/openapi/restapi/operations/server_option"
	"tarsadmin/openapi/restapi/operations/server_pod"
	"tarsadmin/openapi/restapi/operations/server_servant"
	"testing"
)

func TestSelectServerHandler_Handle(t *testing.T) {
	filter, limiter, _ := ConstructSelectParams()
	filter.Like["ServerApp"] = "%Semantics%"
	filter.Like["ServerName"] = "%AnalyserServer%"

	fb, _ := filter.MarshalBinary()
	fs := string(fb)

	lb, _ := limiter.MarshalBinary()
	ls := string(lb)

	var params server.SelectServerParams
	params.Filter = &fs
	params.Limiter = &ls

	handler := k8s.SelectServerHandler{}
	response := handler.Handle(params)
	fmt.Println(response)
}

func TestSelectServicePoolHandler_Handle(t *testing.T) {
	filter, limiter, _ := ConstructSelectParams()
	filter.Eq["ServerApp"] = "Semantics"
	filter.Eq["ServerName"] = "AnalyserServer"

	fb, _ := filter.MarshalBinary()
	fs := string(fb)

	lb, _ := limiter.MarshalBinary()
	ls := string(lb)

	var params release.SelectServicePoolParams
	params.Filter = &fs
	params.Limiter = &ls

	handler := k8s.SelectServicePoolHandler{}
	response := handler.Handle(params)
	fmt.Println(response)
}

func TestCreateServicePoolHandler_Handle(t *testing.T) {
	var params release.CreateServicePoolParams
	serverID := "Semantics-AnalyserServer"
	serverType := "tars.cpp"
	serviceImage := "registry.cn-hangzhou.aliyuncs.com/dtool/semantics.analyserserver:a1600911633405313000"

	params.Params.Metadata = &release.CreateServicePoolParamsBodyMetadata{
		ActivePerson: "jaminzou",
		ActiveReason: "unit test",
		ServerID:     &serverID,
		ServerType:   &serverType,
		ServiceImage: &serviceImage,
		ServiceMark:  "unit test",
	}

	handler := k8s.CreateServicePoolHandler{}
	response := handler.Handle(params)
	fmt.Println(response)
}

func TestDoEnableServiceHandler_Handle(t *testing.T) {
	var params release.DoEnableServiceParams

	replicas := int32(3)
	serverID := "Semantics-AnalyserServer"
	serviceImage := "registry.cn-hangzhou.aliyuncs.com/dtool/semantics.analyserserver:a1600911633405313000"
	serviceVersion := "10002"
	serviceID := fmt.Sprintf("%s|%s", serviceImage, serviceVersion)

	params.Params.Metadata = &release.DoEnableServiceParamsBodyMetadata{
		EnableMark: "unit test",
		Replicas:   &replicas,
		ServerID:   &serverID,
		ServiceID:  &serviceID,
	}

	handler := k8s.DoEnableServiceHandler{}
	response := handler.Handle(params)
	fmt.Println(response)
}

func TestSelectPodAliveHandler_Handle(t *testing.T) {
	filter, limiter, order := ConstructSelectParams()
	filter.Eq["ServerApp"] = "Semantics"
	filter.Eq["ServerName"] = "AnalyserServer"

	fb, _ := filter.MarshalBinary()
	fs := string(fb)

	lb, _ := limiter.MarshalBinary()
	ls := string(lb)

	order = []*models.SelectRequestOrderElem{
		{
			Column: "PodName",
			Order:  "asc",
		},
	}
	ob, _ := json.Marshal(order)
	os := string(ob)

	var params server_pod.SelectPodAliveParams
	params.Filter = &fs
	params.Limiter = &ls
	params.Order = &os

	handler := k8s.SelectPodAliveHandler{}
	response := handler.Handle(params)
	fmt.Println(response)
}

func TestUpdateServerK8SHandler_Handle(t *testing.T) {
	depolyMeta := CreateDeployMeta()
	depolyMeta.ServerK8S.NodeSelector.AbilityPool = &models.NodeSelectorElem{
		Value: make([]string, 0, 1),
	}
	/*
		depolyMeta.ServerK8S.NodeSelector.AbilityPool = nil
		depolyMeta.ServerK8S.NodeSelector.NodeBind = &models.NodeSelectorElem{
			Value: []string{"kube.node117"},
		}
	*/
	notStacked := false
	depolyMeta.ServerK8S.NotStacked = &notStacked
	depolyMeta.ServerK8S.Replicas = 2

	var params server_k8s.UpdateK8SParams
	serverId := "Semantics-AnalyserServer"
	params.Params.Metadata = &server_k8s.UpdateK8SParamsBodyMetadata{ServerID: &serverId}
	params.Params.Target = depolyMeta.ServerK8S

	handler := k8s.UpdateServerK8SHandler{}
	response := handler.Handle(params)
	fmt.Println(response)
}

func TestUpdateServerOptionHandler_Handle(t *testing.T) {
	depolyMeta := CreateDeployMeta()

	serverImportant := int32(5)
	serverProfile := "<tars>\n\n</tars>"
	depolyMeta.ServerOption.ServerTemplate = "tars.default"
	depolyMeta.ServerOption.ServerImportant = &serverImportant
	depolyMeta.ServerOption.ServerProfile = serverProfile

	var params server_option.UpdateServerOptionParams
	serverId := "Semantics-AnalyserServer"
	params.Params.Metadata = &server_option.UpdateServerOptionParamsBodyMetadata{ServerID: &serverId}
	params.Params.Target = depolyMeta.ServerOption

	handler := k8s.UpdateServerOptionHandler{}
	response := handler.Handle(params)
	fmt.Println(response)
}

func TestSelectServerAdapterHandler_Handle(t *testing.T) {
	filter, limiter, _ := ConstructSelectParams()
	filter.Eq["ServerApp"] = "Semantics"
	filter.Eq["ServerName"] = "AnalyserServer"

	fb, _ := filter.MarshalBinary()
	fs := string(fb)

	lb, _ := limiter.MarshalBinary()
	ls := string(lb)

	var params server_servant.SelectServerAdapterParams
	params.Filter = &fs
	params.Limiter = &ls

	handler := k8s.SelectServerAdapterHandler{}
	response := handler.Handle(params)
	fmt.Println(response)
}

func TestDeleteServerAdapterHandler_Handle(t *testing.T) {
	var params server_servant.DeleteServerAdapterParams

	adapterId := "Semantics-AnalyserServer.TestUpdateObj"
	params.Params.Metadata = &server_servant.DeleteServerAdapterParamsBodyMetadata{
		AdapterID: &adapterId,
	}

	handler := k8s.DeleteServerAdapterHandler{}
	response := handler.Handle(params)
	fmt.Println(response)
}

func TestDoListAffinityGroupByNodeHandler_Handle(t *testing.T) {
	var params affinity.DoListAffinityGroupByNodeParams

	nodeName := []string{
		"kube.node118", "kube.node117", "kube.node119", "kube.node68", "kube.node67",
	}
	bs, _ := json.Marshal(nodeName)
	NodeName := string(bs)

	params.NodeName = &NodeName

	handler := compatible.DoListAffinityGroupByNodeHandler{}
	response := handler.Handle(params)
	fmt.Println(response)
}

func init() {
	var err error
	mysql.TafDb, err = loadTafDBDev()
	if err != nil {
		fmt.Println(fmt.Sprintf("load taf_db error: %v", err))
	}

	k8sNamespace, k8sConfig, err := loadK8SDev()
	if err != nil {
		fmt.Println(fmt.Sprintf("load k8s error: %v", err))
	}

	if k8s.K8sOption, k8s.K8sWatcher, err = k8s.StartWatcher(k8sNamespace, k8sConfig); err != nil {
		fmt.Println(fmt.Sprintf("start watcher error: %v", err))
	}

	compatible.StartNodeWatch()
}
