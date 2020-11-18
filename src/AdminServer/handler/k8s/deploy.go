package k8s

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/api/errors"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	crdv1alpha1 "k8s.taf.io/crd/v1alpha1"
	"sort"
	"tafadmin/handler/util"
	"tafadmin/openapi/models"
	"tafadmin/openapi/restapi/operations/deploy"
	"tafadmin/openapi/restapi/operations/server_pod"
)

type CreateDeployHandler struct {}

func (s *CreateDeployHandler) Handle(params deploy.CreateDeployParams) middleware.Responder {
	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	tDeploy := buildTDeploy(namespace, metadata)

	tDeployInterface := K8sOption.CrdClientSet.CrdV1alpha1().TDeploys(namespace)
	if _, err := tDeployInterface.Create(context.TODO(), tDeploy, k8sMetaV1.CreateOptions{}); err != nil && !errors.IsAlreadyExists(err) {
		return deploy.NewCreateDeployInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return deploy.NewCreateDeployOK().WithPayload(&deploy.CreateDeployOKBody{Result: 0})
}

type SelectDeployHandler struct {}

func (s *SelectDeployHandler) Handle(params deploy.SelectDeployParams) middleware.Responder {
	// fetch list
	namespace := K8sOption.Namespace

	// parse query
	selectParams, err := util.ParseSelectQuery(params.Filter, params.Limiter, params.Order)
	if err != nil {
		return deploy.NewSelectDeployInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	requirement, _ := labels.NewRequirement(TDeployApproveLabel, selection.DoubleEquals, []string{"Pending"})
	requirements := []labels.Requirement{*requirement}

	allItems, err := K8sWatcher.tDeployLister.TDeploys(namespace).List(labels.NewSelector().Add(requirements ...))
	if err != nil {
		return deploy.NewSelectDeployInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	// filter
	filterItems := allItems

	// Admin临时版本，TServer特化已有web实现
	if selectParams.Filter != nil {
		filterItems = make([]*crdv1alpha1.TDeploy, 0, len(allItems))
		for _, elem := range allItems {
			if selectParams.Filter.Eq != nil {
				return deploy.NewSelectDeployInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Eq Is Not Supported."})
			}
			if selectParams.Filter.Ne != nil {
				return deploy.NewSelectDeployInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Ne Is Not Supported."})
			}
			if selectParams.Filter.Like != nil {
				m1, err := LikeMatch("ServerApp", selectParams.Filter.Like, elem.Apply.App)
				if err != nil {
					return deploy.NewSelectDeployInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
				}
				m2, err := LikeMatch("ServerName", selectParams.Filter.Like, elem.Apply.Server)
				if err != nil {
					return deploy.NewSelectDeployInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
				}
				if !m1 || !m2 {
					continue
				}
			}
			filterItems = append(filterItems, elem)
		}
	}

	// order
	if selectParams.Order != nil {
		order := ([]*models.SelectRequestOrderElem)(*selectParams.Order)
		if len(order) <= 0 || order[0].Column != "RequestTime" || order[0].Order != "desc" {
			return server_pod.NewSelectPodAliveInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Invalid Select RequestOrder."})
		}
		sort.Sort(TDeployWrapper{Deploy: filterItems, By: func(e1, e2 *crdv1alpha1.TDeploy) bool {
			return !e1.CreationTimestamp.Before(&e2.CreationTimestamp)
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
	result.Count["AllCount"] = int32(len(allItems))
	result.Count["FilterCount"] = int32(len(filterItems))

	// Data填充
	result.Data = make(models.ArrayMapInterface, 0, len(filterItems))
	for _, item := range filterItems {
		elem := make(map[string]interface{})

		elem["DeployId"] = item.Name
		elem["RequestTime"] = item.CreationTimestamp
		elem["ServerApp"] = item.Apply.App
		elem["ServerName"] = item.Apply.Server

		elem["ServerK8S"] = ConvertOperatorK8SToAdminK8S(item.Apply.K8S)
		elem["ServerServant"] = ConvertOperatorServantToAdminK8S(item.Apply.Taf.Servants)
		elem["ServerOption"] = ConvertOperatorOptionToAdminK8S(item.Apply)

		elem["RequestPerson"] = "taf-admin"
		elem["ServerMark"] = "default empty"

		result.Data = append(result.Data, elem)
	}

	return deploy.NewSelectDeployOK().WithPayload(result)
}

type UpdateDeployHandler struct {}

func (s *UpdateDeployHandler) Handle(params deploy.UpdateDeployParams) middleware.Responder {

	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	tDeploy, err := K8sWatcher.tDeployLister.TDeploys(namespace).Get(*metadata.DeployID)
	if err != nil {
		return deploy.NewUpdateDeployInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	// 忽略对比是否相同....
	target := params.Params.Target
	tServer := buildTServer(namespace, tDeploy.Apply.App, tDeploy.Apply.Server, target.ServerServant, target.ServerK8S, target.ServerOption)

	tDeployCopy := tDeploy.DeepCopy()
	tDeployCopy.Apply = tServer.Spec

	tDeployInterface := K8sOption.CrdClientSet.CrdV1alpha1().TDeploys(namespace)
	if _, err := tDeployInterface.Update(context.TODO(), tDeployCopy, k8sMetaV1.UpdateOptions{}); err != nil {
		return deploy.NewUpdateDeployInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return deploy.NewUpdateDeployOK().WithPayload(&deploy.UpdateDeployOKBody{Result: 0})
}

type DeleteDeployHandler struct {}

func (s *DeleteDeployHandler) Handle(params deploy.DeleteDeployParams) middleware.Responder {
	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	tDeployInterface := K8sOption.CrdClientSet.CrdV1alpha1().TDeploys(namespace)
	if err := tDeployInterface.Delete(context.TODO(), *metadata.DeployID, k8sMetaV1.DeleteOptions{}); err != nil {
		return deploy.NewDeleteDeployInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return deploy.NewDeleteDeployOK().WithPayload(&deploy.DeleteDeployOKBody{Result: 0})
}

func buildTDeploy(namespace string, metadata *models.DeployMeta) *crdv1alpha1.TDeploy {

	//部署时 Replicas的值只能为0 ,因为此时没有镜像服务镜像
	metadata.ServerK8S.Replicas = 0
	if metadata.ServerK8S.HostNetwork {
		metadata.ServerK8S.HostPort = make([]*models.HostPortElem, 0, 1)
	}

	// 通过管理平台的部署都是TAF服务
	serverSubType := "taf"
	metadata.ServerOption.ServerSubType = &serverSubType

	deployName := fmt.Sprintf("%s-%s-%s", util.GetTServerName(util.GetServerId(*metadata.ServerApp, *metadata.ServerName)), RandStringRunes(10), RandStringRunes(5))
	tServer := buildTServer(namespace, *metadata.ServerApp, *metadata.ServerName, metadata.ServerServant, metadata.ServerK8S, metadata.ServerOption)

	tDeploy := &crdv1alpha1.TDeploy {
		ObjectMeta: k8sMetaV1.ObjectMeta{
			Name:      deployName,
			Namespace: namespace,
			CreationTimestamp: k8sMetaV1.Now(),
		},
		Apply: tServer.Spec,
	}
	return tDeploy
}

func ConvertOperatorK8SToAdminK8S(operatorK8S crdv1alpha1.TServerK8S) models.ServerK8S {
	var hostPort []*models.HostPortElem
	if operatorK8S.HostPorts != nil {
		hostPort = make([]*models.HostPortElem, len(operatorK8S.HostPorts))
		for i, port := range operatorK8S.HostPorts {
			hostPort[i] = &models.HostPortElem{NameRef: port.NameRef, Port: port.Port}
		}
	}
	var nodeSelector models.NodeSelector
	if operatorK8S.NodeSelector.NodeBind != nil {
		nodeSelector.NodeBind = &models.NodeSelectorElem{}
		nodeSelector.NodeBind.Value = operatorK8S.NodeSelector.NodeBind.Values
	}
	if operatorK8S.NodeSelector.AbilityPool != nil {
		nodeSelector.AbilityPool = &models.NodeSelectorElem{}
		nodeSelector.AbilityPool.Value = operatorK8S.NodeSelector.AbilityPool.Values
	}
	if operatorK8S.NodeSelector.PublicPool != nil {
		nodeSelector.PublicPool = &models.NodeSelectorElem{}
		nodeSelector.PublicPool.Value = operatorK8S.NodeSelector.PublicPool.Values
	}
	if operatorK8S.NodeSelector.DaemonSet != nil {
		nodeSelector.DaemonSet = &models.NodeSelectorElem{}
		nodeSelector.DaemonSet.Value = operatorK8S.NodeSelector.DaemonSet.Values
	}

	return models.ServerK8S{
		HostIpc: operatorK8S.HostIPC,
		HostNetwork: operatorK8S.HostNetwork,
		NotStacked: operatorK8S.NotStacked,
		Replicas: operatorK8S.Replicas,
		HostPort: hostPort,
		NodeSelector: &nodeSelector,
	}
}

func ConvertOperatorServantToAdminK8S(operatorServant []crdv1alpha1.TServant) models.MapServant {
	var serverServant models.MapServant

	if operatorServant != nil {
		serverServant = make(models.MapServant)
		for _, servant := range operatorServant {
			serverServant[servant.Name] = models.ServerServantElem{
				Capacity: servant.Capacity,
				Connections: servant.Connection,
				IsTaf: &servant.IsTaf,
				IsTCP: &servant.IsTcp,
				Name: servant.Name,
				Port: servant.Port,
				Threads: servant.Thread,
				Timeout: servant.Timeout,
			}
		}
	}

	return serverServant
}

func ConvertOperatorOptionToAdminK8S(operatorTServerSpec crdv1alpha1.TServerSpec) models.ServerOption {
	asyncThread := operatorTServerSpec.Taf.AsyncThread
	serverImportant := operatorTServerSpec.Important
	serverSubType := string(operatorTServerSpec.SubType)

	return models.ServerOption{
		AsyncThread: &asyncThread,
		ServerImportant: &serverImportant,
		ServerProfile: operatorTServerSpec.Taf.Profile,
		ServerTemplate: operatorTServerSpec.Taf.Template,
		ServerSubType: &serverSubType,
	}
}
