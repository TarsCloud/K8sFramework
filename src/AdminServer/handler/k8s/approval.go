package k8s

import (
	"context"
	"github.com/go-openapi/runtime/middleware"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	crdv1alpha1 "k8s.taf.io/crd/v1alpha1"
	"sort"
	"tafadmin/handler/util"
	"tafadmin/openapi/models"
	"tafadmin/openapi/restapi/operations/approval"
	"tafadmin/openapi/restapi/operations/server_pod"
)

type CreateApprovalHandler struct {}

func (s *CreateApprovalHandler) Handle(params approval.CreateApprovalParams) middleware.Responder {
	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	tDeploy, err := K8sWatcher.tDeployLister.TDeploys(namespace).Get(*metadata.DeployID)
	if err != nil {
		return approval.NewCreateApprovalInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	if tDeploy.Approve != nil {
		return approval.NewCreateApprovalInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Can Not Do Duplicated Deploy."})
	}

	tDeployCopy := tDeploy.DeepCopy()
	tDeployCopy.Approve = &crdv1alpha1.TDeployApprove{
		Person: "taf-admin",
		Time:   k8sMetaV1.Now(),
		Reason: metadata.ApprovalMark,
		Result: metadata.ApprovalResult,
	}

	if metadata.ApprovalResult {
		serverK8S := ConvertOperatorK8SToAdminK8S(tDeployCopy.Apply.K8S)
		serverServant := ConvertOperatorServantToAdminK8S(tDeployCopy.Apply.Taf.Servants)
		serverOption := ConvertOperatorOptionToAdminK8S(tDeployCopy.Apply)

		if err = CreateServer(tDeployCopy.Apply.App, tDeployCopy.Apply.Server, serverServant, &serverK8S, &serverOption); err != nil {
			return approval.NewCreateApprovalInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
	}

	tDeployInterface := K8sOption.CrdClientSet.CrdV1alpha1().TDeploys(namespace)
	if _, err = tDeployInterface.Update(context.TODO(), tDeployCopy, k8sMetaV1.UpdateOptions{}); err != nil {
		return approval.NewCreateApprovalInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return approval.NewCreateApprovalOK().WithPayload(&approval.CreateApprovalOKBody{Result: 0})
}

type SelectApprovalHandler struct {}

func (s *SelectApprovalHandler) Handle(params approval.SelectApprovalParams) middleware.Responder {
	// fetch list
	namespace := K8sOption.Namespace

	// parse query
	selectParams, err := util.ParseSelectQuery(params.Filter, params.Limiter, params.Order)
	if err != nil {
		return approval.NewSelectApprovalInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	requirement, _ := labels.NewRequirement(TDeployApproveLabel, selection.In, []string{"Approved", "Reject"})
	requirements := []labels.Requirement{*requirement}

	allItems, err := K8sWatcher.tDeployLister.TDeploys(namespace).List(labels.NewSelector().Add(requirements ...))
	if err != nil {
		return approval.NewSelectApprovalInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	// filter
	filterItems := allItems

	// Admin临时版本，TServer特化已有web实现
	if selectParams.Filter != nil {
		filterItems = make([]*crdv1alpha1.TDeploy, 0, len(allItems))
		for _, elem := range allItems {
			if selectParams.Filter.Eq != nil {
				// ignore and not useful
				// return approval.NewSelectApprovalInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Eq Is Not Supported."})
			}
			if selectParams.Filter.Ne != nil {
				return approval.NewSelectApprovalInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Ne Is Not Supported."})
			}
			if selectParams.Filter.Like != nil {
				m1, err := LikeMatch("ServerApp", selectParams.Filter.Like, elem.Apply.App)
				if err != nil {
					return approval.NewSelectApprovalInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
				}
				m2, err := LikeMatch("ServerName", selectParams.Filter.Like, elem.Apply.Server)
				if err != nil {
					return approval.NewSelectApprovalInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
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

		elem["ApprovalTime"] = k8sMetaV1.Now()
		elem["ApprovalPerson"] = item.Approve.Person
		elem["ApprovalResult"] = item.Approve.Result
		elem["ApprovalMark"] = item.Approve.Reason

		result.Data = append(result.Data, elem)
	}

	return approval.NewSelectApprovalOK().WithPayload(result)
}
