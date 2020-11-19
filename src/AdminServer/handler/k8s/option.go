package k8s

import (
	"bytes"
	"fmt"
	tarsConf "github.com/TarsCloud/TarsGo/tars/util/conf"
	"github.com/go-openapi/runtime/middleware"
	"golang.org/x/net/context"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	crdv1alpha1 "k8s.tars.io/crd/v1alpha1"
	"runtime"
	"tarsadmin/handler/util"
	"tarsadmin/openapi/models"
	"tarsadmin/openapi/restapi/operations/server"
	"tarsadmin/openapi/restapi/operations/server_option"
)

type SelectServerOptionHandler struct {}

func (s *SelectServerOptionHandler) Handle(params server_option.SelectServerOptionParams) middleware.Responder {

	// fetch list
	namespace := K8sOption.Namespace

	// parse query
	selectParams, err := util.ParseSelectQuery(params.Filter, params.Limiter, nil)
	if err != nil {
		return server_option.NewSelectServerOptionInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	listAll := true
	if selectParams.Filter != nil && selectParams.Filter.Eq != nil {
		_, idOk := selectParams.Filter.Eq["ServerId"]
		if idOk {
			listAll = false
		}
	}

	filterItems := make([]*crdv1alpha1.TServer, 0, 10)
	if listAll {
		requirements := BuildSubTypeTarsSelector()
		list, err := K8sWatcher.tServerLister.TServers(K8sOption.Namespace).List(labels.NewSelector().Add(requirements ...))
		if err != nil {
			return server_option.NewSelectServerOptionInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
		filterItems = list
	} else {
		serverId, _ := selectParams.Filter.Eq["ServerId"]
		item, err := K8sWatcher.tServerLister.TServers(namespace).Get(util.GetTServerName(serverId.(string)))
		if err != nil {
			return server_option.NewSelectServerOptionInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
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
		elem["ServerImportant"] = item.Spec.Important
		elem["ServerTemplate"] = item.Spec.Tars.Template
		elem["ServerProfile"] = item.Spec.Tars.Profile
		elem["AsyncThread"] = item.Spec.Tars.AsyncThread
		result.Data = append(result.Data, elem)
	}

	return server.NewSelectServerOK().WithPayload(result)
}

type UpdateServerOptionHandler struct {}

func (s *UpdateServerOptionHandler) Handle(params server_option.UpdateServerOptionParams) middleware.Responder {

	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	tServer, err := K8sWatcher.tServerLister.TServers(namespace).Get(util.GetTServerName(*metadata.ServerID))
	if err != nil {
		return server_option.NewUpdateServerOptionInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	if equalServerOption(tServer.Spec.Tars, tServer.Spec.Important, params.Params.Target) {
		return server_option.NewUpdateServerOptionInternalServerError().WithPayload(&models.Error{Code: -1, Message: "No Need To Update Duplicated Option. "})
	}

	tServerCopy := tServer.DeepCopy()
	tServerCopy.Spec.Important = *params.Params.Target.ServerImportant
	tServerCopy.Spec.Tars.Template = params.Params.Target.ServerTemplate
	tServerCopy.Spec.Tars.Profile = params.Params.Target.ServerProfile
	tServerCopy.Spec.Tars.AsyncThread = *params.Params.Target.AsyncThread

	tServerInterface := K8sOption.CrdClientSet.CrdV1alpha1().TServers(namespace)
	if _, err = tServerInterface.Update(context.TODO(), tServerCopy, k8sMetaV1.UpdateOptions{}); err != nil {
		return server_option.NewUpdateServerOptionInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return server_option.NewUpdateServerOptionOK().WithPayload(&server_option.UpdateServerOptionOKBody{Result: 0})
}

type DoPreviewTemplateContentHandler struct {}

func (s *DoPreviewTemplateContentHandler) Handle(params server_option.DoPreviewTemplateContentParams) middleware.Responder {
	namespace := K8sOption.Namespace

	tServer, err := K8sWatcher.tServerLister.TServers(namespace).Get(util.GetTServerName(*params.ServerID))
	if err != nil {
		return server_option.NewDoPreviewTemplateContentInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	profile := []byte(tServer.Spec.Tars.Profile)
	templateName := tServer.Spec.Tars.Template

	allTemplateContent := make([][]byte, 0, 10)
	allTemplateContent = append(allTemplateContent, profile)

	for  {
		curTemplate, err := K8sWatcher.tTemplateLister.TTemplates(namespace).Get(templateName)
		if err != nil {
			return server_option.NewDoPreviewTemplateContentInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
		templateContent := curTemplate.Spec.Content

		parentTemplateName := curTemplate.Spec.Parent

		parTemplate, err := K8sWatcher.tTemplateLister.TTemplates(namespace).Get(parentTemplateName)
		if err != nil {
			return server_option.NewDoPreviewTemplateContentInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
		parentContent := parTemplate.Spec.Content

		if len(allTemplateContent) == 1 {
			allTemplateContent = append(allTemplateContent, []byte(templateContent))
		}

		if parentTemplateName == templateName {
			break
		}

		allTemplateContent = append(allTemplateContent, []byte(parentContent))

		// 目前只有一个递归父模板，同时防止测试阶段随意填导致死循环
		break
	}

	reverseSliceFun := func(s [][]byte) [][]byte {
		for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
			s[i], s[j] = s[j], s[i]
		}
		return s
	}
	allTemplateContent = reverseSliceFun(allTemplateContent)
	conf := tarsConf.New()
	afterJoinTemplateContent := bytes.Join(allTemplateContent, nil)

	if err := conf.InitFromBytes(afterJoinTemplateContent); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))
		return server_option.NewDoPreviewTemplateContentInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return server_option.NewDoPreviewTemplateContentOK().WithPayload(&server_option.DoPreviewTemplateContentOKBody{Result: conf.ToString()})
}

func equalServerOption(serverTars *crdv1alpha1.TServerTars, serverImportant int32, newOption *models.ServerOption) bool {
	if serverTars.Template != newOption.ServerTemplate {
		return false
	}
	if serverTars.Profile != newOption.ServerProfile {
		return false
	}
	if serverImportant != *newOption.ServerImportant {
		return false
	}
	if serverTars.AsyncThread != *newOption.AsyncThread {
		return false
	}
	return true
}
