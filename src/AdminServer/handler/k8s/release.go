package k8s

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/api/errors"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	crdv1alpha1 "k8s.tars.io/crd/v1alpha1"
	"strconv"
	"strings"
	"tarsadmin/handler/util"
	"tarsadmin/openapi/models"
	"tarsadmin/openapi/restapi/operations/release"
)

type SelectServiceEnabledHandler struct {}

func (s *SelectServiceEnabledHandler) Handle(params release.SelectServiceEnabledParams) middleware.Responder {

	// fetch list
	namespace := K8sOption.Namespace

	// parse query
	selectParams, err := util.ParseSelectQuery(params.Filter, params.Limiter, nil)
	if err != nil {
		return release.NewSelectServiceEnabledInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
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
		requirements := BuildSubTypeTarsSelector()
		list, err := K8sWatcher.tServerLister.TServers(K8sOption.Namespace).List(labels.NewSelector().Add(requirements ...))
		if err != nil {
			return release.NewSelectServiceEnabledInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
		allServerItems = list
	} else {
		serverId, _ := selectParams.Filter.Eq["ServerId"]
		tempItems, err := K8sWatcher.tServerLister.TServers(namespace).Get(util.GetTServerName(serverId.(string)))
		if err != nil {
			return release.NewSelectServiceEnabledInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
		allServerItems = append(allServerItems, tempItems)
	}

	filterItems := make([]*crdv1alpha1.TServer, 0, len(allServerItems))
	for _, Server := range allServerItems {
		if Server.Spec.Release == nil {
			continue
		}
		filterItems = append(filterItems, Server)
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
		elem["ServiceId"] = getServiceId(item.Spec.Release.Image, item.Spec.Release.Tag)
		elem["ServerApp"] = item.Spec.App
		elem["ServerName"] = item.Spec.Server
		elem["ServiceVersion"] = item.Spec.Release.Tag
		elem["ServiceImage"] = item.Spec.Release.Image
		elem["EnablePerson"] = item.Spec.Release.ActivePerson
		elem["EnableTime"] = item.Spec.Release.ActiveTime
		elem["EnableMark"] = item.Spec.Release.ActiveReason
		result.Data = append(result.Data, elem)
	}

	return release.NewSelectServiceEnabledOK().WithPayload(result)
}


type SelectServicePoolHandler struct {}

func (s *SelectServicePoolHandler) Handle(params release.SelectServicePoolParams) middleware.Responder {

	// fetch list
	namespace := K8sOption.Namespace

	// parse query
	selectParams, err := util.ParseSelectQuery(params.Filter, params.Limiter, nil)
	if err != nil {
		return release.NewSelectServicePoolInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	listAll := true
	if selectParams.Filter != nil && selectParams.Filter.Eq != nil {
		_, appOk := selectParams.Filter.Eq["ServerApp"]
		_, nameOk := selectParams.Filter.Eq["ServerName"]
		if appOk && nameOk {
			listAll = false
		}
	}

	filterItems := make([]*crdv1alpha1.TRelease, 0, 10)
	if listAll {
		return release.NewSelectServicePoolInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Invalid Select Query Request."})
	} else {
		serverApp, _ := selectParams.Filter.Eq["ServerApp"]
		serverName, _ := selectParams.Filter.Eq["ServerName"]

		serverId := util.GetServerId(serverApp.(string), serverName.(string))
		tServer, err := K8sWatcher.tServerLister.TServers(namespace).Get(util.GetTServerName(serverId))
		if err != nil {
			return release.NewSelectServicePoolInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}

		tRelease, err := K8sWatcher.tReleaseLister.TReleases(namespace).Get(tServer.Name)
		if err != nil {
			if errors.IsNotFound(err) {
				return release.NewSelectServicePoolOK().WithPayload(&models.SelectResult{
					Count: models.MapInt {
						"AllCount": 0,
						"FilterCount": 0,
					},
					Data: make(models.ArrayMapInterface, 0),
				})
			} else {
				return release.NewSelectServicePoolInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
			}
		}
		filterItems = append(filterItems, tRelease)
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
		for _, reVersion := range item.Spec.List {
			elem := make(map[string]interface{})
			if !listAll {
				elem["ServerApp"] = selectParams.Filter.Eq["ServerApp"]
				elem["ServerName"] = selectParams.Filter.Eq["ServerName"]
				elem["ServerId"] = util.GetServerId(elem["ServerApp"].(string), elem["ServerName"].(string))
			}
			elem["ServiceId"] = getServiceId(reVersion.Image, reVersion.Tag)
			elem["ServiceVersion"] = reVersion.Tag
			elem["ServiceImage"] = reVersion.Image
			elem["CreateTime"] = reVersion.CreateTime
			result.Data = append(result.Data, elem)
		}
	}

	return release.NewSelectServicePoolOK().WithPayload(result)
}

type CreateServicePoolHandler struct {}

func (s *CreateServicePoolHandler) Handle(params release.CreateServicePoolParams) middleware.Responder {

	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	tServer, err := K8sWatcher.tServerLister.TServers(namespace).Get(util.GetTServerName(*metadata.ServerID))
	if err != nil {
		return release.NewCreateServicePoolInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}
	if tServer.Spec.SubType != crdv1alpha1.TARS {
		return release.NewCreateServicePoolInternalServerError().WithPayload(&models.Error{Code: -1, Message: "NonTars Is Not Supported."})
	}

	newTServerRelease := &crdv1alpha1.TServerRelease{
		Source: tServer.Name,
		ServerType: *metadata.ServerType,
		Image: *metadata.ServiceImage,
		ImagePullSecret: "tars-image-secret",
		ActivePerson: metadata.ActivePerson,
		ActiveReason: metadata.ActiveReason,
		ActiveTime: k8sMetaV1.Now(),
	}

	// 更新TRelease，Tars服务一一对应
	bCreate := false
	releaseId := tServer.Name
	tRelease, err := K8sWatcher.tReleaseLister.TReleases(namespace).Get(util.GetTServerName(releaseId))
	if err != nil {
		if !errors.IsNotFound(err) {
			return release.NewCreateServicePoolInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
		tRelease = &crdv1alpha1.TRelease{
			ObjectMeta: k8sMetaV1.ObjectMeta{
				Name: releaseId,
				Namespace: namespace,
			},
			Spec: crdv1alpha1.TReleaseSpec{
				List: make([]*crdv1alpha1.TReleaseVersion, 0, 1)},
		}
		bCreate = true
	}

	// 默认tag从10000开始
	maxTag := 9999
	if !bCreate {
		for _, version := range tRelease.Spec.List {
			tag, _ := strconv.Atoi(version.Tag)
			if tag > maxTag {
				maxTag = tag
			}
		}
	}
	newTServerRelease.Tag = strconv.Itoa(maxTag+1)

	tReleaseInterface := K8sOption.CrdClientSet.CrdV1alpha1().TReleases(namespace)
	if bCreate {
		tRelease.Spec.List = append([]*crdv1alpha1.TReleaseVersion{buildTReleaseVersion(newTServerRelease)}, tRelease.Spec.List ...)
		if _, err = tReleaseInterface.Create(context.TODO(), tRelease, k8sMetaV1.CreateOptions{}); err != nil && !errors.IsAlreadyExists(err)  {
			return release.NewCreateServicePoolInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
	} else {
		tReleaseCopy := tRelease.DeepCopy()
		tReleaseCopy.Spec.List = append([]*crdv1alpha1.TReleaseVersion{buildTReleaseVersion(newTServerRelease)}, tReleaseCopy.Spec.List ...)
		if _, err = tReleaseInterface.Update(context.TODO(), tReleaseCopy, k8sMetaV1.UpdateOptions{}); err != nil {
			return release.NewCreateServicePoolInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
	}

	return release.NewCreateServicePoolOK().WithPayload(&release.CreateServicePoolOKBody{Result: 0})
}

type DoEnableServiceHandler struct {}

func (s *DoEnableServiceHandler) Handle(params release.DoEnableServiceParams) middleware.Responder {
	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	tServer, err := K8sWatcher.tServerLister.TServers(namespace).Get(util.GetTServerName(*metadata.ServerID))
	if err != nil {
		return release.NewDoEnableServiceInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}
	if tServer.Spec.SubType != crdv1alpha1.TARS {
		return release.NewDoEnableServiceInternalServerError().WithPayload(&models.Error{Code: -1, Message: "NonTars Is Not Supported."})
	}

	// TARS服务的TRelease和TServer一一对应
	tRelease, err := K8sWatcher.tReleaseLister.TReleases(namespace).Get(util.GetTServerName(tServer.Name))
	if err != nil {
		return release.NewDoEnableServiceInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	serviceField := strings.Split(*metadata.ServiceID, "|")
	if len(serviceField) != 2 {
		return release.NewDoEnableServiceInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Invalid ServiceId"})
	}

	index := -1
	for i, reVersion := range tRelease.Spec.List {
		if reVersion.Image == serviceField[0] && reVersion.Tag == serviceField[1] {
			index = i
			break
		}
	}
	if index == -1 {
		return release.NewDoEnableServiceInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Can Not Find The Release Version."})
	}

	readyActiveRelease := tRelease.Spec.List[index]

	tServerCopy := tServer.DeepCopy()
	if tServerCopy.Spec.Release == nil {
		tServerCopy.Spec.Release = &crdv1alpha1.TServerRelease{}
	}
	tServerCopy.Spec.Release.Source = tRelease.Name
	tServerCopy.Spec.K8S.Replicas = *metadata.Replicas
	tServerCopy.Spec.Release.Image = readyActiveRelease.Image
	tServerCopy.Spec.Release.Tag = readyActiveRelease.Tag
	tServerCopy.Spec.Release.ImagePullSecret = readyActiveRelease.ImagePullSecret
	tServerCopy.Spec.Release.ServerType = readyActiveRelease.ServerType
	tServerCopy.Spec.Release.ActiveReason = metadata.EnableMark
	tServerCopy.Spec.Release.ActiveTime = k8sMetaV1.Now()

	// 更新TServer
	tServerInterface := K8sOption.CrdClientSet.CrdV1alpha1().TServers(namespace)
	if _, err = tServerInterface.Update(context.TODO(), tServerCopy, k8sMetaV1.UpdateOptions{}); err != nil {
		return release.NewDoEnableServiceInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	// 更新TRelease
	if index != -1 {
		tempList := append(tRelease.Spec.List[0:index], tRelease.Spec.List[index+1:] ...)

		tReleaseCopy := tRelease.DeepCopy()
		tReleaseCopy.Spec.List = append([]*crdv1alpha1.TReleaseVersion{readyActiveRelease}, tempList ...)

		tReleaseInterface := K8sOption.CrdClientSet.CrdV1alpha1().TReleases(namespace)
		if _, err = tReleaseInterface.Update(context.TODO(), tReleaseCopy, k8sMetaV1.UpdateOptions{}); err != nil {
			return release.NewDoEnableServiceInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
	}

	return release.NewDoEnableServiceOK().WithPayload(&release.DoEnableServiceOKBody{Result: 0})
}

func buildTReleaseVersion(tServerRelease *crdv1alpha1.TServerRelease) *crdv1alpha1.TReleaseVersion {
	return &crdv1alpha1.TReleaseVersion{
		ServerType: tServerRelease.ServerType,
		Image: tServerRelease.Image,
		Tag: tServerRelease.Tag,
		ImagePullSecret: tServerRelease.ImagePullSecret,
		CreatePerson: tServerRelease.ActivePerson,
		CreateTime: tServerRelease.ActiveTime,
	}
}

func getServiceId(serviceImage, serviceVersion string) string {
	return fmt.Sprintf("%s|%s", serviceImage, serviceVersion)
}