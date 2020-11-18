package k8s

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"golang.org/x/net/context"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	crdv1alpha1 "k8s.taf.io/crd/v1alpha1"
	"strings"
	"tafadmin/handler/util"
	"tafadmin/openapi/models"
	"tafadmin/openapi/restapi/operations/server_servant"
)

type CreateServerAdapterHandler struct {}

func (s *CreateServerAdapterHandler) Handle(params server_servant.CreateServerAdapterParams) middleware.Responder {
	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	tServer, err := K8sWatcher.tServerLister.TServers(namespace).Get(util.GetTServerName(*metadata.ServerID))
	if err != nil {
		return server_servant.NewCreateServerAdapterInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	tServerCopy := tServer.DeepCopy()
	Servants := make([]crdv1alpha1.TServant, 0, len(metadata.Servant))
	for _, target := range metadata.Servant {
		var adapter crdv1alpha1.TServant
		adapter.Name = target.Name
		adapter.IsTaf = *target.IsTaf
		adapter.IsTcp = *target.IsTCP
		adapter.Timeout = target.Timeout
		adapter.Capacity = target.Capacity
		adapter.Port = target.Port
		adapter.Connection = target.Connections
		adapter.Thread = target.Threads
		Servants = append(Servants, adapter)
	}
	tServerCopy.Spec.Taf.Servants = append(tServerCopy.Spec.Taf.Servants, Servants ...)

	tServerInterface := K8sOption.CrdClientSet.CrdV1alpha1().TServers(namespace)
	if _, err = tServerInterface.Update(context.TODO(), tServerCopy, k8sMetaV1.UpdateOptions{}); err != nil {
		return server_servant.NewCreateServerAdapterInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return server_servant.NewCreateServerAdapterOK().WithPayload(&server_servant.CreateServerAdapterOKBody{Result: 0})
}

type SelectServerAdapterHandler struct {}

func (s *SelectServerAdapterHandler) Handle(params server_servant.SelectServerAdapterParams) middleware.Responder {

	// fetch list
	namespace := K8sOption.Namespace

	// parse query
	selectParams, err := util.ParseSelectQuery(params.Filter, params.Limiter, nil)
	if err != nil {
		return server_servant.NewSelectServerAdapterInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	listAll := true
	if selectParams.Filter != nil && selectParams.Filter.Eq != nil {
		_, serverOk := selectParams.Filter.Eq["ServerId"]
		if serverOk {
			listAll = false
		}
	}

	allServerItems := make([]*crdv1alpha1.TServer, 0, 10)
	if listAll {
		requirements := BuildSubTypeTafSelector()
		list, err := K8sWatcher.tServerLister.TServers(K8sOption.Namespace).List(labels.NewSelector().Add(requirements ...))
		if err != nil {
			return server_servant.NewSelectServerAdapterInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
		allServerItems = list
	} else {
		serverId, _ := selectParams.Filter.Eq["ServerId"]
		item, err := K8sWatcher.tServerLister.TServers(namespace).Get(util.GetTServerName(serverId.(string)))
		if err != nil {
			return server_servant.NewSelectServerAdapterInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
		allServerItems = append(allServerItems, item)
	}

	allItems := make([]crdv1alpha1.TServant, 0, len(allServerItems))
	for _, Server := range allServerItems {
		for _, Servant := range Server.Spec.Taf.Servants {
			Servant.Name = getAdapterId(util.GetServerId(Server.Spec.App, Server.Spec.Server), Servant.Name)
			allItems = append(allItems, Servant)
		}
	}

	// filter
	filterItems := allItems

	// Admin临时版本，TServer特化已有web实现
	if selectParams.Filter != nil {
		filterItems = make([]crdv1alpha1.TServant, 0, len(allItems))
		for _, elem := range allItems {
			if selectParams.Filter.Eq != nil {
				pattern, ok := selectParams.Filter.Eq["AdapterId"]
				if ok && pattern != elem.Name{
					continue
				}
				pattern, ok = selectParams.Filter.Eq["IsTaf"]
				if ok && pattern != elem.IsTaf{
					continue
				}
				pattern, ok = selectParams.Filter.Eq["IsTcp"]
				if ok && pattern != elem.IsTcp{
					continue
				}
			}
			if selectParams.Filter.Ne != nil {
				return server_servant.NewSelectServerAdapterInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Ne Is Not Supported."})
			}
			if selectParams.Filter.Like != nil {
				return server_servant.NewSelectServerAdapterInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Like Is Not Supported."})
			}
			filterItems = append(filterItems, elem)
		}
	}

	// limiter
	if selectParams.Limiter != nil {
		start, stop := PageList(len(filterItems), selectParams.Limiter)
		filterItems = filterItems[start:stop]
	}

	// Count填充
	result := &models.SelectResult{}
	result.Count = make(models.MapInt)
	result.Count["AllCount"] = int32(len(allItems))
	result.Count["FilterCount"] = int32(len(filterItems))

	// Data填充
	result.Data = make(models.ArrayMapInterface, 0, len(filterItems))
	for _, item := range filterItems {
		elem := make(map[string]interface{})

		fields := strings.Split(item.Name, ".")
		if len(fields) != 3 {
			continue
		}
		elem["ServerId"] = util.GetServerId(fields[0], fields[1])
		elem["Name"] = fields[2]

		elem["AdapterId"] = item.Name
		elem["Threads"] = item.Thread
		elem["Connections"] = item.Connection
		elem["Port"] = item.Port
		elem["Capacity"] = item.Capacity
		elem["Timeout"] = item.Timeout
		elem["IsTaf"] = item.IsTaf
		elem["IsTcp"] = item.IsTcp
		result.Data = append(result.Data, elem)
	}

	return server_servant.NewSelectServerAdapterOK().WithPayload(result)
}

type UpdateServerAdapterHandler struct {}

func (s *UpdateServerAdapterHandler) Handle(params server_servant.UpdateServerAdapterParams) middleware.Responder {

	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	fields := strings.Split(*metadata.AdapterID, ".")
	if len(fields) != 2 {
		return server_servant.NewUpdateServerAdapterInternalServerError().WithPayload(&models.Error{Code: -1, Message: fmt.Sprintf("Invalid AdapterID syntax: %s", *metadata.AdapterID)})
	}

	ServerId := fields[0]
	AdapterName := fields[1]

	tServer, err := K8sWatcher.tServerLister.TServers(namespace).Get(util.GetTServerName(ServerId))
	if err != nil {
		return server_servant.NewUpdateServerAdapterInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	index := -1
	for i, adapter := range tServer.Spec.Taf.Servants {
		if adapter.Name == AdapterName {
			index = i
		}
	}
	if index == -1 {
		return server_servant.NewUpdateServerAdapterInternalServerError().WithPayload(&models.Error{Code: -1, Message: fmt.Sprintf("Can Not Find AdapterId: %s. ", *metadata.AdapterID)})
	}

	target := params.Params.Target
	if equalServerAdapter(&tServer.Spec.Taf.Servants[index], target) {
		return server_servant.NewUpdateServerAdapterInternalServerError().WithPayload(&models.Error{Code: -1, Message: fmt.Sprintf("No Need To Update Duplicated AdapterId: %s. ", *metadata.AdapterID)})
	}

	tServerCopy := tServer.DeepCopy()
	tServerCopy.Spec.Taf.Servants[index].Name = target.Name
	tServerCopy.Spec.Taf.Servants[index].IsTaf = *target.IsTaf
	tServerCopy.Spec.Taf.Servants[index].IsTcp = *target.IsTCP
	tServerCopy.Spec.Taf.Servants[index].Timeout = target.Timeout
	tServerCopy.Spec.Taf.Servants[index].Capacity = target.Capacity
	tServerCopy.Spec.Taf.Servants[index].Port = target.Port
	tServerCopy.Spec.Taf.Servants[index].Connection = target.Connections
	tServerCopy.Spec.Taf.Servants[index].Thread = target.Threads

	tServerInterface := K8sOption.CrdClientSet.CrdV1alpha1().TServers(namespace)
	if _, err = tServerInterface.Update(context.TODO(), tServerCopy, k8sMetaV1.UpdateOptions{}); err != nil {
		return server_servant.NewUpdateServerAdapterInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return server_servant.NewUpdateServerAdapterOK().WithPayload(&server_servant.UpdateServerAdapterOKBody{Result: 0})
}

type DeleteServerAdapterHandler struct {}

func (s *DeleteServerAdapterHandler) Handle(params server_servant.DeleteServerAdapterParams) middleware.Responder {
	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	fields := strings.Split(*metadata.AdapterID, ".")
	ServerId := fields[0]
	AdapterName := fields[1]

	tServer, err := K8sWatcher.tServerLister.TServers(namespace).Get(util.GetTServerName(ServerId))
	if err != nil {
		return server_servant.NewDeleteServerAdapterInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	index := -1
	for i, adapter := range tServer.Spec.Taf.Servants {
		if adapter.Name == AdapterName {
			index = i
		}
	}
	if index == -1 {
		return server_servant.NewDeleteServerAdapterInternalServerError().WithPayload(&models.Error{Code: -1, Message: fmt.Sprintf("Can Not Find AdapterId: %s. ", *metadata.AdapterID)})
	}

	tServerCopy := tServer.DeepCopy()
	tServerCopy.Spec.Taf.Servants = append(tServerCopy.Spec.Taf.Servants[0:index], tServerCopy.Spec.Taf.Servants[index+1:] ...)

	tServerInterface := K8sOption.CrdClientSet.CrdV1alpha1().TServers(namespace)
	if _, err = tServerInterface.Update(context.TODO(), tServerCopy, k8sMetaV1.UpdateOptions{}); err != nil {
		return server_servant.NewDeleteServerAdapterInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return server_servant.NewDeleteServerAdapterOK().WithPayload(&server_servant.DeleteServerAdapterOKBody{Result: 0})
}

func getAdapterId(serverId, adapterName string) string {
	return fmt.Sprintf("%s.%s", serverId, adapterName)
}

func equalServerAdapter(oldAdapter *crdv1alpha1.TServant, newAdapter *models.ServerServantElem) bool {
	if oldAdapter.Name != newAdapter.Name {
		return false
	}
	if oldAdapter.Thread != newAdapter.Threads {
		return false
	}
	if oldAdapter.Connection != newAdapter.Connections {
		return false
	}
	if oldAdapter.Port != newAdapter.Port {
		return false
	}
	if oldAdapter.Capacity != newAdapter.Capacity {
		return false
	}
	if oldAdapter.Timeout != newAdapter.Timeout {
		return false
	}
	if oldAdapter.IsTcp != *newAdapter.IsTCP {
		return false
	}
	if oldAdapter.IsTaf != *newAdapter.IsTaf {
		return false
	}
	return true
}

