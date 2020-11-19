package k8s

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"golang.org/x/net/context"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	crdv1alpha1 "k8s.tars.io/crd/v1alpha1"
	"tarsadmin/handler/util"
	"tarsadmin/openapi/models"
	"tarsadmin/openapi/restapi/operations/applications"
)

type CreateAppHandler struct {}

func (s *CreateAppHandler) Handle(params applications.CreateAppParams) middleware.Responder {
	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	tTree, index, err := getAppNameIndex(namespace, *metadata.AppName)
	if err != nil {
		return applications.NewCreateAppInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}
	if index != -1 {
		return applications.NewCreateAppInternalServerError().WithPayload(&models.Error{Code: -1, Message: fmt.Sprintf("Duplicated App: %s Already Existed.", *metadata.AppName)})
	}

	tTree.Apps = append(tTree.Apps, crdv1alpha1.TTreeApp{Name: *metadata.AppName, BusinessRef: metadata.BusinessName,
		CreatePerson: metadata.CreatePerson, CreateTime: k8sMetaV1.Now(), Mark: metadata.AppMark})

	if err = updateTTreeSpec(namespace, tTree); err != nil {
		return applications.NewCreateAppInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return applications.NewCreateAppOK().WithPayload(&applications.CreateAppOKBody{Result: 0})
}

type SelectAppHandler struct {}

func (s *SelectAppHandler) Handle(params applications.SelectAppParams) middleware.Responder {
	namespace := K8sOption.Namespace

	tTree, err := K8sWatcher.tTreeLister.TTrees(namespace).Get(TarsTreeName)
	if err != nil {
		return applications.NewSelectAppInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}
	allItems := tTree.Apps

	// parse query
	selectParams, err := util.ParseSelectQuery(params.Filter, params.Limiter, nil)
	if err != nil {
		return applications.NewSelectAppInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	// filter
	filterItems := allItems

	// Admin临时版本，Template特化已有web实现
	if selectParams.Filter != nil {
		filterItems = make([]crdv1alpha1.TTreeApp, 0, len(allItems))
		for _, elem := range allItems {
			if selectParams.Filter.Eq != nil {
				return applications.NewSelectAppInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Eq Is Not Supported."})
			}
			if selectParams.Filter.Ne != nil {
				return applications.NewSelectAppInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Ne Is Not Supported."})
			}
			if selectParams.Filter.Like != nil {
				m1, err := LikeMatch("AppName", selectParams.Filter.Like, elem.Name)
				if err != nil {
					return applications.NewSelectAppInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
				}
				m2, err := LikeMatch("BusinessName", selectParams.Filter.Like, elem.BusinessRef)
				if err != nil {
					return applications.NewSelectAppInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
				}
				if !m1 || !m2 {
					continue
				}
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
		elem["AppName"] = item.Name
		elem["AppMark"] = item.Mark
		elem["BusinessName"] = item.BusinessRef
		elem["CreateTime"] = item.CreateTime
		elem["CreatePerson"] = item.CreatePerson
		result.Data = append(result.Data, elem)
	}

	return applications.NewSelectAppOK().WithPayload(result)
}

type DeleteAppHandler struct {}

func (s *DeleteAppHandler) Handle(params applications.DeleteAppParams) middleware.Responder {

	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	tTree, index, err := getAppNameIndex(namespace, *metadata.AppName)
	if err != nil {
		return applications.NewDeleteAppInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}
	if index == -1 {
		return applications.NewDeleteAppInternalServerError().WithPayload(&models.Error{Code: -1, Message: fmt.Sprintf("Invalid App: %s Is Missing.", *metadata.AppName)})
	}

	// 检测server是否存在，如果存在拒绝删除
	requirements := BuildTarsAppSelector(*metadata.AppName)
	_, err = K8sWatcher.tServerLister.TServers(K8sOption.Namespace).List(labels.NewSelector().Add(requirements ...))
	if err == nil {
		return applications.NewDeleteAppInternalServerError().WithPayload(&models.Error{Code: -1, Message: fmt.Sprintf("Must Clear All TServer With %s.", *metadata.AppName)})
	}

	tTree.Apps = append(tTree.Apps[0:index], tTree.Apps[index+1:len(tTree.Apps)] ...)

	if err = updateTTreeSpec(namespace, tTree.DeepCopy()); err != nil {
		return applications.NewDeleteAppInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return applications.NewDeleteAppOK().WithPayload(&applications.DeleteAppOKBody{Result: 0})
}

type UpdateAppHandler struct {}

func (s *UpdateAppHandler) Handle(params applications.UpdateAppParams) middleware.Responder {

	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	tTree, index, err := getAppNameIndex(namespace, *metadata.AppName)
	if err != nil {
		return applications.NewUpdateAppInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}
	if index == -1 {
		return applications.NewUpdateAppInternalServerError().WithPayload(&models.Error{Code: -1, Message: fmt.Sprintf("Invalid App: %s Is Missing.", *metadata.AppName)})
	}

	target := params.Params.Target
	tTree.Apps[index].BusinessRef = target.BusinessName
	tTree.Apps[index].Mark = target.AppMark

	if err = updateTTreeSpec(namespace, tTree.DeepCopy()); err != nil {
		return applications.NewUpdateAppInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return applications.NewUpdateAppOK().WithPayload(&applications.UpdateAppOKBody{Result: 0})
}

func getAppNameIndex(namespace, appName string) (*crdv1alpha1.TTree, int, error)  {
	tTree, err := K8sWatcher.tTreeLister.TTrees(namespace).Get(TarsTreeName)
	if err != nil {
		return nil, -1, err
	}

	index := -1
	for i, app := range tTree.Apps {
		if appName == app.Name {
			index = i
			break
		}
	}

	return tTree, index, nil
}

func updateTTreeSpec(namespace string, tTree *crdv1alpha1.TTree) error {
	tTreesInterface := K8sOption.CrdClientSet.CrdV1alpha1().TTrees(namespace)
	if _, err := tTreesInterface.Update(context.TODO(), tTree, k8sMetaV1.UpdateOptions{}); err != nil {
		return err
	}
	return nil
}
