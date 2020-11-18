package k8s

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"golang.org/x/net/context"
	k8sCoreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	crdv1alpha1 "k8s.taf.io/crd/v1alpha1"
	"tafadmin/handler/util"
	"tafadmin/openapi/models"
	"tafadmin/openapi/restapi/operations/applications"
	"tafadmin/openapi/restapi/operations/server"
)

type SelectServerHandler struct {}

func (s *SelectServerHandler) Handle(params server.SelectServerParams) middleware.Responder {

	// fetch list
	namespace := K8sOption.Namespace

	// parse query
	selectParams, err := util.ParseSelectQuery(params.Filter, params.Limiter, nil)
	if err != nil {
		return server.NewSelectServerInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	listAll := true
	if selectParams.Filter != nil && selectParams.Filter.Eq != nil {
		_, appOk := selectParams.Filter.Eq["ServerApp"]
		_, nameOk := selectParams.Filter.Eq["ServerName"]
		if appOk || nameOk {
			listAll = false
		}
	}

	allItems := make([]*crdv1alpha1.TServer, 0, 10)
	if listAll {
		requirements := BuildSubTypeTafSelector()
		allItems, err = K8sWatcher.tServerLister.TServers(K8sOption.Namespace).List(labels.NewSelector().Add(requirements ...))
		if err != nil {
			return server.NewSelectServerInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
	} else {
		requirements := BuildDoubleEqualSelector(selectParams.Filter, KeyLabel)
		allItems, err = K8sWatcher.tServerLister.TServers(namespace).List(labels.NewSelector().Add(requirements ...))
		if err != nil {
			return server.NewSelectServerInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
	}

	// filter
	filterItems := allItems

	// Admin临时版本，TServer特化已有web实现
	if selectParams.Filter != nil {
		filterItems = make([]*crdv1alpha1.TServer, 0, len(allItems))
		for _, elem := range allItems {
			if selectParams.Filter.Ne != nil {
				return server.NewSelectServerInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Ne Is Not Supported."})
			}
			if selectParams.Filter.Like != nil {
				m1, err := LikeMatch("ServerApp", selectParams.Filter.Like, elem.Spec.App)
				if err != nil {
					return server.NewSelectServerInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
				}
				m2, err := LikeMatch("ServerName", selectParams.Filter.Like, elem.Spec.Server)
				if err != nil {
					return server.NewSelectServerInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
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
		elem["ServerId"] = util.GetServerId(item.Spec.App, item.Spec.Server)
		elem["ServerApp"] = item.Spec.App
		elem["ServerName"] = item.Spec.Server
		if item.Spec.Release != nil {
			elem["ServerType"] = item.Spec.Release.ServerType
		}
		result.Data = append(result.Data, elem)
	}

	return server.NewSelectServerOK().WithPayload(result)
}

type UpdateServerHandler struct {}

func (s *UpdateServerHandler) Handle(params server.UpdateServerParams) middleware.Responder {

	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	tServer, err := K8sWatcher.tServerLister.TServers(namespace).Get(util.GetTServerName(*metadata.ServerID))
	if err != nil {
		return server.NewUpdateServerInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	tServerCopy := tServer.DeepCopy()
	if tServerCopy.Spec.Release == nil {
		tServerCopy.Spec.Release = &crdv1alpha1.TServerRelease{}
	}

	target := params.Params.Target
	if tServerCopy.Spec.Release.ServerType == target.ServerType {
		// 服务发布时会调用，默认一定是传相同的
		return server.NewUpdateServerOK().WithPayload(&server.UpdateServerOKBody{Result: 0})
	}

	tServerCopy.Spec.Release.ServerType = target.ServerType

	tServerInterface := K8sOption.CrdClientSet.CrdV1alpha1().TServers(namespace)
	if _, err = tServerInterface.Update(context.TODO(), tServerCopy, k8sMetaV1.UpdateOptions{}); err != nil {
		return server.NewUpdateServerInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return server.NewUpdateServerOK().WithPayload(&server.UpdateServerOKBody{Result: 0})
}

type DeleteServerHandler struct {}

func (s *DeleteServerHandler) Handle(params server.DeleteServerParams) middleware.Responder {
	metadata := params.Params.Metadata

	for _, serverId := range metadata.ServerID {
		if err := deleteServer(util.GetTServerName(serverId)); err != nil {
			fmt.Println(fmt.Sprintf("serverId: %s delete error: %s", serverId, err.Error()))
			continue
		}
	}

	return server.NewDeleteServerOK().WithPayload(&server.DeleteServerOKBody{Result: 0})
}

func CreateServer(serverApp, serverName string, serverServant models.MapServant, serverK8S *models.ServerK8S, serverOption *models.ServerOption) error {
	serverId := util.GetServerId(serverApp, serverName)

	namespace := K8sOption.Namespace
	_, err := K8sWatcher.tServerLister.TServers(namespace).Get(util.GetTServerName(serverId))
	if err == nil {
		return fmt.Errorf("Dumplicated %s Already Existed. ", serverId)
	}

	if !errors.IsNotFound(err) {
		return err
	}

	// 检测app是否存在，如果不存在需要加入TTree
	createAppHandler := CreateAppHandler{}
	createAppParams := applications.CreateAppParams{Params: applications.CreateAppBody{Metadata: &applications.CreateAppParamsBodyMetadata{AppName: &serverApp}}}
	_ = createAppHandler.Handle(createAppParams)

	tServer := buildTServer(namespace, serverApp, serverName, serverServant, serverK8S, serverOption)
	tServerInterface := K8sOption.CrdClientSet.CrdV1alpha1().TServers(namespace)
	if _, err = tServerInterface.Create(context.TODO(), tServer, k8sMetaV1.CreateOptions{}); err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

func buildTServer(namespace, serverApp, serverName string, serverServant models.MapServant, serverK8S *models.ServerK8S, serverOption *models.ServerOption) *crdv1alpha1.TServer {
	serverId := util.GetServerId(serverApp, serverName)

	var Servants []crdv1alpha1.TServant
	if serverServant != nil {
		Servants = make([]crdv1alpha1.TServant, 0, len(serverServant))
		for _, obj := range serverServant {
			Servants = append(Servants, crdv1alpha1.TServant{
				Name: obj.Name,
				Port: obj.Port,
				Thread: obj.Threads,
				Connection: obj.Connections,
				Capacity: obj.Capacity,
				IsTaf: *obj.IsTaf,
				IsTcp: *obj.IsTCP,
				Timeout: obj.Timeout,
			})
		}
	}

	Env := []k8sCoreV1.EnvVar{
		{
			Name: "Namespace",
			ValueFrom: &k8sCoreV1.EnvVarSource{
				FieldRef: &k8sCoreV1.ObjectFieldSelector {
					FieldPath: "metadata.namespace",
				},
			},
		},
		{
			Name: "PodName",
			ValueFrom: &k8sCoreV1.EnvVarSource{
				FieldRef: &k8sCoreV1.ObjectFieldSelector {
					FieldPath: "metadata.name",
				},
			},
		},
		{
			Name: "PodIP",
			ValueFrom: &k8sCoreV1.EnvVarSource{
				FieldRef: &k8sCoreV1.ObjectFieldSelector {
					FieldPath: "status.podIP",
				},
			},
		},
		{
			Name: "ServerApp",
			ValueFrom: &k8sCoreV1.EnvVarSource{
				FieldRef: &k8sCoreV1.ObjectFieldSelector {
					APIVersion: "v1",
					FieldPath: "metadata.labels['taf.io/ServerApp']",
				},
			},
		},
	}

	var HostPorts []crdv1alpha1.TK8SHostPort
	if serverK8S.HostPort != nil {
		HostPorts = make([]crdv1alpha1.TK8SHostPort, len(serverK8S.HostPort))
		for i, hostPort := range serverK8S.HostPort {
			HostPorts[i].NameRef = hostPort.NameRef
			HostPorts[i].Port = hostPort.Port
		}
	}

	HostPathType := k8sCoreV1.HostPathDirectoryOrCreate
	Mounts := []crdv1alpha1.TK8SMount{
		{
			Name: "host-log-dir",
			MountPath: "/usr/local/app/taf/app_log",
			SubPathExpr: "$(Namespace).$(PodName)",
			Source: k8sCoreV1.VolumeSource{
				HostPath: &k8sCoreV1.HostPathVolumeSource {
					Path: "/usr/local/app/taf/app_log",
					Type: &HostPathType,
				},
			},
		},
	}

	var NodeSelector crdv1alpha1.TK8SNodeSelector
	if serverK8S.NodeSelector != nil {
		if serverK8S.NodeSelector.NodeBind != nil {
			NodeSelector.NodeBind = &crdv1alpha1.TK8SNodeSelectorKind{}
			NodeSelector.NodeBind.Values = serverK8S.NodeSelector.NodeBind.Value
		}
		if serverK8S.NodeSelector.AbilityPool != nil {
			NodeSelector.AbilityPool = &crdv1alpha1.TK8SNodeSelectorKind{}
			NodeSelector.AbilityPool.Values = serverK8S.NodeSelector.AbilityPool.Value
		}
		if serverK8S.NodeSelector.PublicPool != nil {
			NodeSelector.PublicPool = &crdv1alpha1.TK8SNodeSelectorKind{}
			NodeSelector.PublicPool.Values = serverK8S.NodeSelector.PublicPool.Value
		}
		if serverK8S.NodeSelector.DaemonSet != nil {
			NodeSelector.DaemonSet = &crdv1alpha1.TK8SNodeSelectorKind{}
			NodeSelector.DaemonSet.Values = serverK8S.NodeSelector.DaemonSet.Value
		}
	}

	tServer := &crdv1alpha1.TServer{
		ObjectMeta: k8sMetaV1.ObjectMeta{
			Name:      util.GetTServerName(serverId),
			Namespace: namespace,
		},
		Spec: crdv1alpha1.TServerSpec {
			App: serverApp,
			Server: serverName,
			SubType: crdv1alpha1.TServerSubType(*serverOption.ServerSubType),
			Important: *serverOption.ServerImportant,
			Taf: &crdv1alpha1.TServerTaf{
				Template: serverOption.ServerTemplate,
				Profile: serverOption.ServerProfile,
				Foreground: false,
				AsyncThread: *serverOption.AsyncThread,
				Servants: Servants,
			},
			K8S: crdv1alpha1.TServerK8S{
				ServiceAccount: "",
				Env: Env,
				HostIPC: serverK8S.HostIpc,
				HostNetwork: serverK8S.HostNetwork,
				HostPorts: HostPorts,
				Mounts: Mounts,
				NodeSelector: NodeSelector,
				NotStacked: serverK8S.NotStacked,
				Replicas: 0,
			},
		},
	}

	return tServer
}

func deleteServer(serverId string) error {
	namespace := K8sOption.Namespace
	_, err := K8sWatcher.tServerLister.TServers(namespace).Get(serverId)
	if err != nil {
		return fmt.Errorf("%s Do Not Existed. ", serverId)
	}

	tServerInterface := K8sOption.CrdClientSet.CrdV1alpha1().TServers(namespace)
	if err = tServerInterface.Delete(context.TODO(), serverId, k8sMetaV1.DeleteOptions{}); err != nil {
		return err
	}

	tReleaseInterface := K8sOption.CrdClientSet.CrdV1alpha1().TReleases(namespace)
	return tReleaseInterface.Delete(context.TODO(), serverId, k8sMetaV1.DeleteOptions{})
}

