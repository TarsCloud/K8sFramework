package k8s

import (
	"github.com/go-openapi/runtime/middleware"
	"k8s.io/apimachinery/pkg/labels"
	crdv1alpha1 "k8s.taf.io/crd/v1alpha1"
	"sort"
	"tafadmin/handler/util"
	"tafadmin/openapi/models"
	"tafadmin/openapi/restapi/operations/server_pod"
)


type SelectPodAliveHandler struct {}

func (s *SelectPodAliveHandler) Handle(params server_pod.SelectPodAliveParams) middleware.Responder {

	// fetch list
	namespace := K8sOption.Namespace

	// parse query
	selectParams, err := util.ParseSelectQuery(params.Filter, params.Limiter, params.Order)
	if err != nil {
		return server_pod.NewSelectPodAliveInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	listAll := true
	if selectParams.Filter != nil && selectParams.Filter.Eq != nil {
		if _, ok := selectParams.Filter.Eq["PodId"]; ok {
			return server_pod.NewSelectPodAliveInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Eq PodId Is Not Supported."})
		}
		_, appOk := selectParams.Filter.Eq["ServerApp"]
		_, nameOk := selectParams.Filter.Eq["ServerName"]
		if appOk || nameOk {
			listAll = false
		}
	}

	allEndpointItems := make([]*crdv1alpha1.TEndpoint, 0, 10)
	if listAll {
		return server_pod.NewSelectPodAliveInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Invalid Select Query Request."})
	} else {
		requirements := BuildDoubleEqualSelector(selectParams.Filter, KeyLabel)
		list, err := K8sWatcher.tEndpointLister.TEndpoints(namespace).List(labels.NewSelector().Add(requirements ...))
		if err != nil {
			return server_pod.NewSelectPodAliveInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
		allEndpointItems = list
	}

	filterItems := make([]*crdv1alpha1.TEndpoint, 0, 10)
	for _, Endpoint := range allEndpointItems {
		if Endpoint.Spec.SubType != crdv1alpha1.TAF {
			continue
		}
		if len(Endpoint.Status.PodStatus) <= 0 {
			continue
		}
		filterItems = append(filterItems, Endpoint)
	}

	// Admin临时版本，TServer特化已有web实现
	// 无

	// order
	if selectParams.Order != nil {
		order := ([]*models.SelectRequestOrderElem)(*selectParams.Order)
		if len(order) <= 0 || order[0].Column != "PodName" || order[0].Order != "asc" {
			return server_pod.NewSelectPodAliveInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Invalid Select RequestOrder."})
		}
		sort.Sort(TEndpointWrapper{Endpoint: filterItems, By: func(e1, e2 *crdv1alpha1.TEndpoint) bool {
			return e1.Status.PodStatus[0].Name < e2.Status.PodStatus[0].Name
		}})
	}

	// limiter
	if selectParams.Limiter != nil {
		start, stop := PageList(len(filterItems), selectParams.Limiter)
		filterItems = filterItems[start:stop]
	}

	// Count填充
	result := &models.SelectResult{}
	result.Count = make(models.MapInt)
	result.Count["AllCount"] = int32(len(filterItems))
	result.Count["FilterCount"] = int32(len(filterItems))

	// Data填充
	result.Data = make(models.ArrayMapInterface, 0, len(filterItems))
	for _, item := range filterItems {
		if len(item.Status.PodStatus) <= 0 {
			continue
		}
		for _, podStatus := range item.Status.PodStatus {
			elem := make(map[string]interface{})
			elem["ServerId"] = util.GetServerId(item.Spec.App, item.Spec.Server)
			elem["ServerApp"] = item.Spec.App
			elem["ServerName"] = item.Spec.Server

			elem["PodId"] = podStatus.UID
			elem["PodName"] = podStatus.Name
			elem["PodIp"] = podStatus.PodIP
			elem["NodeIp"] = podStatus.HostIP
			elem["ServiceVersion"] = podStatus.Tag
			elem["SettingState"] = podStatus.SettingState
			elem["PresentState"] = podStatus.PresentState
			elem["PresentMessage"] = podStatus.PresentMessage
			elem["CreateTime"] = podStatus.StartTime
			result.Data = append(result.Data, elem)
		}
	}

	return server_pod.NewSelectPodAliveOK().WithPayload(result)
}

type SelectPodPerishedHandler struct {}

func (s *SelectPodPerishedHandler) Handle(params server_pod.SelectPodPerishedParams) middleware.Responder {

	// fetch list
	namespace := K8sOption.Namespace

	// parse query
	selectParams, err := util.ParseSelectQuery(params.Filter, params.Limiter, params.Order)
	if err != nil {
		return server_pod.NewSelectPodPerishedInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	listAll := true
	if selectParams.Filter != nil && selectParams.Filter.Eq != nil {
		if _, ok := selectParams.Filter.Eq["PodId"]; ok {
			return server_pod.NewSelectPodPerishedInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Eq PodId Is Not Supported."})
		}
		_, appOk := selectParams.Filter.Eq["ServerApp"]
		_, nameOk := selectParams.Filter.Eq["ServerName"]
		if appOk || nameOk {
			listAll = false
		}
	}

	allExitedPodItems := make([]*crdv1alpha1.TExitedRecord, 0, 10)
	if listAll {
		return server_pod.NewSelectPodPerishedInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Invalid Select Query Request."})
	} else {
		requirements := BuildDoubleEqualSelector(selectParams.Filter, KeyLabel)
		list, err := K8sWatcher.tTExitedPod.TExitedRecords(namespace).List(labels.NewSelector().Add(requirements ...))
		if err != nil {
			return server_pod.NewSelectPodPerishedInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
		allExitedPodItems = list
	}

	filterItems := make([]ExitedPodElemEx, 0, 10)
	for _, items := range allExitedPodItems {
		for _, item := range items.Pods {
			filterItems = append(filterItems, ExitedPodElemEx{AppName: items.App, ServerName: items.Server, TExitedPod: item})
		}
	}

	// Admin临时版本，TServer特化已有web实现
	// 无

	// order
	if selectParams.Order != nil {
		order := ([]*models.SelectRequestOrderElem)(*selectParams.Order)
		if len(order) <= 0 || order[0].Column != "DeleteTime" || order[0].Order != "desc" {
			return server_pod.NewSelectPodPerishedInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Invalid Select RequestOrder."})
		}
		sort.Sort(TExitedPodWrapper{ExitedPod: filterItems, By: func(e1, e2 *ExitedPodElemEx) bool {
			return e1.DeleteTime.After(e2.DeleteTime.Time)
		}})
	}

	// limiter
	if selectParams.Limiter != nil {
		start, stop := PageList(len(filterItems), selectParams.Limiter)
		filterItems = filterItems[start:stop]
	}

	// Count填充
	result := &models.SelectResult{}
	result.Count = make(models.MapInt)
	result.Count["AllCount"] = int32(len(filterItems))
	result.Count["FilterCount"] = int32(len(filterItems))

	// Data填充
	result.Data = make(models.ArrayMapInterface, 0, len(filterItems))
	for _, item := range filterItems {
		elem := make(map[string]interface{})
		elem["ServerId"] = util.GetServerId(item.AppName, item.ServerName)
		elem["ServerApp"] = item.AppName
		elem["ServerName"] = item.ServerName
		elem["ServiceVersion"] = item.Tag
		elem["CreateTime"] = item.CreateTime
		elem["DeleteTime"] = item.DeleteTime

		elem["PodId"] = item.UID
		elem["PodName"] = item.Name
		elem["PodIp"] = item.PodIP
		elem["NodeIp"] = item.NodeIP
		result.Data = append(result.Data, elem)
	}

	return server_pod.NewSelectPodPerishedOK().WithPayload(result)
}
