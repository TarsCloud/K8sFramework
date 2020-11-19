package k8s

import (
	"github.com/go-openapi/runtime/middleware"
	"golang.org/x/net/context"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	crdv1alpha1 "k8s.tars.io/crd/v1alpha1"
	"tarsadmin/handler/util"
	"tarsadmin/openapi/models"
	"tarsadmin/openapi/restapi/operations/server_k8s"
)

type SelectServerK8SHandler struct {}

func (s *SelectServerK8SHandler) Handle(params server_k8s.SelectK8SParams) middleware.Responder {

	// fetch list
	namespace := K8sOption.Namespace

	// parse query
	selectParams, err := util.ParseSelectQuery(params.Filter, params.Limiter, nil)
	if err != nil {
		return server_k8s.NewSelectK8SInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	listAll := true
	if selectParams.Filter != nil && selectParams.Filter.Eq != nil {
		_, serverOk := selectParams.Filter.Eq["ServerId"]
		if serverOk {
			listAll = false
		}
	}

	filterItems := make([]*crdv1alpha1.TServer, 0, 10)
	if listAll {
		requirements := BuildSubTypeTarsSelector()
		list, err := K8sWatcher.tServerLister.TServers(K8sOption.Namespace).List(labels.NewSelector().Add(requirements ...))
		if err != nil {
			return server_k8s.NewSelectK8SInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
		filterItems = list
	} else {
		serverId, _ := selectParams.Filter.Eq["ServerId"]
		item, err := K8sWatcher.tServerLister.TServers(namespace).Get(util.GetTServerName(serverId.(string)))
		if err != nil {
			return server_k8s.NewSelectK8SInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
		filterItems = append(filterItems, item)
	}

	// Admin临时版本，TServer特化已有web实现
	// 无

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
		elem["ServerId"] = util.GetServerId(item.Spec.App, item.Spec.Server)
		elem["ServerApp"] = item.Spec.App
		elem["ServerName"] = item.Spec.Server
		elem["Replicas"] = item.Spec.K8S.Replicas
		elem["NodeSelector"] = item.Spec.K8S.NodeSelector
		elem["HostIpc"] = item.Spec.K8S.HostIPC
		elem["HostNetwork"] = item.Spec.K8S.HostNetwork
		elem["HostPort"] = item.Spec.K8S.HostPorts
		result.Data = append(result.Data, elem)
	}

	return server_k8s.NewSelectK8SOK().WithPayload(result)
}

type UpdateServerK8SHandler struct {}

func (s *UpdateServerK8SHandler) Handle(params server_k8s.UpdateK8SParams) middleware.Responder {

	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	tServer, err := K8sWatcher.tServerLister.TServers(namespace).Get(util.GetTServerName(*metadata.ServerID))
	if err != nil {
		return server_k8s.NewUpdateK8SInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	K8S := tServer.Spec.K8S
	target := params.Params.Target
	if equalServerK8S(&K8S, target) {
		return server_k8s.NewUpdateK8SInternalServerError().WithPayload(&models.Error{Code: -1, Message: "No Need To Update Duplicated ServerK8S. "})
	}

	K8S.NotStacked = target.NotStacked
	K8S.Replicas = target.Replicas
	K8S.HostIPC = target.HostIpc
	K8S.HostNetwork = target.HostNetwork
	if target.NodeSelector != nil {
		if target.NodeSelector.NodeBind != nil {
			K8S.NodeSelector.NodeBind = &crdv1alpha1.TK8SNodeSelectorKind{}
			K8S.NodeSelector.NodeBind.Values = target.NodeSelector.NodeBind.Value
		} else {
			K8S.NodeSelector.NodeBind = nil
		}
		if target.NodeSelector.AbilityPool != nil {
			K8S.NodeSelector.AbilityPool = &crdv1alpha1.TK8SNodeSelectorKind{}
			K8S.NodeSelector.AbilityPool.Values = target.NodeSelector.AbilityPool.Value
		} else {
			K8S.NodeSelector.AbilityPool = nil
		}
		if target.NodeSelector.PublicPool != nil {
			K8S.NodeSelector.PublicPool = &crdv1alpha1.TK8SNodeSelectorKind{}
			K8S.NodeSelector.PublicPool.Values = target.NodeSelector.PublicPool.Value
		} else {
			K8S.NodeSelector.PublicPool = nil
		}
		if target.NodeSelector.DaemonSet != nil {
			K8S.NodeSelector.DaemonSet = &crdv1alpha1.TK8SNodeSelectorKind{}
			K8S.NodeSelector.DaemonSet.Values = target.NodeSelector.DaemonSet.Value
		} else {
			K8S.NodeSelector.DaemonSet = nil
		}
	}
	if target.HostPort != nil {
		K8S.HostPorts = make([]crdv1alpha1.TK8SHostPort, len(target.HostPort))
		for i, v := range target.HostPort {
			K8S.HostPorts[i].NameRef = v.NameRef
			K8S.HostPorts[i].Port = v.Port
		}
	}

	// K8S是值...
	tServerCopy := tServer.DeepCopy()
	tServerCopy.Spec.K8S = K8S

	tServerInterface := K8sOption.CrdClientSet.CrdV1alpha1().TServers(namespace)
	if _, err = tServerInterface.Update(context.TODO(), tServerCopy, k8sMetaV1.UpdateOptions{}); err != nil {
		return server_k8s.NewUpdateK8SInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return server_k8s.NewUpdateK8SOK().WithPayload(&server_k8s.UpdateK8SOKBody{Result: 0})
}

func equalServerK8S(oldK8S *crdv1alpha1.TServerK8S, newK8S *models.ServerK8S) bool {
	if oldK8S.Replicas != newK8S.Replicas {
		return false
	}
	if oldK8S.HostNetwork != newK8S.HostNetwork {
		return false
	}
	if oldK8S.HostIPC != newK8S.HostIpc {
		return false
	}
	if len(oldK8S.HostPorts) != len(newK8S.HostPort) {
		return false
	}
	for i, hostPort := range oldK8S.HostPorts {
		if hostPort.NameRef != newK8S.HostPort[i].NameRef {
			return false
		}
		if hostPort.Port != newK8S.HostPort[i].Port {
			return false
		}
	}
	return equalNodeSelector(oldK8S.NodeSelector.AbilityPool, newK8S.NodeSelector.AbilityPool) &&
		equalNodeSelector(oldK8S.NodeSelector.NodeBind, newK8S.NodeSelector.NodeBind) &&
		equalNodeSelector(oldK8S.NodeSelector.PublicPool, newK8S.NodeSelector.PublicPool) &&
		equalNodeSelector(oldK8S.NodeSelector.DaemonSet, newK8S.NodeSelector.DaemonSet)
}

func equalNodeSelector(oldSelector *crdv1alpha1.TK8SNodeSelectorKind, newSelector *models.NodeSelectorElem) bool  {
	if oldSelector == nil && newSelector != nil {
		return false
	}
	if oldSelector != nil && newSelector == nil {
		return false
	}
	if oldSelector != nil && newSelector != nil {
		return false
	}
	if oldSelector != nil && newSelector != nil {
		if len(oldSelector.Values) != len(newSelector.Value) {
			return false
		}
		for i, v := range oldSelector.Values {
			if v != newSelector.Value[i] {
				return false
			}
		}
	}
	return true
}
