package k8s

import (
	"encoding/json"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "k8s.tars.io/crd/v1alpha1"
	"tarsadmin/handler/util"
	"tarsadmin/openapi/models"
	"tarsadmin/openapi/restapi/operations/business"
)

type CreateBusinessHandler struct {}

func (s *CreateBusinessHandler) Handle(params business.CreateBusinessParams) middleware.Responder {

	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	tTree, err := K8sWatcher.tTreeLister.TTrees(namespace).Get(TarsTreeName)
	if err != nil {
		return business.NewCreateBusinessInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	tTree.Businesses = append(tTree.Businesses, crdv1alpha1.TTreeBusiness{Name: *metadata.BusinessName, Show: *metadata.BusinessShow,
		Mark: metadata.BusinessMark, Weight: *metadata.BusinessOrder, CreatePerson: metadata.CreatePerson, CreateTime: k8sMetaV1.Now()})

	if err = updateTTreeSpec(namespace, tTree); err != nil {
		return business.NewCreateBusinessInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return business.NewCreateBusinessOK().WithPayload(&business.CreateBusinessOKBody{Result: 0})
}

type SelectBusinessHandler struct {}

func (s *SelectBusinessHandler) Handle(params business.SelectBusinessParams) middleware.Responder {
	namespace := K8sOption.Namespace

	tTree, err := K8sWatcher.tTreeLister.TTrees(namespace).Get(TarsTreeName)
	if err != nil {
		return business.NewSelectBusinessInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}
	allItems := tTree.Businesses

	// parse query
	selectParams, err := util.ParseSelectQuery(params.Filter, params.Limiter, nil)
	if err != nil {
		return business.NewSelectBusinessInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	// filter
	filterItems := allItems

	// Admin临时版本，Template特化已有web实现
	if selectParams.Filter != nil {
		filterItems = make([]crdv1alpha1.TTreeBusiness, 0, len(allItems))
		for _, elem := range allItems {
			if selectParams.Filter.Eq != nil {
				return business.NewSelectBusinessInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Eq Is Not Supported."})
			}
			if selectParams.Filter.Ne != nil {
				return business.NewSelectBusinessInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Ne Is Not Supported."})
			}
			if selectParams.Filter.Like != nil {
				m1, err := LikeMatch("BusinessName", selectParams.Filter.Like, elem.Name)
				if err != nil {
					return business.NewSelectBusinessInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
				}
				m2, err := LikeMatch("BusinessShow", selectParams.Filter.Like, elem.Show)
				if err != nil {
					return business.NewSelectBusinessInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
				}
				m3, err := LikeMatch("BusinessMark", selectParams.Filter.Like, elem.Mark)
				if err != nil {
					return business.NewSelectBusinessInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
				}
				if !m1 || !m2 || !m3 {
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
		elem["BusinessName"] = item.Name
		elem["BusinessShow"] = item.Show
		elem["BusinessMark"] = item.Mark
		elem["BusinessOrder"] = item.Weight
		elem["CreateTime"] = item.CreateTime
		elem["CreatePerson"] = item.CreatePerson
		result.Data = append(result.Data, elem)
	}

	return business.NewSelectBusinessOK().WithPayload(result)
}

type UpdateBusinessHandler struct {}

func (s *UpdateBusinessHandler) Handle(params business.UpdateBusinessParams) middleware.Responder {

	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	tTree, index, err := getBuzNameIndex(namespace, *metadata.BusinessName)
	if err != nil {
		return business.NewCreateBusinessInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}
	if index == -1 {
		return business.NewCreateBusinessInternalServerError().WithPayload(&models.Error{Code: -1, Message: fmt.Sprintf("Invalid Business: %s Is Missing.", *metadata.BusinessName)})
	}

	target := params.Params.Target
	tTree.Businesses[index].Show = target.BusinessShow
	tTree.Businesses[index].Mark = target.BusinessMark
	tTree.Businesses[index].Weight = target.BusinessOrder

	if err = updateTTreeSpec(namespace, tTree.DeepCopy()); err != nil {
		return business.NewCreateBusinessInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return business.NewUpdateBusinessOK().WithPayload(&business.UpdateBusinessOKBody{Result: 0})
}


type DeleteBusinessHandler struct {}

func (s *DeleteBusinessHandler) Handle(params business.DeleteBusinessParams) middleware.Responder {
	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	tTree, index, err := getBuzNameIndex(namespace, *metadata.BusinessName)
	if err != nil {
		return business.NewDeleteBusinessInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}
	if index == -1 {
		return business.NewDeleteBusinessInternalServerError().WithPayload(&models.Error{Code: -1, Message: fmt.Sprintf("Invalid Business: %s Is Missing.", *metadata.BusinessName)})
	}

	tTree.Businesses = append(tTree.Businesses[0:index], tTree.Businesses[index+1:len(tTree.Businesses)] ...)

	if err = updateTTreeSpec(namespace, tTree.DeepCopy()); err != nil {
		return business.NewDeleteBusinessInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return business.NewDeleteBusinessOK().WithPayload(&business.DeleteBusinessOKBody{Result: 0})
}

type DoListBusinessAppHandler struct {}

func (s *DoListBusinessAppHandler) Handle(params business.DoListBusinessAppParams) middleware.Responder {

	namespace := K8sOption.Namespace

	BusinessName := make([]string, 0, 5)
	if err := json.Unmarshal([]byte(*params.BusinessName), &BusinessName); err != nil {
		return business.NewDoListBusinessAppInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}


	tTree, err := K8sWatcher.tTreeLister.TTrees(namespace).Get(TarsTreeName)
	if err != nil {
		return business.NewDoListBusinessAppInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	// 需要BusinessShow
	buzMap := map[string]*models.BusinessGroupElem {
		"": {BusinessName: "", BusinessShow: "", App: make([]string, 0, 5)},
	}
	for _, buz := range tTree.Businesses {
		buzMap[buz.Name] = &models.BusinessGroupElem{BusinessName: buz.Name, BusinessShow: buz.Show, App: make([]string, 0, 5)}
	}
	for _, app := range tTree.Apps {
		if _, ok := buzMap[app.BusinessRef]; ok {
			buzMap[app.BusinessRef].App = append(buzMap[app.BusinessRef].App, app.Name)
		}
	}

	result := make([]*models.BusinessGroupElem, 0, 5)

	for _, buzName := range BusinessName {
		item, ok := buzMap[buzName]
		if !ok {
			return business.NewDoListBusinessAppInternalServerError().WithPayload(&models.Error{Code: -1, Message: fmt.Sprintf("Invalid Business: %s Is Missing.", buzName)})
		}
		result = append(result, item)
	}

	return business.NewDoListBusinessAppOK().WithPayload(result)
}

type DoAddBusinessAppHandler struct {}

func (s *DoAddBusinessAppHandler) Handle(params business.DoAddBusinessAppParams) middleware.Responder {

	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	tTree, err := K8sWatcher.tTreeLister.TTrees(namespace).Get(TarsTreeName)
	if err != nil {
		return business.NewDoAddBusinessAppInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	appMap := make(map[string]*crdv1alpha1.TTreeApp, len(tTree.Apps))
	for i, app := range tTree.Apps {
		appMap[app.Name] = &tTree.Apps[i]
	}

	for _, app := range metadata.AppName {
		if _, ok := appMap[app]; ok {
			appMap[app].BusinessRef = metadata.BusinessName
		}
	}

	if err = updateTTreeSpec(namespace, tTree.DeepCopy()); err != nil {
		return business.NewDoAddBusinessAppInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return business.NewDoAddBusinessAppOK().WithPayload(&business.DoAddBusinessAppOKBody{Result: 0})
}

func getBuzNameIndex(namespace, buzName string) (*crdv1alpha1.TTree, int, error)  {
	tTree, err := K8sWatcher.tTreeLister.TTrees(namespace).Get(TarsTreeName)
	if err != nil {
		return nil, -1, err
	}

	index := -1
	for i, buz := range tTree.Businesses {
		if buzName == buz.Name {
			index = i
			break
		}
	}

	return tTree, index, nil
}
